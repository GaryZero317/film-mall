package logic

import (
	"context"

	"mall/service/product/api/internal/svc"
	"mall/service/product/api/internal/types"
	"mall/service/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveProductImagesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRemoveProductImagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveProductImagesLogic {
	return &RemoveProductImagesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RemoveProductImagesLogic) RemoveProductImages(req *types.RemoveProductImagesRequest) (resp *types.RemoveProductImagesResponse, err error) {
	_, err = l.svcCtx.ProductRpc.RemoveProductImages(l.ctx, &product.RemoveProductImagesRequest{
		ProductId: req.ProductId,
		ImageUrls: req.ImageUrls,
	})
	if err != nil {
		return nil, err
	}

	return &types.RemoveProductImagesResponse{}, nil
}
