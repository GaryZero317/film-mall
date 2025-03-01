package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CustomerServiceModel = (*customDefaultCustomerServiceModel)(nil)

type (
	// CustomerServiceModel 是customer_service表的操作接口
	CustomerServiceModel interface {
		Insert(ctx context.Context, data *CustomerService) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*CustomerService, error)
		Update(ctx context.Context, data *CustomerService) error
		Delete(ctx context.Context, id int64) error
		FindByUserId(ctx context.Context, userId int64, page, pageSize int64, status int64) ([]*CustomerService, int64, error)
	}

	// CustomerService 客服问题
	CustomerService struct {
		Id         int64     `db:"id"`
		UserId     int64     `db:"user_id"`
		Title      string    `db:"title"`
		Content    string    `db:"content"`
		Type       int64     `db:"type"`
		Status     int64     `db:"status"`
		Reply      string    `db:"reply"`
		ReplyTime  time.Time `db:"reply_time"`
		ContactWay string    `db:"contact_way"`
		CreateTime time.Time `db:"create_time"`
		UpdateTime time.Time `db:"update_time"`
	}

	customDefaultCustomerServiceModel struct {
		conn  sqlx.SqlConn
		table string
	}
)

// NewCustomerServiceModel 返回一个操作customer_service表的模型对象
func NewCustomerServiceModel(conn sqlx.SqlConn) CustomerServiceModel {
	return &customDefaultCustomerServiceModel{
		conn:  conn,
		table: "`customer_service`",
	}
}

// Insert 插入一条新记录
func (m *customDefaultCustomerServiceModel) Insert(ctx context.Context, data *CustomerService) (sql.Result, error) {
	// 使用字段列表和值列表变量，方便调整个别字段
	fields := "user_id, title, content, type, status, contact_way, create_time"
	values := []interface{}{data.UserId, data.Title, data.Content, data.Type, data.Status, data.ContactWay, time.Now()}

	// 如果回复字段有值，则添加到插入字段列表中
	if data.Reply != "" {
		fields += ", reply, reply_time"
		values = append(values, data.Reply, data.ReplyTime)
	} else {
		// 如果没有回复，则只添加reply字段，不添加reply_time
		fields += ", reply"
		values = append(values, data.Reply)
	}

	query := fmt.Sprintf("insert into %s (%s) values (%s)", m.table, fields, createPlaceholders(len(values)))
	ret, err := m.conn.ExecCtx(ctx, query, values...)
	return ret, err
}

// createPlaceholders 创建SQL占位符
func createPlaceholders(count int) string {
	if count <= 0 {
		return ""
	}

	placeholders := "?"
	for i := 1; i < count; i++ {
		placeholders += ", ?"
	}
	return placeholders
}

// FindOne 根据id查询记录
func (m *customDefaultCustomerServiceModel) FindOne(ctx context.Context, id int64) (*CustomerService, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", "id, user_id, title, content, type, status, reply, reply_time, contact_way, create_time, update_time", m.table)
	var resp CustomerService
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
func (m *customDefaultCustomerServiceModel) Update(ctx context.Context, data *CustomerService) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, "user_id = ?, title = ?, content = ?, type = ?, status = ?, reply = ?, reply_time = ?, contact_way = ?, update_time = ?")
	_, err := m.conn.ExecCtx(ctx, query, data.UserId, data.Title, data.Content, data.Type, data.Status, data.Reply, data.ReplyTime, data.ContactWay, time.Now(), data.Id)
	return err
}

// Delete 删除记录
func (m *customDefaultCustomerServiceModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

// FindByUserId 根据用户ID查询记录列表
func (m *customDefaultCustomerServiceModel) FindByUserId(ctx context.Context, userId int64, page, pageSize int64, status int64) ([]*CustomerService, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	whereClause := "where `user_id` = ?"
	args := []interface{}{userId}

	if status > 0 {
		whereClause += " and `status` = ?"
		args = append(args, status)
	}

	// 查询总数
	countQuery := fmt.Sprintf("select count(*) from %s %s", m.table, whereClause)
	var total int64
	err := m.conn.QueryRowCtx(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// 查询列表
	query := fmt.Sprintf("select %s from %s %s order by create_time desc limit ? offset ?",
		"id, user_id, title, content, type, status, reply, reply_time, contact_way, create_time, update_time",
		m.table, whereClause)
	args = append(args, pageSize, offset)

	var resp []*CustomerService
	err = m.conn.QueryRowsCtx(ctx, &resp, query, args...)
	if err != nil {
		return nil, 0, err
	}

	return resp, total, nil
}
