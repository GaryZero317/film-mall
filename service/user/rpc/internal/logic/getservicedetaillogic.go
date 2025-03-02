package logic

import (
	"context"

	"mall/service/user/model"
	"mall/service/user/rpc/internal/svc"
	"mall/service/user/rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/status"
)

type GetServiceDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetServiceDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetServiceDetailLogic {
	return &GetServiceDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetServiceDetailLogic) GetServiceDetail(in *user.ServiceDetailRequest) (*user.ServiceDetailResponse, error) {
	// 参数验证
	if in.UserId <= 0 {
		return nil, status.Error(100, "用户ID无效")
	}
	if in.Id <= 0 {
		return nil, status.Error(100, "问题ID无效")
	}

	// 查询问题详情
	service, err := l.svcCtx.GormCustomerServiceModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, status.Error(100, "问题不存在")
		}
		l.Logger.Errorf("查询客服问题详情失败: %v", err)
		return nil, status.Error(500, "查询客服问题详情失败")
	}

	// 检查问题是否属于当前用户
	if service.UserId != in.UserId {
		return nil, status.Error(100, "无权查看此问题")
	}

	// 组装响应
	replyTime := int64(0)
	if service.ReplyTime != nil && !service.ReplyTime.IsZero() {
		replyTime = service.ReplyTime.Unix()
	}

	return &user.ServiceDetailResponse{
		Id:         service.Id,
		Title:      service.Title,
		Content:    service.Content,
		Type:       service.Type,
		Status:     service.Status,
		Reply:      service.Reply,
		ReplyTime:  replyTime,
		ContactWay: service.ContactWay,
		CreateTime: service.CreateTime.Unix(),
		UpdateTime: service.UpdateTime.Unix(),
	}, nil
}
