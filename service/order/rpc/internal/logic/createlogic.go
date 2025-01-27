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

type CreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateLogic) Create(in *types.CreateRequest) (*types.CreateResponse, error) {
	l.Logger.Infof("RPC创建订单请求: %+v", in)

	// 生成订单号
	now := time.Now()
	oid := fmt.Sprintf("%s%d", now.Format("20060102150405"), in.Uid)

	// 构建订单商品数据
	var orderItems []model.OrderItem
	for _, item := range in.Items {
		orderItems = append(orderItems, model.OrderItem{
			OrderId:      0, // 创建订单后更新
			Pid:          item.Pid,
			ProductName:  item.ProductName,
			ProductImage: item.ProductImage,
			Price:        item.Price,
			Quantity:     item.Quantity,
			Amount:       item.Amount,
		})
	}

	l.Logger.Info("开始创建订单")
	// 创建订单
	order := model.Order{
		Oid:         oid,
		Uid:         in.Uid,
		AddressId:   in.AddressId,
		TotalPrice:  in.TotalPrice,
		ShippingFee: in.ShippingFee,
		Status:      in.Status,
		Remark:      in.Remark,
		CreateTime:  now,
		UpdateTime:  now,
	}

	// 开启事务
	err := l.svcCtx.OrderModel.Trans(l.ctx, func(ctx context.Context, session model.Session) error {
		l.Logger.Info("开始事务")
		// 1. 创建订单
		result, err := l.svcCtx.OrderModel.Insert(ctx, session, &order)
		if err != nil {
			l.Logger.Errorf("创建订单失败: %v", err)
			return err
		}

		orderId, err := result.LastInsertId()
		if err != nil {
			l.Logger.Errorf("获取订单ID失败: %v", err)
			return err
		}

		l.Logger.Info("开始创建订单商品")
		// 2. 创建订单商品
		for i := range orderItems {
			orderItems[i].OrderId = orderId
			_, err := l.svcCtx.OrderItemModel.Insert(ctx, session, &orderItems[i])
			if err != nil {
				l.Logger.Errorf("创建订单商品失败: %v", err)
				return err
			}
		}

		l.Logger.Info("事务完成")
		return nil
	})

	if err != nil {
		l.Logger.Errorf("事务执行失败: %v", err)
		return nil, err
	}

	l.Logger.Infof("创建订单成功: %+v", order)
	return &types.CreateResponse{
		Id:  order.Id,
		Oid: order.Oid,
	}, nil
}
