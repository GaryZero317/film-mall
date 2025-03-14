package logic

import (
	"context"
	"fmt"
	"time"

	"mall/service/order/model"
	"mall/service/order/rpc/internal/svc"
	"mall/service/order/rpc/types"
	"mall/service/product/rpc/pb/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCancelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelLogic {
	return &CancelLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CancelLogic) Cancel(in *types.CancelRequest) (*types.CancelResponse, error) {
	l.Logger.Infof("取消订单: 订单ID=%d", in.Id)

	// 查询订单当前状态
	order, err := l.svcCtx.OrderModel.FindOne(l.ctx, in.Id)
	if err != nil {
		l.Logger.Errorf("查询订单失败: %v", err)
		return nil, fmt.Errorf("查询订单失败: %v", err)
	}

	// 检查订单状态，只有未支付状态的订单才能被取消
	if order.Status != int64(types.OrderStatus_UNPAID) {
		l.Logger.Errorf("订单状态不正确，当前状态: %d", order.Status)
		return nil, fmt.Errorf("订单状态不正确，无法取消")
	}

	// 查询订单商品
	orderItems, err := l.svcCtx.OrderItemModel.FindByOrderId(l.ctx, in.Id)
	if err != nil {
		l.Logger.Errorf("查询订单商品失败: %v", err)
		return nil, fmt.Errorf("查询订单商品失败: %v", err)
	}

	// 更新订单为已取消状态
	err = l.svcCtx.OrderModel.Update(l.ctx, &model.Order{
		Id:         in.Id,
		Status:     int64(types.OrderStatus_CANCELLED), // 已取消状态
		StatusDesc: "已取消",
		UpdateTime: time.Now(),
	})

	if err != nil {
		l.Logger.Errorf("更新订单状态失败: %v", err)
		return nil, fmt.Errorf("更新订单状态失败: %v", err)
	}

	// 恢复库存
	for _, item := range orderItems {
		_, err = l.svcCtx.ProductRpc.RestoreStock(l.ctx, &product.RestoreStockRequest{
			Id:       item.Pid,
			Quantity: item.Quantity,
		})

		if err != nil {
			l.Logger.Errorf("恢复商品库存失败: 商品ID=%d, 数量=%d, 错误=%v",
				item.Pid, item.Quantity, err)
			// 继续处理其他商品，不影响取消流程
		}
	}

	// 从订单超时队列中移除
	timeoutQueueKey := "order:timeout:queue"
	_, err = l.svcCtx.Redis.ZRem(l.ctx, timeoutQueueKey, fmt.Sprintf("%d", in.Id)).Result()
	if err != nil {
		l.Logger.Errorf("从超时队列移除订单失败: %v", err)
		// 不影响取消流程，继续执行
	}

	// 删除订单支付锁
	lockKey := fmt.Sprintf("order:lock:%d", in.Id)
	_, err = l.svcCtx.Redis.Del(l.ctx, lockKey).Result()
	if err != nil {
		l.Logger.Errorf("删除订单支付锁失败: %v", err)
		// 不影响取消流程，继续执行
	}

	l.Logger.Infof("订单取消成功: 订单ID=%d", in.Id)
	return &types.CancelResponse{}, nil
}
