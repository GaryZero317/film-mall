package logic

import (
	"context"
	"fmt"

	"mall/service/statistics/api/internal/svc"
	"mall/service/statistics/api/internal/types"
	"mall/service/statistics/api/internal/utils"

	"github.com/xuri/excelize/v2"
	"github.com/zeromicro/go-zero/core/logx"
)

type ExportStatisticsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExportStatisticsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExportStatisticsLogic {
	return &ExportStatisticsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExportStatisticsLogic) ExportStatistics(req *types.ExportStatisticsReq) (resp *types.ExportStatisticsResp, err error) {
	// 获取时间范围
	start, end := utils.GetTimeRange(req.TimeRange)

	// 创建Excel文件
	f := excelize.NewFile()

	// 导出热门商品数据
	products, err := l.svcCtx.ProductSalesDailyModel.FindHotProducts(l.ctx, start, end, 10)
	if err != nil {
		return nil, err
	}
	sheet := "热门商品"
	f.NewSheet(sheet)
	f.SetCellValue(sheet, "A1", "商品ID")
	f.SetCellValue(sheet, "B1", "销售量")
	for i, p := range products {
		row := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), p.ProductId)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), p.SalesCount)
	}

	// 导出类别销售数据
	categories, err := l.svcCtx.CategorySalesDailyModel.FindCategorySales(l.ctx, start, end)
	if err != nil {
		return nil, err
	}
	sheet = "类别销售"
	f.NewSheet(sheet)
	f.SetCellValue(sheet, "A1", "类别ID")
	f.SetCellValue(sheet, "B1", "销售量")
	f.SetCellValue(sheet, "C1", "销售金额")
	for i, c := range categories {
		row := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), c.CategoryId)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), c.SalesCount)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), c.SalesAmount)
	}

	// 导出用户行为数据
	behaviors, err := l.svcCtx.UserActivityLogModel.FindUserBehaviors(l.ctx, start, end)
	if err != nil {
		return nil, err
	}
	sheet = "用户行为"
	f.NewSheet(sheet)
	f.SetCellValue(sheet, "A1", "日期")
	f.SetCellValue(sheet, "B1", "行为类型")
	f.SetCellValue(sheet, "C1", "数量")
	row := 2
	for _, b := range behaviors {
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), b.ActivityTime.Format("2006-01-02"))
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), b.ActivityType)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), 1)
		row++
	}

	// 删除默认的Sheet1
	f.DeleteSheet("Sheet1")

	// 保存到内存
	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	// 构造返回数据
	resp = &types.ExportStatisticsResp{
		Code: 0,
		Msg:  "success",
		Data: buffer.Bytes(),
	}

	return resp, nil
}
