package logic

import (
	"context"

	"mall/service/community/api/internal/svc"
	"mall/service/community/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminUpdateWorkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminUpdateWorkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUpdateWorkLogic {
	return &AdminUpdateWorkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminUpdateWorkLogic) AdminUpdateWork(req *types.UpdateWorkReq) (resp *types.UpdateWorkResp, err error) {
	// todo: add your logic here and delete this line

	return
}
