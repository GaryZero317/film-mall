package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/syncx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ProductGorm GORM模型
type ProductGorm struct {
	Id         int64          `gorm:"primarykey"`
	Name       string         `gorm:"column:name;not null;default:''" json:"name"`
	Desc       string         `gorm:"column:desc;not null;default:''" json:"desc"`
	Stock      int64          `gorm:"column:stock;not null;default:0" json:"stock"`
	Amount     int64          `gorm:"column:amount;not null;default:0" json:"amount"`
	Status     int64          `gorm:"column:status;not null;default:0" json:"status"`
	MainImage  string         `gorm:"column:main_image;not null;default:''" json:"mainImage"`
	CreateTime time.Time      `gorm:"column:create_time;autoCreateTime" json:"createTime"`
	UpdateTime time.Time      `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`
	Images     []ProductImage `gorm:"foreignKey:ProductId" json:"images,omitempty"`
}

// TableName 设置表名
func (m *ProductGorm) TableName() string {
	return "product"
}

type GormProductModel struct {
	db    *gorm.DB
	cache cache.Cache
}

func NewGormProductModel(sqlDB *sql.DB, c cache.CacheConf) (ProductModel, error) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize GORM: %v", err)
	}

	return &GormProductModel{
		db:    db,
		cache: cache.New(c, syncx.NewSingleFlight(), cache.NewStat("product"), nil),
	}, nil
}

func (m *GormProductModel) Insert(ctx context.Context, data *Product) (sql.Result, error) {
	err := m.db.WithContext(ctx).Create(data).Error
	if err != nil {
		return nil, err
	}
	return &lastInsertIDResult{id: data.Id}, nil
}

func (m *GormProductModel) FindOne(ctx context.Context, id int64) (*Product, error) {
	var product Product
	err := m.db.WithContext(ctx).First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (m *GormProductModel) Update(ctx context.Context, data *Product) error {
	return m.db.WithContext(ctx).Save(data).Error
}

func (m *GormProductModel) Delete(ctx context.Context, id int64) error {
	return m.db.WithContext(ctx).Delete(&Product{}, id).Error
}

func (m *GormProductModel) FindPageListByPage(ctx context.Context, page, pageSize int64) ([]*Product, int64, error) {
	var total int64
	var products []*Product

	if err := m.db.WithContext(ctx).Model(&Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := m.db.WithContext(ctx).Offset(int(offset)).Limit(int(pageSize)).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (m *GormProductModel) DecrStock(ctx context.Context, id int64) error {
	result := m.db.WithContext(ctx).Model(&Product{}).
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

func (m *GormProductModel) Search(ctx context.Context, keyword string, page, pageSize int64) ([]*Product, int64, error) {
	var total int64
	var products []*Product

	query := m.db.WithContext(ctx).Model(&Product{}).
		Where("name LIKE ? OR desc LIKE ?", "%"+keyword+"%", "%"+keyword+"%")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(int(offset)).Limit(int(pageSize)).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

type lastInsertIDResult struct {
	id int64
}

func (r *lastInsertIDResult) LastInsertId() (int64, error) {
	return r.id, nil
}

func (r *lastInsertIDResult) RowsAffected() (int64, error) {
	return 1, nil
}
