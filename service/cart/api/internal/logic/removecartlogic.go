package logic

import (
	"context"
	"encoding/json"
	"mall/service/cart/api/internal/svc"
	"mall/service/cart/api/internal/types"
	"mall/service/cart/model"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveCartLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRemoveCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveCartLogic {
	return &RemoveCartLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RemoveCartLogic) RemoveCart(id string) (resp *types.BaseResp, err error) {
	// 打印完整的context信息
	l.Logger.Infof("DEBUG RemoveCart - Context values: %+v", l.ctx)

	// 从context中获取uid，JWT中的数字类型为json.Number
	uidValue := l.ctx.Value("uid")
	l.Logger.Infof("DEBUG RemoveCart - Raw uid value type: %T", uidValue)

	if uidValue == nil {
		l.Logger.Errorf("DEBUG RemoveCart - uid not found in context")
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
			l.Logger.Errorf("DEBUG RemoveCart - Failed to convert uid to int64: %v", err)
			return &types.BaseResp{
				Code:    401,
				Message: "请先登录",
			}, nil
		}
	default:
		l.Logger.Errorf("DEBUG RemoveCart - Unexpected uid type: %T", v)
		return &types.BaseResp{
			Code:    401,
			Message: "请先登录",
		}, nil
	}

	l.Logger.Infof("DEBUG RemoveCart - Converted userId to int64: %d", userId)

	cartId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return &types.BaseResp{
			Code:    400,
			Message: "无效的购物车ID",
		}, nil
	}

	var cart model.Cart
	result := l.svcCtx.DB.Where("id = ? AND user_id = ?", cartId, userId).First(&cart)
	if result.Error != nil {
		return &types.BaseResp{
			Code:    400,
			Message: "购物车商品不存在",
		}, nil
	}

	if err := l.svcCtx.DB.Delete(&cart).Error; err != nil {
		return &types.BaseResp{
			Code:    500,
			Message: "删除失败",
		}, nil
	}

	return &types.BaseResp{
		Code:    200,
		Message: "删除成功",
	}, nil
}
