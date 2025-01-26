package logic

import (
	"context"

	"mall/service/product/model"
	"mall/service/product/rpc/internal/svc"
	"mall/service/product/rpc/pb/product"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/status"
)

type DetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DetailLogic) Detail(in *product.DetailRequest) (*product.DetailResponse, error) {
	// 查询产品是否存在
	res, err := l.svcCtx.ProductModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, status.Error(100, "产品不存在")
		}
		return nil, status.Error(500, err.Error())
	}

	// 查询商品图片
	var productImages []model.ProductImage
	if err := l.svcCtx.DB.Where("product_id = ?", in.Id).Order("is_main DESC, sort_order ASC").Find(&productImages).Error; err != nil {
		l.Logger.Errorf("查询商品图片失败: %v", err)
		return nil, status.Error(500, err.Error())
	}

	// 提取图片URL列表
	var imageUrls []string
	for _, img := range productImages {
		imageUrls = append(imageUrls, img.ImageUrl)
	}

	l.Logger.Infof("商品 %d 的图片URLs: %v", in.Id, imageUrls)

	return &product.DetailResponse{
		Id:        res.Id,
		Name:      res.Name,
		Desc:      res.Desc,
		Stock:     res.Stock,
		Amount:    res.Amount,
		Status:    res.Status,
		ImageUrls: imageUrls,
	}, nil
}
