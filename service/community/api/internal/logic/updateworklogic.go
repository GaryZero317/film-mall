package logic

import (
	"context"

	"mall/common/ctxdata"
	"mall/service/community/api/internal/svc"
	"mall/service/community/api/internal/types"
	"mall/service/community/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateWorkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateWorkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateWorkLogic {
	return &UpdateWorkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateWorkLogic) UpdateWork(req *types.UpdateWorkReq) (resp *types.UpdateWorkResp, err error) {
	resp = &types.UpdateWorkResp{
		Code: 0,
		Msg:  "更新成功",
	}

	// 获取用户ID
	uid, ok := ctxdata.GetUserIdFromCtx(l.ctx)
	if !ok || uid <= 0 {
		resp.Code = 401
		resp.Msg = "未登录或登录已过期"
		return resp, nil
	}

	// 查询作品是否存在且是否是用户自己的作品
	work, err := l.svcCtx.WorkModel.FindOneByUid(l.ctx, uid, req.Id)
	if err != nil {
		if err == model.ErrNotFound {
			resp.Code = 404
			resp.Msg = "作品不存在或不是您的作品"
			return resp, nil
		}
		resp.Code = 500
		resp.Msg = "查询作品失败: " + err.Error()
		return resp, nil
	}

	// 更新作品信息
	if req.Title != "" {
		work.Title = req.Title
	}
	if req.Description != "" {
		work.Description = req.Description
	}
	if req.CoverUrl != "" {
		work.CoverUrl = req.CoverUrl
	}
	if req.FilmType != "" {
		work.FilmType = req.FilmType
	}
	if req.FilmBrand != "" {
		work.FilmBrand = req.FilmBrand
	}
	if req.Camera != "" {
		work.Camera = req.Camera
	}
	if req.Lens != "" {
		work.Lens = req.Lens
	}
	if req.ExifInfo != "" {
		work.ExifInfo = req.ExifInfo
	}
	if req.Status > 0 {
		work.Status = req.Status
	}

	// 执行更新
	err = l.svcCtx.WorkModel.Update(l.ctx, work)
	if err != nil {
		resp.Code = 500
		resp.Msg = "更新作品失败: " + err.Error()
		return resp, nil
	}

	l.Logger.Infof("用户[%d]更新了作品[%d]的信息", uid, req.Id)
	return resp, nil
}
