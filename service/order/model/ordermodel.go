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
		// base methods
		Insert(ctx context.Context, session Session, data *Order) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*Order, error)
		Update(ctx context.Context, data *Order) error
		Delete(ctx context.Context, id int64) error
		Trans(ctx context.Context, fn func(ctx context.Context, session Session) error) error

		// custom methods
		FindByUid(ctx context.Context, uid, status, page, pageSize int64) ([]*Order, int64, error)
		FindAllByUid(ctx context.Context, uid int64) ([]*Order, error)
		FindPageListByPage(ctx context.Context, page, pageSize int64) ([]*Order, int64, error)
		FindAll(ctx context.Context, status, page, pageSize int64) ([]*Order, int64, error)
	}

	Session interface {
		Commit() error
		Rollback() error
	}

	customOrderModel struct {
		*GormOrderModel
	}
)

// NewOrderModel returns a model for the database table.
func NewOrderModel(conn sqlx.SqlConn, c cache.CacheConf) OrderModel {
	// Get the raw database connection
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

// FindPageListByPage returns a page of orders
func (m *customOrderModel) FindPageListByPage(ctx context.Context, page, pageSize int64) ([]*Order, int64, error) {
	return m.GormOrderModel.FindPageListByPage(ctx, page, pageSize)
}

// FindAll returns all orders with pagination and status filter
func (m *customOrderModel) FindAll(ctx context.Context, status, page, pageSize int64) ([]*Order, int64, error) {
	return m.GormOrderModel.FindAll(ctx, status, page, pageSize)
}
