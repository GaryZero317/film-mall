package logic

import (
	"context"

	"mall/service/product/api/internal/svc"
	"mall/service/product/api/internal/types"
	"mall/service/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddProductImagesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddProductImagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddProductImagesLogic {
	return &AddProductImagesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddProductImagesLogic) AddProductImages(req *types.AddProductImagesRequest) (resp *types.AddProductImagesResponse, err error) {
	// 调用RPC服务添加商品图片
	l.Logger.Infof("添加商品图片请求: productId=%d, imageUrls=%v", req.ProductId, req.ImageUrls)

	_, err = l.svcCtx.ProductRpc.AddProductImages(l.ctx, &product.AddProductImagesRequest{
		ProductId: req.ProductId,
		ImageUrls: req.ImageUrls,
	})
	if err != nil {
		l.Logger.Errorf("添加商品图片失败: %v", err)
		return nil, err
	}

	l.Logger.Info("添加商品图片成功")
	return &types.AddProductImagesResponse{}, nil
}
