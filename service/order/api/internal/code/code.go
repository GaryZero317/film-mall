package code

// 定义业务状态码
const (
	Success       = 200
	Error         = 500
	InvalidParams = 400

	// 订单相关错误码 (2000-2999)
	OrderNotExist           = 2001
	OrderCreateFailed       = 2002
	OrderUpdateFailed       = 2003
	OrderDeleteFailed       = 2004
	OrderStatusInvalid      = 2005
	OrderItemCreateFailed   = 2006
	OrderItemNotExist       = 2007
	OrderItemUpdateFailed   = 2008
	OrderItemDeleteFailed   = 2009
	OrderAddressNotExist    = 2010
	OrderUserNotExist       = 2011
	OrderProductNotExist    = 2012
	OrderProductStockLack   = 2013
	OrderAmountInvalid      = 2014
	OrderStatusNotAllowPaid = 2015
)

// 获取错误码对应的消息
var codeMsg = map[int]string{
	Success:       "成功",
	Error:         "服务器内部错误",
	InvalidParams: "参数错误",

	OrderNotExist:           "订单不存在",
	OrderCreateFailed:       "创建订单失败",
	OrderUpdateFailed:       "更新订单失败",
	OrderDeleteFailed:       "删除订单失败",
	OrderStatusInvalid:      "订单状态无效",
	OrderItemCreateFailed:   "创建订单商品失败",
	OrderItemNotExist:       "订单商品不存在",
	OrderItemUpdateFailed:   "更新订单商品失败",
	OrderItemDeleteFailed:   "删除订单商品失败",
	OrderAddressNotExist:    "收货地址不存在",
	OrderUserNotExist:       "用户不存在",
	OrderProductNotExist:    "商品不存在",
	OrderProductStockLack:   "商品库存不足",
	OrderAmountInvalid:      "订单金额无效",
	OrderStatusNotAllowPaid: "订单状态不允许支付",
}

// GetMsg 获取错误码对应的消息
func GetMsg(code int) string {
	msg, ok := codeMsg[code]
	if ok {
		return msg
	}
	return codeMsg[Error]
}
