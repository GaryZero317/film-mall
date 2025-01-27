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

// GormOrderItem represents the order item model for GORM
type GormOrderItem struct {
	ID           int64     `gorm:"primaryKey;column:id;autoIncrement"`
	OrderId      int64     `gorm:"column:order_id;index"`
	Pid          int64     `gorm:"column:pid;index"`
	ProductName  string    `gorm:"column:product_name;type:varchar(255)"`
	ProductImage string    `gorm:"column:product_image;type:varchar(255)"`
	Price        int64     `gorm:"column:price"`
	Quantity     int64     `gorm:"column:quantity"`
	Amount       int64     `gorm:"column:amount"`
	CreateTime   time.Time `gorm:"column:create_time;type:datetime;not null;default:CURRENT_TIMESTAMP"`
	UpdateTime   time.Time `gorm:"column:update_time;type:datetime;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

// TableName specifies the table name for GORM
func (GormOrderItem) TableName() string {
	return "order_item"
}

// GormOrderItemModel represents the GORM implementation of OrderItemModel
type GormOrderItemModel struct {
	db    *gorm.DB
	cache cache.Cache
}

// NewGormOrderItemModel creates a new instance of GormOrderItemModel
func NewGormOrderItemModel(sqlDB *sql.DB, c cache.CacheConf) (*GormOrderItemModel, error) {
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
	err = db.AutoMigrate(&GormOrderItem{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate schema: %v", err)
	}

	cacheInstance := cache.New(c, syncx.NewSingleFlight(), cache.NewStat("order_item"), nil)
	return &GormOrderItemModel{
		db:    db,
		cache: cacheInstance,
	}, nil
}

// Insert adds a new order item to the database
func (m *GormOrderItemModel) Insert(ctx context.Context, session Session, data *OrderItem) (sql.Result, error) {
	gormOrderItem := &GormOrderItem{
		OrderId:      data.OrderId,
		Pid:          data.Pid,
		ProductName:  data.ProductName,
		ProductImage: data.ProductImage,
		Price:        data.Price,
		Quantity:     data.Quantity,
		Amount:       data.Amount,
	}

	db := m.db
	if session != nil {
		db = session.(*GormSession).DB
	}

	err := db.WithContext(ctx).Create(gormOrderItem).Error
	if err != nil {
		return nil, err
	}

	// Update the original order item with the new ID and timestamps
	data.Id = gormOrderItem.ID

	return &lastInsertIDResult{id: gormOrderItem.ID}, nil
}

// FindOne retrieves an order item by ID
func (m *GormOrderItemModel) FindOne(ctx context.Context, id int64) (*OrderItem, error) {
	var gormOrderItem GormOrderItem
	err := m.db.WithContext(ctx).First(&gormOrderItem, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, sqlx.ErrNotFound
		}
		return nil, err
	}
	return convertToOrderItem(&gormOrderItem), nil
}

// Update modifies an existing order item
func (m *GormOrderItemModel) Update(ctx context.Context, data *OrderItem) error {
	result := m.db.WithContext(ctx).Model(&GormOrderItem{}).Where("id = ?", data.Id).Updates(map[string]interface{}{
		"order_id":      data.OrderId,
		"pid":           data.Pid,
		"product_name":  data.ProductName,
		"product_image": data.ProductImage,
		"price":         data.Price,
		"quantity":      data.Quantity,
		"amount":        data.Amount,
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return sqlx.ErrNotFound
	}
	return nil
}

// Delete removes an order item from the database
func (m *GormOrderItemModel) Delete(ctx context.Context, id int64) error {
	result := m.db.WithContext(ctx).Delete(&GormOrderItem{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return sqlx.ErrNotFound
	}
	return nil
}

// FindByOrderId retrieves all items for a given order ID
func (m *GormOrderItemModel) FindByOrderId(ctx context.Context, orderId int64) ([]*OrderItem, error) {
	var gormOrderItems []GormOrderItem
	err := m.db.WithContext(ctx).Where("order_id = ?", orderId).Find(&gormOrderItems).Error
	if err != nil {
		return nil, err
	}

	items := make([]*OrderItem, len(gormOrderItems))
	for i := range gormOrderItems {
		items[i] = convertToOrderItem(&gormOrderItems[i])
	}
	return items, nil
}

// convertToOrderItem converts a GormOrderItem to an OrderItem
func convertToOrderItem(gormOrderItem *GormOrderItem) *OrderItem {
	return &OrderItem{
		Id:           gormOrderItem.ID,
		OrderId:      gormOrderItem.OrderId,
		Pid:          gormOrderItem.Pid,
		ProductName:  gormOrderItem.ProductName,
		ProductImage: gormOrderItem.ProductImage,
		Price:        gormOrderItem.Price,
		Quantity:     gormOrderItem.Quantity,
		Amount:       gormOrderItem.Amount,
	}
}
