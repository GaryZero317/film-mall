package logic

import (
	"context"

	"mall/service/product/api/internal/svc"
	"mall/service/product/api/internal/types"
	"mall/service/product/model"
	"mall/service/product/rpc/pb/product"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.UpdateRequest) (resp *types.UpdateResponse, err error) {
	// 开启事务
	err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 更新商品基本信息
		_, err = l.svcCtx.ProductRpc.Update(l.ctx, &product.UpdateRequest{
			Id:     req.Id,
			Name:   req.Name,
			Desc:   req.Desc,
			Stock:  req.Stock,
			Amount: req.Amount,
			Status: req.Status,
		})
		if err != nil {
			return err
		}

		// 2. 如果有新的图片列表，更新图片
		if len(req.ImageUrls) > 0 {
			// 2.1 删除旧的图片记录
			productImage := &model.ProductImage{}
			if err := productImage.DeleteByProductId(l.ctx, tx, req.Id); err != nil {
				return err
			}

			// 2.2 添加新的图片记录
			var images []*model.ProductImage
			for i, url := range req.ImageUrls {
				images = append(images, &model.ProductImage{
					ProductId: req.Id,
					ImageUrl:  url,
					IsMain:    url == req.MainImage, // 使用指定的主图
					SortOrder: i,
				})
			}

			// 如果没有指定主图，则使用第一张图片作为主图
			if req.MainImage == "" && len(images) > 0 {
				images[0].IsMain = true
			}

			if err := productImage.BatchInsert(l.ctx, tx, images); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &types.UpdateResponse{}, nil
}
