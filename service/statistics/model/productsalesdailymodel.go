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
	productSalesDailyFieldNames          = builder.RawFieldNames(&ProductSalesDaily{})
	productSalesDailyRows                = strings.Join(productSalesDailyFieldNames, ",")
	productSalesDailyRowsExpectAutoSet   = strings.Join(stringx.Remove(productSalesDailyFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	productSalesDailyRowsWithPlaceHolder = strings.Join(stringx.Remove(productSalesDailyFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"
)

type (
	ProductSalesDailyModel interface {
		Insert(ctx context.Context, data *ProductSalesDaily) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*ProductSalesDaily, error)
		FindOneByProductIdSalesDate(ctx context.Context, productId int64, salesDate time.Time) (*ProductSalesDaily, error)
		Update(ctx context.Context, data *ProductSalesDaily) error
		Delete(ctx context.Context, id int64) error
		FindHotProducts(ctx context.Context, start, end time.Time, limit int) ([]*ProductSalesDaily, error)
	}

	defaultProductSalesDailyModel struct {
		conn  sqlx.SqlConn
		table string
	}

	ProductSalesDaily struct {
		Id          int64     `db:"id"`
		ProductId   int64     `db:"product_id"`   // 商品ID
		CategoryId  int64     `db:"category_id"`  // 类别ID
		SalesDate   time.Time `db:"sales_date"`   // 销售日期
		SalesCount  int64     `db:"sales_count"`  // 销售数量
		SalesAmount float64   `db:"sales_amount"` // 销售金额
		CreatedAt   time.Time `db:"created_at"`   // 创建时间
		UpdatedAt   time.Time `db:"updated_at"`   // 更新时间
	}
)

func NewProductSalesDailyModel(conn sqlx.SqlConn) ProductSalesDailyModel {
	return &defaultProductSalesDailyModel{
		conn:  conn,
		table: "`product_sales_daily`",
	}
}

func (m *defaultProductSalesDailyModel) Insert(ctx context.Context, data *ProductSalesDaily) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?)", m.table, productSalesDailyRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.ProductId, data.CategoryId, data.SalesDate, data.SalesCount, data.SalesAmount, data.CreatedAt, data.UpdatedAt)
	return ret, err
}

func (m *defaultProductSalesDailyModel) FindOne(ctx context.Context, id int64) (*ProductSalesDaily, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", productSalesDailyRows, m.table)
	var resp ProductSalesDaily
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

func (m *defaultProductSalesDailyModel) FindOneByProductIdSalesDate(ctx context.Context, productId int64, salesDate time.Time) (*ProductSalesDaily, error) {
	var resp ProductSalesDaily
	query := fmt.Sprintf("select %s from %s where `product_id` = ? and `sales_date` = ? limit 1", productSalesDailyRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, productId, salesDate)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultProductSalesDailyModel) Update(ctx context.Context, data *ProductSalesDaily) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, productSalesDailyRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.ProductId, data.CategoryId, data.SalesDate, data.SalesCount, data.SalesAmount, data.CreatedAt, data.UpdatedAt, data.Id)
	return err
}

func (m *defaultProductSalesDailyModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

// FindHotProducts 查询指定时间范围内的热门商品
func (m *defaultProductSalesDailyModel) FindHotProducts(ctx context.Context, start, end time.Time, limit int) ([]*ProductSalesDaily, error) {
	query := fmt.Sprintf("select %s from %s where sales_date between ? and ? group by product_id order by sum(sales_count) desc limit ?", productSalesDailyRows, m.table)
	var resp []*ProductSalesDaily
	err := m.conn.QueryRowsCtx(ctx, &resp, query, start.Format("2006-01-02"), end.Format("2006-01-02"), limit)
	return resp, err
}
