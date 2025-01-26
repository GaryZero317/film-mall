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
	// 打印完整的context信息
	l.Logger.Infof("DEBUG AddCart - Context values: %+v", l.ctx)

	// 从context中获取uid，JWT中的数字类型为json.Number
	uidValue := l.ctx.Value("uid")
	l.Logger.Infof("DEBUG AddCart - Raw uid value type: %T", uidValue)

	if uidValue == nil {
		l.Logger.Errorf("DEBUG AddCart - uid not found in context")
		return &types.BaseResp{
			Code:    401,
			Message: "请先登录",
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
			l.Logger.Errorf("DEBUG AddCart - Failed to convert uid to int64: %v", err)
			return &types.BaseResp{
				Code:    401,
				Message: "请先登录",
			}, nil
		}
	default:
		l.Logger.Errorf("DEBUG AddCart - Unexpected uid type: %T", v)
		return &types.BaseResp{
			Code:    401,
			Message: "请先登录",
		}, nil
	}

	l.Logger.Infof("DEBUG AddCart - Converted userId to int64: %d", userId)

	// 检查商品是否存在
	product, err := l.svcCtx.ProductClient.Detail(l.ctx, &product.DetailRequest{
		Id: req.ProductId,
	})
	if err != nil {
		l.Logger.Errorf("DEBUG AddCart - Failed to get product details: %v", err)
		return &types.BaseResp{
			Code:    1,
			Message: "商品不存在",
		}, nil
	}

	// 检查购物车中是否已存在该商品
	var cart model.Cart
	result := l.svcCtx.DB.Where("user_id = ? AND product_id = ?", userId, req.ProductId).First(&cart)

	quantity := req.Quantity
	if quantity == 0 {
		quantity = 1
	}

	if result.Error == nil {
		// 商品已存在，更新数量
		cart.Quantity += quantity
		if err := l.svcCtx.DB.Save(&cart).Error; err != nil {
			return &types.BaseResp{
				Code:    1,
				Message: "更新购物车失败",
			}, nil
		}
	} else {
		// 商品不存在，新增
		cart = model.Cart{
			UserId:    userId,
			ProductId: product.Id,
			Quantity:  quantity,
			Selected:  true,
		}
		if err := l.svcCtx.DB.Create(&cart).Error; err != nil {
			return &types.BaseResp{
				Code:    1,
				Message: "添加购物车失败",
			}, nil
		}
	}

	return &types.BaseResp{
		Code:    0,
		Message: "添加成功",
	}, nil
}
