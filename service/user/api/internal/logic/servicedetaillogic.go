package logic

import (
	"context"
	"encoding/json"

	"mall/common/errorx"
	"mall/service/user/api/internal/svc"
	"mall/service/user/api/internal/types"
	"mall/service/user/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type ServiceDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewServiceDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ServiceDetailLogic {
	return &ServiceDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ServiceDetailLogic) ServiceDetail(req *types.ServiceDetailRequest) (resp *types.ServiceDetailResponse, err error) {
	// 从上下文获取用户ID
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()

	// 查询问题详情
	service, err := l.svcCtx.CustomerServiceModel.FindOne(l.ctx, req.Id)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, errorx.NewDefaultError("问题不存在")
		}
		l.Logger.Errorf("查询客服问题详情失败: %v", err)
		return nil, err
	}

	// 检查问题是否属于当前用户
	if service.UserId != uid {
		return nil, errorx.NewDefaultError("无权查看此问题")
	}

	// 组装响应
	replyTime := int64(0)
	if !service.ReplyTime.IsZero() {
		replyTime = service.ReplyTime.Unix()
	}

	return &types.ServiceDetailResponse{
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
