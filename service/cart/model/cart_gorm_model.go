package model

import "time"

type Cart struct {
	Id        int64     `gorm:"primaryKey;column:id" json:"id"`
	UserId    int64     `gorm:"column:user_id" json:"userId"`
	ProductId int64     `gorm:"column:product_id" json:"productId"`
	Quantity  int       `gorm:"column:quantity" json:"quantity"`
	Selected  bool      `gorm:"column:selected" json:"selected"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (m *Cart) TableName() string {
	return "cart"
}
