package logic

import (
	"context"

	"mall/service/user/rpc/internal/svc"
	"mall/service/user/rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/status"
)

type GetServiceListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetServiceListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetServiceListLogic {
	return &GetServiceListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetServiceListLogic) GetServiceList(in *user.ServiceListRequest) (*user.ServiceListResponse, error) {
	// 参数验证
	if in.UserId <= 0 {
		return nil, status.Error(100, "用户ID无效")
	}

	// 设置默认分页参数
	page := in.Page
	if page < 1 {
		page = 1
	}
	pageSize := in.PageSize
	if pageSize < 1 {
		pageSize = 10
	}

	// 查询数据
	list, total, err := l.svcCtx.GormCustomerServiceModel.FindByUserId(l.ctx, in.UserId, page, pageSize, in.Status)
	if err != nil {
		l.Logger.Errorf("查询客服问题列表失败: %v", err)
		return nil, status.Error(500, "查询客服问题列表失败")
	}

	// 组装响应
	var serviceItems []*user.ServiceItem
	for _, item := range list {
		serviceItems = append(serviceItems, &user.ServiceItem{
			Id:         item.Id,
			Title:      item.Title,
			Type:       item.Type,
			Status:     item.Status,
			CreateTime: item.CreateTime.Unix(),
			UpdateTime: item.UpdateTime.Unix(),
		})
	}

	return &user.ServiceListResponse{
		Total: total,
		List:  serviceItems,
	}, nil
}
