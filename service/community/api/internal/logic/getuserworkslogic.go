package logic

import (
	"context"

	"mall/service/community/api/internal/svc"
	"mall/service/community/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserWorksLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserWorksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserWorksLogic {
	return &GetUserWorksLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserWorksLogic) GetUserWorks(req *types.ListWorkReq) (resp *types.ListWorkResp, err error) {
	// todo: add your logic here and delete this line

	return
}
