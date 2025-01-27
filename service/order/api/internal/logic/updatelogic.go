package logic

import (
	"context"

	"mall/service/order/api/internal/svc"
	"mall/service/order/api/internal/types"
	"mall/service/order/rpc/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.UpdateOrderReq) (resp *types.UpdateOrderResp, err error) {
	// 调用 RPC 更新订单
	_, err = l.svcCtx.OrderRpc.Update(l.ctx, &order.UpdateRequest{
		Id:     req.Id,
		Status: req.Status,
	})
	if err != nil {
		return nil, err
	}

	return &types.UpdateOrderResp{}, nil
}
