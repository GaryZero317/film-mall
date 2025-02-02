package logic

import (
	"context"

	"mall/service/order/model"
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
	// 获取订单列表
	var orders []*model.Order
	var total int64
	var err error

	if in.Uid == 0 {
		// 管理员查询所有订单
		orders, total, err = l.svcCtx.OrderModel.FindAll(l.ctx, in.Status, in.Page, in.PageSize)
	} else {
		// 普通用户查询自己的订单
		orders, total, err = l.svcCtx.OrderModel.FindByUid(l.ctx, in.Uid, in.Status, in.Page, in.PageSize)
	}

	if err != nil {
		return nil, err
	}

	var orderList []*types.DetailResponse
	for _, order := range orders {
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

		// 使用model层已经设置好的StatusDesc
		orderList = append(orderList, &types.DetailResponse{
			Id:          order.Id,
			Oid:         order.Oid,
			Uid:         order.Uid,
			AddressId:   order.AddressId,
			TotalPrice:  order.TotalPrice,
			ShippingFee: order.ShippingFee,
			Status:      order.Status,
			StatusDesc:  order.StatusDesc, // 直接使用model层设置的状态描述
			Remark:      order.Remark,
			Items:       orderItems,
			CreateTime:  order.CreateTime.Format("2006-01-02 15:04:05"),
			UpdateTime:  order.UpdateTime.Format("2006-01-02 15:04:05"),
		})
	}

	return &types.ListResponse{
		Total: total,
		Data:  orderList,
	}, nil
}
