package logic

import (
	"context"
	"fmt"

	"mall/service/order/rpc/order"
	"mall/service/pay/model"
	"mall/service/pay/rpc/internal/svc"
	"mall/service/pay/rpc/types/pay"
	"mall/service/user/rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/status"
)

type CallbackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CallbackLogic {
	return &CallbackLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CallbackLogic) Callback(in *pay.CallbackRequest) (*pay.CallbackResponse, error) {
	l.Logger.Info("支付回调请求", logx.Field("request", in))
	fmt.Printf("支付回调请求: Oid=%d, Status=%d\n", in.Oid, in.Status)

	// 查询用户是否存在（如果提供了用户ID）
	if in.Uid > 0 {
		_, err := l.svcCtx.UserRpc.UserInfo(l.ctx, &user.UserInfoRequest{
			Id: in.Uid,
		})
		if err != nil {
			l.Logger.Errorf("查询用户失败: %v", err)
			return nil, status.Error(100, "用户不存在")
		}
	}

	// 查询订单是否存在（如果不是取消操作）
	if in.Oid > 0 && in.Status != 4 && in.Status != 9 {
		_, err := l.svcCtx.OrderRpc.Detail(l.ctx, &order.DetailRequest{
			Id: in.Oid,
		})
		if err != nil {
			l.Logger.Errorf("查询订单失败: %v", err)
			return nil, status.Error(100, "订单不存在")
		}
	}

	var res *model.Pay
	var err error

	// 取消订单操作特殊处理
	if (in.Status == 4 || in.Status == 9) && in.Oid > 0 {
		// 对取消操作进行特殊处理
		l.Logger.Info("处理取消订单操作", logx.Field("orderId", in.Oid))
		fmt.Printf("处理取消订单操作: Oid=%d\n", in.Oid)

		// 尝试查询支付记录
		res, err = l.svcCtx.PayModel.FindOneByOid(l.ctx, in.Oid)
		if err != nil {
			if err == model.ErrNotFound {
				l.Logger.Info("未找到支付记录，但这是取消操作，视为成功", logx.Field("orderId", in.Oid))
				fmt.Printf("未找到订单ID=%d的支付记录，但这是取消操作，视为成功\n", in.Oid)
				return &pay.CallbackResponse{}, nil
			}
			l.Logger.Error("查询支付记录出错", logx.Field("error", err))
			return nil, status.Error(500, err.Error())
		}

		l.Logger.Info("找到支付记录，更新为已取消", logx.Field("orderId", in.Oid), logx.Field("payId", res.Id))
		fmt.Printf("找到订单ID=%d的支付记录(ID=%d)，更新为已取消\n", in.Oid, res.Id)

		// 更新为取消状态
		res.Status = in.Status
		err = l.svcCtx.PayModel.Update(l.ctx, res)
		if err != nil {
			l.Logger.Error("更新支付记录失败", logx.Field("error", err))
			return nil, status.Error(500, err.Error())
		}

		l.Logger.Info("支付记录已成功更新为取消状态", logx.Field("orderId", in.Oid))
		fmt.Printf("订单ID=%d的支付记录已成功更新为取消状态\n", in.Oid)
		return &pay.CallbackResponse{}, nil
	}

	// 非取消操作，正常查询支付记录
	if in.Id > 0 {
		// 按支付ID查询
		l.Logger.Info("按支付ID查询", logx.Field("payId", in.Id))
		res, err = l.svcCtx.PayModel.FindOne(l.ctx, in.Id)
	} else if in.Oid > 0 {
		// 按订单ID查询
		l.Logger.Info("按订单ID查询", logx.Field("orderId", in.Oid))
		res, err = l.svcCtx.PayModel.FindOneByOid(l.ctx, in.Oid)
	}

	// 处理查询结果
	if err != nil {
		if err == model.ErrNotFound {
			l.Logger.Error("支付记录不存在", logx.Field("orderId", in.Oid))
			return nil, status.Error(100, "支付记录不存在")
		}
		l.Logger.Error("查询支付记录失败", logx.Field("error", err))
		return nil, status.Error(500, err.Error())
	}

	// 如果是支付成功，验证支付金额
	if in.Status == 1 && in.Amount > 0 && in.Amount != res.Amount {
		l.Logger.Error("支付金额与订单金额不符",
			logx.Field("expected", res.Amount),
			logx.Field("actual", in.Amount))
		return nil, status.Error(100, "支付金额与订单金额不符")
	}

	// 更新支付记录状态
	res.Status = in.Status
	if in.Source > 0 {
		res.Source = in.Source
	}

	err = l.svcCtx.PayModel.Update(l.ctx, res)
	if err != nil {
		l.Logger.Error("更新支付记录失败", logx.Field("error", err))
		return nil, status.Error(500, err.Error())
	}

	// 仅当支付成功时，调用订单服务更新订单状态为已支付
	if in.Status == 1 {
		l.Logger.Info("支付成功，更新订单状态", logx.Field("orderId", in.Oid))
		_, err = l.svcCtx.OrderRpc.Paid(l.ctx, &order.PaidRequest{
			Id: in.Oid,
		})
		if err != nil {
			l.Logger.Error("更新订单状态失败", logx.Field("error", err))
			return nil, status.Error(500, err.Error())
		}
	}

	l.Logger.Info("支付记录状态更新成功",
		logx.Field("orderId", in.Oid),
		logx.Field("status", in.Status))
	return &pay.CallbackResponse{}, nil
}
