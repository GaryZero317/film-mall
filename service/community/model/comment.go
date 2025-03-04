package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CommentModel = (*customCommentModel)(nil)

type (
	// CommentModel 评论模型接口
	CommentModel interface {
		Insert(ctx context.Context, data *Comment) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*Comment, error)
		FindOneByUidAndId(ctx context.Context, uid, id int64) (*Comment, error)
		Delete(ctx context.Context, id int64) error
		List(ctx context.Context, workId, page, pageSize int64) ([]*Comment, int64, error)
		FindReplies(ctx context.Context, commentId int64) ([]*Comment, error)
	}

	// 数据库评论表映射结构体
	Comment struct {
		Id         int64     `db:"id"`          // 评论ID
		Uid        int64     `db:"uid"`         // 用户ID
		WorkId     int64     `db:"work_id"`     // 作品ID
		Content    string    `db:"content"`     // 评论内容
		ReplyId    int64     `db:"reply_id"`    // 回复的评论ID(为0表示顶级评论)
		Status     int64     `db:"status"`      // 状态:0正常,1已删除
		CreateTime time.Time `db:"create_time"` // 创建时间
	}

	customCommentModel struct {
		conn  sqlx.SqlConn
		table string
	}
)

// NewCommentModel 创建评论模型
func NewCommentModel(conn sqlx.SqlConn) CommentModel {
	return &customCommentModel{
		conn:  conn,
		table: "comments",
	}
}

// Insert 插入评论
func (m *customCommentModel) Insert(ctx context.Context, data *Comment) (sql.Result, error) {
	data.CreateTime = time.Now()
	query := fmt.Sprintf("INSERT INTO %s (uid, work_id, content, reply_id, status, create_time) VALUES (?, ?, ?, ?, ?, ?)", m.table)
	return m.conn.ExecCtx(ctx, query, data.Uid, data.WorkId, data.Content, data.ReplyId, data.Status, data.CreateTime)
}

// FindOne 查询单个评论
func (m *customCommentModel) FindOne(ctx context.Context, id int64) (*Comment, error) {
	var resp Comment
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = ? LIMIT 1", m.table)
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

// FindOneByUidAndId 查询用户的指定评论
func (m *customCommentModel) FindOneByUidAndId(ctx context.Context, uid, id int64) (*Comment, error) {
	var resp Comment
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = ? AND uid = ? LIMIT 1", m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, id, uid)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// Delete 删除评论
func (m *customCommentModel) Delete(ctx context.Context, id int64) error {
	// 软删除
	query := fmt.Sprintf("UPDATE %s SET status = 1 WHERE id = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

// List 查询评论列表(仅顶级评论)
func (m *customCommentModel) List(ctx context.Context, workId, page, pageSize int64) ([]*Comment, int64, error) {
	conditions := []string{}
	args := []interface{}{}

	conditions = append(conditions, "work_id = ?")
	args = append(args, workId)

	conditions = append(conditions, "reply_id = 0") // 仅查询顶级评论
	conditions = append(conditions, "status = 0")   // 仅查询正常评论

	conditionStr := "WHERE " + strings.Join(conditions, " AND ")

	// 查询总数
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s %s", m.table, conditionStr)
	var count int64
	err := m.conn.QueryRowCtx(ctx, &count, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// 查询列表
	offset := (page - 1) * pageSize
	query := fmt.Sprintf("SELECT * FROM %s %s ORDER BY create_time DESC LIMIT ?, ?",
		m.table, conditionStr)
	args = append(args, offset, pageSize)

	var resp []*Comment
	err = m.conn.QueryRowsCtx(ctx, &resp, query, args...)
	if err != nil {
		return nil, 0, err
	}

	return resp, count, nil
}

// FindReplies 查询评论的回复
func (m *customCommentModel) FindReplies(ctx context.Context, commentId int64) ([]*Comment, error) {
	var resp []*Comment
	query := fmt.Sprintf("SELECT * FROM %s WHERE reply_id = ? AND status = 0 ORDER BY create_time ASC", m.table)
	err := m.conn.QueryRowsCtx(ctx, &resp, query, commentId)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
