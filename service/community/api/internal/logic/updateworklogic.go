package logic

import (
	"context"

	"mall/service/community/api/internal/svc"
	"mall/service/community/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateWorkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateWorkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateWorkLogic {
	return &UpdateWorkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateWorkLogic) UpdateWork(req *types.UpdateWorkReq) (resp *types.UpdateWorkResp, err error) {
	// todo: add your logic here and delete this line

	return
}
