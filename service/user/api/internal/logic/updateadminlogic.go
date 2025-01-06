package logic

import (
	"context"
	"mall/common/errorx"
	"mall/service/user/api/internal/svc"
	"mall/service/user/api/internal/types"

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
	// 检查管理员是否存在
	admin, err := l.svcCtx.AdminModel.FindOne(l.ctx, req.Id)
	if err != nil {
		logx.Errorf("查询管理员失败: %v", err)
		return nil, errorx.NewDefaultError("管理员不存在")
	}

	// 更新管理员信息
	if req.Password != "" {
		admin.Password = req.Password
	}
	admin.Level = int(req.Level)

	// 保存更新
	err = l.svcCtx.AdminModel.Update(l.ctx, admin)
	if err != nil {
		logx.Errorf("更新管理员失败: %v", err)
		return nil, errorx.NewDefaultError("更新管理员失败")
	}

	return &types.UpdateAdminResponse{
		Id:       admin.ID,
		Username: admin.Username,
		Level:    int32(admin.Level),
	}, nil
}
