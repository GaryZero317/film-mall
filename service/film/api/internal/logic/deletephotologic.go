package logic

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"mall/service/film/api/internal/svc"
	"mall/service/film/api/internal/types"
	"mall/service/film/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletePhotoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeletePhotoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletePhotoLogic {
	return &DeletePhotoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeletePhotoLogic) DeletePhoto(req *types.DeleteFilmPhotoReq) (resp *types.DeleteFilmPhotoResp, err error) {
	l.Infof("尝试删除照片ID: %d", req.Id)

	// 查询照片信息
	photo, err := l.svcCtx.FilmPhotoModel.FindOne(l.ctx, req.Id)
	if err != nil {
		if err == model.ErrNotFound {
			l.Errorf("照片 %d 不存在", req.Id)
			return &types.DeleteFilmPhotoResp{
				Code: 404,
				Msg:  fmt.Sprintf("照片 %d 不存在", req.Id),
			}, nil
		}
		l.Errorf("查询照片信息失败: %v", err)
		return &types.DeleteFilmPhotoResp{
			Code: 500,
			Msg:  "查询照片信息失败: " + err.Error(),
		}, nil
	}

	// 删除文件系统中的照片文件
	if photo.Url != "" {
		// 提取文件名
		filename := ""
		if idx := strings.LastIndex(photo.Url, "/"); idx >= 0 {
			filename = photo.Url[idx+1:]
		}

		if filename != "" {
			// 构建文件路径
			filePath := filepath.Join("uploads", "film", "photos", filename)

			// 检查文件是否存在
			if _, err := os.Stat(filePath); err == nil {
				// 删除文件
				if err := os.Remove(filePath); err != nil {
					l.Errorf("删除文件失败 %s: %v", filePath, err)
					// 不中断流程，继续删除数据库记录
				} else {
					l.Infof("成功删除文件: %s", filePath)
				}
			} else {
				l.Infof("文件不存在或无法访问: %s", filePath)
				// 不中断流程
			}
		}
	}

	// 从数据库中删除照片记录
	err = l.svcCtx.FilmPhotoModel.Delete(l.ctx, req.Id)
	if err != nil {
		l.Errorf("从数据库删除照片记录失败: %v", err)
		return &types.DeleteFilmPhotoResp{
			Code: 500,
			Msg:  "删除照片失败: " + err.Error(),
		}, nil
	}

	l.Infof("成功删除照片ID: %d", req.Id)
	return &types.DeleteFilmPhotoResp{
		Code: 0,
		Msg:  "删除成功",
	}, nil
}
