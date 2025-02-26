package model

// 胶片冲洗订单状态
const (
	FilmOrderStatusWaitPay     = 0 // 待付款
	FilmOrderStatusProcessing  = 1 // 冲洗处理中
	FilmOrderStatusWaitReceive = 2 // 待收货
	FilmOrderStatusFinished    = 3 // 已完成
)

// 胶片冲洗订单状态描述
var FilmOrderStatusText = map[int64]string{
	FilmOrderStatusWaitPay:     "待付款",
	FilmOrderStatusProcessing:  "冲洗处理中",
	FilmOrderStatusWaitReceive: "待收货",
	FilmOrderStatusFinished:    "已完成",
}

// 获取胶片冲洗订单状态文本
func GetFilmOrderStatusText(status int64) string {
	if text, ok := FilmOrderStatusText[status]; ok {
		return text
	}
	return "未知状态"
}
