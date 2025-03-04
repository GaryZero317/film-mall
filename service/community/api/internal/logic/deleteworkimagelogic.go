package logic

import (
	"context"

	"mall/service/community/api/internal/svc"
	"mall/service/community/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteWorkImageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteWorkImageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteWorkImageLogic {
	return &DeleteWorkImageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteWorkImageLogic) DeleteWorkImage(req *types.DeleteWorkImageReq) (resp *types.DeleteWorkImageResp, err error) {
	// todo: add your logic here and delete this line

	return
}
