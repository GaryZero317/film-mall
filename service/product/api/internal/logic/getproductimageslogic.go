package logic

import (
	"context"
	"mall/service/product/api/internal/svc"
	"mall/service/product/api/internal/types"
	"mall/service/product/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProductImagesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProductImagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductImagesLogic {
	return &GetProductImagesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProductImagesLogic) GetProductImages(req *types.GetProductImagesRequest) (resp *types.GetProductImagesResponse, err error) {
	l.Logger.Infof("获取商品图片, 商品ID: %d", req.ProductId)

	// 创建ProductImage实例
	productImage := &model.ProductImage{}

	// 获取商品图片列表
	images, err := productImage.FindByProductId(l.ctx, l.svcCtx.DB, req.ProductId)
	if err != nil {
		l.Logger.Errorf("获取商品图片失败: %v", err)
		return &types.GetProductImagesResponse{
			Code: 1,
			Msg:  "获取商品图片失败",
		}, nil
	}

	// 转换为响应格式
	var imageList []types.ProductImageInfo
	for _, img := range images {
		imageList = append(imageList, types.ProductImageInfo{
			Id:        img.Id,
			ProductId: img.ProductId,
			ImageUrl:  img.ImageUrl,
			IsMain:    img.IsMain,
		})
	}

	l.Logger.Infof("获取商品图片成功, 图片数量: %d", len(imageList))
	return &types.GetProductImagesResponse{
		Code: 0,
		Msg:  "success",
		Data: imageList,
	}, nil
}
