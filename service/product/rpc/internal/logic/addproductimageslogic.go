package logic

import (
	"context"

	"mall/service/product/rpc/internal/svc"
	"mall/service/product/rpc/pb/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddProductImagesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddProductImagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddProductImagesLogic {
	return &AddProductImagesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddProductImagesLogic) AddProductImages(in *product.AddProductImagesRequest) (*product.AddProductImagesResponse, error) {
	// todo: add your logic here and delete this line

	return &product.AddProductImagesResponse{}, nil
}
