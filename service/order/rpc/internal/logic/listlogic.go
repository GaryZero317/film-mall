package logic

import (
	"context"

	"mall/service/order/rpc/internal/svc"
	"mall/service/order/rpc/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLogic {
	return &ListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListLogic) List(in *types.ListRequest) (*types.ListResponse, error) {
	// 查询用户订单列表
	orders, err := l.svcCtx.OrderModel.FindAllByUid(l.ctx, in.Uid)
	if err != nil {
		return nil, err
	}

	var orderList []*types.DetailResponse
	for _, order := range orders {
		orderList = append(orderList, &types.DetailResponse{
			Id:         order.Id,
			Oid:        order.Oid,
			Uid:        order.Uid,
			Pid:        order.Pid,
			Amount:     order.Amount,
			Status:     order.Status,
			CreateTime: order.CreateTime.Format("2006-01-02 15:04:05"),
			UpdateTime: order.UpdateTime.Format("2006-01-02 15:04:05"),
		})
	}

	return &types.ListResponse{
		Data: orderList,
	}, nil
}
