package logic

import (
	"context"
	"mall/common/errorx"
	"mall/service/user/api/internal/svc"
	"mall/service/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteAdminLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAdminLogic {
	return &DeleteAdminLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteAdminLogic) DeleteAdmin(req *types.DeleteAdminRequest) (resp *types.DeleteAdminResponse, err error) {
	// 获取要删除的管理员信息
	admin, err := l.svcCtx.AdminModel.FindOne(l.ctx, req.Id)
	if err != nil {
		logx.Errorf("查询管理员失败: %v", err)
		return nil, errorx.NewDefaultError("管理员不存在")
	}

	// 执行删除操作
	err = l.svcCtx.AdminModel.Delete(l.ctx, req.Id)
	if err != nil {
		logx.Errorf("删除管理员失败: %v", err)
		return nil, errorx.NewDefaultError("删除管理员失败")
	}

	return &types.DeleteAdminResponse{
		Id:       admin.ID,
		Username: admin.Username,
		Level:    int32(admin.Level),
	}, nil
}
