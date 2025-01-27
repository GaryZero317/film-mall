package model

import (
	"context"
	"database/sql"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ OrderItemModel = (*customOrderItemModel)(nil)

type (
	// OrderItemModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOrderItemModel.
	OrderItemModel interface {
		Insert(ctx context.Context, session Session, data *OrderItem) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*OrderItem, error)
		Update(ctx context.Context, data *OrderItem) error
		Delete(ctx context.Context, id int64) error
		FindByOrderId(ctx context.Context, orderId int64) ([]*OrderItem, error)
	}

	customOrderItemModel struct {
		*GormOrderItemModel
	}
)

// NewOrderItemModel returns a model for the database table.
func NewOrderItemModel(conn sqlx.SqlConn, c cache.CacheConf) OrderItemModel {
	// Get the raw database connection
	rawDB, err := conn.RawDB()
	if err != nil {
		panic(err)
	}

	gormModel, err := NewGormOrderItemModel(rawDB, c)
	if err != nil {
		panic(err)
	}
	return &customOrderItemModel{
		GormOrderItemModel: gormModel,
	}
}
