package logic

import (
	"context"
	"mall/common/errorx"
	"mall/service/user/api/internal/svc"
	"mall/service/user/api/internal/types"
	"mall/service/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAdminLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAdminLogic {
	return &UpdateAdminLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateAdminLogic) UpdateAdmin(req *types.UpdateAdminRequest) (resp *types.UpdateAdminResponse, err error) {
	// 调用 RPC 服务更新管理员
	res, err := l.svcCtx.UserRpc.UpdateAdmin(l.ctx, &user.UpdateAdminRequest{
		Id:       req.Id,
		Password: req.Password,
		Level:    req.Level,
	})
	if err != nil {
		return nil, errorx.NewDefaultError(err.Error())
	}

	return &types.UpdateAdminResponse{
		Id:       res.Id,
		Username: res.Username,
		Level:    res.Level,
	}, nil
}
