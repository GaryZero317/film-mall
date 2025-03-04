package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserModel = (*customUserModel)(nil)

type (
	// UserModel 用户模型接口
	UserModel interface {
		FindOne(ctx context.Context, id int64) (*User, error)
		FindBatch(ctx context.Context, ids []int64) (map[int64]*User, error)
	}

	// 用户信息
	User struct {
		Id       int64  `db:"id"`       // 用户ID
		Nickname string `db:"nickname"` // 昵称
		Avatar   string `db:"avatar"`   // 头像
	}

	customUserModel struct {
		conn  sqlx.SqlConn
		table string
	}
)

// NewUserModel 创建用户模型
func NewUserModel(conn sqlx.SqlConn) UserModel {
	return &customUserModel{
		conn:  conn,
		table: "user",
	}
}

// FindOne 查询单个用户信息
func (m *customUserModel) FindOne(ctx context.Context, id int64) (*User, error) {
	var resp User
	query := fmt.Sprintf("SELECT id, nickname, avatar FROM %s WHERE id = ? LIMIT 1", m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// FindBatch 批量查询用户信息
func (m *customUserModel) FindBatch(ctx context.Context, ids []int64) (map[int64]*User, error) {
	if len(ids) == 0 {
		return map[int64]*User{}, nil
	}

	// 构建IN查询参数
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf("SELECT id, nickname, avatar FROM %s WHERE id IN (%s)",
		m.table, strings.Join(placeholders, ","))

	var users []*User
	err := m.conn.QueryRowsCtx(ctx, &users, query, args...)
	if err != nil {
		return nil, err
	}

	// 转换为map方便查找
	userMap := make(map[int64]*User)
	for _, user := range users {
		userMap[user.Id] = user
	}

	return userMap, nil
}
