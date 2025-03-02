package logic

import (
	"context"

	"mall/service/user/api/internal/svc"
	"mall/service/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminServiceListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminServiceListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminServiceListLogic {
	return &AdminServiceListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminServiceListLogic) AdminServiceList(req *types.AdminServiceListRequest) (resp *types.AdminServiceListResponse, err error) {
	// 设置默认分页参数
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 10
	}

	// 查询数据（注意：管理员可以查看所有客服问题，无需用户ID筛选）
	list, total, err := l.svcCtx.GormCustomerServiceModel.FindAll(l.ctx, page, pageSize, req.Status, req.Type)
	if err != nil {
		l.Logger.Errorf("查询客服问题列表失败: %v", err)
		return nil, err
	}

	// 组装响应
	var itemList []types.ServiceItemResponse
	for _, item := range list {
		itemList = append(itemList, types.ServiceItemResponse{
			Id:         item.Id,
			Title:      item.Title,
			Type:       item.Type,
			Status:     item.Status,
			CreateTime: item.CreateTime.Unix(),
			UpdateTime: item.UpdateTime.Unix(),
		})
	}

	return &types.AdminServiceListResponse{
		Total: total,
		List:  itemList,
	}, nil
}
