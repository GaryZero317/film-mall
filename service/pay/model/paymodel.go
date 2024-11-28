package model

import (
	"context"
	"database/sql"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ PayModel = (*customPayModel)(nil)

type (
	// PayModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPayModel.
	PayModel interface {
		Insert(ctx context.Context, data *Pay) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*Pay, error)
		Update(ctx context.Context, data *Pay) error
		Delete(ctx context.Context, id int64) error
		FindOneByOid(ctx context.Context, oid int64) (*Pay, error)
	}

	customPayModel struct {
		*GormPayModel
	}
)

var (
	cachePayOidPrefix = "cache:pay:oid:"
)

// NewPayModel returns a model for the database table.
func NewPayModel(conn sqlx.SqlConn, c cache.CacheConf) PayModel {
	// 获取原始数据库连接
	rawDB, err := conn.RawDB()
	if err != nil {
		panic(err)
	}
	
	gormModel, err := NewGormPayModel(rawDB, c)
	if err != nil {
		panic(err)
	}
	return &customPayModel{
		GormPayModel: gormModel,
	}
}

// FindOneByOid returns a pay record by order ID
func (m *customPayModel) FindOneByOid(ctx context.Context, oid int64) (*Pay, error) {
	return m.GormPayModel.FindOneByOid(ctx, oid)
}
