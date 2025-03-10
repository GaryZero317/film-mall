package model

import (
	"context"
	"database/sql"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserModel = (*customUserModel)(nil)

type (
	// UserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserModel.
	UserModel interface {
		Insert(ctx context.Context, data *User) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*User, error)
		FindOneByMobile(ctx context.Context, mobile string) (*User, error)
		Update(ctx context.Context, data *User) error
		Delete(ctx context.Context, id int64) error
	}

	customUserModel struct {
		*GormUserModel
	}
)

// NewUserModel 返回数据库表的模型。
func NewUserModel(conn sqlx.SqlConn, c cache.CacheConf) UserModel {
	// 获取原始数据库连接
	rawDB, err := conn.RawDB()
	if err != nil {
		panic(err)
	}
	
	gormModel, err := NewGormUserModel(rawDB, c)
	if err != nil {
		panic(err)
	}
	return &customUserModel{
		GormUserModel: gormModel,
	}
}

// Insert implements UserModel.Insert
func (m *customUserModel) Insert(ctx context.Context, data *User) (sql.Result, error) {
	return m.GormUserModel.Insert(ctx, data)
}

// FindOne implements UserModel.FindOne
func (m *customUserModel) FindOne(ctx context.Context, id int64) (*User, error) {
	return m.GormUserModel.FindOne(ctx, id)
}

// FindOneByMobile implements UserModel.FindOneByMobile
func (m *customUserModel) FindOneByMobile(ctx context.Context, mobile string) (*User, error) {
	return m.GormUserModel.FindOneByMobile(ctx, mobile)
}

// Update implements UserModel.Update
func (m *customUserModel) Update(ctx context.Context, data *User) error {
	return m.GormUserModel.Update(ctx, data)
}

// Delete implements UserModel.Delete
func (m *customUserModel) Delete(ctx context.Context, id int64) error {
	return m.GormUserModel.Delete(ctx, id)
}
