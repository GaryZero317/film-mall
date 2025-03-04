package logic

import (
	"context"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"mall/service/community/api/internal/svc"
	"mall/service/community/api/internal/types"
	"mall/service/community/model"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type UploadWorkImageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
}

func NewUploadWorkImageLogic(ctx context.Context, svcCtx *svc.ServiceContext, r *http.Request) *UploadWorkImageLogic {
	return &UploadWorkImageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		r:      r,
	}
}

func (l *UploadWorkImageLogic) UploadWorkImage(req *types.UploadWorkImageReq) (resp *types.UploadWorkImageResp, err error) {
	resp = &types.UploadWorkImageResp{
		Code: 0,
		Msg:  "上传成功",
	}

	// 获取用户ID
	uid, ok := l.ctx.Value("uid").(int64)
	if !ok || uid <= 0 {
		resp.Code = 401
		resp.Msg = "未登录或登录已过期"
		return resp, nil
	}

	// 检查作品是否存在
	work, err := l.svcCtx.WorkModel.FindOne(l.ctx, req.WorkId)
	if err != nil {
		resp.Code = 400
		resp.Msg = "作品不存在"
		return resp, nil
	}

	// 检查作品是否属于当前用户
	if work.Uid != uid {
		resp.Code = 403
		resp.Msg = "无权为该作品上传图片"
		return resp, nil
	}

	// 处理上传文件
	file, header, err := l.r.FormFile("file")
	if err != nil {
		resp.Code = 400
		resp.Msg = "获取上传文件失败: " + err.Error()
		return resp, nil
	}
	defer file.Close()

	// 验证文件类型
	fileExt := strings.ToLower(path.Ext(header.Filename))
	if !isAllowedFileType(fileExt) {
		resp.Code = 400
		resp.Msg = "不支持的文件类型，仅支持图片格式(jpg, jpeg, png, gif)"
		return resp, nil
	}

	// 确保上传目录存在
	uploadPath := l.svcCtx.Config.FileUpload.UploadPath + "/works"
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		resp.Code = 500
		resp.Msg = "创建上传目录失败: " + err.Error()
		return resp, nil
	}

	// 生成文件名
	fileName := generateFileName(fileExt)
	filePath := path.Join(uploadPath, fileName)

	// 保存文件
	dst, err := os.Create(filePath)
	if err != nil {
		resp.Code = 500
		resp.Msg = "创建文件失败: " + err.Error()
		return resp, nil
	}
	defer dst.Close()

	if _, err = io.Copy(dst, file); err != nil {
		resp.Code = 500
		resp.Msg = "保存文件失败: " + err.Error()
		return resp, nil
	}

	// 构建文件URL
	fileUrl := l.svcCtx.Config.FileUpload.UrlPrefix + "works/" + fileName

	// 设置排序值
	sort := req.Sort
	if sort <= 0 {
		// 如果未指定排序值，查询当前最大排序值并加1
		images, err := l.svcCtx.WorkImageModel.FindByWorkId(l.ctx, req.WorkId)
		if err == nil && len(images) > 0 {
			maxSort := int64(0)
			for _, img := range images {
				if img.Sort > maxSort {
					maxSort = img.Sort
				}
			}
			sort = maxSort + 1
		} else {
			sort = 1 // 默认第一个位置
		}
	}

	// 保存到数据库
	workImage := &model.WorkImage{
		WorkId: req.WorkId,
		Url:    fileUrl,
		Sort:   sort,
	}

	result, err := l.svcCtx.WorkImageModel.Insert(l.ctx, workImage)
	if err != nil {
		// 如果数据库插入失败，删除已上传的文件
		_ = os.Remove(filePath)
		resp.Code = 500
		resp.Msg = "保存图片信息失败: " + err.Error()
		return resp, nil
	}

	// 如果是第一张图片且作品没有封面，则自动设置为封面
	if len(fileUrl) > 0 && (work.CoverUrl == "" || work.CoverUrl == "null") {
		work.CoverUrl = fileUrl
		_ = l.svcCtx.WorkModel.Update(l.ctx, work)
	}

	imageId, _ := result.LastInsertId()
	resp.Data = types.UploadWorkImageData{
		Id:  imageId,
		Url: fileUrl,
	}

	return resp, nil
}

// isAllowedFileType 判断文件类型是否允许上传
func isAllowedFileType(fileExt string) bool {
	allowedTypes := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
	}
	return allowedTypes[fileExt]
}

// generateFileName 生成唯一的文件名
func generateFileName(fileExt string) string {
	// 生成UUID
	u := uuid.New()
	// 时间戳 + UUID + 扩展名
	return strconv.FormatInt(time.Now().UnixNano(), 10) + "-" + u.String() + fileExt
}
