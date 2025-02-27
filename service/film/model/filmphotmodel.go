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
	filmPhotoFieldNames          = builder.RawFieldNames(&FilmPhoto{})
	filmPhotoRows                = strings.Join(filmPhotoFieldNames, ",")
	filmPhotoRowsExpectAutoSet   = strings.Join(stringx.Remove(filmPhotoFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	filmPhotoRowsWithPlaceHolder = strings.Join(stringx.Remove(filmPhotoFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"
)

type (
	FilmPhotoModel interface {
		Insert(ctx context.Context, data *FilmPhoto) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*FilmPhoto, error)
		FindByFilmOrderId(ctx context.Context, filmOrderId int64) ([]*FilmPhoto, error)
		Update(ctx context.Context, data *FilmPhoto) error
		Delete(ctx context.Context, id int64) error
	}

	defaultFilmPhotoModel struct {
		conn  sqlx.SqlConn
		table string
	}

	FilmPhoto struct {
		Id          int64     `db:"id"`
		FilmOrderId int64     `db:"film_order_id"`
		Url         string    `db:"url"`
		Sort        int64     `db:"sort"`
		CreateTime  time.Time `db:"create_time"`
		UpdateTime  time.Time `db:"update_time"`
	}
)

func NewFilmPhotoModel(conn sqlx.SqlConn) FilmPhotoModel {
	return &defaultFilmPhotoModel{
		conn:  conn,
		table: "`film_photo`",
	}
}

func (m *defaultFilmPhotoModel) Insert(ctx context.Context, data *FilmPhoto) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?)", m.table, filmPhotoRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.FilmOrderId, data.Url, data.Sort)
	return ret, err
}

func (m *defaultFilmPhotoModel) FindOne(ctx context.Context, id int64) (*FilmPhoto, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", filmPhotoRows, m.table)
	var resp FilmPhoto
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

func (m *defaultFilmPhotoModel) FindByFilmOrderId(ctx context.Context, filmOrderId int64) ([]*FilmPhoto, error) {
	query := fmt.Sprintf("select %s from %s where `film_order_id` = ? order by `sort` asc", filmPhotoRows, m.table)
	var resp []*FilmPhoto
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

func (m *defaultFilmPhotoModel) Update(ctx context.Context, data *FilmPhoto) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, filmPhotoRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.FilmOrderId, data.Url, data.Sort, data.Id)
	return err
}

func (m *defaultFilmPhotoModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}
