package logic

import (
	"context"

	"mall/service/community/api/internal/svc"
	"mall/service/community/api/internal/types"
	"mall/service/community/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type LikeWorkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLikeWorkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikeWorkLogic {
	return &LikeWorkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LikeWorkLogic) LikeWork(req *types.LikeWorkReq) (resp *types.LikeWorkResp, err error) {
	resp = &types.LikeWorkResp{
		Code: 0,
		Msg:  "操作成功",
	}

	// 1. 检查作品是否存在
	work, err := l.svcCtx.WorkModel.FindOne(l.ctx, req.WorkId)
	if err != nil {
		if err == model.ErrNotFound {
			resp.Code = 404
			resp.Msg = "作品不存在"
			return resp, nil
		}
		resp.Code = 500
		resp.Msg = "系统错误: " + err.Error()
		return resp, nil
	}

	// 2. 检查作品状态
	if work.Status != 1 {
		resp.Code = 403
		resp.Msg = "该作品不可点赞"
		return resp, nil
	}

	// 3. 检查是否已经点赞
	isLiked, err := l.svcCtx.LikeModel.IsLiked(l.ctx, req.Uid, req.WorkId)
	if err != nil {
		resp.Code = 500
		resp.Msg = "系统错误: " + err.Error()
		return resp, nil
	}

	if !isLiked {
		// 未点赞，添加点赞记录
		like := &model.Like{
			Uid:    req.Uid,
			WorkId: req.WorkId,
		}
		_, err = l.svcCtx.LikeModel.Insert(l.ctx, like)
		if err != nil {
			resp.Code = 500
			resp.Msg = "点赞失败: " + err.Error()
			return resp, nil
		}

		// 更新作品点赞数量+1
		err = l.svcCtx.WorkModel.IncrLikeCount(l.ctx, req.WorkId, 1)
		if err != nil {
			// 点赞数更新失败，但不影响用户操作，记录日志
			l.Error("更新作品点赞数失败:", err)
		}

		resp.Data = types.LikeWorkData{
			IsLiked: true,
			Count:   work.LikeCount + 1,
		}
	} else {
		// 已点赞，取消点赞
		err = l.svcCtx.LikeModel.Delete(l.ctx, req.Uid, req.WorkId)
		if err != nil {
			resp.Code = 500
			resp.Msg = "取消点赞失败: " + err.Error()
			return resp, nil
		}

		// 更新作品点赞数量-1（防止数量为负）
		newLikeCount := work.LikeCount - 1
		if newLikeCount < 0 {
			newLikeCount = 0
		}
		err = l.svcCtx.WorkModel.IncrLikeCount(l.ctx, req.WorkId, -1)
		if err != nil {
			// 点赞数更新失败，但不影响用户操作，记录日志
			l.Error("更新作品点赞数失败:", err)
		}

		resp.Data = types.LikeWorkData{
			IsLiked: false,
			Count:   newLikeCount,
		}
	}

	return resp, nil
}
