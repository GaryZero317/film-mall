package logic

import (
	"context"

	"mall/service/order/api/internal/code"
	"mall/service/order/api/internal/svc"
	"mall/service/order/api/internal/types"
	"mall/service/order/rpc/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLogic {
	return &GetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLogic) Get(req *types.GetOrderReq) (resp *types.GetOrderResp, err error) {
	l.Logger.Infof("API获取订单请求: %+v", req)

	// 调用 RPC 获取订单
	res, err := l.svcCtx.OrderRpc.Detail(l.ctx, &order.DetailRequest{
		Id: req.Id,
	})

	if err != nil {
		l.Logger.Errorf("RPC获取订单失败: %v", err)
		return nil, code.NewCodeError(code.OrderNotExist)
	}

	// 转换订单商品列表
	var items []types.OrderItem
	for _, item := range res.Items {
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

	l.Logger.Info("获取订单成功")
	return &types.GetOrderResp{
		Code: 0,
		Msg:  "success",
		Data: types.Order{
			Id:          res.Id,
			Oid:         res.Oid,
			Uid:         res.Uid,
			AddressId:   res.AddressId,
			TotalPrice:  res.TotalPrice,
			ShippingFee: res.ShippingFee,
			Status:      res.Status,
			StatusDesc:  res.StatusDesc,
			Remark:      res.Remark,
			Items:       items,
			CreateTime:  res.CreateTime,
			UpdateTime:  res.UpdateTime,
		},
	}, nil
}

// 获取订单状态描述
func getStatusDesc(status int64) string {
	switch status {
	case 0:
		return "待付款"
	case 1:
		return "待发货"
	case 2:
		return "待收货"
	case 3:
		return "已完成"
	default:
		return "未知状态"
	}
}
