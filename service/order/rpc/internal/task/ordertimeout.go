package task

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"mall/service/order/rpc/internal/logic"
	"mall/service/order/rpc/internal/svc"
	"mall/service/order/rpc/types"

	"github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
)

// 订单超时处理
type OrderTimeoutTask struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

// 创建订单超时处理器
func NewOrderTimeoutTask(ctx context.Context, svcCtx *svc.ServiceContext) *OrderTimeoutTask {
	return &OrderTimeoutTask{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 开始订单超时检查
func (t *OrderTimeoutTask) Start() {
	t.Logger.Info("订单超时处理任务启动")
	threading.GoSafe(func() {
		t.run()
	})
}

// 运行超时订单处理
func (t *OrderTimeoutTask) run() {
	for {
		// 处理超时订单
		t.processTimeoutOrders()
		// 每30秒检查一次超时订单，降低资源占用
		time.Sleep(30 * time.Second)
	}
}

// 处理超时订单
func (t *OrderTimeoutTask) processTimeoutOrders() {
	t.Logger.Info("开始检查超时订单")

	// 超时队列的key
	timeoutQueueKey := "order:timeout:queue"

	// 获取当前时间
	now := time.Now().Unix()

	// 获取所有已超时的订单（score <= 当前时间）
	orders, err := t.svcCtx.Redis.ZRangeByScore(t.ctx, timeoutQueueKey, &redis.ZRangeBy{
		Min:    "0",
		Max:    fmt.Sprintf("%d", now),
		Offset: 0,
		Count:  100, // 每次最多处理100个
	}).Result()

	if err != nil {
		t.Logger.Errorf("查询超时订单失败: %v", err)
		return
	}

	if len(orders) == 0 {
		t.Logger.Info("没有超时订单需要处理")
		return
	}

	t.Logger.Infof("发现 %d 个超时订单需要取消", len(orders))

	// 创建订单取消逻辑
	cancelLogic := logic.NewCancelLogic(t.ctx, t.svcCtx)

	// 处理每个超时订单
	for _, orderIdStr := range orders {
		// 转换为订单ID
		orderId, err := strconv.ParseInt(orderIdStr, 10, 64)
		if err != nil {
			t.Logger.Errorf("解析订单ID失败: %v", err)
			continue
		}

		// 查询订单状态
		order, err := t.svcCtx.OrderModel.FindOne(t.ctx, orderId)
		if err != nil {
			t.Logger.Errorf("查询订单失败: %v", err)
			continue
		}

		// 只处理未支付的订单
		if order.Status == int64(types.OrderStatus_UNPAID) {
			// 调用取消订单逻辑
			_, err = cancelLogic.Cancel(&types.CancelRequest{
				Id: orderId,
			})

			if err != nil {
				t.Logger.Errorf("取消超时订单失败: 订单ID=%d, 错误=%v", orderId, err)
			} else {
				t.Logger.Infof("成功取消超时订单: 订单ID=%d", orderId)
			}
		} else {
			// 已支付或已取消的订单，从超时队列中移除
			_, err = t.svcCtx.Redis.ZRem(t.ctx, timeoutQueueKey, orderIdStr).Result()
			if err != nil {
				t.Logger.Errorf("从超时队列移除订单失败: %v", err)
			}
		}
	}
}
