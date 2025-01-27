package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	addressFieldNames          = builder.RawFieldNames(&Address{})
	addressRows                = strings.Join(addressFieldNames, ",")
	addressRowsExpectAutoSet   = strings.Join(stringx.Remove(addressFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	addressRowsWithPlaceHolder = strings.Join(stringx.Remove(addressFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheAddressIdPrefix = "cache:address:id:"
)

type (
	AddressModel interface {
		Insert(ctx context.Context, data *Address) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*Address, error)
		Update(ctx context.Context, data *Address) error
		Delete(ctx context.Context, id int64) error
		FindByUserId(ctx context.Context, userId int64) ([]*Address, error)
		UpdateDefaultByUserId(ctx context.Context, userId int64, isDefault bool) error
	}

	defaultAddressModel struct {
		sqlc.CachedConn
		table string
	}

	Address struct {
		Id         int64     `db:"id"`
		UserId     int64     `db:"user_id"`     // 用户ID
		Name       string    `db:"name"`        // 收货人姓名
		Phone      string    `db:"phone"`       // 联系电话
		Province   string    `db:"province"`    // 省份
		City       string    `db:"city"`        // 城市
		District   string    `db:"district"`    // 区/县
		DetailAddr string    `db:"detail_addr"` // 详细地址
		IsDefault  int64     `db:"is_default"`  // 是否为默认地址：0-否，1-是
		CreateTime time.Time `db:"create_time"`
		UpdateTime time.Time `db:"update_time"`
	}
)

func NewAddressModel(conn sqlx.SqlConn, c cache.CacheConf) AddressModel {
	return &defaultAddressModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`address`",
	}
}

func (m *defaultAddressModel) Insert(ctx context.Context, data *Address) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?)", m.table, addressRowsExpectAutoSet)
	ret, err := m.ExecNoCacheCtx(ctx, query, data.UserId, data.Name, data.Phone, data.Province, data.City, data.District, data.DetailAddr, data.IsDefault)

	return ret, err
}

func (m *defaultAddressModel) FindOne(ctx context.Context, id int64) (*Address, error) {
	addressIdKey := fmt.Sprintf("%s%v", cacheAddressIdPrefix, id)
	var resp Address
	err := m.QueryRowCtx(ctx, &resp, addressIdKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", addressRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, id)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultAddressModel) Update(ctx context.Context, data *Address) error {
	addressIdKey := fmt.Sprintf("%s%v", cacheAddressIdPrefix, data.Id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, addressRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, data.UserId, data.Name, data.Phone, data.Province, data.City, data.District, data.DetailAddr, data.IsDefault, data.Id)
	}, addressIdKey)
	return err
}

func (m *defaultAddressModel) Delete(ctx context.Context, id int64) error {
	addressIdKey := fmt.Sprintf("%s%v", cacheAddressIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, addressIdKey)
	return err
}

func (m *defaultAddressModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheAddressIdPrefix, primary)
}

func (m *defaultAddressModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", addressRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultAddressModel) FindByUserId(ctx context.Context, userId int64) ([]*Address, error) {
	var resp []*Address
	query := fmt.Sprintf("select %s from %s where `user_id` = ? order by `is_default` desc, `update_time` desc", addressRows, m.table)
	err := m.CachedConn.QueryRowsNoCacheCtx(ctx, &resp, query, userId)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultAddressModel) UpdateDefaultByUserId(ctx context.Context, userId int64, isDefault bool) error {
	query := fmt.Sprintf("update %s set `is_default` = ? where `user_id` = ?", m.table)
	_, err := m.ExecNoCacheCtx(ctx, query, isDefault, userId)
	return err
}
