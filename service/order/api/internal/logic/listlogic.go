package logic

import (
	"context"

	"mall/service/order/api/internal/code"
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
	l.Logger.Infof("API获取订单列表请求: %+v", req)

	// 参数校验
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 调用 RPC 获取订单列表
	res, err := l.svcCtx.OrderRpc.List(l.ctx, &order.ListRequest{
		Uid:      req.Uid, // 如果uid为0，则查询所有订单
		Status:   req.Status,
		Page:     req.Page,
		PageSize: req.PageSize,
	})

	if err != nil {
		l.Logger.Errorf("RPC获取订单列表失败: %v", err)
		return nil, code.NewCodeError(code.Error)
	}

	// 转换订单列表
	var orders []types.Order
	for _, o := range res.Data {
		// 转换订单商品列表
		var items []types.OrderItem
		for _, item := range o.Items {
			items = append(items, types.OrderItem{
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
			Id:          o.Id,
			Oid:         o.Oid,
			Uid:         o.Uid,
			AddressId:   o.AddressId,
			TotalPrice:  o.TotalPrice,
			ShippingFee: o.ShippingFee,
			Status:      o.Status,
			StatusDesc:  o.StatusDesc,
			Remark:      o.Remark,
			Items:       items,
			CreateTime:  o.CreateTime,
			UpdateTime:  o.UpdateTime,
		})
	}

	l.Logger.Info("获取订单列表成功")
	return &types.ListOrderResp{
		Code: 0,
		Msg:  "success",
		Data: types.ListOrderData{
			Total: res.Total,
			List:  orders,
		},
	}, nil
}
