package logic

import (
	"context"

	"mall/service/film/api/internal/svc"
	"mall/service/film/api/internal/types"
	"mall/service/film/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminDeleteLogic {
	return &AdminDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminDeleteLogic) AdminDelete(req *types.DeleteFilmOrderReq) (resp *types.DeleteFilmOrderResp, err error) {
	l.Logger.Infof("管理员删除胶片冲洗订单: %+v", req)

	// 获取订单信息，确认存在
	filmOrder, err := l.svcCtx.FilmOrderModel.FindOne(l.ctx, req.Id)
	if err != nil {
		if err == model.ErrNotFound {
			return &types.DeleteFilmOrderResp{
				Code: 404,
				Msg:  "订单不存在",
			}, nil
		}

		l.Logger.Errorf("获取胶片冲洗订单失败: %v", err)
		return &types.DeleteFilmOrderResp{
			Code: 500,
			Msg:  "获取订单失败",
		}, nil
	}

	// 只能删除已完成的订单
	if filmOrder.Status != model.FilmOrderStatusFinished {
		return &types.DeleteFilmOrderResp{
			Code: 400,
			Msg:  "只能删除已完成的订单",
		}, nil
	}

	// 删除订单项
	items, err := l.svcCtx.FilmOrderItemModel.FindByFilmOrderId(l.ctx, req.Id)
	if err == nil {
		for _, item := range items {
			err := l.svcCtx.FilmOrderItemModel.Delete(l.ctx, item.Id)
			if err != nil {
				l.Logger.Errorf("删除胶片冲洗订单项失败: %v", err)
				// 继续尝试删除其他项
			}
		}
	}

	// 删除订单
	err = l.svcCtx.FilmOrderModel.Delete(l.ctx, req.Id)
	if err != nil {
		l.Logger.Errorf("删除胶片冲洗订单失败: %v", err)
		return &types.DeleteFilmOrderResp{
			Code: 500,
			Msg:  "删除订单失败",
		}, nil
	}

	return &types.DeleteFilmOrderResp{
		Code: 0,
		Msg:  "success",
	}, nil
}
