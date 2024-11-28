package model

import (
	"context"
	"database/sql"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ OrderModel = (*customOrderModel)(nil)

type (
	// OrderModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOrderModel.
	OrderModel interface {
		Insert(ctx context.Context, data *Order) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*Order, error)
		Update(ctx context.Context, data *Order) error
		Delete(ctx context.Context, id int64) error
		FindAllByUid(ctx context.Context, uid int64) ([]*Order, error)
	}

	customOrderModel struct {
		*GormOrderModel
	}
)

// NewOrderModel 返回数据库表的模型。
func NewOrderModel(conn sqlx.SqlConn, c cache.CacheConf) OrderModel {
	// 获取原始数据库连接
	rawDB, err := conn.RawDB()
	if err != nil {
		panic(err)
	}
	
	gormModel, err := NewGormOrderModel(rawDB, c)
	if err != nil {
		panic(err)
	}
	return &customOrderModel{
		GormOrderModel: gormModel,
	}
}

// FindAllByUid returns all orders for a given user ID
func (m *customOrderModel) FindAllByUid(ctx context.Context, uid int64) ([]*Order, error) {
	return m.GormOrderModel.FindAllByUid(ctx, uid)
}
