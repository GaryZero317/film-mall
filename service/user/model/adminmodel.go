package model

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AdminModel = (*customAdminModel)(nil)

type (
	// AdminModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAdminModel.
	AdminModel interface {
		Insert(ctx context.Context, data *GormAdmin) error
		FindOne(ctx context.Context, id int64) (*GormAdmin, error)
		FindOneByUsername(ctx context.Context, username string) (*GormAdmin, error)
		Update(ctx context.Context, data *GormAdmin) error
		Delete(ctx context.Context, id int64) error
	}

	customAdminModel struct {
		*GormAdminModel
	}
)

// NewAdminModel 返回数据库表的模型。
func NewAdminModel(conn sqlx.SqlConn, c cache.CacheConf) AdminModel {
	// 获取原始数据库连接
	rawDB, err := conn.RawDB()
	if err != nil {
		panic(err)
	}

	gormModel, err := NewGormAdminModel(rawDB, c)
	if err != nil {
		panic(err)
	}
	return &customAdminModel{
		GormAdminModel: gormModel,
	}
}

// Insert implements AdminModel.Insert
func (m *customAdminModel) Insert(ctx context.Context, data *GormAdmin) error {
	return m.GormAdminModel.Insert(ctx, data)
}

// FindOne implements AdminModel.FindOne
func (m *customAdminModel) FindOne(ctx context.Context, id int64) (*GormAdmin, error) {
	return m.GormAdminModel.FindOne(ctx, id)
}

// FindOneByUsername implements AdminModel.FindOneByUsername
func (m *customAdminModel) FindOneByUsername(ctx context.Context, username string) (*GormAdmin, error) {
	return m.GormAdminModel.FindOneByUsername(ctx, username)
}

// Update implements AdminModel.Update
func (m *customAdminModel) Update(ctx context.Context, data *GormAdmin) error {
	return m.GormAdminModel.Update(ctx, data)
}

// Delete implements AdminModel.Delete
func (m *customAdminModel) Delete(ctx context.Context, id int64) error {
	return m.GormAdminModel.Delete(ctx, id)
}
