package logic

import (
	"context"

	"mall/service/order/api/internal/code"
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
	l.Logger.Infof("API更新订单请求: %+v", req)

	// 参数校验
	if req.Status < 0 || req.Status > 4 {
		return nil, code.NewCodeError(code.OrderStatusInvalid)
	}

	// 调用 RPC 更新订单
	_, err = l.svcCtx.OrderRpc.Update(l.ctx, &order.UpdateRequest{
		Id:     req.Id,
		Status: req.Status,
	})

	if err != nil {
		l.Logger.Errorf("RPC更新订单失败: %v", err)
		return nil, code.NewCodeError(code.OrderUpdateFailed)
	}

	l.Logger.Info("更新订单成功")
	return &types.UpdateOrderResp{
		Code: 0,
		Msg:  "success",
	}, nil
}
