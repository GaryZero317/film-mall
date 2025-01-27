package logic

import (
	"context"

	"mall/service/order/api/internal/svc"
	"mall/service/order/api/internal/types"
	"mall/service/order/rpc/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLogic {
	return &ListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListLogic) List(req *types.ListOrderReq) (resp *types.ListOrderResp, err error) {
	// 调用 RPC 获取订单列表
	res, err := l.svcCtx.OrderRpc.List(l.ctx, &order.ListRequest{
		Uid:      req.Uid,
		Status:   req.Status,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, err
	}

	var orders []types.Order
	for _, orderDetail := range res.Data {
		var orderItems []types.OrderItem
		for _, item := range orderDetail.Items {
			orderItems = append(orderItems, types.OrderItem{
				Id:           item.Id,
				OrderId:      item.OrderId,
				Pid:          item.Pid,
				ProductName:  item.ProductName,
				ProductImage: item.ProductImage,
				Price:        item.Price,
				Quantity:     item.Quantity,
				Amount:       item.Amount,
			})
		}

		orders = append(orders, types.Order{
			Id:          orderDetail.Id,
			Oid:         orderDetail.Oid,
			Uid:         orderDetail.Uid,
			AddressId:   orderDetail.AddressId,
			TotalPrice:  orderDetail.TotalPrice,
			ShippingFee: orderDetail.ShippingFee,
			Status:      orderDetail.Status,
			StatusDesc:  orderDetail.StatusDesc,
			Remark:      orderDetail.Remark,
			Items:       orderItems,
			CreateTime:  orderDetail.CreateTime,
			UpdateTime:  orderDetail.UpdateTime,
		})
	}

	return &types.ListOrderResp{
		Total: res.Total,
		Data:  orders,
	}, nil
}
