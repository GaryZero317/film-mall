package logic

import (
	"context"

	"mall/service/community/api/internal/svc"
	"mall/service/community/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminGetWorkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminGetWorkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminGetWorkLogic {
	return &AdminGetWorkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminGetWorkLogic) AdminGetWork(req *types.GetWorkReq) (resp *types.GetWorkResp, err error) {
	// todo: add your logic here and delete this line

	return
}
