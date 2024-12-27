package logic

import (
	"context"

	"mall/service/product/model"
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
	l.Logger.Infof("开始添加商品图片: productId=%d, imageUrls=%v", in.ProductId, in.ImageUrls)

	// 创建商品图片记录
	images := make([]*model.ProductImage, 0, len(in.ImageUrls))
	for _, url := range in.ImageUrls {
		images = append(images, &model.ProductImage{
			ProductId: in.ProductId,
			ImageUrl:  url,
			IsMain:    false,
		})
	}

	l.Logger.Info("准备批量插入图片记录")
	// 批量插入图片记录
	productImage := &model.ProductImage{}
	if err := productImage.BatchInsert(l.ctx, l.svcCtx.DB, images); err != nil {
		l.Logger.Errorf("批量插入图片记录失败: %v", err)
		return nil, err
	}

	l.Logger.Info("添加商品图片成功")
	return &product.AddProductImagesResponse{}, nil
}
