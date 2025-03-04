package logic

import (
	"context"

	"mall/service/community/api/internal/svc"
	"mall/service/community/api/internal/types"
	"mall/service/community/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateWorkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateWorkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateWorkLogic {
	return &CreateWorkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateWorkLogic) CreateWork(req *types.CreateWorkReq) (resp *types.CreateWorkResp, err error) {
	resp = &types.CreateWorkResp{
		Code: 0,
		Msg:  "创建成功",
	}

	// 获取用户ID
	uid, ok := l.ctx.Value("uid").(int64)
	if !ok || uid <= 0 {
		resp.Code = 401
		resp.Msg = "未登录或登录已过期"
		return resp, nil
	}

	// 创建作品
	work := &model.Work{
		Uid:          uid,
		Title:        req.Title,
		Description:  req.Description,
		CoverUrl:     req.CoverUrl,
		FilmType:     req.FilmType,
		FilmBrand:    req.FilmBrand,
		Camera:       req.Camera,
		Lens:         req.Lens,
		ExifInfo:     req.ExifInfo,
		ViewCount:    0,
		LikeCount:    0,
		CommentCount: 0,
		Status:       req.Status,
	}

	result, err := l.svcCtx.WorkModel.Insert(l.ctx, work)
	if err != nil {
		resp.Code = 500
		resp.Msg = "创建作品失败: " + err.Error()
		return resp, nil
	}

	insertId, err := result.LastInsertId()
	if err != nil {
		resp.Code = 500
		resp.Msg = "获取作品ID失败: " + err.Error()
		return resp, nil
	}

	resp.Data = types.CreateWorkData{
		Id: insertId,
	}

	return resp, nil
}
