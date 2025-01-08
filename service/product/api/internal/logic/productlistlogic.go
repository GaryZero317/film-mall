package logic

import (
	"context"

	"mall/service/product/api/internal/svc"
	"mall/service/product/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

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

	// 使用 GORM 查询商品列表
	var products []types.Product
	result := l.svcCtx.DB.Offset(int((page - 1) * pageSize)).Limit(int(pageSize)).Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}

	// 使用 GORM 查询商品总数
	var total int64
	result = l.svcCtx.DB.Model(&types.Product{}).Count(&total)
	if result.Error != nil {
		return nil, result.Error
	}

	// 组装响应数据
	data := &types.ProductListData{
		Total: total,
		List:  products,
	}
	resp = &types.ProductListResponse{
		Code: 0,
		Msg:  "success",
		Data: data,
	}
	return resp, nil
}
