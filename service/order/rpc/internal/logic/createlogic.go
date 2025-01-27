package logic

import (
	"context"
	"fmt"
	"math/rand"
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

// 生成订单号
func generateOrderID(uid int64) string {
	// 获取当前时间
	now := time.Now()
	// 生成4位随机数
	random := rand.Intn(10000)
	// 格式化订单号：时间戳(14位) + 随机数(4位) + 用户ID
	return fmt.Sprintf("%s%04d%d",
		now.Format("20060102150405"),
		random,
		uid)
}

func (l *CreateLogic) Create(in *types.CreateRequest) (*types.CreateResponse, error) {
	// 生成订单号
	oid := generateOrderID(in.Uid)

	// 创建订单记录
	order := &model.Order{
		Uid:    in.Uid,
		Pid:    in.Pid,
		Amount: in.Amount,
		Status: in.Status,
		Oid:    oid,
	}

	// 插入数据库
	result, err := l.svcCtx.OrderModel.Insert(l.ctx, order)
	if err != nil {
		return nil, err
	}

	// 获取自增ID
	newOrder, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &types.CreateResponse{
		Id:  newOrder,
		Oid: oid,
	}, nil
}
