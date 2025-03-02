package logic

import (
	"context"
	"time"

	"mall/service/user/api/internal/svc"
	"mall/service/user/api/internal/types"
	"mall/service/user/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminFaqAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminFaqAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminFaqAddLogic {
	return &AdminFaqAddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminFaqAddLogic) AdminFaqAdd(req *types.AdminFaqAddRequest) (resp *types.AdminFaqAddResponse, err error) {
	// 参数校验
	if req.Question == "" || req.Answer == "" {
		return &types.AdminFaqAddResponse{
			Code: 400,
			Msg:  "问题和答案不能为空",
		}, nil
	}

	// 创建FAQ记录
	faq := &model.Faq{
		Question:   req.Question,
		Answer:     req.Answer,
		Category:   req.Type, // 使用Type字段作为Category
		Priority:   req.Sort, // 使用Sort字段作为Priority
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	// 插入数据库
	_, err = l.svcCtx.FaqModel.Insert(l.ctx, faq)
	if err != nil {
		l.Logger.Errorf("添加FAQ失败: %v", err)
		return &types.AdminFaqAddResponse{
			Code: 500,
			Msg:  "添加FAQ失败，系统错误",
		}, nil
	}

	return &types.AdminFaqAddResponse{
		Code: 0,
		Msg:  "添加FAQ成功",
	}, nil
}
