// Code generated by goctl. DO NOT EDIT.
package types

type HotProductsData struct {
	Products []string `json:"products"` // 商品名称列表
	Sales    []int    `json:"sales"`    // 销售量列表
}

type HotProductsReq struct {
	TimeRange string `json:"timeRange"` // 时间范围：7days, 30days, month, quarter, year
}

type HotProductsResp struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data HotProductsData `json:"data"`
}

type CategoryStatsData struct {
	Categories []string `json:"categories"` // 类别名称列表
	Sales      []int    `json:"sales"`      // 销售量列表
}

type CategoryStatsReq struct {
	TimeRange string `json:"timeRange"`
}

type CategoryStatsResp struct {
	Code int               `json:"code"`
	Msg  string            `json:"msg"`
	Data CategoryStatsData `json:"data"`
}

type UserBehaviorData struct {
	Dates  []string `json:"dates"`  // 日期列表
	Views  []int    `json:"views"`  // 浏览量列表
	Carts  []int    `json:"carts"`  // 加购量列表
	Orders []int    `json:"orders"` // 下单量列表
}

type UserBehaviorReq struct {
	TimeRange string `json:"timeRange"`
}

type UserBehaviorResp struct {
	Code int              `json:"code"`
	Msg  string           `json:"msg"`
	Data UserBehaviorData `json:"data"`
}

type UserActivityData struct {
	Hours    []string `json:"hours"`    // 小时列表 (0-23)
	Activity [][]int  `json:"activity"` // 活跃度数据 [hour, weekday, value]
}

type UserActivityReq struct {
	TimeRange string `json:"timeRange"`
}

type UserActivityResp struct {
	Code int              `json:"code"`
	Msg  string           `json:"msg"`
	Data UserActivityData `json:"data"`
}

type ExportStatisticsReq struct {
	TimeRange string `json:"timeRange"`
}

type ExportStatisticsResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data []byte `json:"data"` // Excel文件二进制数据
}
