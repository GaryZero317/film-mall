package logic

import (
	"context"
	"mall/common/jwtx"
	"mall/service/user/api/internal/svc"
	"mall/service/user/api/internal/types"
	"mall/service/user/rpc/pb/user"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminLoginLogic {
	return &AdminLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminLoginLogic) AdminLogin(req *types.AdminLoginRequest) (resp *types.AdminLoginResponse, err error) {
	// 调用RPC服务进行管理员登录
	adminResp, err := l.svcCtx.UserRpc.AdminLogin(l.ctx, &user.AdminLoginRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	// 生成JWT token
	auth := l.svcCtx.Config.AdminAuth
	now := time.Now().Unix()
	accessExpire := auth.AccessExpire
	
	l.Logger.Infof("=== DEBUG Token Generation ===")
	l.Logger.Infof("Secret: %s", auth.AccessSecret)
	l.Logger.Infof("Now: %d", now)
	l.Logger.Infof("Expire: %d", accessExpire)
	l.Logger.Infof("AdminID: %d", adminResp.Id)
	
	token, err := jwtx.GetToken(auth.AccessSecret, now, accessExpire, adminResp.Id)
	if err != nil {
		return nil, err
	}
	
	l.Logger.Infof("=== Generated Token ===")
	l.Logger.Infof("Token: %s", token)
	l.Logger.Infof("Use this header in your next request:")
	l.Logger.Infof("Authorization: Bearer %s", token)
	l.Logger.Infof("=====================")

	return &types.AdminLoginResponse{
		AccessToken:  token,
		AccessExpire: now + accessExpire,
	}, nil
}
