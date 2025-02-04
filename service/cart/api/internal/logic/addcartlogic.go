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

type AddCartLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddCartLogic {
	return &AddCartLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddCartLogic) AddCart(req *types.AddCartReq) (resp *types.BaseResp, err error) {
	// 从context中获取uid
	uidValue := l.ctx.Value("uid")
	if uidValue == nil {
		return &types.BaseResp{
			Code:    401,
			Message: "请先登录",
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
			return &types.BaseResp{
				Code:    401,
				Message: "请先登录",
			}, nil
		}
	default:
		return &types.BaseResp{
			Code:    401,
			Message: "请先登录",
		}, nil
	}

	// 检查商品是否存在
	l.Logger.Infof("检查商品是否存在, 商品ID: %d", req.ProductId)
	product, err := l.svcCtx.ProductClient.Detail(l.ctx, &product.DetailRequest{
		Id: req.ProductId,
	})
	if err != nil {
		l.Logger.Errorf("获取商品详情失败: %v", err)
		return &types.BaseResp{
			Code:    1,
			Message: "商品不存在",
		}, nil
	}
	l.Logger.Infof("商品详情: %+v", product)

	// 检查购物车中是否已存在该商品
	l.Logger.Infof("检查购物车是否已存在该商品, 用户ID: %d, 商品ID: %d", userId, req.ProductId)
	var cart model.Cart
	result := l.svcCtx.DB.Where("user_id = ? AND product_id = ?", userId, req.ProductId).First(&cart)

	quantity := req.Quantity
	if quantity == 0 {
		quantity = 1
	}

	if result.Error == nil {
		// 商品已存在，更新数量
		l.Logger.Infof("商品已存在购物车，更新数量 原数量: %d, 新增数量: %d", cart.Quantity, quantity)
		cart.Quantity += quantity
		if err := l.svcCtx.DB.Save(&cart).Error; err != nil {
			l.Logger.Errorf("更新购物车失败: %v", err)
			return &types.BaseResp{
				Code:    1,
				Message: "更新购物车失败",
			}, nil
		}
		l.Logger.Info("更新购物车成功")
	} else {
		// 商品不存在，新增
		l.Logger.Info("商品不在购物车中，新增记录")
		cart = model.Cart{
			UserId:    userId,
			ProductId: product.Id,
			Quantity:  quantity,
			Selected:  true,
		}
		if err := l.svcCtx.DB.Create(&cart).Error; err != nil {
			l.Logger.Errorf("添加购物车失败: %v", err)
			return &types.BaseResp{
				Code:    1,
				Message: "添加购物车失败",
			}, nil
		}
		l.Logger.Info("添加购物车成功")
	}

	// 删除缓存，让下次获取购物车列表时重新加载
	cacheKey := fmt.Sprintf("cart:%d", userId)
	l.Logger.Infof("删除购物车缓存, key: %s", cacheKey)
	_, err = l.svcCtx.Cache.Del(cacheKey)
	if err != nil {
		l.Logger.Errorf("删除缓存失败: %v", err)
	}

	return &types.BaseResp{
		Code:    0,
		Message: "添加成功",
	}, nil
}
