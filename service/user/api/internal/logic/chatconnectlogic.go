package logic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"mall/service/user/api/internal/svc"
)

type ChatConnectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChatConnectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatConnectLogic {
	return &ChatConnectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChatConnectLogic) ChatConnect() error {
	// todo: add your logic here and delete this line

	return nil
}
