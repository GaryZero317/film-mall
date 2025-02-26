package logic

import (
	"context"
	"fmt"
	"time"

	"mall/service/film/api/internal/svc"
	"mall/service/film/api/internal/types"
	"mall/service/film/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 生成订单号
func generateOrderNo() string {
	now := time.Now()
	return fmt.Sprintf("FM%s%d", now.Format("20060102150405"), now.UnixNano()%1000)
}

func (l *CreateLogic) Create(req *types.CreateFilmOrderReq) (resp *types.CreateFilmOrderResp, err error) {
	l.Logger.Infof("创建胶片冲洗订单: %+v", req)

	// 参数校验
	if req.Uid <= 0 {
		return &types.CreateFilmOrderResp{
			Code: 400,
			Msg:  "用户ID不能为空",
		}, nil
	}

	if req.AddressId <= 0 {
		return &types.CreateFilmOrderResp{
			Code: 400,
			Msg:  "收货地址不能为空",
		}, nil
	}

	if req.TotalPrice <= 0 {
		return &types.CreateFilmOrderResp{
			Code: 400,
			Msg:  "订单总价不能为空",
		}, nil
	}

	if len(req.Items) == 0 {
		return &types.CreateFilmOrderResp{
			Code: 400,
			Msg:  "订单项不能为空",
		}, nil
	}

	// 生成订单号
	foid := generateOrderNo()

	// 创建订单
	filmOrder := &model.FilmOrder{
		Foid:        foid,
		Uid:         req.Uid,
		AddressId:   req.AddressId,
		ReturnFilm:  req.ReturnFilm,
		TotalPrice:  req.TotalPrice,
		ShippingFee: req.ShippingFee,
		Status:      req.Status,
		Remark:      req.Remark,
	}

	// 插入订单
	result, err := l.svcCtx.FilmOrderModel.Insert(l.ctx, filmOrder)
	if err != nil {
		l.Logger.Errorf("插入胶片冲洗订单失败: %v", err)
		return &types.CreateFilmOrderResp{
			Code: 500,
			Msg:  "创建订单失败",
		}, nil
	}

	// 获取订单ID
	filmOrderId, err := result.LastInsertId()
	if err != nil {
		l.Logger.Errorf("获取胶片冲洗订单ID失败: %v", err)
		return &types.CreateFilmOrderResp{
			Code: 500,
			Msg:  "创建订单失败",
		}, nil
	}

	// 批量添加订单项
	for _, item := range req.Items {
		filmOrderItem := &model.FilmOrderItem{
			FilmOrderId: filmOrderId,
			FilmType:    item.FilmType,
			FilmBrand:   item.FilmBrand,
			Size:        item.Size,
			Quantity:    item.Quantity,
			Price:       item.Price,
			Amount:      item.Amount,
			Remark:      item.Remark,
		}

		_, err = l.svcCtx.FilmOrderItemModel.Insert(l.ctx, filmOrderItem)
		if err != nil {
			l.Logger.Errorf("插入胶片冲洗订单项失败: %v", err)
			// 这里应该有事务回滚机制，简化处理
			return &types.CreateFilmOrderResp{
				Code: 500,
				Msg:  "创建订单项失败",
			}, nil
		}
	}

	return &types.CreateFilmOrderResp{
		Code: 0,
		Msg:  "success",
		Data: types.CreateFilmOrderData{
			Id:   filmOrderId,
			Foid: foid,
		},
	}, nil
}
