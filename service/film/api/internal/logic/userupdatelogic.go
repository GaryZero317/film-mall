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
	l.Logger.Infof("请求中的ReturnFilm值: %v (类型: %T)", req.ReturnFilm, req.ReturnFilm)

	// 获取用户ID
	uid, ok := ctxdata.GetUserIdFromCtx(l.ctx)
	if !ok {
		return &types.UpdateFilmOrderResp{
			Code: 401,
			Msg:  "未认证的用户",
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

	l.Logger.Infof("更新前的订单数据: %+v", filmOrder)
	l.Logger.Infof("更新前的ReturnFilm值: %v (类型: %T)", filmOrder.ReturnFilm, filmOrder.ReturnFilm)

	// 检查订单是否属于当前用户
	if filmOrder.Uid != uid {
		return &types.UpdateFilmOrderResp{
			Code: 403,
			Msg:  "无权修改该订单",
		}, nil
	}

	// 修改订单状态逻辑：
	// 1. 允许用户取消订单（将状态设为0）
	// 2. 允许用户将待付款(0)订单更新为冲洗处理中(1)，模拟支付成功
	// 3. 其他状态更新不允许
	if req.Status >= 0 {
		// 待付款 -> 冲洗处理中 (模拟支付成功)
		if req.Status == 1 && filmOrder.Status == 0 {
			filmOrder.Status = 1
			l.Logger.Infof("用户支付订单：状态从待付款(0)更新为冲洗处理中(1), 订单ID: %d", req.Id)

			// 支付操作时不修改ReturnFilm值，保持原值
			l.Logger.Infof("支付操作，保持原始ReturnFilm值: %v", filmOrder.ReturnFilm)

			// 如果不回寄底片，则不收取邮费
			if !filmOrder.ReturnFilm && filmOrder.ShippingFee > 0 {
				l.Logger.Infof("不回寄底片，邮费从 %d 调整为 0", filmOrder.ShippingFee)
				filmOrder.ShippingFee = 0
			} else if filmOrder.ReturnFilm {
				l.Logger.Infof("回寄底片，保留邮费 %d", filmOrder.ShippingFee)
			}
		} else if req.Status == 0 {
			// 取消订单 (任何状态 -> 待付款)
			filmOrder.Status = 0

			// 取消订单时可以修改ReturnFilm
			if req.ReturnFilm != filmOrder.ReturnFilm {
				l.Logger.Infof("更新ReturnFilm: %v -> %v", filmOrder.ReturnFilm, req.ReturnFilm)
				filmOrder.ReturnFilm = req.ReturnFilm
			}
		} else {
			// 其他状态变更不允许
			return &types.UpdateFilmOrderResp{
				Code: 400,
				Msg:  "用户只能取消待付款的订单或将待付款订单更新为冲洗处理中",
			}, nil
		}
	} else {
		// 仅更新ReturnFilm，不改变状态
		if req.ReturnFilm != filmOrder.ReturnFilm {
			l.Logger.Infof("更新ReturnFilm: %v -> %v", filmOrder.ReturnFilm, req.ReturnFilm)
			filmOrder.ReturnFilm = req.ReturnFilm
		} else {
			l.Logger.Infof("ReturnFilm值未变更: %v", filmOrder.ReturnFilm)
		}
	}

	// 更新收货地址
	if req.AddressId > 0 {
		filmOrder.AddressId = req.AddressId
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
