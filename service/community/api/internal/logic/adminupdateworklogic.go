package logic

import (
	"context"

	"mall/common/ctxdata"
	"mall/service/community/api/internal/svc"
	"mall/service/community/api/internal/types"
	"mall/service/community/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminUpdateWorkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminUpdateWorkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUpdateWorkLogic {
	return &AdminUpdateWorkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminUpdateWorkLogic) AdminUpdateWork(req *types.UpdateWorkReq) (resp *types.UpdateWorkResp, err error) {
	resp = &types.UpdateWorkResp{
		Code: 0,
		Msg:  "更新成功",
	}

	// 获取管理员ID
	adminId, ok := ctxdata.GetUserIdFromCtx(l.ctx)
	if !ok || adminId <= 0 {
		resp.Code = 401
		resp.Msg = "管理员未登录或登录已过期"
		return resp, nil
	}

	// 查询作品是否存在
	work, err := l.svcCtx.WorkModel.FindOne(l.ctx, req.Id)
	if err != nil {
		if err == model.ErrNotFound {
			resp.Code = 404
			resp.Msg = "作品不存在"
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

	l.Logger.Infof("管理员[%d]更新了作品[%d]的信息，状态: %d", adminId, req.Id, work.Status)
	return resp, nil
}
