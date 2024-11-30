package logic

import (
	"context"
	"errors"

	"mall/service/user/api/internal/svc"
	"mall/service/user/api/internal/types"
	"mall/service/user/rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminInfoLogic {
	return &AdminInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminInfoLogic) AdminInfo() (resp *types.AdminInfoResponse, err error) {
	// 从ctx中获取当前登录的管理员ID
	adminId, ok := l.ctx.Value("id").(int64)
	if !ok {
		return nil, errors.New("invalid admin id")
	}

	// 调用RPC服务获取管理员信息
	adminInfo, err := l.svcCtx.UserRpc.AdminInfo(l.ctx, &user.AdminInfoRequest{
		Id: adminId,
	})
	if err != nil {
		return nil, err
	}

	return &types.AdminInfoResponse{
		Id:       adminInfo.Id,
		Username: adminInfo.Username,
		Level:    adminInfo.Level,
	}, nil
}
