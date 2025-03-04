package logic

import (
	"context"

	"mall/service/community/api/internal/svc"
	"mall/service/community/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminDeleteCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminDeleteCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminDeleteCommentLogic {
	return &AdminDeleteCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminDeleteCommentLogic) AdminDeleteComment(req *types.DeleteCommentReq) (resp *types.DeleteCommentResp, err error) {
	// todo: add your logic here and delete this line

	return
}
