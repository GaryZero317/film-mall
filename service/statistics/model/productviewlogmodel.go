package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	productViewLogFieldNames          = builder.RawFieldNames(&ProductViewLog{})
	productViewLogRows                = strings.Join(productViewLogFieldNames, ",")
	productViewLogRowsExpectAutoSet   = strings.Join(stringx.Remove(productViewLogFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	productViewLogRowsWithPlaceHolder = strings.Join(stringx.Remove(productViewLogFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"
)

type (
	ProductViewLogModel interface {
		Insert(ctx context.Context, data *ProductViewLog) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*ProductViewLog, error)
		Update(ctx context.Context, data *ProductViewLog) error
		Delete(ctx context.Context, id int64) error
	}

	defaultProductViewLogModel struct {
		conn  sqlx.SqlConn
		table string
	}

	ProductViewLog struct {
		Id        int64          `db:"id"`
		UserId    int64          `db:"user_id"`    // 用户ID
		ProductId int64          `db:"product_id"` // 商品ID
		ViewTime  time.Time      `db:"view_time"`  // 浏览时间
		Ip        sql.NullString `db:"ip"`         // 访问IP
		UserAgent sql.NullString `db:"user_agent"` // 用户代理
	}
)

func NewProductViewLogModel(conn sqlx.SqlConn) ProductViewLogModel {
	return &defaultProductViewLogModel{
		conn:  conn,
		table: "`product_view_log`",
	}
}

func (m *defaultProductViewLogModel) Insert(ctx context.Context, data *ProductViewLog) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?)", m.table, productViewLogRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.UserId, data.ProductId, data.ViewTime, data.Ip, data.UserAgent)
	return ret, err
}

func (m *defaultProductViewLogModel) FindOne(ctx context.Context, id int64) (*ProductViewLog, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", productViewLogRows, m.table)
	var resp ProductViewLog
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

func (m *defaultProductViewLogModel) Update(ctx context.Context, data *ProductViewLog) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, productViewLogRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.UserId, data.ProductId, data.ViewTime, data.Ip, data.UserAgent, data.Id)
	return err
}

func (m *defaultProductViewLogModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}
