package logic

import (
	"context"
	"time"

	"mall/service/user/model"
	"mall/service/user/rpc/internal/svc"
	"mall/service/user/rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/status"
)

type SubmitCustomerServiceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSubmitCustomerServiceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubmitCustomerServiceLogic {
	return &SubmitCustomerServiceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 客服服务
func (l *SubmitCustomerServiceLogic) SubmitCustomerService(in *user.CustomerServiceRequest) (*user.CustomerServiceResponse, error) {
	// 参数验证
	if in.Title == "" {
		return nil, status.Error(100, "问题标题不能为空")
	}
	if in.Content == "" {
		return nil, status.Error(100, "问题内容不能为空")
	}
	if in.Type <= 0 || in.Type > 4 {
		return nil, status.Error(100, "问题类型无效")
	}
	if in.UserId <= 0 {
		return nil, status.Error(100, "用户ID无效")
	}

	// 创建客服问题记录
	now := time.Now()
	customerService := &model.CustomerService{
		UserId:     in.UserId,
		Title:      in.Title,
		Content:    in.Content,
		Type:       in.Type,
		Status:     1, // 1-待处理
		ContactWay: in.ContactWay,
		CreateTime: now,
		UpdateTime: now,
	}

	// 插入数据库
	result, err := l.svcCtx.CustomerServiceModel.Insert(l.ctx, customerService)
	if err != nil {
		l.Logger.Errorf("创建客服问题失败：%v", err)
		return nil, status.Error(500, "创建客服问题失败")
	}

	// 获取插入的ID
	id, err := result.LastInsertId()
	if err != nil {
		l.Logger.Errorf("获取插入ID失败：%v", err)
		return nil, status.Error(500, "创建客服问题失败")
	}

	return &user.CustomerServiceResponse{
		Id:         id,
		Status:     1, // 1-待处理
		CreateTime: now.Unix(),
	}, nil
}
