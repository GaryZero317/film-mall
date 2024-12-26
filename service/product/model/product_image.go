package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type ProductImage struct {
	Id         int64     `gorm:"primarykey"`
	ProductId  int64     `gorm:"column:product_id;not null"`
	ImageUrl   string    `gorm:"column:image_url;not null"`
	IsMain     bool      `gorm:"column:is_main;not null;default:0"`
	SortOrder  int       `gorm:"column:sort_order;not null;default:0"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime"`
	UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime"`
}

func (m *ProductImage) TableName() string {
	return "product_image"
}

// 添加商品图片
func (m *ProductImage) Insert(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).Create(m).Error
}

// 批量添加商品图片
func (m *ProductImage) BatchInsert(ctx context.Context, db *gorm.DB, images []*ProductImage) error {
	return db.WithContext(ctx).Create(images).Error
}

// 删除商品图片
func (m *ProductImage) Delete(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).Delete(m).Error
}

// 批量删除商品图片
func (m *ProductImage) BatchDelete(ctx context.Context, db *gorm.DB, productId int64, imageUrls []string) error {
	return db.WithContext(ctx).Where("product_id = ? AND image_url IN ?", productId, imageUrls).Delete(&ProductImage{}).Error
}

// 设置商品主图
func (m *ProductImage) SetMainImage(ctx context.Context, db *gorm.DB, productId int64, imageUrl string) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 先将所有图片设置为非主图
		if err := tx.Model(&ProductImage{}).Where("product_id = ?", productId).Update("is_main", false).Error; err != nil {
			return err
		}
		// 设置新的主图
		return tx.Model(&ProductImage{}).Where("product_id = ? AND image_url = ?", productId, imageUrl).Update("is_main", true).Error
	})
}

// 获取商品所有图片
func (m *ProductImage) FindByProductId(ctx context.Context, db *gorm.DB, productId int64) ([]*ProductImage, error) {
	var images []*ProductImage
	err := db.WithContext(ctx).Where("product_id = ?", productId).Order("is_main DESC, sort_order ASC").Find(&images).Error
	return images, err
}

// 获取商品主图
func (m *ProductImage) FindMainImage(ctx context.Context, db *gorm.DB, productId int64) (*ProductImage, error) {
	var image ProductImage
	err := db.WithContext(ctx).Where("product_id = ? AND is_main = ?", productId, true).First(&image).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &image, err
}

// 删除商品所有图片
func (m *ProductImage) DeleteByProductId(ctx context.Context, db *gorm.DB, productId int64) error {
	return db.WithContext(ctx).Where("product_id = ?", productId).Delete(&ProductImage{}).Error
}
