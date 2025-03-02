package logic

import (
	"context"
	"encoding/json"
	"errors"

	"mall/service/user/api/internal/svc"
	"mall/service/user/api/internal/types"
	"mall/service/user/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminGetChatHistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminGetChatHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminGetChatHistoryLogic {
	return &AdminGetChatHistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminGetChatHistoryLogic) AdminGetChatHistory(req *types.ChatHistoryRequest) (resp *types.ChatHistoryResponse, err error) {
	logx.Info("开始处理获取聊天历史请求")
	logx.Infof("请求参数: %+v", req)

	// 获取请求参数
	userId := req.UserId
	page := req.Page
	pageSize := req.PageSize

	// 参数验证
	if userId <= 0 {
		logx.Error("用户ID无效")
		return nil, errors.New("用户ID无效")
	}

	// 设置默认分页参数
	if page <= 0 {
		page = 1
		logx.Info("未指定页码，使用默认值1")
	}
	if pageSize <= 0 {
		pageSize = 20
		logx.Info("未指定每页数量，使用默认值20")
	}

	// 尝试获取管理员ID
	adminId, ok := l.ctx.Value("adminId").(int64)
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
	}

	if !ok || adminId <= 0 {
		logx.Error("获取管理员ID失败或ID无效")
		return nil, errors.New("管理员未授权")
	}

	logx.Infof("成功获取管理员ID: %d", adminId)

	// 获取聊天历史记录
	messages, totalCount, err := l.svcCtx.GetChatHistory(l.ctx, userId, adminId, int(page), int(pageSize))
	if err != nil {
		logx.Errorf("获取聊天历史失败: %v", err)
		return nil, err
	}

	logx.Infof("成功获取聊天历史，共 %d 条消息", len(messages))

	// 转换为响应格式
	items := make([]types.ChatMessageItem, 0, len(messages))
	for _, msg := range messages {
		direction := int64(1) // 默认为用户到管理员
		if msg.SenderType == model.SenderTypeAdmin {
			direction = 2 // 管理员到用户
		}

		items = append(items, types.ChatMessageItem{
			Id:         msg.Id,
			UserId:     msg.SessionId, // 会话ID即用户ID
			AdminId:    msg.SenderId,  // 发送者ID
			Direction:  direction,
			Content:    msg.Content,
			ReadStatus: int64(msg.ReadStatus),
			CreateTime: msg.CreateTime.Unix(),
		})
	}

	// 标记消息为已读
	if len(messages) > 0 {
		// 找出所有未读的用户消息
		var unreadMsgIds []int64
		for _, msg := range messages {
			if msg.SenderType == model.SenderTypeUser && msg.ReadStatus == 0 {
				unreadMsgIds = append(unreadMsgIds, msg.Id)
			}
		}

		// 如果有未读消息，标记为已读
		if len(unreadMsgIds) > 0 {
			logx.Infof("标记 %d 条消息为已读", len(unreadMsgIds))
			err = l.svcCtx.MarkMessagesAsRead(l.ctx, unreadMsgIds)
			if err != nil {
				logx.Errorf("标记消息为已读失败: %v", err)
				// 不中断处理，继续返回消息
			}
		}
	}

	logx.Info("聊天历史处理完成")
	return &types.ChatHistoryResponse{
		Total: totalCount,
		List:  items,
	}, nil
}
