package logic

import (
	"context"
	"mall/service/product/api/internal/svc"
	"mall/service/product/api/internal/types"
	"mall/service/product/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type DetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailLogic) Detail(req *types.DetailRequest) (resp *types.DetailResponse, err error) {
	l.Logger.Infof("获取商品详情, 商品ID: %d", req.Id)

	// 获取商品基本信息
	product, err := l.svcCtx.ProductModel.FindOne(l.ctx, req.Id)
	if err != nil {
		l.Logger.Errorf("获取商品信息失败: %v", err)
		return nil, err
	}

	// 获取商品图片
	productImage := &model.ProductImage{}
	images, err := productImage.FindByProductId(l.ctx, l.svcCtx.DB, req.Id)
	if err != nil {
		l.Logger.Errorf("获取商品图片失败: %v", err)
		return nil, err
	}

	// 处理图片信息
	var imageUrls []string
	var mainImage string
	for _, img := range images {
		imageUrls = append(imageUrls, img.ImageUrl)
		if img.IsMain {
			mainImage = img.ImageUrl
		}
	}

	// 如果没有主图但有其他图片，使用第一张图片作为主图
	if mainImage == "" && len(imageUrls) > 0 {
		mainImage = imageUrls[0]
	}

	l.Logger.Info("获取商品详情成功")
	return &types.DetailResponse{
		Id:         product.Id,
		Name:       product.Name,
		Desc:       product.Desc,
		Stock:      product.Stock,
		Amount:     product.Amount,
		Status:     product.Status,
		ImageUrls:  imageUrls,
		MainImage:  mainImage,
		CategoryId: product.CategoryId,
	}, nil
}
