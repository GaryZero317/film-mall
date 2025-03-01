package logic

import (
	"context"
	"encoding/json"
	"time"

	"mall/common/errorx"
	"mall/service/user/api/internal/svc"
	"mall/service/user/api/internal/types"
	"mall/service/user/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type SubmitServiceRequestLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSubmitServiceRequestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubmitServiceRequestLogic {
	return &SubmitServiceRequestLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SubmitServiceRequestLogic) SubmitServiceRequest(req *types.SubmitServiceRequest) (resp *types.SubmitServiceResponse, err error) {
	// 从上下文获取用户ID
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()

	l.Logger.Infof("提交客服问题: 用户ID=%d, 标题=%s, 类型=%d", uid, req.Title, req.Type)

	// 验证请求参数
	if req.Title == "" {
		return nil, errorx.NewDefaultError("问题标题不能为空")
	}
	if req.Content == "" {
		return nil, errorx.NewDefaultError("问题内容不能为空")
	}
	if req.Type <= 0 || req.Type > 4 {
		return nil, errorx.NewDefaultError("问题类型无效")
	}

	// 创建问题记录
	now := time.Now()
	customerService := &model.CustomerService{
		UserId:     uid,
		Title:      req.Title,
		Content:    req.Content,
		Type:       req.Type,
		Status:     1,  // 1-待处理
		Reply:      "", // 空回复
		ContactWay: req.ContactWay,
		CreateTime: now,
		UpdateTime: now,
		// 注意：不设置ReplyTime字段，由模型层处理
	}

	// 插入数据库
	result, err := l.svcCtx.CustomerServiceModel.Insert(l.ctx, customerService)
	if err != nil {
		l.Logger.Errorf("插入客服问题失败: %v", err)
		return nil, errorx.NewDefaultError("提交问题失败，请稍后重试")
	}

	// 获取插入的ID
	id, err := result.LastInsertId()
	if err != nil {
		l.Logger.Errorf("获取插入ID失败: %v", err)
		return nil, errorx.NewDefaultError("提交问题失败，请稍后重试")
	}

	return &types.SubmitServiceResponse{
		Id:         id,
		Status:     1, // 1-待处理
		CreateTime: now.Unix(),
	}, nil
}
