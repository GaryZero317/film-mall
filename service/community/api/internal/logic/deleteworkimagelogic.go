package logic

import (
	"context"
	"os"
	"path/filepath"

	"mall/common/ctxdata"
	"mall/service/community/api/internal/svc"
	"mall/service/community/api/internal/types"
	"mall/service/community/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteWorkImageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteWorkImageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteWorkImageLogic {
	return &DeleteWorkImageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteWorkImageLogic) DeleteWorkImage(req *types.DeleteWorkImageReq) (resp *types.DeleteWorkImageResp, err error) {
	resp = &types.DeleteWorkImageResp{
		Code: 0,
		Msg:  "删除成功",
	}

	// 获取用户ID
	uid, ok := ctxdata.GetUserIdFromCtx(l.ctx)
	if !ok || uid <= 0 {
		resp.Code = 401
		resp.Msg = "未登录或登录已过期"
		return resp, nil
	}

	// 查询作品图片是否存在
	workImage, err := l.svcCtx.WorkImageModel.FindOne(l.ctx, req.Id)
	if err != nil {
		if err == model.ErrNotFound {
			resp.Code = 404
			resp.Msg = "作品图片不存在"
			return resp, nil
		}
		resp.Code = 500
		resp.Msg = "查询作品图片失败: " + err.Error()
		return resp, nil
	}

	// 查询作品是否存在，以及是否是用户自己的作品
	work, err := l.svcCtx.WorkModel.FindOneByUid(l.ctx, uid, workImage.WorkId)
	if err != nil {
		if err == model.ErrNotFound {
			resp.Code = 403
			resp.Msg = "无权删除此图片"
			return resp, nil
		}
		resp.Code = 500
		resp.Msg = "查询作品失败: " + err.Error()
		return resp, nil
	}

	// 执行删除数据库记录操作
	err = l.svcCtx.WorkImageModel.Delete(l.ctx, req.Id)
	if err != nil {
		resp.Code = 500
		resp.Msg = "删除作品图片失败: " + err.Error()
		return resp, nil
	}

	// 尝试删除物理文件（可选，根据业务需求决定）
	// 从URL中提取文件名
	fileName := filepath.Base(workImage.Url)
	filePath := filepath.Join(l.svcCtx.Config.FileUpload.UploadPath, "works", fileName)

	// 尝试删除物理文件，忽略错误（文件可能不存在或已被删除）
	_ = os.Remove(filePath)

	// 如果删除的是封面图片，判断是否需要更新作品封面
	if work.CoverUrl == workImage.Url {
		// 查询该作品的其他图片
		images, err := l.svcCtx.WorkImageModel.FindByWorkId(l.ctx, work.Id)
		if err == nil && len(images) > 0 {
			// 设置第一张图片为新的封面
			work.CoverUrl = images[0].Url
			err = l.svcCtx.WorkModel.Update(l.ctx, work)
			if err != nil {
				// 更新封面失败，记录日志但不影响删除操作
				l.Logger.Errorf("更新作品[%d]封面失败: %v", work.Id, err)
			}
		} else {
			// 没有其他图片，清空封面
			work.CoverUrl = ""
			err = l.svcCtx.WorkModel.Update(l.ctx, work)
			if err != nil {
				// 更新封面失败，记录日志但不影响删除操作
				l.Logger.Errorf("清空作品[%d]封面失败: %v", work.Id, err)
			}
		}
	}

	l.Logger.Infof("用户[%d]删除了作品[%d]的图片[%d]", uid, workImage.WorkId, req.Id)
	return resp, nil
}
