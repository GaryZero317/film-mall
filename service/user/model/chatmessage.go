package model

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// 发送者类型常量
const (
	SenderTypeUser  = 1 // 用户
	SenderTypeAdmin = 2 // 管理员
)

// ChatMessage 聊天消息结构 - 匹配数据库结构
type ChatMessage struct {
	Id         int64     `db:"id"`
	SessionId  int64     `db:"session_id"`  // 会话ID（对应用户ID）
	SenderType int64     `db:"sender_type"` // 发送者类型：1-用户，2-管理员
	SenderId   int64     `db:"sender_id"`   // 发送者ID
	Content    string    `db:"content"`     // 消息内容
	ReadStatus int64     `db:"read_status"` // 读取状态：0-未读，1-已读
	CreateTime time.Time `db:"create_time"` // 创建时间
}

// ChatMessageModel 聊天消息模型接口
type ChatMessageModel interface {
	Insert(context.Context, *ChatMessage) (sql.Result, error)
	FindBySession(context.Context, int64, int, int) ([]*ChatMessage, int64, error)
	FindByUserAndAdmin(context.Context, int64, int64, int, int) ([]*ChatMessage, int64, error)
	UpdateReadStatus(context.Context, []int64, int64) error
}

type defaultChatMessageModel struct {
	conn sqlx.SqlConn
}

// NewChatMessageModel 创建聊天消息模型
func NewChatMessageModel(conn sqlx.SqlConn) ChatMessageModel {
	return &defaultChatMessageModel{conn: conn}
}

// Insert 插入聊天消息
func (m *defaultChatMessageModel) Insert(ctx context.Context, data *ChatMessage) (sql.Result, error) {
	query := `insert into chat_message (session_id, sender_type, sender_id, content, read_status, create_time) 
              values (?, ?, ?, ?, ?, ?)`
	return m.conn.ExecCtx(ctx, query, data.SessionId, data.SenderType, data.SenderId, data.Content, data.ReadStatus, data.CreateTime)
}

// FindBySession 根据会话ID查询消息
func (m *defaultChatMessageModel) FindBySession(ctx context.Context, sessionId int64, page, pageSize int) ([]*ChatMessage, int64, error) {
	var messages []*ChatMessage
	query := `select id, session_id, sender_type, sender_id, content, read_status, create_time 
              from chat_message where session_id = ? 
              order by create_time desc limit ?, ?`
	err := m.conn.QueryRowsCtx(ctx, &messages, query, sessionId, (page-1)*pageSize, pageSize)

	var count int64
	countQuery := `select count(*) from chat_message where session_id = ?`
	_ = m.conn.QueryRowCtx(ctx, &count, countQuery, sessionId)

	return messages, count, err
}

// FindByUserAndAdmin 根据用户ID和管理员ID查询消息（向后兼容）
func (m *defaultChatMessageModel) FindByUserAndAdmin(ctx context.Context, userId, adminId int64, page, pageSize int) ([]*ChatMessage, int64, error) {
	// 使用userId作为会话ID
	var messages []*ChatMessage
	query := `select id, session_id, sender_type, sender_id, content, read_status, create_time 
              from chat_message where session_id = ? 
              order by create_time desc limit ?, ?`
	err := m.conn.QueryRowsCtx(ctx, &messages, query, userId, (page-1)*pageSize, pageSize)

	var count int64
	countQuery := `select count(*) from chat_message where session_id = ?`
	_ = m.conn.QueryRowCtx(ctx, &count, countQuery, userId)

	return messages, count, err
}

// UpdateReadStatus 批量更新消息读取状态
func (m *defaultChatMessageModel) UpdateReadStatus(ctx context.Context, messageIds []int64, status int64) error {
	if len(messageIds) == 0 {
		return nil
	}

	// 构建问号占位符
	placeholders := make([]string, len(messageIds))
	args := make([]interface{}, len(messageIds)+1)
	args[0] = status

	for i := range messageIds {
		placeholders[i] = "?"
		args[i+1] = messageIds[i]
	}

	// 使用IN子句批量更新
	query := `update chat_message set read_status = ? where id in (` + strings.Join(placeholders, ",") + `)`
	_, err := m.conn.ExecCtx(ctx, query, args...)
	return err
}
