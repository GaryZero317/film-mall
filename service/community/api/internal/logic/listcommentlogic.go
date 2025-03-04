package logic

import (
	"context"

	"mall/service/community/api/internal/svc"
	"mall/service/community/api/internal/types"
	"mall/service/community/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCommentLogic {
	return &ListCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListCommentLogic) ListComment(req *types.ListCommentReq) (resp *types.ListCommentResp, err error) {
	resp = &types.ListCommentResp{
		Code: 0,
		Msg:  "获取成功",
	}

	// 获取顶级评论
	comments, total, err := l.svcCtx.CommentModel.List(l.ctx, req.WorkId, req.Page, req.PageSize)
	if err != nil {
		resp.Code = 500
		resp.Msg = "获取评论失败: " + err.Error()
		return resp, nil
	}

	// 收集所有需要获取详情的用户ID
	userIds := make([]int64, 0)
	for _, comment := range comments {
		userIds = append(userIds, comment.Uid)
	}

	// 查询用户信息
	userMap, err := l.svcCtx.UserModel.FindBatch(l.ctx, userIds)
	if err != nil {
		// 忽略错误，使用默认用户信息
		userMap = make(map[int64]*model.User)
	}

	// 构建评论列表数据
	commentList := make([]types.CommentDetail, 0, len(comments))
	for _, comment := range comments {
		// 获取评论者信息
		userInfo := types.UserSimple{
			Uid:      comment.Uid,
			Nickname: "未知用户",
			Avatar:   "/static/images/default_avatar.png",
		}

		// 如果能找到用户信息，则使用实际用户信息
		if user, ok := userMap[comment.Uid]; ok {
			userInfo.Nickname = user.Nickname
		}

		// 获取回复列表
		replies, err := l.getCommentReplies(comment.Id)
		if err != nil {
			// 忽略错误，使用空回复列表
			replies = []types.CommentDetail{}
		}

		// 构建评论详情
		commentDetail := types.CommentDetail{
			Comment: types.Comment{
				Id:         comment.Id,
				Uid:        comment.Uid,
				WorkId:     comment.WorkId,
				Content:    comment.Content,
				ReplyId:    comment.ReplyId,
				Status:     comment.Status,
				CreateTime: comment.CreateTime.Format("2006-01-02 15:04:05"),
			},
			User:    userInfo,
			Replies: replies,
		}

		commentList = append(commentList, commentDetail)
	}

	resp.Data = types.ListCommentData{
		Total: total,
		List:  commentList,
	}

	return resp, nil
}

// getCommentReplies 获取评论的回复列表
func (l *ListCommentLogic) getCommentReplies(commentId int64) ([]types.CommentDetail, error) {
	// 获取回复
	replies, err := l.svcCtx.CommentModel.FindReplies(l.ctx, commentId)
	if err != nil {
		return nil, err
	}

	if len(replies) == 0 {
		return []types.CommentDetail{}, nil
	}

	// 收集所有回复的用户ID
	userIds := make([]int64, 0, len(replies))
	for _, reply := range replies {
		userIds = append(userIds, reply.Uid)
	}

	// 查询用户信息
	userMap, err := l.svcCtx.UserModel.FindBatch(l.ctx, userIds)
	if err != nil {
		// 忽略错误，使用默认用户信息
		userMap = make(map[int64]*model.User)
	}

	// 构建回复列表
	replyList := make([]types.CommentDetail, 0, len(replies))
	for _, reply := range replies {
		// 获取回复者信息
		userInfo := types.UserSimple{
			Uid:      reply.Uid,
			Nickname: "未知用户",
			Avatar:   "/static/images/default_avatar.png",
		}

		// 如果能找到用户信息，则使用实际用户信息
		if user, ok := userMap[reply.Uid]; ok {
			userInfo.Nickname = user.Nickname
		}

		// 构建回复详情
		replyDetail := types.CommentDetail{
			Comment: types.Comment{
				Id:         reply.Id,
				Uid:        reply.Uid,
				WorkId:     reply.WorkId,
				Content:    reply.Content,
				ReplyId:    reply.ReplyId,
				Status:     reply.Status,
				CreateTime: reply.CreateTime.Format("2006-01-02 15:04:05"),
			},
			User: userInfo,
			// 不再递归获取回复的回复
		}

		replyList = append(replyList, replyDetail)
	}

	return replyList, nil
}
