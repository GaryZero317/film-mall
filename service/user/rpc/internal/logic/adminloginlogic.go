package logic

import (
	"context"

	"mall/service/user/rpc/internal/svc"
	"mall/service/user/rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminLoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAdminLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminLoginLogic {
	return &AdminLoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  管理员服务
func (l *AdminLoginLogic) AdminLogin(in *user.AdminLoginRequest) (*user.AdminLoginResponse, error) {
	// todo: add your logic here and delete this line

	return &user.AdminLoginResponse{}, nil
}
