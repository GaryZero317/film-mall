package logic

import (
	"context"

	"mall/service/user/api/internal/svc"
	"mall/service/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FaqListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFaqListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FaqListLogic {
	return &FaqListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FaqListLogic) FaqList(req *types.FaqListRequest) (resp *types.FaqListResponse, err error) {
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
	list, total, err := l.svcCtx.FaqModel.FindList(l.ctx, page, pageSize, req.Category)
	if err != nil {
		l.Logger.Errorf("查询常见问题列表失败: %v", err)
		return nil, err
	}

	// 组装响应
	var itemList []types.FaqItemResponse
	for _, item := range list {
		itemList = append(itemList, types.FaqItemResponse{
			Id:       item.Id,
			Question: item.Question,
			Answer:   item.Answer,
			Category: item.Category,
		})
	}

	return &types.FaqListResponse{
		Total: total,
		List:  itemList,
	}, nil
}
