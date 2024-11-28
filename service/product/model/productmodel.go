package model

import (
	"context"
	"database/sql"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ProductModel = (*customProductModel)(nil)

type (
	// ProductModel is an interface to be customized, add more methods here,
	// and implement the added methods in customProductModel.
	ProductModel interface {
		Insert(ctx context.Context, data *Product) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*Product, error)
		Update(ctx context.Context, data *Product) error
		Delete(ctx context.Context, id int64) error
	}

	customProductModel struct {
		*GormProductModel
	}
)

// NewProductModel 返回数据库表的模型。
func NewProductModel(conn sqlx.SqlConn, c cache.CacheConf) ProductModel {
	// 获取原始数据库连接
	rawDB, err := conn.RawDB()
	if err != nil {
		panic(err)
	}
	
	gormModel, err := NewGormProductModel(rawDB, c)
	if err != nil {
		panic(err)
	}
	return &customProductModel{
		GormProductModel: gormModel,
	}
}
