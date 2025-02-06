package logic

import (
	"context"
	"mall/service/product/api/internal/svc"
	"mall/service/product/api/internal/types"
	"mall/service/product/model"

	"github.com/zeromicro/go-zero/core/logx"
)

// Product 数据库模型
type Product struct {
	Id         int64    `gorm:"column:id;primaryKey"`
	Name       string   `gorm:"column:name"`
	Desc       string   `gorm:"column:desc"`
	Stock      int64    `gorm:"column:stock"`
	Amount     int64    `gorm:"column:amount"`
	Status     int64    `gorm:"column:status"`
	MainImage  string   `gorm:"column:main_image"`
	Images     []string `gorm:"column:images;type:json"`
	CategoryId int64    `gorm:"column:category_id"`
}

// TableName 指定表名
func (Product) TableName() string {
	return "product"
}

type ProductListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProductListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductListLogic {
	return &ProductListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProductListLogic) ProductList(req *types.ProductListRequest) (resp *types.ProductListResponse, err error) {
	l.Logger.Infof("获取商品列表, page: %d, pageSize: %d", req.Page, req.PageSize)

	// 设置默认分页参数
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 10
	}

	// 获取商品列表
	products, total, err := l.svcCtx.ProductModel.FindPageListByPage(l.ctx, req.Page, req.PageSize)
	if err != nil {
		l.Logger.Errorf("获取商品列表失败: %v", err)
		return nil, err
	}

	// 转换数据
	var list []types.Product
	for _, p := range products {
		// 获取商品图片
		var images []*model.ProductImage
		if err := l.svcCtx.DB.WithContext(l.ctx).
			Where("product_id = ?", p.Id).
			Order("is_main DESC, sort_order ASC").
			Find(&images).Error; err != nil {
			l.Logger.Errorf("获取商品[%d]图片失败: %v", p.Id, err)
			continue
		}

		// 提取图片URL列表和主图
		var imageUrls []string
		var mainImage string
		for _, img := range images {
			imageUrls = append(imageUrls, img.ImageUrl)
			if img.IsMain {
				mainImage = img.ImageUrl
			}
		}

		// 如果没有主图，但有其他图片，则使用第一张图片作为主图
		if mainImage == "" && len(imageUrls) > 0 {
			mainImage = imageUrls[0]
		}

		// 添加到列表
		list = append(list, types.Product{
			Id:         p.Id,
			Name:       p.Name,
			Desc:       p.Desc,
			Stock:      p.Stock,
			Amount:     p.Amount,
			Status:     p.Status,
			CategoryId: p.CategoryId,
			Images:     imageUrls,
			MainImage:  mainImage,
		})
	}

	l.Logger.Infof("获取商品列表成功: total=%d, list=%d", total, len(list))
	return &types.ProductListResponse{
		Code: 0,
		Msg:  "success",
		Data: &types.ProductListData{
			Total: total,
			List:  list,
		},
	}, nil
}
