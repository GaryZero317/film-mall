package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// ChatMessage 聊天消息结构
type ChatMessage struct {
	Id         int64     `db:"id"`
	UserId     int64     `db:"user_id"`
	AdminId    int64     `db:"admin_id"`
	Direction  int64     `db:"direction"`
	Content    string    `db:"content"`
	ReadStatus int64     `db:"read_status"`
	CreateTime time.Time `db:"create_time"`
}

// ChatMessageModel 聊天消息模型接口
type ChatMessageModel interface {
	Insert(context.Context, *ChatMessage) (sql.Result, error)
	FindByUserAndAdmin(context.Context, int64, int64, int, int) ([]*ChatMessage, int64, error)
}

type defaultChatMessageModel struct {
	conn sqlx.SqlConn
}

// NewChatMessageModel 创建聊天消息模型
func NewChatMessageModel(conn sqlx.SqlConn) ChatMessageModel {
	return &defaultChatMessageModel{conn: conn}
}

func (m *defaultChatMessageModel) Insert(ctx context.Context, data *ChatMessage) (sql.Result, error) {
	query := `insert into chat_message (user_id, admin_id, direction, content, read_status, create_time) values (?, ?, ?, ?, ?, ?)`
	return m.conn.ExecCtx(ctx, query, data.UserId, data.AdminId, data.Direction, data.Content, data.ReadStatus, data.CreateTime)
}

func (m *defaultChatMessageModel) FindByUserAndAdmin(ctx context.Context, userId, adminId int64, page, pageSize int) ([]*ChatMessage, int64, error) {
	var messages []*ChatMessage
	query := `select id, user_id, admin_id, direction, content, read_status, create_time from chat_message where user_id = ? and admin_id = ? order by create_time desc limit ?, ?`
	err := m.conn.QueryRowsCtx(ctx, &messages, query, userId, adminId, (page-1)*pageSize, pageSize)

	var count int64
	countQuery := `select count(*) from chat_message where user_id = ? and admin_id = ?`
	_ = m.conn.QueryRowCtx(ctx, &count, countQuery, userId, adminId)

	return messages, count, err
}
