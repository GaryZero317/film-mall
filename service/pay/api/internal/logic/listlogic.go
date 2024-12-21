package logic

import (
	"context"

	"mall/service/pay/api/internal/svc"
	"mall/service/pay/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLogic {
	return &ListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListLogic) List(req *types.ListRequest) (resp *types.ListResponse, err error) {
	// 处理分页参数
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	// 查询数据库
	pays, total, err := l.svcCtx.PayModel.FindList(l.ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	// 转换数据格式
	var list []types.DetailResponse
	for _, pay := range pays {
		list = append(list, types.DetailResponse{
			Id:         pay.Id,
			Uid:        pay.Uid,
			Oid:        pay.Oid,
			Amount:     pay.Amount,
			Source:     pay.Source,
			Status:     pay.Status,
			CreateTime: pay.CreateTime.Format("2006-01-02 15:04:05"),
			UpdateTime: pay.UpdateTime.Format("2006-01-02 15:04:05"),
		})
	}

	return &types.ListResponse{
		Total: total,
		List:  list,
	}, nil
}
