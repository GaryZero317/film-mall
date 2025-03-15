package logic

import (
	"context"
	"fmt"
	"time"

	"mall/service/order/model"
	"mall/service/order/rpc/internal/svc"
	"mall/service/order/rpc/types"
	payclient "mall/service/pay/rpc/pay"
	productclient "mall/service/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateLogic) Update(in *types.UpdateRequest) (*types.UpdateResponse, error) {
	l.Logger.Infof("更新订单状态请求: %+v", in)

	// 根据状态自动设置状态描述
	statusDesc := ""
	switch in.Status {
	case 0:
		statusDesc = "待付款"
	case 1:
		statusDesc = "待发货"
	case 2:
		statusDesc = "待收货"
	case 3:
		statusDesc = "已完成"
	case 4:
		statusDesc = "已取消"
	case 9:
		statusDesc = "已取消(超时)"
	default:
		return nil, fmt.Errorf("无效的订单状态: %d", in.Status)
	}

	// 查询订单，如果状态为0且要取消(4或9)，则需要恢复库存
	var needRestoreStock bool = false
	var originalStatus int64 = -1

	if in.Status == 4 || in.Status == 9 {
		order, err := l.svcCtx.OrderModel.FindOne(l.ctx, in.Id)
		if err == nil {
			originalStatus = order.Status
			// 如果是从待支付状态变更为取消状态，需要恢复库存
			if originalStatus == 0 {
				needRestoreStock = true
			}
		}
	}

	// 更新订单状态
	err := l.svcCtx.OrderModel.Update(l.ctx, &model.Order{
		Id:         in.Id,
		Status:     in.Status,
		StatusDesc: statusDesc,
		UpdateTime: time.Now(),
	})

	if err != nil {
		l.Logger.Errorf("更新订单失败: %v", err)
		return nil, err
	}

	// 如果需要恢复库存
	if needRestoreStock {
		// 查询订单商品
		orderItems, err := l.svcCtx.OrderItemModel.FindByOrderId(l.ctx, in.Id)
		if err != nil {
			l.Logger.Errorf("查询订单商品失败: %v", err)
			// 继续流程，不影响订单状态更新
		} else {
			// 恢复库存
			for _, item := range orderItems {
				_, err = l.svcCtx.ProductRpc.RestoreStock(l.ctx, &productclient.RestoreStockRequest{
					Id:       item.Pid,
					Quantity: item.Quantity,
				})

				if err != nil {
					l.Logger.Errorf("恢复商品库存失败: 商品ID=%d, 数量=%d, 错误=%v",
						item.Pid, item.Quantity, err)
					// 继续处理其他商品，不影响取消流程
				}
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
	}

	// 如果订单状态为取消且支付服务可用，更新支付记录状态
	// 注意：这里不再依赖needRestoreStock条件
	if (in.Status == 4 || in.Status == 9) && l.svcCtx.PayRpc != nil {
		l.Logger.Info("准备更新支付记录状态为已取消", logx.Field("orderId", in.Id))

		// 查询订单获取用户ID
		var userId int64 = 0
		order, err := l.svcCtx.OrderModel.FindOne(l.ctx, in.Id)
		if err == nil {
			userId = order.Uid
		}

		// 更新支付记录状态为已取消
		_, callbackErr := l.svcCtx.PayRpc.Callback(l.ctx, &payclient.CallbackRequest{
			Id:     0,      // 不指定支付ID，让支付服务根据订单ID查找
			Oid:    in.Id,  // 订单ID
			Uid:    userId, // 用户ID
			Status: 9,      // 支付状态：9表示已取消
		})

		if callbackErr != nil {
			l.Logger.Errorf("更新支付记录状态失败: 订单ID=%d, 错误=%v", in.Id, callbackErr)
			// 不影响订单取消流程，继续执行
		} else {
			l.Logger.Infof("更新支付记录状态成功: 订单ID=%d, 状态=已取消", in.Id)
		}
	}

	l.Logger.Info("更新订单成功")
	return &types.UpdateResponse{}, nil
}
