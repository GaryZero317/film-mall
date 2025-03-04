package logic

import (
	"context"

	"mall/common/ctxdata"
	"mall/service/community/api/internal/svc"
	"mall/service/community/api/internal/types"
	"mall/service/community/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminDeleteWorkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminDeleteWorkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminDeleteWorkLogic {
	return &AdminDeleteWorkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminDeleteWorkLogic) AdminDeleteWork(req *types.DeleteWorkReq) (resp *types.DeleteWorkResp, err error) {
	resp = &types.DeleteWorkResp{
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

	// 如果作品已经是删除状态，直接返回成功
	if work.Status == 2 {
		return resp, nil
	}

	// 执行删除操作（软删除，将状态设置为2）
	err = l.svcCtx.WorkModel.Delete(l.ctx, req.Id)
	if err != nil {
		resp.Code = 500
		resp.Msg = "删除作品失败: " + err.Error()
		return resp, nil
	}

	l.Logger.Infof("管理员[%d]删除了作品[%d]", adminId, req.Id)
	return resp, nil
}
