package logic

import (
	"context"

	"mall/service/product/api/internal/svc"
	"mall/service/product/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// Product 数据库模型
type Product struct {
	Id        int64    `gorm:"column:id;primaryKey"`
	Name      string   `gorm:"column:name"`
	Desc      string   `gorm:"column:desc"`
	Stock     int64    `gorm:"column:stock"`
	Amount    int64    `gorm:"column:amount"`
	Status    int64    `gorm:"column:status"`
	MainImage string   `gorm:"column:main_image"`
	Images    []string `gorm:"column:images;type:json"`
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
	// 设置默认分页参数
	page := req.Page
	if page == 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize == 0 {
		pageSize = 10
	}

	// 构建查询条件
	query := l.svcCtx.DB.Model(&Product{}).Where("status = ?", 1) // 只查询上架商品

	// 使用 GORM 查询商品列表
	var dbProducts []Product
	result := query.Order("id DESC"). // 按ID降序排序，最新商品排在前面
						Offset(int((page - 1) * pageSize)).
						Limit(int(pageSize)).
						Find(&dbProducts)
	if result.Error != nil {
		return &types.ProductListResponse{
			Code: 500,
			Msg:  "获取商品列表失败",
		}, result.Error
	}

	// 使用相同的查询条件获取总数
	var total int64
	result = query.Count(&total)
	if result.Error != nil {
		return &types.ProductListResponse{
			Code: 500,
			Msg:  "获取商品总数失败",
		}, result.Error
	}

	// 转换为API响应类型
	products := make([]types.Product, len(dbProducts))
	for i, p := range dbProducts {
		products[i] = types.Product{
			Id:        p.Id,
			Name:      p.Name,
			Desc:      p.Desc,
			Stock:     p.Stock,
			Amount:    p.Amount,
			Status:    p.Status,
			Images:    p.Images,
			MainImage: p.MainImage,
		}
	}

	// 组装响应数据
	return &types.ProductListResponse{
		Code: 0,
		Msg:  "success",
		Data: &types.ProductListData{
			Total: total,
			List:  products,
		},
	}, nil
}
