package logic

import (
	"context"

	"mall/service/order/api/internal/svc"
	"mall/service/order/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminOrderListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminOrderListLogic {
	return &AdminOrderListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminOrderListLogic) AdminOrderList(req *types.AdminOrderListRequest) (resp *types.AdminOrderListResponse, err error) {
	// 获取分页参数
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 10
	}

	// 调用model层获取数据
	orders, total, err := l.svcCtx.OrderModel.FindPageListByPage(l.ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	// 转换数据
	var list []types.Order
	for _, o := range orders {
		list = append(list, types.Order{
			Id:         o.Id,
			Uid:        o.Uid,
			Pid:        o.Pid,
			Amount:     o.Amount,
			Status:     o.Status,
			CreateTime: o.CreateTime.Format("2006-01-02 15:04:05"),
			UpdateTime: o.UpdateTime.Format("2006-01-02 15:04:05"),
		})
	}

	return &types.AdminOrderListResponse{
		Total: total,
		List:  list,
	}, nil
}
