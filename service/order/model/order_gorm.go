package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/syncx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// GormOrder represents the order model for GORM
type GormOrder struct {
	ID         int64     `gorm:"primaryKey;column:id;autoIncrement"`
	Oid        string    `gorm:"column:oid;type:varchar(32);uniqueIndex:uk_oid"`
	Uid        int64     `gorm:"column:uid;index"`
	Pid        int64     `gorm:"column:pid;index"`
	Amount     int64     `gorm:"column:amount"`
	Status     int64     `gorm:"column:status"`
	CreateTime time.Time `gorm:"column:create_time;type:datetime;not null;default:CURRENT_TIMESTAMP"`
	UpdateTime time.Time `gorm:"column:update_time;type:datetime;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

// TableName specifies the table name for GORM
func (GormOrder) TableName() string {
	return "order"
}

// GormOrderModel represents the GORM implementation of OrderModel
type GormOrderModel struct {
	db    *gorm.DB
	cache cache.Cache
}

// NewGormOrderModel creates a new instance of GormOrderModel
func NewGormOrderModel(sqlDB *sql.DB, c cache.CacheConf) (*GormOrderModel, error) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
		DefaultDatetimePrecision:  nil,
		DisableDatetimePrecision:  true,
	}), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
		CreateBatchSize: 1000,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize GORM: %v", err)
	}

	sqlDB2, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying DB: %v", err)
	}

	// Set connection pool settings
	sqlDB2.SetMaxIdleConns(10)
	sqlDB2.SetMaxOpenConns(100)
	sqlDB2.SetConnMaxLifetime(time.Hour)

	// Auto migrate the schema
	err = db.AutoMigrate(&GormOrder{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate schema: %v", err)
	}

	cacheInstance := cache.New(c, syncx.NewSingleFlight(), cache.NewStat("order"), nil)
	return &GormOrderModel{
		db:    db,
		cache: cacheInstance,
	}, nil
}

// Insert adds a new order to the database
func (m *GormOrderModel) Insert(ctx context.Context, data *Order) (sql.Result, error) {
	gormOrder := &GormOrder{
		Uid:    data.Uid,
		Pid:    data.Pid,
		Amount: data.Amount,
		Status: data.Status,
		Oid:    data.Oid,
	}

	err := m.db.WithContext(ctx).Create(gormOrder).Error
	if err != nil {
		return nil, err
	}

	// Update the original order with the new ID and timestamps
	data.Id = gormOrder.ID
	data.Oid = gormOrder.Oid
	data.CreateTime = gormOrder.CreateTime
	data.UpdateTime = gormOrder.UpdateTime

	return &lastInsertIDResult{id: gormOrder.ID}, nil
}

// FindOne retrieves an order by ID
func (m *GormOrderModel) FindOne(ctx context.Context, id int64) (*Order, error) {
	var gormOrder GormOrder
	err := m.db.WithContext(ctx).First(&gormOrder, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, sqlx.ErrNotFound
		}
		return nil, err
	}

	order := &Order{
		Id:         gormOrder.ID,
		Uid:        gormOrder.Uid,
		Pid:        gormOrder.Pid,
		Amount:     gormOrder.Amount,
		Status:     gormOrder.Status,
		CreateTime: gormOrder.CreateTime,
		UpdateTime: gormOrder.UpdateTime,
	}

	return order, nil
}

// FindAllByUid retrieves all orders for a given user ID
func (m *GormOrderModel) FindAllByUid(ctx context.Context, uid int64) ([]*Order, error) {
	var gormOrders []GormOrder
	err := m.db.WithContext(ctx).Where("uid = ?", uid).Find(&gormOrders).Error
	if err != nil {
		return nil, err
	}

	orders := make([]*Order, len(gormOrders))
	for i, gormOrder := range gormOrders {
		orders[i] = &Order{
			Id:         gormOrder.ID,
			Uid:        gormOrder.Uid,
			Pid:        gormOrder.Pid,
			Amount:     gormOrder.Amount,
			Status:     gormOrder.Status,
			CreateTime: gormOrder.CreateTime,
			UpdateTime: gormOrder.UpdateTime,
		}
	}

	return orders, nil
}

// Update modifies an existing order
func (m *GormOrderModel) Update(ctx context.Context, data *Order) error {
	// 只更新需要修改的字段，不包括创建时间
	updates := map[string]interface{}{
		"uid":         data.Uid,
		"pid":         data.Pid,
		"amount":      data.Amount,
		"status":      data.Status,
		"update_time": time.Now(),
	}

	result := m.db.WithContext(ctx).Model(&GormOrder{}).Where("id = ?", data.Id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return sqlx.ErrNotFound
	}

	return nil
}

// Delete removes an order from the database
func (m *GormOrderModel) Delete(ctx context.Context, id int64) error {
	return m.db.WithContext(ctx).Delete(&GormOrder{}, id).Error
}

// lastInsertIDResult implements sql.Result interface
type lastInsertIDResult struct {
	id int64
}

func (r *lastInsertIDResult) LastInsertId() (int64, error) {
	return r.id, nil
}

func (r *lastInsertIDResult) RowsAffected() (int64, error) {
	return 1, nil
}

// FindPageListByPage 分页获取订单列表
func (m *GormOrderModel) FindPageListByPage(ctx context.Context, page, pageSize int64) ([]*Order, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	var gormOrders []*GormOrder
	var total int64

	// 查询总数
	if err := m.db.Table("order").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询数据
	if err := m.db.Table("order").Offset(int(offset)).Limit(int(pageSize)).Find(&gormOrders).Error; err != nil {
		return nil, 0, err
	}

	// 转换为Order类型
	var orders []*Order
	for _, gorder := range gormOrders {
		orders = append(orders, &Order{
			Id:         gorder.ID,
			Uid:        gorder.Uid,
			Pid:        gorder.Pid,
			Amount:     gorder.Amount,
			Status:     gorder.Status,
			CreateTime: gorder.CreateTime,
			UpdateTime: gorder.UpdateTime,
		})
	}

	return orders, total, nil
}
