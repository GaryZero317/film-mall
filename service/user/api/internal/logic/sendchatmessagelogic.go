package logic

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"mall/service/user/api/internal/svc"
	"mall/service/user/api/internal/types"
	"mall/service/user/api/internal/ws"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendChatMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendChatMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendChatMessageLogic {
	return &SendChatMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendChatMessageLogic) SendChatMessage(req *types.SendChatMessageRequest) (resp *types.SendChatMessageResponse, err error) {
	// 记录请求的详细信息
	logx.Infof("开始处理发送消息请求: adminId=%d, content=%s, type=%d",
		req.AdminId, req.Content, req.Type)

	// 尝试从上下文中获取用户ID
	logx.Info("尝试从上下文中获取用户ID")

	// 修改: 从 "userId" 改为 "uid"，go-zero JWT中间件将用户ID存储为"uid"
	userId, ok := l.ctx.Value("uid").(int64)
	if !ok {
		// 尝试其他可能的键名
		logx.Error("从上下文中获取uid失败，尝试其他可能的键名")
		var uid interface{}
		uid = l.ctx.Value("uid")
		logx.Infof("ctx.Value(\"uid\")的值类型: %T, 值: %v", uid, uid)

		// 如果是数字类型但不是int64，尝试转换
		if uid != nil {
			switch v := uid.(type) {
			case float64:
				userId = int64(v)
				ok = true
				logx.Infof("成功将float64类型的uid转换为int64: %d", userId)
			case int:
				userId = int64(v)
				ok = true
				logx.Infof("成功将int类型的uid转换为int64: %d", userId)
			case json.Number:
				userId, err = v.Int64()
				if err == nil {
					ok = true
					logx.Infof("成功将json.Number类型的uid转换为int64: %d", userId)
				}
			default:
				logx.Errorf("无法转换uid，类型: %T", v)
			}
		}

		if !ok {
			logx.Error("从上下文中获取用户ID失败")
			return nil, errors.New("用户未授权")
		}
	}

	if userId <= 0 {
		logx.Errorf("无效的用户ID: %d", userId)
		return nil, errors.New("用户未授权")
	}

	logx.Infof("成功获取用户ID: %d", userId)

	// 默认文本消息
	msgType := int(req.Type)
	if msgType == 0 {
		msgType = ws.TextMessage
		logx.Info("未指定消息类型，使用默认文本消息类型")
	}

	// 创建消息
	message := &ws.Message{
		UserId:    userId,
		AdminId:   req.AdminId,
		Direction: ws.UserToAdmin,
		Content:   req.Content,
		Type:      msgType,
		Timestamp: time.Now().Unix(),
	}

	logx.Infof("创建消息对象: %+v", message)

	// 检查WebSocket管理器状态
	if l.svcCtx.WSManager == nil {
		logx.Error("WebSocket管理器未初始化")
		return nil, errors.New("WebSocket服务不可用")
	}

	logx.Info("WebSocket管理器状态正常，开始处理消息")

	// 广播消息
	// 存储消息获取ID
	id, err := l.svcCtx.WSManager.StoreMessage(message)
	if err != nil {
		logx.Errorf("存储消息失败: %v", err)
	} else {
		message.ID = id
		logx.Infof("消息存储成功，ID: %d", id)
	}

	// 广播消息
	logx.Info("发送消息到广播通道")
	l.svcCtx.WSManager.Broadcast <- message
	logx.Info("消息已发送到广播通道")

	// 返回响应
	resp = &types.SendChatMessageResponse{
		Id:         message.ID,
		Content:    message.Content,
		CreateTime: message.Timestamp,
	}

	logx.Infof("发送消息处理完成，返回响应: %+v", resp)
	return resp, nil
}
