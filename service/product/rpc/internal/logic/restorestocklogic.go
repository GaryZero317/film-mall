package logic

import (
	"context"

	"mall/service/product/rpc/internal/svc"
	"mall/service/product/rpc/pb/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type RestoreStockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRestoreStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RestoreStockLogic {
	return &RestoreStockLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RestoreStockLogic) RestoreStock(in *product.RestoreStockRequest) (*product.RestoreStockResponse, error) {
	l.Logger.Infof("恢复商品库存: 商品ID=%d, 数量=%d", in.Id, in.Quantity)

	// 调用模型恢复库存
	err := l.svcCtx.ProductModel.RestoreStock(l.ctx, in.Id, in.Quantity)
	if err != nil {
		l.Logger.Errorf("恢复库存失败: %v", err)
		return &product.RestoreStockResponse{Success: false}, err
	}

	return &product.RestoreStockResponse{Success: true}, nil
}
