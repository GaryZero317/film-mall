package logic

import (
	"context"
	"encoding/json"
	"fmt"
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
	// 从context中获取uid
	uidValue := l.ctx.Value("uid")
	if uidValue == nil {
		return &types.CartListResp{
			List: make([]types.CartItem, 0),
		}, nil
	}

	// 转换uid为int64
	var userId int64
	switch v := uidValue.(type) {
	case float64:
		userId = int64(v)
	case json.Number:
		userId, err = v.Int64()
		if err != nil {
			return &types.CartListResp{
				List: make([]types.CartItem, 0),
			}, nil
		}
	default:
		return &types.CartListResp{
			List: make([]types.CartItem, 0),
		}, nil
	}

	// 尝试从缓存获取购物车数据
	cacheKey := fmt.Sprintf("cart:%d", userId)
	l.Logger.Infof("正在从缓存获取购物车数据, key: %s", cacheKey)

	var cartItems []types.CartItem
	val, err := l.svcCtx.Cache.Get(cacheKey)
	if err != nil {
		l.Logger.Infof("缓存未命中或发生错误: %v", err)
	} else {
		l.Logger.Infof("获取到缓存数据: %s", val)
	}

	if err == nil && val != "" {
		err = json.Unmarshal([]byte(val), &cartItems)
		if err == nil {
			l.Logger.Info("成功解析缓存数据，返回购物车列表")
			return &types.CartListResp{
				List: cartItems,
			}, nil
		}
		l.Logger.Errorf("解析缓存数据失败: %v", err)
	}

	// 缓存未命中，从数据库获取
	l.Logger.Info("从数据库获取购物车数据")
	var carts []model.Cart
	if err := l.svcCtx.DB.Where("user_id = ?", userId).Find(&carts).Error; err != nil {
		l.Logger.Errorf("数据库查询失败: %v", err)
		return &types.CartListResp{
			List: make([]types.CartItem, 0),
		}, nil
	}
	l.Logger.Infof("数据库查询到 %d 条购物车记录", len(carts))

	// 获取商品详情并组装数据
	cartItems = make([]types.CartItem, 0)
	for _, cart := range carts {
		l.Logger.Infof("获取商品详情, 商品ID: %d", cart.ProductId)
		product, err := l.svcCtx.ProductClient.Detail(l.ctx, &product.DetailRequest{
			Id: cart.ProductId,
		})
		if err != nil {
			l.Logger.Errorf("获取商品详情失败: %v", err)
			continue
		}
		l.Logger.Infof("商品详情: %+v", product)

		cartItems = append(cartItems, types.CartItem{
			Id:           cart.Id,
			ProductId:    cart.ProductId,
			Quantity:     cart.Quantity,
			Selected:     cart.Selected,
			ProductName:  product.Name,
			ProductImage: product.ImageUrls[0],
			Price:        float64(product.Amount) / 100,
			Stock:        int(product.Stock),
		})
	}

	// 将数据写入缓存
	if len(cartItems) > 0 {
		jsonData, _ := json.Marshal(cartItems)
		l.Logger.Infof("写入缓存, key: %s, data: %s", cacheKey, string(jsonData))
		err = l.svcCtx.Cache.Set(cacheKey, string(jsonData))
		if err != nil {
			l.Logger.Errorf("写入缓存失败: %v", err)
		}
	} else {
		l.Logger.Info("购物车为空，不写入缓存")
	}

	return &types.CartListResp{
		List: cartItems,
	}, nil
}
