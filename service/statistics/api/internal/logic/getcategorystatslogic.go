package logic

import (
	"context"
	"fmt"
	"mall/service/statistics/api/internal/svc"
	"mall/service/statistics/api/internal/types"
	"mall/service/statistics/api/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCategoryStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCategoryStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCategoryStatsLogic {
	return &GetCategoryStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCategoryStatsLogic) GetCategoryStats(req *types.CategoryStatsReq) (resp *types.CategoryStatsResp, err error) {
	// 获取时间范围
	start, end := utils.GetTimeRange(req.TimeRange)

	// 构造返回数据
	resp = &types.CategoryStatsResp{
		Code: 0,
		Msg:  "success",
		Data: types.CategoryStatsData{
			Categories: make([]string, 0),
			Sales:      make([]int, 0),
		},
	}

	// 查询类别销售数据
	categories, err := l.svcCtx.CategorySalesDailyModel.FindCategorySales(l.ctx, start, end)
	if err != nil {
		l.Logger.Errorf("查询类别销售数据失败: %v", err)
		// 即使查询失败，也返回空数据而不是错误
		return resp, nil
	}

	// 填充数据
	for _, c := range categories {
		resp.Data.Categories = append(resp.Data.Categories, fmt.Sprintf("%d", c.CategoryId))
		resp.Data.Sales = append(resp.Data.Sales, int(c.SalesCount))
	}

	// 如果没有数据，添加一些默认数据供前端展示
	if len(resp.Data.Categories) == 0 {
		resp.Data.Categories = []string{"暂无数据"}
		resp.Data.Sales = []int{0}
	}

	return resp, nil
}
