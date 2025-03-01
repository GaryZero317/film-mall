package logic

import (
	"context"
	"errors"
	"mall/service/user/model"

	"mall/service/user/api/internal/svc"
	"mall/service/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetChatHistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetChatHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatHistoryLogic {
	return &GetChatHistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetChatHistoryLogic) GetChatHistory(req *types.ChatHistoryRequest) (resp *types.ChatHistoryResponse, err error) {
	userId, ok := l.ctx.Value("userId").(int64)
	if !ok || userId <= 0 {
		return nil, errors.New("用户未授权")
	}

	// 设置默认分页参数
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 20
	}

	// 查询聊天历史
	// 注意：这里需要聊天消息模型支持相关查询方法
	messages, total, err := l.svcCtx.ChatMessageModel.FindByUserAndAdmin(
		l.ctx,
		userId,
		req.AdminId,
		int(page),
		int(pageSize),
	)
	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	items := make([]types.ChatMessageItem, 0, len(messages))
	for _, msg := range messages {
		var userId, adminId int64
		var direction int64

		// 根据发送者类型确定用户ID和管理员ID
		if msg.SenderType == model.SenderTypeUser {
			userId = msg.SenderId
			adminId = 0   // 可能需要从其他地方获取adminId
			direction = 1 // 用户到管理员
		} else {
			userId = msg.SessionId // 会话ID即为用户ID
			adminId = msg.SenderId
			direction = 2 // 管理员到用户
		}

		items = append(items, types.ChatMessageItem{
			Id:         msg.Id,
			UserId:     userId,
			AdminId:    adminId,
			Direction:  direction,
			Content:    msg.Content,
			ReadStatus: msg.ReadStatus,
			CreateTime: msg.CreateTime.Unix(),
		})
	}

	return &types.ChatHistoryResponse{
		Total: total,
		List:  items,
	}, nil
}
