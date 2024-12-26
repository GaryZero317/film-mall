package logic

import (
	"context"

	"mall/service/product/api/internal/svc"
	"mall/service/product/api/internal/types"
	"mall/service/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetMainImageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetMainImageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetMainImageLogic {
	return &SetMainImageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetMainImageLogic) SetMainImage(req *types.SetMainImageRequest) (resp *types.SetMainImageResponse, err error) {
	_, err = l.svcCtx.ProductRpc.SetMainImage(l.ctx, &product.SetMainImageRequest{
		ProductId: req.ProductId,
		ImageUrl:  req.ImageUrl,
	})
	if err != nil {
		return nil, err
	}

	return &types.SetMainImageResponse{}, nil
}
