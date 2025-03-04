package logic

import (
	"context"

	"mall/common/ctxdata"
	"mall/service/community/api/internal/svc"
	"mall/service/community/api/internal/types"
	"mall/service/community/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserWorksLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserWorksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserWorksLogic {
	return &GetUserWorksLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserWorksLogic) GetUserWorks(req *types.ListWorkReq) (resp *types.ListWorkResp, err error) {
	resp = &types.ListWorkResp{
		Code: 0,
		Msg:  "获取成功",
	}

	// 获取用户ID
	uid, ok := ctxdata.GetUserIdFromCtx(l.ctx)
	if !ok || uid <= 0 {
		resp.Code = 401
		resp.Msg = "未登录或登录已过期"
		return resp, nil
	}

	// 强制设置查询条件为当前用户
	req.Uid = uid

	// 查询作品列表，包括草稿和已发布的作品
	works, total, err := l.svcCtx.WorkModel.List(l.ctx, req.Uid, req.Keyword, req.FilmType, req.FilmBrand, req.Status, req.Page, req.PageSize)
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

	// 获取用户信息
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, uid)
	if err != nil {
		// 用户信息获取失败，使用默认值
		user = &model.User{
			Id:       uid,
			Nickname: "未知用户",
		}
	}

	// 构建作品列表数据
	workList := make([]types.WorkDetail, 0, len(works))
	for _, work := range works {
		// 构建作者信息
		authorInfo := types.UserSimple{
			Uid:      user.Id,
			Nickname: user.Nickname,
			Avatar:   "/static/images/default_avatar.png",
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
