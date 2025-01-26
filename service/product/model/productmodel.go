package model

import (
	"context"
	"database/sql"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ProductModel interface {
	Insert(ctx context.Context, data *Product) (sql.Result, error)
	FindOne(ctx context.Context, id int64) (*Product, error)
	Update(ctx context.Context, data *Product) error
	Delete(ctx context.Context, id int64) error
	FindPageListByPage(ctx context.Context, page, pageSize int64) ([]*Product, int64, error)
	DecrStock(ctx context.Context, id int64) error
	Search(ctx context.Context, keyword string, page, pageSize int64) ([]*Product, int64, error)
}

type customProductModel struct {
	*defaultGormProductModel
}

// NewProductModel returns a model for the database table.
func NewProductModel(conn sqlx.SqlConn, c cache.CacheConf) ProductModel {
	rawDB, err := conn.RawDB()
	if err != nil {
		panic(err)
	}

	gormModel, err := NewDefaultGormProductModel(rawDB, c)
	if err != nil {
		panic(err)
	}
	return &customProductModel{
		defaultGormProductModel: gormModel,
	}
}

// FindPageListByPage 分页获取商品列表
func (m *customProductModel) FindPageListByPage(ctx context.Context, page, pageSize int64) ([]*Product, int64, error) {
	return m.defaultGormProductModel.FindPageListByPage(ctx, page, pageSize)
}

// DecrStock 减少库存
func (m *customProductModel) DecrStock(ctx context.Context, id int64) error {
	return m.defaultGormProductModel.DecrStock(ctx, id)
}

// Search 搜索商品
func (m *customProductModel) Search(ctx context.Context, keyword string, page, pageSize int64) ([]*Product, int64, error) {
	return m.defaultGormProductModel.Search(ctx, keyword, page, pageSize)
}
