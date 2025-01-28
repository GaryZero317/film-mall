package logic

import (
	"context"

	"mall/service/order/api/internal/code"
	"mall/service/order/api/internal/svc"
	"mall/service/order/api/internal/types"
	"mall/service/order/rpc/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(req *types.CreateOrderReq) (resp *types.CreateOrderResp, err error) {
	l.Logger.Infof("API创建订单请求: %+v", req)

	// 参数校验
	if req.Uid == 0 {
		return nil, code.NewCodeError(code.InvalidParams)
	}
	if req.AddressId == 0 {
		return nil, code.NewCodeError(code.InvalidParams)
	}
	if req.TotalPrice <= 0 {
		return nil, code.NewCodeError(code.OrderAmountInvalid)
	}
	if len(req.Items) == 0 {
		return nil, code.NewCodeError(code.InvalidParams)
	}

	// 构建订单商品列表
	var items []*order.OrderItem
	for _, item := range req.Items {
		items = append(items, &order.OrderItem{
			Pid:          item.Pid,
			ProductName:  item.ProductName,
			ProductImage: item.ProductImage,
			Price:        item.Price,
			Quantity:     item.Quantity,
			Amount:       item.Amount,
		})
	}

	l.Logger.Info("调用RPC创建订单")
	// 调用RPC创建订单
	res, err := l.svcCtx.OrderRpc.Create(l.ctx, &order.CreateRequest{
		Uid:         req.Uid,
		AddressId:   req.AddressId,
		TotalPrice:  req.TotalPrice,
		ShippingFee: req.ShippingFee,
		Status:      1, // 待支付
		Remark:      req.Remark,
		Items:       items,
	})

	if err != nil {
		l.Logger.Errorf("RPC创建订单失败: %v", err)
		return nil, code.NewCodeError(code.OrderCreateFailed)
	}

	l.Logger.Infof("创建订单成功: %+v", res)
	return &types.CreateOrderResp{
		Code: 0,
		Msg:  "success",
		Data: types.CreateOrderData{
			Id:  res.Id,
			Oid: res.Oid,
		},
	}, nil
}
