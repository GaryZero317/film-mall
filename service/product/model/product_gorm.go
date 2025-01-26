package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/syncx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// TableName specifies the table name for GORM
func (Product) TableName() string {
	return "product"
}

type defaultGormProductModel struct {
	db    *gorm.DB
	cache cache.Cache
}

func NewDefaultGormProductModel(sqlDB *sql.DB, c cache.CacheConf) (*defaultGormProductModel, error) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &defaultGormProductModel{
		db:    db,
		cache: cache.New(c, syncx.NewSingleFlight(), cache.NewStat("product"), nil),
	}, nil
}

func (m *defaultGormProductModel) Insert(ctx context.Context, data *Product) (sql.Result, error) {
	err := m.db.WithContext(ctx).Create(data).Error
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (m *defaultGormProductModel) FindOne(ctx context.Context, id int64) (*Product, error) {
	var product Product
	err := m.db.WithContext(ctx).First(&product, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, sqlx.ErrNotFound
		}
		return nil, err
	}
	return &product, nil
}

func (m *defaultGormProductModel) Update(ctx context.Context, data *Product) error {
	return m.db.WithContext(ctx).Save(data).Error
}

func (m *defaultGormProductModel) Delete(ctx context.Context, id int64) error {
	return m.db.WithContext(ctx).Delete(&Product{}, id).Error
}

func (m *defaultGormProductModel) FindPageListByPage(ctx context.Context, page, pageSize int64) ([]*Product, int64, error) {
	var total int64
	var products []*Product

	if err := m.db.WithContext(ctx).Model(&Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := m.db.WithContext(ctx).
		Model(&Product{}).
		Order("id DESC").
		Offset(int(offset)).
		Limit(int(pageSize)).
		Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (m *defaultGormProductModel) DecrStock(ctx context.Context, id int64) error {
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

func (m *defaultGormProductModel) Search(ctx context.Context, keyword string, page, pageSize int64) ([]*Product, int64, error) {
	var total int64
	var products []*Product

	query := m.db.WithContext(ctx).Model(&Product{}).
		Where("name LIKE ? OR `desc` LIKE ?", "%"+keyword+"%", "%"+keyword+"%")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(int(offset)).Limit(int(pageSize)).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}
