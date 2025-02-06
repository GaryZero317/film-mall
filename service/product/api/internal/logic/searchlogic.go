package logic

import (
	"context"

	"mall/service/product/api/internal/svc"
	"mall/service/product/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchLogic {
	return &SearchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchLogic) Search(req *types.SearchRequest) (resp *types.SearchResponse, err error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 10
	}

	products, total, err := l.svcCtx.ProductModel.Search(l.ctx, req.Keyword, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	var productList []types.Product
	for _, p := range products {
		// 获取商品主图
		images, err := l.svcCtx.ProductImage.FindByProductId(l.ctx, l.svcCtx.DB, p.Id)
		if err != nil {
			return nil, err
		}

		mainImage := ""
		for _, img := range images {
			if img.IsMain {
				mainImage = img.ImageUrl
				break
			}
		}

		productList = append(productList, types.Product{
			Id:         p.Id,
			Name:       p.Name,
			Desc:       p.Desc,
			Stock:      p.Stock,
			Amount:     p.Amount,
			Status:     p.Status,
			CategoryId: p.CategoryId,
			MainImage:  mainImage,
		})
	}

	return &types.SearchResponse{
		Code: 0,
		Msg:  "success",
		Data: &types.ProductListData{
			Total: total,
			List:  productList,
		},
	}, nil
}
