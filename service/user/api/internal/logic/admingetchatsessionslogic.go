package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"mall/service/user/api/internal/svc"
	"mall/service/user/api/internal/types"
	"mall/service/user/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminGetChatSessionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminGetChatSessionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminGetChatSessionsLogic {
	return &AdminGetChatSessionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminGetChatSessionsLogic) AdminGetChatSessions() (resp *types.ChatSessionListResponse, err error) {
	logx.Info("开始处理获取管理员聊天会话列表请求")

	// 尝试获取管理员ID，首先尝试adminId
	logx.Info("尝试从上下文中获取管理员ID")

	// 先尝试获取adminId
	adminId, ok := l.ctx.Value("adminId").(int64)

	// 如果adminId获取失败，尝试uid
	if !ok {
		logx.Error("从上下文中获取adminId失败，尝试获取uid")

		// 获取uid并记录类型信息
		var uid interface{}
		uid = l.ctx.Value("uid")
		logx.Infof("ctx.Value(\"uid\")的值类型: %T, 值: %v", uid, uid)

		// 尝试类型转换
		if uid != nil {
			switch v := uid.(type) {
			case float64:
				adminId = int64(v)
				ok = true
				logx.Infof("成功将float64类型的uid转换为管理员ID: %d", adminId)
			case int:
				adminId = int64(v)
				ok = true
				logx.Infof("成功将int类型的uid转换为管理员ID: %d", adminId)
			case json.Number:
				adminId, err = v.Int64()
				if err == nil {
					ok = true
					logx.Infof("成功将json.Number类型的uid转换为管理员ID: %d", adminId)
				}
			default:
				logx.Errorf("无法转换uid，类型: %T", v)
			}
		}
	} else {
		logx.Infof("成功从adminId获取管理员ID: %d", adminId)
	}

	if !ok || adminId <= 0 {
		logx.Error("获取管理员ID失败或ID无效")
		return nil, errors.New("管理员未授权")
	}

	logx.Infof("成功获取管理员ID: %d", adminId)

	// 获取连接对象
	// 反射获取实际实现类型
	modelValue := reflect.ValueOf(l.svcCtx.ChatMessageModel)
	modelType := modelValue.Type()
	logx.Infof("ChatMessageModel的类型: %v", modelType)

	// 直接使用sql.Conn对象进行原生SQL查询
	db := l.svcCtx.Config.Mysql.DataSource
	logx.Infof("使用数据库连接: %s", db)

	// 使用简化的查询方式，直接使用已有的服务上下文方法
	messages, _, err := l.svcCtx.GetChatHistory(l.ctx, 0, adminId, 1, 9999)
	if err != nil {
		logx.Errorf("获取聊天历史失败: %v", err)
		return nil, err
	}

	// 提取所有不同的会话ID (这里使用的是用户ID)
	sessionIds := make(map[int64]bool)
	for _, msg := range messages {
		sessionIds[msg.SessionId] = true
	}

	logx.Infof("找到 %d 个聊天会话", len(sessionIds))

	// 转换为响应格式
	items := make([]types.ChatSessionItem, 0, len(sessionIds))
	for sessionId := range sessionIds {
		logx.Infof("处理会话 %d 的信息", sessionId)

		// 筛选属于此会话的消息
		var sessionMessages []*model.ChatMessage
		for _, msg := range messages {
			if msg.SessionId == sessionId {
				sessionMessages = append(sessionMessages, msg)
			}
		}

		if len(sessionMessages) == 0 {
			logx.Infof("会话 %d 没有消息记录，跳过", sessionId)
			continue
		}

		// 获取最后一条消息（首条为最新的，因为是按时间倒序）
		lastMessage := sessionMessages[0]

		// 计算未读消息数
		var unreadCount int64
		for _, msg := range sessionMessages {
			if msg.SenderType == model.SenderTypeUser && msg.ReadStatus == 0 {
				unreadCount++
			}
		}

		// 用户名（简化处理，使用会话ID）
		username := fmt.Sprintf("用户%d", sessionId)

		items = append(items, types.ChatSessionItem{
			Id:           sessionId,
			UserName:     username,
			LastMessage:  lastMessage.Content,
			UnreadCount:  unreadCount,
			LastActivity: lastMessage.CreateTime.Unix(),
		})

		logx.Infof("会话 %d 信息处理完成", sessionId)
	}

	logx.Infof("管理员聊天会话列表生成完成，共 %d 个会话", len(items))
	return &types.ChatSessionListResponse{
		Total: int64(len(items)),
		List:  items,
	}, nil
}
