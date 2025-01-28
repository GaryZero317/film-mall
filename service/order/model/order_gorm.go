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
	ID          int64     `gorm:"primaryKey;column:id;autoIncrement"`
	Oid         string    `gorm:"column:oid;type:varchar(32);uniqueIndex:uk_oid"`
	Uid         int64     `gorm:"column:uid;index"`
	AddressId   int64     `gorm:"column:address_id"`
	TotalPrice  int64     `gorm:"column:total_price"`
	ShippingFee int64     `gorm:"column:shipping_fee"`
	Status      int64     `gorm:"column:status"`
	StatusDesc  string    `gorm:"column:status_desc;type:varchar(32)"`
	Remark      string    `gorm:"column:remark;type:varchar(255)"`
	CreateTime  time.Time `gorm:"column:create_time;type:datetime;not null;default:CURRENT_TIMESTAMP"`
	UpdateTime  time.Time `gorm:"column:update_time;type:datetime;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
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
func (m *GormOrderModel) Insert(ctx context.Context, session Session, data *Order) (sql.Result, error) {
	gormOrder := &GormOrder{
		Uid:         data.Uid,
		AddressId:   data.AddressId,
		TotalPrice:  data.TotalPrice,
		ShippingFee: data.ShippingFee,
		Status:      data.Status,
		StatusDesc:  data.StatusDesc,
		Remark:      data.Remark,
		Oid:         data.Oid,
	}

	db := m.db
	if session != nil {
		db = session.(*GormSession).DB
	}

	err := db.WithContext(ctx).Create(gormOrder).Error
	if err != nil {
		return nil, err
	}

	// Update the original order with the new ID and timestamps
	data.Id = gormOrder.ID
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

	return convertToOrder(&gormOrder), nil
}

// Update modifies an existing order
func (m *GormOrderModel) Update(ctx context.Context, data *Order) error {
	// 先获取现有订单
	var existingOrder GormOrder
	if err := m.db.WithContext(ctx).First(&existingOrder, data.Id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return sqlx.ErrNotFound
		}
		return err
	}

	updates := make(map[string]interface{})

	// 只更新提供的字段，保持其他字段不变
	if data.Status != 0 {
		updates["status"] = data.Status
	}
	if data.StatusDesc != "" {
		updates["status_desc"] = data.StatusDesc
	}
	if data.Remark != "" {
		updates["remark"] = data.Remark
	}
	if data.Uid != 0 {
		updates["uid"] = data.Uid
	}
	if data.AddressId != 0 {
		updates["address_id"] = data.AddressId
	}
	if data.TotalPrice != 0 {
		updates["total_price"] = data.TotalPrice
	}
	if data.ShippingFee != 0 {
		updates["shipping_fee"] = data.ShippingFee
	}
	if data.Oid != "" {
		updates["oid"] = data.Oid
	}

	// 如果没有任何字段需要更新，直接返回
	if len(updates) == 0 {
		return nil
	}

	// 添加更新时间
	updates["update_time"] = time.Now()

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
	result := m.db.WithContext(ctx).Delete(&GormOrder{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return sqlx.ErrNotFound
	}

	return nil
}

// FindAllByUid retrieves all orders for a given user ID
func (m *GormOrderModel) FindAllByUid(ctx context.Context, uid int64) ([]*Order, error) {
	var gormOrders []GormOrder
	err := m.db.WithContext(ctx).Where("uid = ?", uid).Find(&gormOrders).Error
	if err != nil {
		return nil, err
	}

	orders := make([]*Order, len(gormOrders))
	for i := range gormOrders {
		orders[i] = convertToOrder(&gormOrders[i])
	}

	return orders, nil
}

// FindPageListByPage returns a page of orders
func (m *GormOrderModel) FindPageListByPage(ctx context.Context, page, pageSize int64) ([]*Order, int64, error) {
	var total int64
	err := m.db.WithContext(ctx).Model(&GormOrder{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	var gormOrders []GormOrder
	err = m.db.WithContext(ctx).Offset(int((page - 1) * pageSize)).Limit(int(pageSize)).Find(&gormOrders).Error
	if err != nil {
		return nil, 0, err
	}

	orders := make([]*Order, len(gormOrders))
	for i := range gormOrders {
		orders[i] = convertToOrder(&gormOrders[i])
	}

	return orders, total, nil
}

// FindByUid returns orders for a given user ID with pagination
func (m *GormOrderModel) FindByUid(ctx context.Context, uid, status, page, pageSize int64) ([]*Order, int64, error) {
	var total int64
	query := m.db.WithContext(ctx).Model(&GormOrder{}).Where("uid = ?", uid)

	if status != 0 {
		query = query.Where("status = ?", status)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	var gormOrders []GormOrder
	err = query.Offset(int((page - 1) * pageSize)).Limit(int(pageSize)).Find(&gormOrders).Error
	if err != nil {
		return nil, 0, err
	}

	orders := make([]*Order, len(gormOrders))
	for i := range gormOrders {
		orders[i] = convertToOrder(&gormOrders[i])
	}

	return orders, total, nil
}

// FindAll 查询所有订单
func (m *GormOrderModel) FindAll(ctx context.Context, status, page, pageSize int64) ([]*Order, int64, error) {
	var gormOrders []GormOrder
	var total int64

	query := m.db.WithContext(ctx).Model(&GormOrder{})

	// 如果指定了状态，则按状态筛选
	if status > 0 {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err = query.Order("create_time DESC").Offset(int((page - 1) * pageSize)).Limit(int(pageSize)).Find(&gormOrders).Error
	if err != nil {
		return nil, 0, err
	}

	orders := make([]*Order, len(gormOrders))
	for i := range gormOrders {
		orders[i] = convertToOrder(&gormOrders[i])
	}

	return orders, total, nil
}

// convertToOrder converts a GormOrder to an Order
func convertToOrder(gormOrder *GormOrder) *Order {
	return &Order{
		Id:          gormOrder.ID,
		Oid:         gormOrder.Oid,
		Uid:         gormOrder.Uid,
		AddressId:   gormOrder.AddressId,
		TotalPrice:  gormOrder.TotalPrice,
		ShippingFee: gormOrder.ShippingFee,
		Status:      gormOrder.Status,
		StatusDesc:  gormOrder.StatusDesc,
		Remark:      gormOrder.Remark,
		CreateTime:  gormOrder.CreateTime,
		UpdateTime:  gormOrder.UpdateTime,
	}
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

// GormSession wraps gorm.DB to implement Session interface
type GormSession struct {
	*gorm.DB
}

func (s *GormSession) Commit() error {
	return s.DB.Commit().Error
}

func (s *GormSession) Rollback() error {
	return s.DB.Rollback().Error
}

// Trans executes given function in a transaction
func (m *GormOrderModel) Trans(ctx context.Context, fn func(ctx context.Context, session Session) error) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(ctx, &GormSession{DB: tx})
	})
}
