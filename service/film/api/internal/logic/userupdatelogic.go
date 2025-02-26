package logic

import (
	"context"

	"mall/common/ctxdata"
	"mall/service/film/api/internal/svc"
	"mall/service/film/api/internal/types"
	"mall/service/film/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserUpdateLogic {
	return &UserUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserUpdateLogic) UserUpdate(req *types.UpdateFilmOrderReq) (resp *types.UpdateFilmOrderResp, err error) {
	l.Logger.Infof("用户更新胶片冲洗订单: %+v", req)

	// 从上下文获取当前登录用户ID
	uid, ok := ctxdata.GetUserIdFromCtx(l.ctx)
	if !ok {
		return &types.UpdateFilmOrderResp{
			Code: 401,
			Msg:  "请先登录",
		}, nil
	}

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

	// 检查订单是否属于当前用户
	if filmOrder.Uid != uid {
		return &types.UpdateFilmOrderResp{
			Code: 403,
			Msg:  "无权修改该订单",
		}, nil
	}

	// 用户只能取消订单（更新状态为待付款的订单）
	if req.Status != 0 || filmOrder.Status != 0 {
		return &types.UpdateFilmOrderResp{
			Code: 400,
			Msg:  "用户只能取消待付款的订单",
		}, nil
	}

	// 更新收货地址
	if req.AddressId > 0 {
		filmOrder.AddressId = req.AddressId
	}

	// 更新是否回寄底片
	if req.ReturnFilm != filmOrder.ReturnFilm {
		filmOrder.ReturnFilm = req.ReturnFilm
	}

	// 更新备注
	if req.Remark != "" {
		filmOrder.Remark = req.Remark
	}

	// 保存更新
	err = l.svcCtx.FilmOrderModel.Update(l.ctx, filmOrder)
	if err != nil {
		l.Logger.Errorf("更新胶片冲洗订单失败: %v", err)
		return &types.UpdateFilmOrderResp{
			Code: 500,
			Msg:  "更新订单失败",
		}, nil
	}

	return &types.UpdateFilmOrderResp{
		Code: 0,
		Msg:  "success",
	}, nil
}
