package logic

import (
	"context"
	"reflect"

	"mall/common/ctxdata"
	"mall/service/film/api/internal/svc"
	"mall/service/film/api/internal/types"
	"mall/service/film/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserListLogic {
	return &UserListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserListLogic) UserList(req *types.ListFilmOrderReq) (resp *types.ListFilmOrderResp, err error) {
	l.Logger.Infof("用户获取胶片冲洗订单列表: %+v", req)

	// 打印上下文中所有可用的键值
	l.Logger.Info("开始调试上下文信息...")

	// 尝试获取"uid"键
	l.Logger.Info("尝试从上下文获取uid...")
	uid, ok := ctxdata.GetUserIdFromCtx(l.ctx)

	// 记录调试信息
	l.Logger.Infof("uid获取结果: ok=%v, uid=%v, 类型=%v", ok, l.ctx.Value("uid"), reflect.TypeOf(l.ctx.Value("uid")))

	// 尝试获取可能的其他键名
	possibleKeys := []string{"userId", "user_id", "userID", "id", "UID", "ID"}
	for _, key := range possibleKeys {
		if l.ctx.Value(key) != nil {
			l.Logger.Infof("找到可能的用户ID: key=%s, value=%v, 类型=%v", key, l.ctx.Value(key), reflect.TypeOf(l.ctx.Value(key)))
		}
	}

	// 检查上下文中是否有JWT claims
	if l.ctx.Value("JWT_CLAIMS") != nil {
		l.Logger.Infof("找到JWT_CLAIMS: %+v", l.ctx.Value("JWT_CLAIMS"))
	}

	if !ok {
		l.Logger.Error("无法从上下文获取用户ID，认证失败")
		return &types.ListFilmOrderResp{
			Code: 401,
			Msg:  "请先登录",
		}, nil
	}

	// 使用当前用户ID，忽略请求中的uid
	status := req.Status
	if status == 0 {
		status = -1 // -1表示查询所有状态
	}

	// 获取订单列表
	l.Logger.Infof("开始查询用户(uid=%d)的订单列表，状态=%d", uid, status)
	filmOrders, total, err := l.svcCtx.FilmOrderModel.FindByUid(l.ctx, uid, status, req.Page, req.PageSize)
	if err != nil {
		l.Logger.Errorf("获取胶片冲洗订单列表失败: %v", err)
		return &types.ListFilmOrderResp{
			Code: 500,
			Msg:  "获取订单列表失败",
		}, nil
	}
	l.Logger.Infof("成功获取订单列表，共%d条记录", total)

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
