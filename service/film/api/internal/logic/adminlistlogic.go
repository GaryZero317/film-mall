package logic

import (
	"context"

	"mall/service/film/api/internal/svc"
	"mall/service/film/api/internal/types"
	"mall/service/film/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListLogic {
	return &AdminListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminListLogic) AdminList(req *types.ListFilmOrderReq) (resp *types.ListFilmOrderResp, err error) {
	l.Logger.Infof("管理员获取胶片冲洗订单列表: %+v", req)

	var filmOrders []*model.FilmOrder
	var total int64
	var queryErr error

	status := req.Status
	if status == 0 {
		status = -1 // -1表示查询所有状态
	}

	// 如果指定了用户ID，则查询指定用户的订单
	if req.Uid > 0 {
		filmOrders, total, queryErr = l.svcCtx.FilmOrderModel.FindByUid(l.ctx, req.Uid, status, req.Page, req.PageSize)
	} else {
		// 否则查询所有订单
		filmOrders, total, queryErr = l.svcCtx.FilmOrderModel.FindAll(l.ctx, status, req.Page, req.PageSize)
	}

	if queryErr != nil {
		l.Logger.Errorf("获取胶片冲洗订单列表失败: %v", queryErr)
		return &types.ListFilmOrderResp{
			Code: 500,
			Msg:  "获取订单列表失败",
		}, nil
	}

	// 组装返回结果
	var orderList []types.FilmOrder
	for _, order := range filmOrders {
		// 获取订单项
		items, _ := l.svcCtx.FilmOrderItemModel.FindByFilmOrderId(l.ctx, order.Id)

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

		orderList = append(orderList, types.FilmOrder{
			Id:          order.Id,
			Foid:        order.Foid,
			Uid:         order.Uid,
			AddressId:   order.AddressId,
			ReturnFilm:  order.ReturnFilm,
			TotalPrice:  order.TotalPrice,
			ShippingFee: order.ShippingFee,
			Status:      order.Status,
			StatusDesc:  model.GetFilmOrderStatusText(order.Status),
			Remark:      order.Remark,
			Items:       itemList,
			CreateTime:  order.CreateTime.Format("2006-01-02 15:04:05"),
			UpdateTime:  order.UpdateTime.Format("2006-01-02 15:04:05"),
		})
	}

	return &types.ListFilmOrderResp{
		Code: 0,
		Msg:  "success",
		Data: types.ListFilmOrderData{
			Total: total,
			List:  orderList,
		},
	}, nil
}
