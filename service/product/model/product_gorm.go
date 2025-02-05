package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/syncx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	defaultExpiry  = 24 * time.Hour       // 默认缓存过期时间为24小时
	cacheKeyPrefix = "mall:product:gorm:" // 缓存键前缀
)

// TableName specifies the table name for GORM
func (Product) TableName() string {
	return "product"
}

type defaultGormProductModel struct {
	db    *gorm.DB
	cache cache.Cache
}

func NewDefaultGormProductModel(sqlDB *sql.DB, c cache.CacheConf) (*defaultGormProductModel, error) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &defaultGormProductModel{
		db:    db,
		cache: cache.New(c, syncx.NewSingleFlight(), cache.NewStat("product"), nil),
	}, nil
}

func (m *defaultGormProductModel) Insert(ctx context.Context, data *Product) (sql.Result, error) {
	now := time.Now()
	data.CreateTime = now
	data.UpdateTime = now

	// 开启事务
	err := m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 创建商品
		if err := tx.Create(data).Error; err != nil {
			return err
		}

		// 如果有图片，保存图片信息
		if len(data.Images) > 0 {
			var images []*ProductImage
			for i, url := range data.Images {
				images = append(images, &ProductImage{
					ProductId: data.Id,
					ImageUrl:  url,
					IsMain:    i == 0 || url == data.MainImage, // 第一张图片或指定的主图为主图
					SortOrder: i,
				})
			}
			if err := tx.Create(&images).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &lastInsertIDResult{id: data.Id}, nil
}

func (m *defaultGormProductModel) FindOne(ctx context.Context, id int64) (*Product, error) {
	// 先尝试从缓存获取
	productIdKey := fmt.Sprintf("%s%v", cacheKeyPrefix, id)
	fmt.Printf("正在尝试从缓存获取数据，key: %s\n", productIdKey)

	var productMap map[string]interface{}
	err := m.cache.GetCtx(ctx, productIdKey, &productMap)
	if err != nil {
		fmt.Printf("从缓存获取数据失败，错误: %v\n", err)
	} else {
		fmt.Printf("从缓存获取的数据: %+v\n", productMap)
	}

	var product Product
	if err == nil && productMap != nil {
		// 将map转换回结构体
		if id, ok := productMap["id"].(float64); ok && id > 0 {
			fmt.Printf("缓存命中，正在转换数据\n")
			product = Product{
				Id:         int64(id),
				Name:       productMap["name"].(string),
				Desc:       productMap["desc"].(string),
				Stock:      int64(productMap["stock"].(float64)),
				Amount:     int64(productMap["amount"].(float64)),
				Status:     int64(productMap["status"].(float64)),
				CategoryId: int64(productMap["category_id"].(float64)),
			}
		}
	} else {
		fmt.Printf("缓存未命中或无效，从数据库查询\n")
		// 缓存未命中或无效，从数据库查询
		err = m.db.WithContext(ctx).First(&product, id).Error
		if err != nil {
			fmt.Printf("数据库查询失败，错误: %v\n", err)
			if err == gorm.ErrRecordNotFound {
				return nil, sqlx.ErrNotFound
			}
			return nil, err
		}
		fmt.Printf("从数据库查询到的数据: %+v\n", product)

		// 将基本数据存入缓存
		productMap = map[string]interface{}{
			"id":          product.Id,
			"name":        product.Name,
			"desc":        product.Desc,
			"stock":       product.Stock,
			"amount":      product.Amount,
			"status":      product.Status,
			"category_id": product.CategoryId,
		}
		fmt.Printf("准备写入缓存的数据: %+v\n", productMap)

		// 写入缓存，设置过期时间
		err = m.cache.SetWithExpireCtx(ctx, productIdKey, productMap, defaultExpiry)
		if err != nil {
			fmt.Printf("写入缓存失败，错误: %v\n", err)
		} else {
			fmt.Printf("成功写入缓存\n")
		}
	}

	// 加载商品图片（不缓存图片数据）
	var images []*ProductImage
	if err := m.db.WithContext(ctx).
		Where("product_id = ?", product.Id).
		Order("is_main DESC, sort_order ASC").
		Find(&images).Error; err != nil {
		return nil, err
	}

	// 提取图片URL列表和主图
	var imageUrls []string
	for _, img := range images {
		imageUrls = append(imageUrls, img.ImageUrl)
		if img.IsMain {
			product.MainImage = img.ImageUrl
		}
	}
	product.Images = imageUrls

	// 如果没有主图但有其他图片，使用第一张图片作为主图
	if product.MainImage == "" && len(imageUrls) > 0 {
		product.MainImage = imageUrls[0]
	}

	return &product, nil
}

func (m *defaultGormProductModel) Update(ctx context.Context, data *Product) error {
	productIdKey := fmt.Sprintf("%s%v", cacheKeyPrefix, data.Id)
	fmt.Printf("正在更新商品，ID: %d\n", data.Id)

	err := m.db.WithContext(ctx).Save(data).Error
	if err != nil {
		fmt.Printf("更新数据库失败，错误: %v\n", err)
		return err
	}
	fmt.Printf("数据库更新成功\n")

	// 删除缓存
	err = m.cache.DelCtx(ctx, productIdKey)
	if err != nil {
		fmt.Printf("删除缓存失败，错误: %v\n", err)
	} else {
		fmt.Printf("成功删除缓存\n")
	}
	return nil
}

func (m *defaultGormProductModel) Delete(ctx context.Context, id int64) error {
	productIdKey := fmt.Sprintf("%s%v", cacheKeyPrefix, id)
	fmt.Printf("正在删除商品，ID: %d\n", id)

	// 删除商品图片
	if err := m.db.WithContext(ctx).Table("product_image").Where("product_id = ?", id).Delete(&ProductImage{}).Error; err != nil {
		fmt.Printf("删除商品图片失败，错误: %v\n", err)
		return err
	}

	// 删除商品
	if err := m.db.WithContext(ctx).Table("product").Where("id = ?", id).Delete(&Product{}).Error; err != nil {
		fmt.Printf("删除商品失败，错误: %v\n", err)
		return err
	}
	fmt.Printf("数据库记录删除成功\n")

	// 删除缓存
	err := m.cache.DelCtx(ctx, productIdKey)
	if err != nil {
		fmt.Printf("删除缓存失败，错误: %v\n", err)
	} else {
		fmt.Printf("成功删除缓存\n")
	}
	return nil
}

func (m *defaultGormProductModel) FindPageListByPage(ctx context.Context, page, pageSize int64) ([]*Product, int64, error) {
	var total int64
	var products []*Product

	if err := m.db.WithContext(ctx).Model(&Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := m.db.WithContext(ctx).
		Model(&Product{}).
		Order("id DESC").
		Offset(int(offset)).
		Limit(int(pageSize)).
		Find(&products).Error; err != nil {
		return nil, 0, err
	}

	// 加载每个商品的图片
	for _, p := range products {
		var images []*ProductImage
		if err := m.db.WithContext(ctx).
			Where("product_id = ?", p.Id).
			Order("is_main DESC, sort_order ASC").
			Find(&images).Error; err != nil {
			return nil, 0, err
		}

		// 提取图片URL列表和主图
		var imageUrls []string
		for _, img := range images {
			imageUrls = append(imageUrls, img.ImageUrl)
			if img.IsMain {
				p.MainImage = img.ImageUrl
			}
		}
		p.Images = imageUrls

		// 如果没有主图但有其他图片，使用第一张图片作为主图
		if p.MainImage == "" && len(imageUrls) > 0 {
			p.MainImage = imageUrls[0]
		}
	}

	return products, total, nil
}

func (m *defaultGormProductModel) DecrStock(ctx context.Context, id int64) error {
	result := m.db.WithContext(ctx).Model(&Product{}).
		Where("id = ? AND stock > 0", id).
		UpdateColumn("stock", gorm.Expr("stock - ?", 1))

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("product %d stock not available", id)
	}

	return nil
}

func (m *defaultGormProductModel) Search(ctx context.Context, keyword string, page, pageSize int64) ([]*Product, int64, error) {
	var total int64
	var products []*Product

	query := m.db.WithContext(ctx).Model(&Product{}).
		Where("name LIKE ? OR `desc` LIKE ?", "%"+keyword+"%", "%"+keyword+"%")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(int(offset)).Limit(int(pageSize)).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	// 加载每个商品的图片
	for _, p := range products {
		var images []*ProductImage
		if err := m.db.WithContext(ctx).
			Where("product_id = ?", p.Id).
			Order("is_main DESC, sort_order ASC").
			Find(&images).Error; err != nil {
			return nil, 0, err
		}

		// 提取图片URL列表和主图
		var imageUrls []string
		for _, img := range images {
			imageUrls = append(imageUrls, img.ImageUrl)
			if img.IsMain {
				p.MainImage = img.ImageUrl
			}
		}
		p.Images = imageUrls

		// 如果没有主图但有其他图片，使用第一张图片作为主图
		if p.MainImage == "" && len(imageUrls) > 0 {
			p.MainImage = imageUrls[0]
		}
	}

	return products, total, nil
}

type lastInsertIDResult struct {
	id int64
}

func (r *lastInsertIDResult) LastInsertId() (int64, error) {
	return r.id, nil
}

func (r *lastInsertIDResult) RowsAffected() (int64, error) {
	return 1, nil
}
