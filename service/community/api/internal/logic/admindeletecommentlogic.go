package logic

import (
	"context"

	"mall/common/ctxdata"
	"mall/service/community/api/internal/svc"
	"mall/service/community/api/internal/types"
	"mall/service/community/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminDeleteCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminDeleteCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminDeleteCommentLogic {
	return &AdminDeleteCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminDeleteCommentLogic) AdminDeleteComment(req *types.DeleteCommentReq) (resp *types.DeleteCommentResp, err error) {
	resp = &types.DeleteCommentResp{
		Code: 0,
		Msg:  "删除成功",
	}

	// 获取管理员ID
	adminId, ok := ctxdata.GetUserIdFromCtx(l.ctx)
	if !ok || adminId <= 0 {
		resp.Code = 401
		resp.Msg = "管理员未登录或登录已过期"
		return resp, nil
	}

	// 查询评论是否存在
	comment, err := l.svcCtx.CommentModel.FindOne(l.ctx, req.Id)
	if err != nil {
		if err == model.ErrNotFound {
			resp.Code = 404
			resp.Msg = "评论不存在"
			return resp, nil
		}
		resp.Code = 500
		resp.Msg = "查询评论失败: " + err.Error()
		return resp, nil
	}

	// 如果评论已经是删除状态，直接返回成功
	if comment.Status == 1 {
		return resp, nil
	}

	// 执行删除操作（软删除，将状态设置为1）
	err = l.svcCtx.CommentModel.Delete(l.ctx, req.Id)
	if err != nil {
		resp.Code = 500
		resp.Msg = "删除评论失败: " + err.Error()
		return resp, nil
	}

	// 减少作品的评论计数
	err = l.svcCtx.WorkModel.IncrCommentCount(l.ctx, comment.WorkId, -1)
	if err != nil {
		// 评论数减少失败，仅记录日志，不影响主流程
		l.Logger.Errorf("减少作品[%d]的评论计数失败: %v", comment.WorkId, err)
	}

	l.Logger.Infof("管理员[%d]删除了评论[%d]", adminId, req.Id)
	return resp, nil
}
