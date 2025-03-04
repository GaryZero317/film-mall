package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ LikeModel = (*customLikeModel)(nil)

type (
	// LikeModel 点赞模型接口
	LikeModel interface {
		Insert(ctx context.Context, data *Like) (sql.Result, error)
		Delete(ctx context.Context, uid, workId int64) error
		IsLiked(ctx context.Context, uid, workId int64) (bool, error)
	}

	// 数据库点赞表映射结构体
	Like struct {
		Id         int64     `db:"id"`          // 点赞ID
		Uid        int64     `db:"uid"`         // 用户ID
		WorkId     int64     `db:"work_id"`     // 作品ID
		CreateTime time.Time `db:"create_time"` // 创建时间
	}

	customLikeModel struct {
		conn  sqlx.SqlConn
		table string
	}
)

// NewLikeModel 创建点赞模型
func NewLikeModel(conn sqlx.SqlConn) LikeModel {
	return &customLikeModel{
		conn:  conn,
		table: "likes",
	}
}

// Insert 插入点赞记录
func (m *customLikeModel) Insert(ctx context.Context, data *Like) (sql.Result, error) {
	data.CreateTime = time.Now()
	query := fmt.Sprintf("INSERT INTO %s (uid, work_id, create_time) VALUES (?, ?, ?)", m.table)
	return m.conn.ExecCtx(ctx, query, data.Uid, data.WorkId, data.CreateTime)
}

// Delete 删除点赞记录
func (m *customLikeModel) Delete(ctx context.Context, uid, workId int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE uid = ? AND work_id = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, uid, workId)
	return err
}

// IsLiked 检查用户是否已点赞
func (m *customLikeModel) IsLiked(ctx context.Context, uid, workId int64) (bool, error) {
	var count int64
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE uid = ? AND work_id = ?", m.table)
	err := m.conn.QueryRowCtx(ctx, &count, query, uid, workId)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
