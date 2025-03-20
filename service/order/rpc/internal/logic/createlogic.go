package logic

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"mall/service/order/model"
	"mall/service/order/rpc/internal/svc"
	"mall/service/order/rpc/types"
	"mall/service/product/rpc/pb/product"

	"github.com/go-redis/redis/v8"
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

	// 生成订单号 - 添加随机数防止并发冲突
	now := time.Now()
	// 初始化随机数生成器
	rand.Seed(now.UnixNano())
	// 生成6位随机数
	randomNum := rand.Intn(1000000)
	// 新的订单号格式: 年月日时分秒+用户ID+6位随机数
	oid := fmt.Sprintf("%s%d%06d", now.Format("20060102150405"), in.Uid, randomNum)

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
		Status:      int64(types.OrderStatus_UNPAID), // 设置为待支付状态
		StatusDesc:  "待支付",
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
		order.Id = orderId

		l.Logger.Info("开始创建订单商品和减库存")
		// 2. 创建订单商品和减库存
		for i := range orderItems {
			orderItems[i].OrderId = orderId
			_, err := l.svcCtx.OrderItemModel.Insert(ctx, session, &orderItems[i])
			if err != nil {
				l.Logger.Errorf("创建订单商品失败: %v", err)
				return err
			}

			// 调用产品服务减少库存
			_, err = l.svcCtx.ProductRpc.DecrStock(ctx, &product.DecrStockRequest{
				Id:       orderItems[i].Pid,
				Quantity: orderItems[i].Quantity,
			})

			if err != nil {
				l.Logger.Errorf("减少商品库存失败: 商品ID=%d, 数量=%d, 错误=%v",
					orderItems[i].Pid, orderItems[i].Quantity, err)
				return fmt.Errorf("减少商品库存失败: %v", err)
			}
		}

		l.Logger.Info("事务完成")
		return nil
	})

	if err != nil {
		l.Logger.Errorf("事务执行失败: %v", err)
		return nil, err
	}

	// 设置订单支付超时时间，默认15分钟
	expiry := time.Duration(l.svcCtx.OrderLockExpiry) * time.Second
	if expiry == 0 {
		expiry = 15 * time.Minute
		l.Logger.Infof("未配置订单锁过期时间，使用默认值15分钟")
	} else {
		l.Logger.Infof("订单锁过期时间设置为: %v秒", l.svcCtx.OrderLockExpiry)
	}

	// 使用Redis设置订单支付超时锁
	lockKey := fmt.Sprintf("order:lock:%d", order.Id)
	err = l.svcCtx.Redis.Set(l.ctx, lockKey, "1", expiry).Err()
	if err != nil {
		l.Logger.Errorf("设置订单支付超时锁失败: %v", err)
		// 不影响订单创建，继续执行
	} else {
		l.Logger.Infof("成功设置订单支付超时锁，订单ID=%d，过期时间=%v", order.Id, expiry)
	}

	// 添加到订单超时队列，用于后续处理超时订单
	timeoutQueueKey := "order:timeout:queue"
	// 使用ZSet存储，score为过期时间的时间戳
	expireAt := now.Add(expiry).Unix()
	err = l.svcCtx.Redis.ZAdd(l.ctx, timeoutQueueKey, &redis.Z{
		Score:  float64(expireAt),
		Member: strconv.FormatInt(order.Id, 10),
	}).Err()

	if err != nil {
		l.Logger.Errorf("添加订单到超时队列失败: %v", err)
		// 不影响订单创建，继续执行
	}

	l.Logger.Infof("创建订单成功: %+v", order)
	return &types.CreateResponse{
		Id:  order.Id,
		Oid: order.Oid,
	}, nil
}
