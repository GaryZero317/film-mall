package logic

import (
	"context"
	"fmt"
	"mall/service/statistics/api/internal/svc"
	"mall/service/statistics/api/internal/types"

	"mall/service/statistics/api/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetHotProductsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetHotProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetHotProductsLogic {
	return &GetHotProductsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetHotProductsLogic) GetHotProducts(req *types.HotProductsReq) (resp *types.HotProductsResp, err error) {
	// 获取时间范围
	start, end := utils.GetTimeRange(req.TimeRange)

	// 构造返回数据
	resp = &types.HotProductsResp{
		Code: 0,
		Msg:  "success",
		Data: types.HotProductsData{
			Products: make([]string, 0),
			Sales:    make([]int, 0),
		},
	}

	// 查询热门商品
	products, err := l.svcCtx.ProductSalesDailyModel.FindHotProducts(l.ctx, start, end, 10)
	if err != nil {
		l.Logger.Errorf("查询热门商品失败: %v", err)
		// 即使查询失败，也返回空数据而不是错误
		return resp, nil
	}

	// 填充数据
	for _, p := range products {
		resp.Data.Products = append(resp.Data.Products, fmt.Sprintf("%d", p.ProductId))
		resp.Data.Sales = append(resp.Data.Sales, int(p.SalesCount))
	}

	// 如果没有数据，添加一些默认数据供前端展示
	if len(resp.Data.Products) == 0 {
		resp.Data.Products = []string{"暂无数据"}
		resp.Data.Sales = []int{0}
	}

	return resp, nil
}
