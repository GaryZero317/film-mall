package logic

import (
	"context"

	"mall/service/film/api/internal/svc"
	"mall/service/film/api/internal/types"
	"mall/service/film/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUpdateLogic {
	return &AdminUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminUpdateLogic) AdminUpdate(req *types.UpdateFilmOrderReq) (resp *types.UpdateFilmOrderResp, err error) {
	l.Logger.Infof("管理员更新胶片冲洗订单: %+v", req)
	l.Logger.Infof("请求中的ReturnFilm值: %v (类型: %T)", req.ReturnFilm, req.ReturnFilm)

	// 获取订单信息
	filmOrder, err := l.svcCtx.FilmOrderModel.FindOne(l.ctx, req.Id)
	if err != nil {
		if err == model.ErrNotFound {
			return &types.UpdateFilmOrderResp{
				Code: 404,
				Msg:  "订单不存在",
			}, nil
		}

		l.Logger.Errorf("获取胶片冲洗订单失败: %v", err)
		return &types.UpdateFilmOrderResp{
			Code: 500,
			Msg:  "获取订单失败",
		}, nil
	}

	l.Logger.Infof("更新前的订单数据: %+v", filmOrder)
	l.Logger.Infof("更新前的ReturnFilm值: %v (类型: %T)", filmOrder.ReturnFilm, filmOrder.ReturnFilm)

	// 管理员可以更新订单状态
	if req.Status >= 0 && req.Status <= 3 {
		filmOrder.Status = req.Status
	}

	// 更新收货地址
	if req.AddressId > 0 {
		filmOrder.AddressId = req.AddressId
	}

	// 更新是否回寄底片
	if req.ReturnFilm != filmOrder.ReturnFilm {
		l.Logger.Infof("更新ReturnFilm: %v -> %v", filmOrder.ReturnFilm, req.ReturnFilm)
		filmOrder.ReturnFilm = req.ReturnFilm
	} else {
		l.Logger.Infof("ReturnFilm值未变更: %v", filmOrder.ReturnFilm)
	}

	// 更新备注
	if req.Remark != "" {
		filmOrder.Remark = req.Remark
	}

	l.Logger.Infof("更新后的订单数据: %+v", filmOrder)
	l.Logger.Infof("即将保存的ReturnFilm值: %v (类型: %T)", filmOrder.ReturnFilm, filmOrder.ReturnFilm)

	// 保存更新
	err = l.svcCtx.FilmOrderModel.Update(l.ctx, filmOrder)
	if err != nil {
		l.Logger.Errorf("更新胶片冲洗订单失败: %v", err)
		return &types.UpdateFilmOrderResp{
			Code: 500,
			Msg:  "更新订单失败",
		}, nil
	}

	// 再次查询订单，验证更新是否成功
	updatedFilmOrder, err := l.svcCtx.FilmOrderModel.FindOne(l.ctx, req.Id)
	if err == nil {
		l.Logger.Infof("更新后查询的订单数据: %+v", updatedFilmOrder)
		l.Logger.Infof("更新后的ReturnFilm值: %v (类型: %T)", updatedFilmOrder.ReturnFilm, updatedFilmOrder.ReturnFilm)
	}

	return &types.UpdateFilmOrderResp{
		Code: 0,
		Msg:  "success",
	}, nil
}
