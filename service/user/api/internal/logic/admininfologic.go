package logic

import (
	"context"
	"encoding/json"

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
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	res, err := l.svcCtx.UserRpc.AdminInfo(l.ctx, &user.AdminInfoRequest{
		Id: uid,
	})
	if err != nil {
		return nil, err
	}

	return &types.AdminInfoResponse{
		Id:       res.Id,
		Username: res.Username,
		Level:    res.Level,
	}, nil
}
