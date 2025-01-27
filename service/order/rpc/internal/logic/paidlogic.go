package logic

import (
	"context"
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
	// 更新订单为已支付状态
	err := l.svcCtx.OrderModel.Update(l.ctx, &model.Order{
		Id:         in.Id,
		Status:     1, // 已支付状态
		StatusDesc: "已支付",
		UpdateTime: time.Now(),
	})

	if err != nil {
		return nil, err
	}

	return &types.PaidResponse{}, nil
}
