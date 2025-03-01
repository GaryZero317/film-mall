package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ FaqModel = (*customDefaultFaqModel)(nil)

type (
	// FaqModel 是faq表的操作接口
	FaqModel interface {
		Insert(ctx context.Context, data *Faq) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*Faq, error)
		Update(ctx context.Context, data *Faq) error
		Delete(ctx context.Context, id int64) error
		FindList(ctx context.Context, page, pageSize int64, category int64) ([]*Faq, int64, error)
	}

	// Faq 常见问题
	Faq struct {
		Id         int64     `db:"id"`
		Question   string    `db:"question"`
		Answer     string    `db:"answer"`
		Category   int64     `db:"category"`
		Priority   int64     `db:"priority"`
		CreateTime time.Time `db:"create_time"`
		UpdateTime time.Time `db:"update_time"`
	}

	customDefaultFaqModel struct {
		conn  sqlx.SqlConn
		table string
	}
)

// NewFaqModel 返回一个操作faq表的模型对象
func NewFaqModel(conn sqlx.SqlConn) FaqModel {
	return &customDefaultFaqModel{
		conn:  conn,
		table: "`faq`",
	}
}

// Insert 插入一条新记录
func (m *customDefaultFaqModel) Insert(ctx context.Context, data *Faq) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?)", m.table, "question, answer, category, priority, create_time")
	ret, err := m.conn.ExecCtx(ctx, query, data.Question, data.Answer, data.Category, data.Priority, time.Now())
	return ret, err
}

// FindOne 根据id查询记录
func (m *customDefaultFaqModel) FindOne(ctx context.Context, id int64) (*Faq, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", "id, question, answer, category, priority, create_time, update_time", m.table)
	var resp Faq
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

// Update 更新记录
func (m *customDefaultFaqModel) Update(ctx context.Context, data *Faq) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, "question = ?, answer = ?, category = ?, priority = ?, update_time = ?")
	_, err := m.conn.ExecCtx(ctx, query, data.Question, data.Answer, data.Category, data.Priority, time.Now(), data.Id)
	return err
}

// Delete 删除记录
func (m *customDefaultFaqModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

// FindList 查询记录列表
func (m *customDefaultFaqModel) FindList(ctx context.Context, page, pageSize int64, category int64) ([]*Faq, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	whereClause := ""
	args := []interface{}{}

	if category > 0 {
		whereClause = "where `category` = ?"
		args = append(args, category)
	}

	// 查询总数
	countQuery := fmt.Sprintf("select count(*) from %s %s", m.table, whereClause)
	var total int64
	err := m.conn.QueryRowCtx(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// 查询列表
	query := fmt.Sprintf("select %s from %s %s order by priority desc, id asc limit ? offset ?",
		"id, question, answer, category, priority, create_time, update_time",
		m.table, whereClause)
	args = append(args, pageSize, offset)

	var resp []*Faq
	err = m.conn.QueryRowsCtx(ctx, &resp, query, args...)
	if err != nil {
		return nil, 0, err
	}

	return resp, total, nil
}
