package model

import "time"

// Order represents the business model
type Order struct {
	Id          int64     `json:"id"`           // 订单ID
	Oid         string    `json:"oid"`          // 订单号
	Uid         int64     `json:"uid"`          // 用户ID
	AddressId   int64     `json:"address_id"`   // 收货地址ID
	TotalPrice  int64     `json:"total_price"`  // 订单总价(分)
	ShippingFee int64     `json:"shipping_fee"` // 运费(分)
	Status      int64     `json:"status"`       // 订单状态
	StatusDesc  string    `json:"status_desc"`  // 状态描述
	Remark      string    `json:"remark"`       // 订单备注
	CreateTime  time.Time `json:"create_time"`  // 创建时间
	UpdateTime  time.Time `json:"update_time"`  // 更新时间
}

// OrderItem represents the order item model
type OrderItem struct {
	Id           int64  `json:"id"`            // 订单商品ID
	OrderId      int64  `json:"order_id"`      // 订单ID
	Pid          int64  `json:"pid"`           // 商品ID
	ProductName  string `json:"product_name"`  // 商品名称
	ProductImage string `json:"product_image"` // 商品图片
	Price        int64  `json:"price"`         // 商品单价(分)
	Quantity     int64  `json:"quantity"`      // 购买数量
	Amount       int64  `json:"amount"`        // 商品总价(分)
}
