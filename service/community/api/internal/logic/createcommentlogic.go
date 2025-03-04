package logic

import (
	"context"

	"mall/common/ctxdata"
	"mall/service/community/api/internal/svc"
	"mall/service/community/api/internal/types"
	"mall/service/community/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCommentLogic {
	return &CreateCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateCommentLogic) CreateComment(req *types.CreateCommentReq) (resp *types.CreateCommentResp, err error) {
	resp = &types.CreateCommentResp{
		Code: 0,
		Msg:  "评论成功",
	}

	// 获取用户ID
	uid, ok := ctxdata.GetUserIdFromCtx(l.ctx)
	if !ok || uid <= 0 {
		resp.Code = 401
		resp.Msg = "未登录或登录已过期"
		return resp, nil
	}

	// 校验参数
	if req.Content == "" {
		resp.Code = 400
		resp.Msg = "评论内容不能为空"
		return resp, nil
	}

	// 检查作品是否存在
	work, err := l.svcCtx.WorkModel.FindOne(l.ctx, req.WorkId)
	if err != nil {
		resp.Code = 400
		resp.Msg = "作品不存在"
		return resp, nil
	}

	// 检查作品状态是否正常
	if work.Status != 1 {
		resp.Code = 400
		resp.Msg = "作品状态异常，无法评论"
		return resp, nil
	}

	// 如果是回复评论，检查被回复的评论是否存在
	if req.ReplyId > 0 {
		parentComment, err := l.svcCtx.CommentModel.FindOne(l.ctx, req.ReplyId)
		if err != nil {
			resp.Code = 400
			resp.Msg = "回复的评论不存在"
			return resp, nil
		}

		// 确保回复的评论属于同一个作品
		if parentComment.WorkId != req.WorkId {
			resp.Code = 400
			resp.Msg = "回复的评论不属于该作品"
			return resp, nil
		}

		// 确保回复的评论状态正常
		if parentComment.Status != 0 {
			resp.Code = 400
			resp.Msg = "回复的评论已被删除"
			return resp, nil
		}
	}

	// 创建评论
	comment := &model.Comment{
		Uid:     uid,
		WorkId:  req.WorkId,
		Content: req.Content,
		ReplyId: req.ReplyId,
		Status:  0, // 正常状态
	}

	result, err := l.svcCtx.CommentModel.Insert(l.ctx, comment)
	if err != nil {
		resp.Code = 500
		resp.Msg = "评论失败: " + err.Error()
		return resp, nil
	}

	// 更新作品评论数
	_ = l.svcCtx.WorkModel.IncrCommentCount(l.ctx, req.WorkId, 1)

	// 获取评论ID
	commentId, err := result.LastInsertId()
	if err != nil {
		resp.Code = 500
		resp.Msg = "获取评论ID失败，但评论已创建成功"
		return resp, nil
	}

	resp.Data = types.CreateCommentData{
		Id: commentId,
	}

	return resp, nil
}
