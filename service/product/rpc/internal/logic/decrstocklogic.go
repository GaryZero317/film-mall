package logic

import (
	"context"

	"mall/service/product/rpc/internal/svc"
	"mall/service/product/rpc/pb/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type DecrStockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDecrStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DecrStockLogic {
	return &DecrStockLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DecrStockLogic) DecrStock(in *product.DecrStockRequest) (*product.DecrStockResponse, error) {
	l.Logger.Infof("减少商品库存: 商品ID=%d, 数量=%d", in.Id, in.Quantity)

	// 调用模型减少库存
	err := l.svcCtx.ProductModel.DecrStock(l.ctx, in.Id, in.Quantity)
	if err != nil {
		l.Logger.Errorf("减少库存失败: %v", err)
		return &product.DecrStockResponse{Success: false}, err
	}

	return &product.DecrStockResponse{Success: true}, nil
}
