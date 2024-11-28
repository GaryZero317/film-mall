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

// GormProduct represents the product model for GORM
type GormProduct struct {
	ID         int64     `gorm:"primaryKey;column:id;autoIncrement"`
	Name       string    `gorm:"column:name;type:varchar(255)"`
	Desc       string    `gorm:"column:desc;type:varchar(255)"`
	Stock      int64     `gorm:"column:stock"`
	Amount     int64     `gorm:"column:amount"`
	Status     int64     `gorm:"column:status"`
	CreateTime time.Time `gorm:"column:create_time;type:datetime;not null;default:CURRENT_TIMESTAMP"`
	UpdateTime time.Time `gorm:"column:update_time;type:datetime;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

// TableName specifies the table name for GORM
func (GormProduct) TableName() string {
	return "product"
}

// GormProductModel represents the GORM implementation of ProductModel
type GormProductModel struct {
	db    *gorm.DB
	cache cache.Cache
}

// NewGormProductModel creates a new instance of GormProductModel
func NewGormProductModel(sqlDB *sql.DB, c cache.CacheConf) (*GormProductModel, error) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
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
	err = db.AutoMigrate(&GormProduct{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate schema: %v", err)
	}

	cacheInstance := cache.New(c, syncx.NewSingleFlight(), cache.NewStat("product"), nil)
	return &GormProductModel{
		db:    db,
		cache: cacheInstance,
	}, nil
}

// Insert adds a new product to the database
func (m *GormProductModel) Insert(ctx context.Context, data *Product) (sql.Result, error) {
	gormProduct := &GormProduct{
		Name:   data.Name,
		Desc:   data.Desc,
		Stock:  data.Stock,
		Amount: data.Amount,
		Status: data.Status,
	}

	err := m.db.WithContext(ctx).Create(gormProduct).Error
	if err != nil {
		return nil, err
	}

	// Update the original product with the new ID and timestamps
	data.Id = gormProduct.ID
	data.CreateTime = gormProduct.CreateTime
	data.UpdateTime = gormProduct.UpdateTime

	return &lastInsertIDResult{id: gormProduct.ID}, nil
}

// FindOne retrieves a product by ID
func (m *GormProductModel) FindOne(ctx context.Context, id int64) (*Product, error) {
	var gormProduct GormProduct
	err := m.db.WithContext(ctx).First(&gormProduct, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, sqlx.ErrNotFound
		}
		return nil, err
	}

	product := &Product{
		Id:         gormProduct.ID,
		Name:       gormProduct.Name,
		Desc:       gormProduct.Desc,
		Stock:      gormProduct.Stock,
		Amount:     gormProduct.Amount,
		Status:     gormProduct.Status,
		CreateTime: gormProduct.CreateTime,
		UpdateTime: gormProduct.UpdateTime,
	}

	return product, nil
}

// Update modifies an existing product
func (m *GormProductModel) Update(ctx context.Context, data *Product) error {
	// 只更新需要修改的字段，不包括创建时间
	updates := map[string]interface{}{
		"name":        data.Name,
		"desc":        data.Desc,
		"stock":       data.Stock,
		"amount":      data.Amount,
		"status":      data.Status,
		"update_time": time.Now(),
	}

	result := m.db.WithContext(ctx).Model(&GormProduct{}).Where("id = ?", data.Id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return sqlx.ErrNotFound
	}

	return nil
}

// Delete removes a product from the database
func (m *GormProductModel) Delete(ctx context.Context, id int64) error {
	return m.db.WithContext(ctx).Delete(&GormProduct{}, id).Error
}

// DecrStock decrements the stock of a product
func (m *GormProductModel) DecrStock(ctx context.Context, id int64) error {
	result := m.db.WithContext(ctx).Model(&GormProduct{}).
		Where("id = ? AND stock > 0", id).
		UpdateColumn("stock", gorm.Expr("stock - ?", 1))
	
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return fmt.Errorf("product %d stock not available", id)
	}
	
	return nil
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
