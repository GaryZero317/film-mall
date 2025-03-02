package logic

import (
	"context"

	"mall/service/user/api/internal/svc"
	"mall/service/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminServiceDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminServiceDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminServiceDetailLogic {
	return &AdminServiceDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminServiceDetailLogic) AdminServiceDetail(req *types.AdminServiceDetailRequest) (resp *types.AdminServiceDetailResponse, err error) {
	// 查询客服问题详情
	serviceQuestion, err := l.svcCtx.GormCustomerServiceModel.FindOne(l.ctx, req.Id)
	if err != nil {
		l.Logger.Errorf("查询客服问题详情失败: %v", err)
		return nil, err
	}

	// 转换回复时间
	var replyTime int64
	if serviceQuestion.ReplyTime != nil && !serviceQuestion.ReplyTime.IsZero() {
		replyTime = serviceQuestion.ReplyTime.Unix()
	} else {
		replyTime = 0
	}

	// 返回响应
	return &types.AdminServiceDetailResponse{
		Id:         serviceQuestion.Id,
		UserId:     serviceQuestion.UserId,
		Title:      serviceQuestion.Title,
		Content:    serviceQuestion.Content,
		Type:       serviceQuestion.Type,
		Status:     serviceQuestion.Status,
		Reply:      serviceQuestion.Reply,
		ReplyTime:  replyTime,
		ContactWay: serviceQuestion.ContactWay,
		CreateTime: serviceQuestion.CreateTime.Unix(),
		UpdateTime: serviceQuestion.UpdateTime.Unix(),
	}, nil
}
