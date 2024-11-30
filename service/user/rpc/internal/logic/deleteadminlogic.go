package logic

import (
	"context"

	"mall/service/user/model"
	"mall/service/user/rpc/internal/svc"
	"mall/service/user/rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/status"
)

type DeleteAdminLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAdminLogic {
	return &DeleteAdminLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteAdminLogic) DeleteAdmin(in *user.DeleteAdminRequest) (*user.DeleteAdminResponse, error) {
	// 检查管理员是否存在
	admin, err := l.svcCtx.AdminModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, status.Error(100, "管理员不存在")
		}
		return nil, status.Error(500, err.Error())
	}

	l.Logger.Infof("DEBUG DeleteAdmin - Deleting admin: ID=%d, Username=%s", 
		admin.ID, admin.Username)

	// 删除管理员
	err = l.svcCtx.AdminModel.Delete(l.ctx, in.Id)
	if err != nil {
		return nil, status.Error(500, err.Error())
	}

	response := &user.DeleteAdminResponse{
		Success: true,
	}

	l.Logger.Info("DEBUG DeleteAdmin - Admin deleted successfully")

	return response, nil
}
