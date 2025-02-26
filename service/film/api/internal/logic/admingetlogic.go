package logic

import (
	"context"

	"mall/service/film/api/internal/svc"
	"mall/service/film/api/internal/types"
	"mall/service/film/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminGetLogic {
	return &AdminGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminGetLogic) AdminGet(req *types.GetFilmOrderReq) (resp *types.GetFilmOrderResp, err error) {
	l.Logger.Infof("管理员获取胶片冲洗订单: %+v", req)

	// 获取订单信息
	filmOrder, err := l.svcCtx.FilmOrderModel.FindOne(l.ctx, req.Id)
	if err != nil {
		if err == model.ErrNotFound {
			return &types.GetFilmOrderResp{
				Code: 404,
				Msg:  "订单不存在",
			}, nil
		}

		l.Logger.Errorf("获取胶片冲洗订单失败: %v", err)
		return &types.GetFilmOrderResp{
			Code: 500,
			Msg:  "获取订单失败",
		}, nil
	}

	// 获取订单项
	items, err := l.svcCtx.FilmOrderItemModel.FindByFilmOrderId(l.ctx, filmOrder.Id)
	if err != nil && err != model.ErrNotFound {
		l.Logger.Errorf("获取胶片冲洗订单项失败: %v", err)
		return &types.GetFilmOrderResp{
			Code: 500,
			Msg:  "获取订单项失败",
		}, nil
	}

	// 组装返回结果
	var itemList []types.FilmOrderItem
	for _, item := range items {
		itemList = append(itemList, types.FilmOrderItem{
			Id:          item.Id,
			FilmOrderId: item.FilmOrderId,
			FilmType:    item.FilmType,
			FilmBrand:   item.FilmBrand,
			Size:        item.Size,
			Quantity:    item.Quantity,
			Price:       item.Price,
			Amount:      item.Amount,
			Remark:      item.Remark,
		})
	}

	return &types.GetFilmOrderResp{
		Code: 0,
		Msg:  "success",
		Data: types.FilmOrder{
			Id:          filmOrder.Id,
			Foid:        filmOrder.Foid,
			Uid:         filmOrder.Uid,
			AddressId:   filmOrder.AddressId,
			ReturnFilm:  filmOrder.ReturnFilm,
			TotalPrice:  filmOrder.TotalPrice,
			ShippingFee: filmOrder.ShippingFee,
			Status:      filmOrder.Status,
			StatusDesc:  model.GetFilmOrderStatusText(filmOrder.Status),
			Remark:      filmOrder.Remark,
			Items:       itemList,
			CreateTime:  filmOrder.CreateTime.Format("2006-01-02 15:04:05"),
			UpdateTime:  filmOrder.UpdateTime.Format("2006-01-02 15:04:05"),
		},
	}, nil
}
