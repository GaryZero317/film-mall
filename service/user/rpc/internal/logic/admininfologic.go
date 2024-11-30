package logic

import (
	"context"

	"mall/service/user/model"
	"mall/service/user/rpc/internal/svc"
	"mall/service/user/rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAdminInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminInfoLogic {
	return &AdminInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AdminInfoLogic) AdminInfo(in *user.AdminInfoRequest) (*user.AdminInfoResponse, error) {
	// 从数据库中查询管理员信息
	admin, err := l.svcCtx.AdminModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, err
		}
		return nil, err
	}

	return &user.AdminInfoResponse{
		Id:       admin.ID,
		Username: admin.Username,
		Level:    int32(admin.Level),
	}, nil
}
