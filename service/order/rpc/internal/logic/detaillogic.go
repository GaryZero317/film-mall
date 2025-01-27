package logic

import (
	"context"

	"mall/service/order/rpc/internal/svc"
	"mall/service/order/rpc/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DetailLogic) Detail(in *types.DetailRequest) (*types.DetailResponse, error) {
	// 获取订单信息
	order, err := l.svcCtx.OrderModel.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	// 获取订单商品列表
	items, err := l.svcCtx.OrderItemModel.FindByOrderId(l.ctx, order.Id)
	if err != nil {
		return nil, err
	}

	// 构建商品列表响应
	var orderItems []*types.OrderItem
	for _, item := range items {
		orderItems = append(orderItems, &types.OrderItem{
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

	return &types.DetailResponse{
		Id:          order.Id,
		Oid:         order.Oid,
		Uid:         order.Uid,
		AddressId:   order.AddressId,
		TotalPrice:  order.TotalPrice,
		ShippingFee: order.ShippingFee,
		Status:      order.Status,
		StatusDesc:  order.StatusDesc,
		Remark:      order.Remark,
		Items:       orderItems,
		CreateTime:  order.CreateTime.Format("2006-01-02 15:04:05"),
		UpdateTime:  order.UpdateTime.Format("2006-01-02 15:04:05"),
	}, nil
}
