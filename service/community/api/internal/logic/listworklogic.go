package logic

import (
	"context"

	"mall/service/community/api/internal/svc"
	"mall/service/community/api/internal/types"
	"mall/service/community/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListWorkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListWorkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListWorkLogic {
	return &ListWorkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListWorkLogic) ListWork(req *types.ListWorkReq) (resp *types.ListWorkResp, err error) {
	resp = &types.ListWorkResp{
		Code: 0,
		Msg:  "获取成功",
	}

	// 默认只显示已发布的作品
	status := int64(1)
	if req.Status > 0 {
		status = req.Status
	}

	// 查询作品列表
	works, total, err := l.svcCtx.WorkModel.List(l.ctx, req.Uid, req.Keyword, req.FilmType, req.FilmBrand, status, req.Page, req.PageSize)
	if err != nil {
		resp.Code = 500
		resp.Msg = "获取作品列表失败: " + err.Error()
		return resp, nil
	}

	// 如果没有作品，直接返回空列表
	if len(works) == 0 {
		resp.Data = types.ListWorkData{
			Total: 0,
			List:  []types.WorkDetail{},
		}
		return resp, nil
	}

	// 收集所有作者的用户ID
	userIds := make([]int64, 0, len(works))
	for _, work := range works {
		userIds = append(userIds, work.Uid)
	}

	// 批量查询用户信息
	userMap, err := l.svcCtx.UserModel.FindBatch(l.ctx, userIds)
	if err != nil {
		// 忽略错误，使用默认用户信息
		userMap = make(map[int64]*model.User)
	}

	// 构建作品列表数据
	workList := make([]types.WorkDetail, 0, len(works))
	for _, work := range works {
		// 获取作者信息
		authorInfo := types.UserSimple{
			Uid:      work.Uid,
			Nickname: "未知用户",
			Avatar:   "/static/images/default_avatar.png",
		}

		// 如果能找到用户信息，则使用实际用户信息
		if user, ok := userMap[work.Uid]; ok {
			authorInfo.Nickname = user.Nickname
		}

		// 构建作品详情
		workDetail := types.WorkDetail{
			Work: types.Work{
				Id:           work.Id,
				Uid:          work.Uid,
				Title:        work.Title,
				Description:  work.Description,
				CoverUrl:     work.CoverUrl,
				FilmType:     work.FilmType,
				FilmBrand:    work.FilmBrand,
				Camera:       work.Camera,
				Lens:         work.Lens,
				ExifInfo:     work.ExifInfo,
				ViewCount:    work.ViewCount,
				LikeCount:    work.LikeCount,
				CommentCount: work.CommentCount,
				Status:       work.Status,
				CreateTime:   work.CreateTime.Format("2006-01-02 15:04:05"),
				UpdateTime:   work.UpdateTime.Format("2006-01-02 15:04:05"),
			},
			Author: authorInfo,
		}

		workList = append(workList, workDetail)
	}

	resp.Data = types.ListWorkData{
		Total: total,
		List:  workList,
	}

	return resp, nil
}
