package logic

import (
	"context"

	"mall/service/user/api/internal/svc"
	"mall/service/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminFaqDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminFaqDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminFaqDeleteLogic {
	return &AdminFaqDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminFaqDeleteLogic) AdminFaqDelete(req *types.AdminFaqDeleteRequest) (resp *types.AdminFaqDeleteResponse, err error) {
	// 参数校验
	if req.Id <= 0 {
		return &types.AdminFaqDeleteResponse{
			Code: 400,
			Msg:  "参数错误：ID不能为空",
		}, nil
	}

	// 查询FAQ是否存在
	_, err = l.svcCtx.FaqModel.FindOne(l.ctx, req.Id)
	if err != nil {
		l.Logger.Errorf("查询FAQ失败: %v", err)
		return &types.AdminFaqDeleteResponse{
			Code: 500,
			Msg:  "FAQ不存在或系统错误",
		}, nil
	}

	// 删除FAQ
	err = l.svcCtx.FaqModel.Delete(l.ctx, req.Id)
	if err != nil {
		l.Logger.Errorf("删除FAQ失败: %v", err)
		return &types.AdminFaqDeleteResponse{
			Code: 500,
			Msg:  "删除失败，系统错误",
		}, nil
	}

	return &types.AdminFaqDeleteResponse{
		Code: 0,
		Msg:  "删除FAQ成功",
	}, nil
}
