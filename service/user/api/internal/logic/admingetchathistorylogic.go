package logic

import (
	"context"

	"mall/service/user/api/internal/svc"
	"mall/service/user/api/internal/types"

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
	// todo: add your logic here and delete this line

	return
}
