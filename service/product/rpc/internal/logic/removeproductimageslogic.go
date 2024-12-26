package logic

import (
	"context"

	"mall/service/product/rpc/internal/svc"
	"mall/service/product/rpc/pb/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveProductImagesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRemoveProductImagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveProductImagesLogic {
	return &RemoveProductImagesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RemoveProductImagesLogic) RemoveProductImages(in *product.RemoveProductImagesRequest) (*product.RemoveProductImagesResponse, error) {
	// todo: add your logic here and delete this line

	return &product.RemoveProductImagesResponse{}, nil
}
