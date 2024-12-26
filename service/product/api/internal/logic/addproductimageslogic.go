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
	_, err = l.svcCtx.ProductRpc.AddProductImages(l.ctx, &product.AddProductImagesRequest{
		ProductId: req.ProductId,
		ImageUrls: req.ImageUrls,
	})
	if err != nil {
		return nil, err
	}

	return &types.AddProductImagesResponse{}, nil
}
