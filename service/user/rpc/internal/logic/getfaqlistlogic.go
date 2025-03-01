package logic

import (
	"context"

	"mall/service/user/rpc/internal/svc"
	"mall/service/user/rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/status"
)

type GetFaqListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFaqListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFaqListLogic {
	return &GetFaqListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFaqListLogic) GetFaqList(in *user.FaqListRequest) (*user.FaqListResponse, error) {
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
	list, total, err := l.svcCtx.FaqModel.FindList(l.ctx, page, pageSize, in.Category)
	if err != nil {
		l.Logger.Errorf("查询常见问题列表失败: %v", err)
		return nil, status.Error(500, "查询常见问题列表失败")
	}

	// 组装响应
	var faqItems []*user.FaqItem
	for _, item := range list {
		faqItems = append(faqItems, &user.FaqItem{
			Id:       item.Id,
			Question: item.Question,
			Answer:   item.Answer,
			Category: item.Category,
		})
	}

	return &user.FaqListResponse{
		Total: total,
		List:  faqItems,
	}, nil
}
