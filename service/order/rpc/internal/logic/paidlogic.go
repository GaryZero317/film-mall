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

type PaidLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPaidLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PaidLogic {
	return &PaidLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PaidLogic) Paid(in *types.PaidRequest) (*types.PaidResponse, error) {
	l.Logger.Infof("处理订单支付: 订单ID=%d", in.Id)

	// 查询订单当前状态
	order, err := l.svcCtx.OrderModel.FindOne(l.ctx, in.Id)
	if err != nil {
		l.Logger.Errorf("查询订单失败: %v", err)
		return nil, fmt.Errorf("查询订单失败: %v", err)
	}

	// 检查订单状态，只有未支付状态的订单才能被支付
	if order.Status != int64(types.OrderStatus_UNPAID) {
		l.Logger.Errorf("订单状态不正确，当前状态: %d", order.Status)
		return nil, fmt.Errorf("订单状态不正确，无法支付")
	}

	// 检查支付锁是否存在
	lockKey := fmt.Sprintf("order:lock:%d", in.Id)
	exists, err := l.svcCtx.Redis.Exists(l.ctx, lockKey).Result()
	if err != nil {
		l.Logger.Errorf("检查支付锁失败: %v", err)
		// 继续处理，不影响支付流程
	} else if exists == 0 {
		l.Logger.Errorf("订单支付锁已过期: 订单ID=%d", in.Id)
		// 继续处理，允许支付
	}

	// 更新订单为已支付状态
	err = l.svcCtx.OrderModel.Update(l.ctx, &model.Order{
		Id:         in.Id,
		Status:     int64(types.OrderStatus_PAID), // 已支付状态
		StatusDesc: "已支付",
		UpdateTime: time.Now(),
	})

	if err != nil {
		l.Logger.Errorf("更新订单状态失败: %v", err)
		return nil, fmt.Errorf("更新订单状态失败: %v", err)
	}

	// 从订单超时队列中移除
	timeoutQueueKey := "order:timeout:queue"
	_, err = l.svcCtx.Redis.ZRem(l.ctx, timeoutQueueKey, fmt.Sprintf("%d", in.Id)).Result()
	if err != nil {
		l.Logger.Errorf("从超时队列移除订单失败: %v", err)
		// 不影响支付流程，继续执行
	}

	// 删除订单支付锁
	_, err = l.svcCtx.Redis.Del(l.ctx, lockKey).Result()
	if err != nil {
		l.Logger.Errorf("删除订单支付锁失败: %v", err)
		// 不影响支付流程，继续执行
	}

	l.Logger.Infof("订单支付成功: 订单ID=%d", in.Id)
	return &types.PaidResponse{}, nil
}
