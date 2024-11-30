package logic

import (
	"context"

	"mall/common/cryptx"
	"mall/service/user/model"
	"mall/service/user/rpc/internal/svc"
	"mall/service/user/rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/status"
)

type UpdateAdminLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAdminLogic {
	return &UpdateAdminLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateAdminLogic) UpdateAdmin(in *user.UpdateAdminRequest) (*user.UpdateAdminResponse, error) {
	// 检查管理员是否存在
	admin, err := l.svcCtx.AdminModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, status.Error(100, "管理员不存在")
		}
		return nil, status.Error(500, err.Error())
	}

	l.Logger.Infof("DEBUG UpdateAdmin - Updating admin: ID=%d, Username=%s", 
		admin.ID, admin.Username)

	// 更新管理员信息
	if in.Password != "" {
		admin.Password = cryptx.PasswordEncrypt(l.svcCtx.Config.Salt, in.Password)
	}
	admin.Level = int(in.Level)

	err = l.svcCtx.AdminModel.Update(l.ctx, admin)
	if err != nil {
		return nil, status.Error(500, err.Error())
	}

	response := &user.UpdateAdminResponse{
		Id:       admin.ID,
		Username: admin.Username,
		Level:    int32(admin.Level),
	}

	l.Logger.Infof("DEBUG UpdateAdmin - Admin updated successfully: ID=%d, Username=%s, Level=%d",
		response.Id, response.Username, response.Level)

	return response, nil
}
