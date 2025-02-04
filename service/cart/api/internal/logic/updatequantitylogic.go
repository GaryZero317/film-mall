package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"mall/service/cart/api/internal/svc"
	"mall/service/cart/api/internal/types"
	"mall/service/cart/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateQuantityLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateQuantityLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateQuantityLogic {
	return &UpdateQuantityLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateQuantityLogic) UpdateQuantity(req *types.UpdateQuantityReq) (resp *types.BaseResp, err error) {
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

	var cart model.Cart
	result := l.svcCtx.DB.Where("id = ? AND user_id = ?", req.Id, userId).First(&cart)
	if result.Error != nil {
		return &types.BaseResp{
			Code:    1,
			Message: "购物车商品不存在",
		}, nil
	}

	cart.Quantity = req.Quantity
	if err := l.svcCtx.DB.Save(&cart).Error; err != nil {
		return &types.BaseResp{
			Code:    1,
			Message: "更新失败",
		}, nil
	}

	// 删除缓存，让下次获取购物车列表时重新加载
	cacheKey := fmt.Sprintf("cart:%d", userId)
	_, _ = l.svcCtx.Cache.Del(cacheKey)

	return &types.BaseResp{
		Code:    0,
		Message: "更新成功",
	}, nil
}
