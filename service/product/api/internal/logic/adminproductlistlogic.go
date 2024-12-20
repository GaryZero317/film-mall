package logic

import (
	"context"

	"mall/service/product/api/internal/svc"
	"mall/service/product/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminProductListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminProductListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminProductListLogic {
	return &AdminProductListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminProductListLogic) AdminProductList(req *types.AdminProductListRequest) (resp *types.AdminProductListResponse, err error) {
	// 获取分页参数
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 10
	}

	// 调用model层获取数据
	products, total, err := l.svcCtx.ProductModel.FindPageListByPage(l.ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	// 转换数据
	var list []types.Product
	for _, p := range products {
		list = append(list, types.Product{
			Id:     p.Id,
			Name:   p.Name,
			Desc:   p.Desc,
			Stock:  p.Stock,
			Amount: p.Amount,
			Status: p.Status,
		})
	}

	return &types.AdminProductListResponse{
		Total: total,
		List:  list,
	}, nil
}
