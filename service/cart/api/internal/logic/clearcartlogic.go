package logic

import (
	"context"
	"encoding/json"
	"mall/service/cart/api/internal/svc"
	"mall/service/cart/api/internal/types"
	"mall/service/cart/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type ClearCartLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewClearCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClearCartLogic {
	return &ClearCartLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ClearCartLogic) ClearCart() (resp *types.BaseResp, err error) {
	// 打印完整的context信息
	l.Logger.Infof("DEBUG ClearCart - Context values: %+v", l.ctx)

	// 从context中获取uid，JWT中的数字类型为json.Number
	uidValue := l.ctx.Value("uid")
	l.Logger.Infof("DEBUG ClearCart - Raw uid value type: %T", uidValue)

	if uidValue == nil {
		l.Logger.Errorf("DEBUG ClearCart - uid not found in context")
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
			l.Logger.Errorf("DEBUG ClearCart - Failed to convert uid to int64: %v", err)
			return &types.BaseResp{
				Code:    401,
				Message: "请先登录",
			}, nil
		}
	default:
		l.Logger.Errorf("DEBUG ClearCart - Unexpected uid type: %T", v)
		return &types.BaseResp{
			Code:    401,
			Message: "请先登录",
		}, nil
	}

	l.Logger.Infof("DEBUG ClearCart - Converted userId to int64: %d", userId)

	result := l.svcCtx.DB.Where("user_id = ?", userId).Delete(&model.Cart{})
	if result.Error != nil {
		return &types.BaseResp{
			Code:    1,
			Message: "清空购物车失败",
		}, nil
	}

	return &types.BaseResp{
		Code:    0,
		Message: "清空成功",
	}, nil
}
