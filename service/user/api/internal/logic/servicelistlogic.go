package logic

import (
	"context"
	"encoding/json"

	"mall/service/user/api/internal/svc"
	"mall/service/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ServiceListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewServiceListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ServiceListLogic {
	return &ServiceListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ServiceListLogic) ServiceList(req *types.ServiceListRequest) (resp *types.ServiceListResponse, err error) {
	// 从上下文获取用户ID
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()

	// 设置默认分页参数
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 10
	}

	// 查询数据
	list, total, err := l.svcCtx.CustomerServiceModel.FindByUserId(l.ctx, uid, page, pageSize, req.Status)
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

	return &types.ServiceListResponse{
		Total: total,
		List:  itemList,
	}, nil
}
