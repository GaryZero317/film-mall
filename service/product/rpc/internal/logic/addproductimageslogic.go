package logic

import (
	"context"

	"mall/service/product/model"
	"mall/service/product/rpc/internal/svc"
	"mall/service/product/rpc/pb/product"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
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

	// 开启事务
	err := l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		// 检查是否已有主图
		productImage := &model.ProductImage{}
		mainImage, err := productImage.FindMainImage(l.ctx, tx, in.ProductId)
		if err != nil {
			return err
		}

		// 获取当前最大的排序值
		var maxSortOrder int
		if err := tx.Model(&model.ProductImage{}).
			Where("product_id = ?", in.ProductId).
			Select("COALESCE(MAX(sort_order), -1)").
			Scan(&maxSortOrder).Error; err != nil {
			return err
		}

		// 创建商品图片记录
		images := make([]*model.ProductImage, 0, len(in.ImageUrls))
		for i, url := range in.ImageUrls {
			images = append(images, &model.ProductImage{
				ProductId: in.ProductId,
				ImageUrl:  url,
				IsMain:    mainImage == nil && i == 0, // 如果没有主图，第一张设为主图
				SortOrder: maxSortOrder + i + 1,
			})
		}

		l.Logger.Info("准备批量插入图片记录")
		// 批量插入图片记录
		if err := tx.Create(&images).Error; err != nil {
			l.Logger.Errorf("批量插入图片记录失败: %v", err)
			return err
		}

		return nil
	})

	if err != nil {
		l.Logger.Errorf("添加商品图片失败: %v", err)
		return nil, err
	}

	l.Logger.Info("添加商品图片成功")
	return &product.AddProductImagesResponse{}, nil
}
