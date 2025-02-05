package logic

import (
	"context"

	"mall/service/product/rpc/internal/svc"
	"mall/service/product/rpc/pb/product"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/status"
)

type RemoveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRemoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveLogic {
	return &RemoveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RemoveLogic) Remove(in *product.RemoveRequest) (*product.RemoveResponse, error) {
	l.Logger.Infof("接收到删除商品请求，ID: %d", in.Id)

	// 直接调用删除方法
	err := l.svcCtx.ProductModel.Delete(l.ctx, in.Id)
	if err != nil {
		l.Logger.Errorf("删除商品失败: %v", err)
		return nil, status.Error(500, err.Error())
	}

	l.Logger.Infof("商品删除成功，ID: %d", in.Id)
	return &product.RemoveResponse{}, nil
}
