package logic

import (
	"context"
	"time"

	"mall/service/user/api/internal/svc"
	"mall/service/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminFaqUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminFaqUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminFaqUpdateLogic {
	return &AdminFaqUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminFaqUpdateLogic) AdminFaqUpdate(req *types.AdminFaqUpdateRequest) (resp *types.AdminFaqUpdateResponse, err error) {
	// 参数校验
	if req.Id <= 0 || req.Question == "" || req.Answer == "" {
		return &types.AdminFaqUpdateResponse{
			Code: 400,
			Msg:  "参数错误：ID、问题和答案不能为空",
		}, nil
	}

	// 查询FAQ是否存在
	faq, err := l.svcCtx.FaqModel.FindOne(l.ctx, req.Id)
	if err != nil {
		l.Logger.Errorf("查询FAQ失败: %v", err)
		return &types.AdminFaqUpdateResponse{
			Code: 500,
			Msg:  "FAQ不存在或系统错误",
		}, nil
	}

	// 更新FAQ信息
	faq.Question = req.Question
	faq.Answer = req.Answer
	faq.Category = req.Type // 使用Type字段作为Category
	faq.Priority = req.Sort // 使用Sort字段作为Priority
	faq.UpdateTime = time.Now()

	// 保存到数据库
	err = l.svcCtx.FaqModel.Update(l.ctx, faq)
	if err != nil {
		l.Logger.Errorf("更新FAQ失败: %v", err)
		return &types.AdminFaqUpdateResponse{
			Code: 500,
			Msg:  "更新失败，系统错误",
		}, nil
	}

	return &types.AdminFaqUpdateResponse{
		Code: 0,
		Msg:  "更新FAQ成功",
	}, nil
}
