package logic

import (
	"context"
	"encoding/json"
	"mall/service/cart/api/internal/svc"
	"mall/service/cart/api/internal/types"
	"mall/service/cart/model"
	"mall/service/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCartLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCartLogic {
	return &GetCartLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCartLogic) GetCart() (resp *types.CartListResp, err error) {
	// 打印完整的context信息
	l.Logger.Infof("DEBUG GetCart - Context values: %+v", l.ctx)

	// 从context中获取uid，JWT中的数字类型为json.Number
	uidValue := l.ctx.Value("uid")
	l.Logger.Infof("DEBUG GetCart - Raw uid value type: %T", uidValue)

	if uidValue == nil {
		l.Logger.Errorf("DEBUG GetCart - uid not found in context")
		return &types.CartListResp{
			List: make([]types.CartItem, 0),
		}, nil
	}

	// 尝试将uid转换为int64
	var userId int64
	switch v := uidValue.(type) {
	case float64:
		userId = int64(v)
	case json.Number:
		userId, err = v.Int64()
		if err != nil {
			l.Logger.Errorf("DEBUG GetCart - Failed to convert uid to int64: %v", err)
			return &types.CartListResp{
				List: make([]types.CartItem, 0),
			}, nil
		}
	default:
		l.Logger.Errorf("DEBUG GetCart - Unexpected uid type: %T", v)
		return &types.CartListResp{
			List: make([]types.CartItem, 0),
		}, nil
	}

	l.Logger.Infof("DEBUG GetCart - Converted userId to int64: %d", userId)

	var carts []model.Cart
	if err := l.svcCtx.DB.Where("user_id = ?", userId).Find(&carts).Error; err != nil {
		l.Logger.Errorf("DEBUG GetCart - Failed to get cart items: %v", err)
		return &types.CartListResp{
			List: make([]types.CartItem, 0),
		}, nil
	}

	var result []types.CartItem
	for _, cart := range carts {
		// 获取商品信息
		product, err := l.svcCtx.ProductClient.Detail(l.ctx, &product.DetailRequest{
			Id: cart.ProductId,
		})
		if err != nil {
			l.Logger.Errorf("DEBUG GetCart - Failed to get product details for id %d: %v", cart.ProductId, err)
			continue
		}

		l.Logger.Infof("DEBUG GetCart - Product details for id %d: %+v", cart.ProductId, product)
		l.Logger.Infof("DEBUG GetCart - Product image URLs: %v", product.ImageUrls)

		// 使用商品主图
		productImage := ""
		if len(product.ImageUrls) > 0 {
			for _, url := range product.ImageUrls {
				if url != "" {
					productImage = url
					l.Logger.Infof("DEBUG GetCart - Using product image URL: %s", url)
					break
				}
			}
		}
		if productImage == "" {
			l.Logger.Infof("DEBUG GetCart - No valid image URL found for product %d", cart.ProductId)
		}

		result = append(result, types.CartItem{
			Id:           cart.Id,
			ProductId:    cart.ProductId,
			Quantity:     cart.Quantity,
			Selected:     cart.Selected,
			ProductName:  product.Name,
			ProductImage: productImage,
			Price:        float64(product.Amount) / 100,
			Stock:        int(product.Stock),
		})
	}

	l.Logger.Infof("DEBUG GetCart - Found %d cart items", len(result))
	return &types.CartListResp{
		List: result,
	}, nil
}
