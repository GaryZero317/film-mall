package logic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"mall/service/user/api/internal/svc"
)

type AdminChatConnectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminChatConnectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminChatConnectLogic {
	return &AdminChatConnectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminChatConnectLogic) AdminChatConnect() error {
	// todo: add your logic here and delete this line

	return nil
}
