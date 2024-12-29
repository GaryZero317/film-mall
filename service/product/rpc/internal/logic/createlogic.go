package logic

import (
	"context"
	"time"

	"mall/service/product/model"
	"mall/service/product/rpc/internal/svc"
	"mall/service/product/rpc/pb/product"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type CreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateLogic) Create(in *product.CreateRequest) (*product.CreateResponse, error) {
	now := time.Now()
	newProduct := &model.Product{
		Name:       in.Name,
		Desc:       in.Desc,
		Stock:      in.Stock,
		Amount:     in.Amount,
		Status:     in.Status,
		CreateTime: now,
		UpdateTime: now,
	}

	// 开启事务
	err := l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		// 创建商品
		if err := tx.Create(newProduct).Error; err != nil {
			return err
		}

		// 如果有图片，添加图片记录
		if len(in.ImageUrls) > 0 {
			var images []*model.ProductImage
			for i, url := range in.ImageUrls {
				images = append(images, &model.ProductImage{
					ProductId: newProduct.Id,
					ImageUrl:  url,
					IsMain:    i == 0, // 第一张图片设为主图
					SortOrder: i,
				})
			}

			if err := tx.Create(&images).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, status.Error(500, err.Error())
	}

	return &product.CreateResponse{
		Id: newProduct.Id,
	}, nil
}
