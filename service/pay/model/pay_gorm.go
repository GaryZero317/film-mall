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

// GormPay represents the pay model for GORM
type GormPay struct {
	ID         int64     `gorm:"primaryKey;column:id;autoIncrement"`
	Uid        int64     `gorm:"column:uid;index"`
	Oid        int64     `gorm:"column:oid;index"`
	Amount     int64     `gorm:"column:amount"`
	Source     int64     `gorm:"column:source"`
	Status     int64     `gorm:"column:status"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime"`
	UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime"`
}

// TableName specifies the table name for GORM
func (GormPay) TableName() string {
	return "pay"
}

// GormPayModel represents the GORM implementation of PayModel
type GormPayModel struct {
	db    *gorm.DB
	cache cache.Cache
}

// NewGormPayModel creates a new instance of GormPayModel
func NewGormPayModel(sqlDB *sql.DB, c cache.CacheConf) (*GormPayModel, error) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize GORM: %v", err)
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&GormPay{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate schema: %v", err)
	}

	cacheInstance := cache.New(c, syncx.NewSingleFlight(), cache.NewStat("pay"), nil)
	return &GormPayModel{
		db:    db,
		cache: cacheInstance,
	}, nil
}

// Insert adds a new payment to the database
func (m *GormPayModel) Insert(ctx context.Context, data *Pay) (sql.Result, error) {
	gormPay := &GormPay{
		Uid:    data.Uid,
		Oid:    data.Oid,
		Amount: data.Amount,
		Source: data.Source,
		Status: data.Status,
	}

	err := m.db.WithContext(ctx).Create(gormPay).Error
	if err != nil {
		return nil, err
	}

	// Update the original payment with the new ID and timestamps
	data.Id = gormPay.ID
	data.CreateTime = gormPay.CreateTime
	data.UpdateTime = gormPay.UpdateTime

	return &lastInsertIDResult{id: gormPay.ID}, nil
}

// FindOne retrieves a payment by ID
func (m *GormPayModel) FindOne(ctx context.Context, id int64) (*Pay, error) {
	var gormPay GormPay
	err := m.db.WithContext(ctx).First(&gormPay, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, sqlx.ErrNotFound
		}
		return nil, err
	}

	pay := &Pay{
		Id:         gormPay.ID,
		Uid:        gormPay.Uid,
		Oid:        gormPay.Oid,
		Amount:     gormPay.Amount,
		Source:     gormPay.Source,
		Status:     gormPay.Status,
		CreateTime: gormPay.CreateTime,
		UpdateTime: gormPay.UpdateTime,
	}

	return pay, nil
}

// FindOneByOid retrieves a payment by order ID
func (m *GormPayModel) FindOneByOid(ctx context.Context, oid int64) (*Pay, error) {
	var gormPay GormPay
	err := m.db.WithContext(ctx).Where("oid = ?", oid).First(&gormPay).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, sqlx.ErrNotFound
		}
		return nil, err
	}

	pay := &Pay{
		Id:         gormPay.ID,
		Uid:        gormPay.Uid,
		Oid:        gormPay.Oid,
		Amount:     gormPay.Amount,
		Source:     gormPay.Source,
		Status:     gormPay.Status,
		CreateTime: gormPay.CreateTime,
		UpdateTime: gormPay.UpdateTime,
	}

	return pay, nil
}

// Update modifies an existing payment
func (m *GormPayModel) Update(ctx context.Context, data *Pay) error {
	// 只更新需要修改的字段，不包括创建时间
	updates := map[string]interface{}{
		"uid":         data.Uid,
		"oid":         data.Oid,
		"amount":      data.Amount,
		"source":      data.Source,
		"status":      data.Status,
		"update_time": time.Now(),
	}

	result := m.db.WithContext(ctx).Model(&GormPay{}).Where("id = ?", data.Id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return sqlx.ErrNotFound
	}

	return nil
}

// Delete removes a payment from the database
func (m *GormPayModel) Delete(ctx context.Context, id int64) error {
	return m.db.WithContext(ctx).Delete(&GormPay{}, id).Error
}

// FindList retrieves a paginated list of payments
func (m *GormPayModel) FindList(ctx context.Context, page, pageSize int64) ([]*Pay, int64, error) {
	var gormPays []*GormPay
	var total int64

	// Get total count
	if err := m.db.WithContext(ctx).Model(&GormPay{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	err := m.db.WithContext(ctx).
		Offset(int(offset)).
		Limit(int(pageSize)).
		Order("id DESC").
		Find(&gormPays).Error
	if err != nil {
		return nil, 0, err
	}

	// Convert to Pay model
	var pays []*Pay
	for _, gormPay := range gormPays {
		pays = append(pays, &Pay{
			Id:         gormPay.ID,
			Uid:        gormPay.Uid,
			Oid:        gormPay.Oid,
			Amount:     gormPay.Amount,
			Source:     gormPay.Source,
			Status:     gormPay.Status,
			CreateTime: gormPay.CreateTime,
			UpdateTime: gormPay.UpdateTime,
		})
	}

	return pays, total, nil
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
