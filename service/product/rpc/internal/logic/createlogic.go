package logic

import (
	"context"

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
	newProduct := model.Product{
		Name:   in.Name,
		Desc:   in.Desc,
		Stock:  in.Stock,
		Amount: in.Amount,
		Status: in.Status,
	}

	// 开启事务
	err := l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		// 创建商品
		result, err := l.svcCtx.ProductModel.Insert(l.ctx, &newProduct)
		if err != nil {
			return err
		}

		productId, err := result.LastInsertId()
		if err != nil {
			return err
		}

		// 如果有图片，添加图片记录
		if len(in.ImageUrls) > 0 {
			var images []*model.ProductImage
			for i, url := range in.ImageUrls {
				images = append(images, &model.ProductImage{
					ProductId: productId,
					ImageUrl:  url,
					IsMain:    i == 0, // 第一张图片设为主图
					SortOrder: i,
				})
			}

			productImage := &model.ProductImage{}
			if err := productImage.BatchInsert(l.ctx, tx, images); err != nil {
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
