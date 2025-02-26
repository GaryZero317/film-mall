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
	filmOrderItemFieldNames          = builder.RawFieldNames(&FilmOrderItem{})
	filmOrderItemRows                = strings.Join(filmOrderItemFieldNames, ",")
	filmOrderItemRowsExpectAutoSet   = strings.Join(stringx.Remove(filmOrderItemFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	filmOrderItemRowsWithPlaceHolder = strings.Join(stringx.Remove(filmOrderItemFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"
)

type (
	FilmOrderItemModel interface {
		Insert(ctx context.Context, data *FilmOrderItem) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*FilmOrderItem, error)
		FindByFilmOrderId(ctx context.Context, filmOrderId int64) ([]*FilmOrderItem, error)
		Update(ctx context.Context, data *FilmOrderItem) error
		Delete(ctx context.Context, id int64) error
	}

	defaultFilmOrderItemModel struct {
		conn  sqlx.SqlConn
		table string
	}

	FilmOrderItem struct {
		Id          int64     `db:"id"`
		FilmOrderId int64     `db:"film_order_id"`
		FilmType    string    `db:"film_type"`
		FilmBrand   string    `db:"film_brand"`
		Size        string    `db:"size"`
		Quantity    int64     `db:"quantity"`
		Price       int64     `db:"price"`
		Amount      int64     `db:"amount"`
		Remark      string    `db:"remark"`
		CreateTime  time.Time `db:"create_time"`
		UpdateTime  time.Time `db:"update_time"`
	}
)

func NewFilmOrderItemModel(conn sqlx.SqlConn) FilmOrderItemModel {
	return &defaultFilmOrderItemModel{
		conn:  conn,
		table: "`film_order_item`",
	}
}

func (m *defaultFilmOrderItemModel) Insert(ctx context.Context, data *FilmOrderItem) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?)", m.table, filmOrderItemRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.FilmOrderId, data.FilmType, data.FilmBrand, data.Size, data.Quantity, data.Price, data.Amount, data.Remark)
	return ret, err
}

func (m *defaultFilmOrderItemModel) FindOne(ctx context.Context, id int64) (*FilmOrderItem, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", filmOrderItemRows, m.table)
	var resp FilmOrderItem
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

func (m *defaultFilmOrderItemModel) FindByFilmOrderId(ctx context.Context, filmOrderId int64) ([]*FilmOrderItem, error) {
	query := fmt.Sprintf("select %s from %s where `film_order_id` = ?", filmOrderItemRows, m.table)
	var resp []*FilmOrderItem
	err := m.conn.QueryRowsCtx(ctx, &resp, query, filmOrderId)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultFilmOrderItemModel) Update(ctx context.Context, data *FilmOrderItem) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, filmOrderItemRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.FilmOrderId, data.FilmType, data.FilmBrand, data.Size, data.Quantity, data.Price, data.Amount, data.Remark, data.Id)
	return err
}

func (m *defaultFilmOrderItemModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}
