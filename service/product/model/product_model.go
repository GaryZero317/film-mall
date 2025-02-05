package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// Product 商品模型
type Product struct {
	Id         int64     `db:"id" gorm:"column:id"`
	Name       string    `db:"name" gorm:"column:name"`               // 商品名称
	Desc       string    `db:"desc" gorm:"column:desc"`               // 商品描述
	Stock      int64     `db:"stock" gorm:"column:stock"`             // 商品库存
	Amount     int64     `db:"amount" gorm:"column:amount"`           // 商品金额
	Status     int64     `db:"status" gorm:"column:status"`           // 商品状态
	CategoryId int64     `db:"category_id" gorm:"column:category_id"` // 分类ID
	CreateTime time.Time `db:"create_time" gorm:"column:create_time"` // 创建时间
	UpdateTime time.Time `db:"update_time" gorm:"column:update_time"` // 更新时间
	MainImage  string    `db:"-" gorm:"-"`                            // 商品主图
	Images     []string  `db:"-" gorm:"-"`                            // 商品图片列表
}

// ProductModel 定义商品模型接口
type ProductModel interface {
	Insert(ctx context.Context, data *Product) (sql.Result, error)
	FindOne(ctx context.Context, id int64) (*Product, error)
	Update(ctx context.Context, data *Product) error
	Delete(ctx context.Context, id int64) error
	FindPageListByPage(ctx context.Context, page, pageSize int64) ([]*Product, int64, error)
	DecrStock(ctx context.Context, id int64) error
	Search(ctx context.Context, keyword string, page, pageSize int64) ([]*Product, int64, error)
}

// NewProductModel 创建商品模型实例
func NewProductModel(conn sqlx.SqlConn, c cache.CacheConf) ProductModel {
	sqlDB, err := conn.RawDB()
	if err != nil {
		panic(err)
	}
	model, err := NewDefaultGormProductModel(sqlDB, c)
	if err != nil {
		panic(err)
	}
	return model
}
