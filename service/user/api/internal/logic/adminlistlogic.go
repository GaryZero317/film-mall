package logic

import (
	"context"

	"mall/service/user/api/internal/svc"
	"mall/service/user/api/internal/types"

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

func (l *AdminListLogic) AdminList(req *types.AdminListRequest) (resp *types.AdminListResponse, err error) {
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
	admins, total, err := l.svcCtx.AdminModel.FindList(l.ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	// 转换数据格式
	var list []types.AdminInfoResponse
	for _, admin := range admins {
		list = append(list, types.AdminInfoResponse{
			Id:       admin.ID,
			Username: admin.Username,
			Level:    int32(admin.Level),
		})
	}

	return &types.AdminListResponse{
		Total: total,
		List:  list,
	}, nil
}
