package logic

import (
	"context"

	"mall/service/product/model"
	"mall/service/product/rpc/internal/svc"
	"mall/service/product/rpc/pb/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetMainImageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetMainImageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetMainImageLogic {
	return &SetMainImageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SetMainImageLogic) SetMainImage(in *product.SetMainImageRequest) (*product.SetMainImageResponse, error) {
	// 设置主图
	productImage := &model.ProductImage{}
	if err := productImage.SetMainImage(l.ctx, l.svcCtx.DB, in.ProductId, in.ImageUrl); err != nil {
		return nil, err
	}

	return &product.SetMainImageResponse{}, nil
}
