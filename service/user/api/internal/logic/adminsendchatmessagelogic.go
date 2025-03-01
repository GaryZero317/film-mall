package logic

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"mall/service/user/api/internal/svc"
	"mall/service/user/api/internal/types"
	"mall/service/user/api/internal/ws"
	"mall/service/user/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminSendChatMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminSendChatMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminSendChatMessageLogic {
	return &AdminSendChatMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminSendChatMessageLogic) AdminSendChatMessage(req *types.AdminSendChatMessageRequest) (resp *types.AdminSendChatMessageResponse, err error) {
	logx.Infof("接收到管理员发送消息请求：用户ID=%d, 内容=%s, 消息类型=%d", req.UserId, req.Content, req.Type)

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

	logx.Infof("成功获取管理员ID: %d，准备创建消息", adminId)

	// 默认文本消息
	msgType := int(req.Type)
	if msgType == 0 {
		msgType = ws.TextMessage
		logx.Info("未指定消息类型，使用默认文本消息类型")
	}

	// 创建ChatMessage对象用于数据库存储 - 使用新的数据库结构
	chatMsg := &model.ChatMessage{
		SessionId:  req.UserId,            // 使用用户ID作为会话ID
		SenderType: model.SenderTypeAdmin, // 管理员发送
		SenderId:   adminId,               // 管理员ID
		Content:    req.Content,           // 消息内容
		ReadStatus: 0,                     // 初始为未读
		CreateTime: time.Now(),            // 当前时间
	}
	logx.Infof("创建消息对象: %+v", chatMsg)

	// 检查WebSocket管理器状态
	if l.svcCtx.WSManager == nil {
		logx.Error("WebSocket管理器未初始化")
		return nil, errors.New("WebSocket服务不可用")
	}

	logx.Info("WebSocket管理器状态正常，开始处理消息")

	// 使用ChatMessageModel存储消息
	result, err := l.svcCtx.ChatMessageModel.Insert(l.ctx, chatMsg)
	if err != nil {
		logx.Errorf("存储消息失败: %v", err)
		return nil, err
	}

	// 获取消息ID
	id, err := result.LastInsertId()
	if err != nil {
		logx.Errorf("获取消息ID失败: %v", err)
		return nil, err
	}
	logx.Infof("消息存储成功，ID: %d", id)

	// 创建WebSocket消息并发送到广播通道
	wsMessage := &ws.Message{
		ID:        id,
		UserId:    req.UserId,
		AdminId:   adminId,
		Direction: ws.AdminToUser,
		Content:   req.Content,
		Type:      msgType,
		Timestamp: time.Now().Unix(),
	}

	// 广播消息 - 注意这里是往通道发送，不是函数调用
	logx.Info("发送消息到广播通道")
	l.svcCtx.WSManager.Broadcast <- wsMessage
	logx.Info("消息已发送到广播通道")

	// 返回响应
	resp = &types.AdminSendChatMessageResponse{
		Id:         id,
		Content:    req.Content,
		CreateTime: time.Now().Unix(),
	}

	logx.Infof("管理员发送消息处理完成，返回响应: %+v", resp)
	return resp, nil
}
