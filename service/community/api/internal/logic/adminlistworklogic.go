package logic

import (
	"context"

	"mall/service/community/api/internal/svc"
	"mall/service/community/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListWorkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminListWorkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListWorkLogic {
	return &AdminListWorkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminListWorkLogic) AdminListWork(req *types.ListWorkReq) (resp *types.ListWorkResp, err error) {
	// todo: add your logic here and delete this line

	return
}
