package logic

import (
	"context"
	"time"

	"mall/service/user/api/internal/svc"
	"mall/service/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminServiceUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminServiceUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminServiceUpdateLogic {
	return &AdminServiceUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminServiceUpdateLogic) AdminServiceUpdate(req *types.AdminServiceUpdateRequest) (resp *types.AdminServiceUpdateResponse, err error) {
	// 查询问题是否存在
	serviceQuestion, err := l.svcCtx.GormCustomerServiceModel.FindOne(l.ctx, req.Id)
	if err != nil {
		l.Logger.Errorf("查询客服问题失败: %v", err)
		return &types.AdminServiceUpdateResponse{
			Code: 500,
			Msg:  "问题不存在或系统错误",
		}, nil
	}

	// 更新回复内容
	serviceQuestion.Reply = req.Reply
	now := time.Now()
	serviceQuestion.ReplyTime = &now

	// 更新状态，如果未指定则默认为已回复(1)
	if req.Status > 0 {
		serviceQuestion.Status = req.Status
	} else {
		serviceQuestion.Status = 1 // 默认已回复状态
	}

	// 更新时间
	serviceQuestion.UpdateTime = time.Now()

	// 保存更新
	err = l.svcCtx.GormCustomerServiceModel.Update(l.ctx, serviceQuestion)
	if err != nil {
		l.Logger.Errorf("更新客服问题失败: %v", err)
		return &types.AdminServiceUpdateResponse{
			Code: 500,
			Msg:  "更新失败，系统错误",
		}, nil
	}

	return &types.AdminServiceUpdateResponse{
		Code: 0,
		Msg:  "回复成功",
	}, nil
}
