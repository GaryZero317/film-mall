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

func (l *RemoveCartLogic) RemoveCart(req *types.RemoveCartReq) (resp *types.BaseResp, err error) {
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

	result := l.svcCtx.DB.Where("id = ? AND user_id = ?", req.Id, userId).Delete(&model.Cart{})
	if result.Error != nil {
		return &types.BaseResp{
			Code:    1,
			Message: "删除失败",
		}, nil
	}

	if result.RowsAffected == 0 {
		return &types.BaseResp{
			Code:    1,
			Message: "购物车商品不存在",
		}, nil
	}

	// 删除缓存，让下次获取购物车列表时重新加载
	cacheKey := fmt.Sprintf("cart:%d", userId)
	_, _ = l.svcCtx.Cache.Del(cacheKey)

	return &types.BaseResp{
		Code:    0,
		Message: "删除成功",
	}, nil
}
