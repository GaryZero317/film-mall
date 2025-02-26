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
	filmOrderFieldNames          = builder.RawFieldNames(&FilmOrder{})
	filmOrderRows                = strings.Join(filmOrderFieldNames, ",")
	filmOrderRowsExpectAutoSet   = strings.Join(stringx.Remove(filmOrderFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	filmOrderRowsWithPlaceHolder = strings.Join(stringx.Remove(filmOrderFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"
)

type (
	FilmOrderModel interface {
		Insert(ctx context.Context, data *FilmOrder) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*FilmOrder, error)
		FindOneByFoid(ctx context.Context, foid string) (*FilmOrder, error)
		FindByUid(ctx context.Context, uid int64, status int64, page, pageSize int64) ([]*FilmOrder, int64, error)
		FindAll(ctx context.Context, status int64, page, pageSize int64) ([]*FilmOrder, int64, error)
		Update(ctx context.Context, data *FilmOrder) error
		Delete(ctx context.Context, id int64) error
	}

	defaultFilmOrderModel struct {
		conn  sqlx.SqlConn
		table string
	}

	FilmOrder struct {
		Id          int64     `db:"id"`
		Foid        string    `db:"foid"`
		Uid         int64     `db:"uid"`
		AddressId   int64     `db:"address_id"`
		ReturnFilm  bool      `db:"return_film"`
		TotalPrice  int64     `db:"total_price"`
		ShippingFee int64     `db:"shipping_fee"`
		Status      int64     `db:"status"`
		Remark      string    `db:"remark"`
		CreateTime  time.Time `db:"create_time"`
		UpdateTime  time.Time `db:"update_time"`
	}
)

func NewFilmOrderModel(conn sqlx.SqlConn) FilmOrderModel {
	return &defaultFilmOrderModel{
		conn:  conn,
		table: "`film_order`",
	}
}

func (m *defaultFilmOrderModel) Insert(ctx context.Context, data *FilmOrder) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?)", m.table, filmOrderRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.Foid, data.Uid, data.AddressId, data.ReturnFilm, data.TotalPrice, data.ShippingFee, data.Status, data.Remark)
	return ret, err
}

func (m *defaultFilmOrderModel) FindOne(ctx context.Context, id int64) (*FilmOrder, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", filmOrderRows, m.table)
	var resp FilmOrder
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

func (m *defaultFilmOrderModel) FindOneByFoid(ctx context.Context, foid string) (*FilmOrder, error) {
	query := fmt.Sprintf("select %s from %s where `foid` = ? limit 1", filmOrderRows, m.table)
	var resp FilmOrder
	err := m.conn.QueryRowCtx(ctx, &resp, query, foid)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultFilmOrderModel) FindByUid(ctx context.Context, uid int64, status int64, page, pageSize int64) ([]*FilmOrder, int64, error) {
	var condition string
	var args []interface{}

	condition = "`uid` = ?"
	args = append(args, uid)

	if status != -1 {
		condition += " AND `status` = ?"
		args = append(args, status)
	}

	countQuery := fmt.Sprintf("select count(*) from %s where %s", m.table, condition)
	var count int64
	err := m.conn.QueryRowCtx(ctx, &count, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	query := fmt.Sprintf("select %s from %s where %s order by `id` desc limit ?, ?", filmOrderRows, m.table, condition)
	args = append(args, offset, pageSize)

	var resp []*FilmOrder
	err = m.conn.QueryRowsCtx(ctx, &resp, query, args...)
	if err != nil {
		return nil, 0, err
	}

	return resp, count, nil
}

func (m *defaultFilmOrderModel) FindAll(ctx context.Context, status int64, page, pageSize int64) ([]*FilmOrder, int64, error) {
	var condition string
	var args []interface{}

	if status != -1 {
		condition = "`status` = ?"
		args = append(args, status)
	}

	var countQuery string
	if condition != "" {
		countQuery = fmt.Sprintf("select count(*) from %s where %s", m.table, condition)
	} else {
		countQuery = fmt.Sprintf("select count(*) from %s", m.table)
	}

	var count int64
	err := m.conn.QueryRowCtx(ctx, &count, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	var query string
	if condition != "" {
		query = fmt.Sprintf("select %s from %s where %s order by `id` desc limit ?, ?", filmOrderRows, m.table, condition)
		args = append(args, offset, pageSize)
	} else {
		query = fmt.Sprintf("select %s from %s order by `id` desc limit ?, ?", filmOrderRows, m.table)
		args = append(args, offset, pageSize)
	}

	var resp []*FilmOrder
	err = m.conn.QueryRowsCtx(ctx, &resp, query, args...)
	if err != nil {
		return nil, 0, err
	}

	return resp, count, nil
}

func (m *defaultFilmOrderModel) Update(ctx context.Context, data *FilmOrder) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, filmOrderRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.Foid, data.Uid, data.AddressId, data.ReturnFilm, data.TotalPrice, data.ShippingFee, data.Status, data.Remark, data.Id)
	return err
}

func (m *defaultFilmOrderModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}
