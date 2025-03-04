package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ WorkImageModel = (*customWorkImageModel)(nil)

type (
	// WorkImageModel 作品图片模型接口
	WorkImageModel interface {
		Insert(ctx context.Context, data *WorkImage) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*WorkImage, error)
		FindByWorkId(ctx context.Context, workId int64) ([]*WorkImage, error)
		Delete(ctx context.Context, id int64) error
		DeleteByWorkId(ctx context.Context, workId int64) error
	}

	// 数据库作品图片表映射结构体
	WorkImage struct {
		Id         int64     `db:"id"`          // 图片ID
		WorkId     int64     `db:"work_id"`     // 作品ID
		Url        string    `db:"url"`         // 图片URL
		Sort       int64     `db:"sort"`        // 排序
		CreateTime time.Time `db:"create_time"` // 创建时间
	}

	customWorkImageModel struct {
		conn  sqlx.SqlConn
		table string
	}
)

// NewWorkImageModel 创建作品图片模型
func NewWorkImageModel(conn sqlx.SqlConn) WorkImageModel {
	return &customWorkImageModel{
		conn:  conn,
		table: "work_images",
	}
}

// Insert 插入作品图片
func (m *customWorkImageModel) Insert(ctx context.Context, data *WorkImage) (sql.Result, error) {
	data.CreateTime = time.Now()
	query := fmt.Sprintf("INSERT INTO %s (work_id, url, sort, create_time) VALUES (?, ?, ?, ?)", m.table)
	return m.conn.ExecCtx(ctx, query, data.WorkId, data.Url, data.Sort, data.CreateTime)
}

// FindOne 查询单个作品图片
func (m *customWorkImageModel) FindOne(ctx context.Context, id int64) (*WorkImage, error) {
	var resp WorkImage
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

// FindByWorkId 查询指定作品的所有图片
func (m *customWorkImageModel) FindByWorkId(ctx context.Context, workId int64) ([]*WorkImage, error) {
	var resp []*WorkImage
	query := fmt.Sprintf("SELECT * FROM %s WHERE work_id = ? ORDER BY sort ASC", m.table)
	err := m.conn.QueryRowsCtx(ctx, &resp, query, workId)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// Delete 删除作品图片
func (m *customWorkImageModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

// DeleteByWorkId 删除指定作品的所有图片
func (m *customWorkImageModel) DeleteByWorkId(ctx context.Context, workId int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE work_id = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, workId)
	return err
}
