package logic

import (
	"context"
	"fmt"
	"time"

	"mall/service/order/model"
	"mall/service/order/rpc/internal/svc"
	"mall/service/order/rpc/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateLogic) Update(in *types.UpdateRequest) (*types.UpdateResponse, error) {
	l.Logger.Infof("更新订单状态请求: %+v", in)

	// 根据状态自动设置状态描述
	statusDesc := ""
	switch in.Status {
	case 0:
		statusDesc = "待付款"
	case 1:
		statusDesc = "待发货"
	case 2:
		statusDesc = "待收货"
	case 3:
		statusDesc = "已完成"
	case 4:
		statusDesc = "已取消"
	default:
		return nil, fmt.Errorf("无效的订单状态: %d", in.Status)
	}

	// 更新订单状态
	err := l.svcCtx.OrderModel.Update(l.ctx, &model.Order{
		Id:         in.Id,
		Status:     in.Status,
		StatusDesc: statusDesc,
		UpdateTime: time.Now(),
	})

	if err != nil {
		l.Logger.Errorf("更新订单失败: %v", err)
		return nil, err
	}

	l.Logger.Info("更新订单成功")
	return &types.UpdateResponse{}, nil
}
