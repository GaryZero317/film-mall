package logic

import (
	"context"
	"time"

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

	// 使用异步库存处理器处理请求
	resultChan := l.svcCtx.StockProcessor.AsyncRestoreStock(in.Id, in.Quantity)

	// 设置超时，避免长时间等待
	select {
	case err := <-resultChan:
		if err != nil {
			l.Logger.Errorf("恢复库存失败: %v", err)
			return &product.RestoreStockResponse{Success: false}, err
		}
		return &product.RestoreStockResponse{Success: true}, nil
	case <-time.After(2 * time.Second): // 设置2秒超时
		l.Logger.Infof("恢复库存请求超时，但会在后台继续处理")
		// 超时但不返回错误，库存操作会在后台继续执行
		return &product.RestoreStockResponse{Success: true}, nil
	}
}
