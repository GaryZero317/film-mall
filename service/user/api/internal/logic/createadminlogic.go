package logic

import (
	"context"
	"mall/common/errorx"
	"mall/service/user/api/internal/svc"
	"mall/service/user/api/internal/types"
	"mall/service/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateAdminLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAdminLogic {
	return &CreateAdminLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateAdminLogic) CreateAdmin(req *types.CreateAdminRequest) (resp *types.CreateAdminResponse, err error) {
	// 调用 RPC 服务创建管理员
	res, err := l.svcCtx.UserRpc.CreateAdmin(l.ctx, &user.CreateAdminRequest{
		Username: req.Username,
		Password: req.Password,
		Level:    1, // 直接设置默认值为1
	})
	if err != nil {
		return nil, errorx.NewDefaultError(err.Error())
	}

	if res == nil {
		return nil, errorx.NewDefaultError("创建管理员失败")
	}

	return &types.CreateAdminResponse{
		Id:       res.Id,
		Username: res.Username,
		Level:    res.Level,
	}, nil
}
