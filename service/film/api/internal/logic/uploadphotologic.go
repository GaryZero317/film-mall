package logic

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"mall/service/film/api/internal/svc"
	"mall/service/film/api/internal/types"
	"mall/service/film/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadPhotoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadPhotoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadPhotoLogic {
	return &UploadPhotoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UploadPhoto 处理照片上传
func (l *UploadPhotoLogic) UploadPhoto(req *types.UploadFilmPhotoReq) (resp *types.UploadFilmPhotoResp, err error) {
	// 调试信息：打印当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		logx.Errorf("获取当前工作目录失败: %v", err)
	} else {
		logx.Infof("当前工作目录: %s", currentDir)
	}

	// 获取上传的文件
	r := l.ctx.Value("request").(*http.Request)
	if r == nil {
		return &types.UploadFilmPhotoResp{
			Code: 400,
			Msg:  "无法获取HTTP请求",
		}, nil
	}

	// 解析表单数据
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		return &types.UploadFilmPhotoResp{
			Code: 400,
			Msg:  "解析表单失败: " + err.Error(),
		}, nil
	}

	// 获取文件
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		logx.Errorf("打开上传文件失败: %v", err)
		return &types.UploadFilmPhotoResp{
			Code: 400,
			Msg:  "上传文件打开失败: " + err.Error(),
		}, nil
	}
	defer file.Close()

	// 检查文件类型
	contentType := fileHeader.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return &types.UploadFilmPhotoResp{
			Code: 400,
			Msg:  "上传的文件不是图片",
		}, nil
	}

	// 获取文件扩展名
	fileName := fileHeader.Filename
	logx.Infof("上传的原始文件名: %s, 文件类型: %s", fileName, contentType)

	// 检查胶卷订单ID
	if req.FilmOrderId <= 0 {
		return &types.UploadFilmPhotoResp{
			Code: 400,
			Msg:  "无效的胶卷订单ID",
		}, nil
	}

	// 验证胶卷订单是否存在
	_, err = l.svcCtx.FilmOrderModel.FindOne(l.ctx, req.FilmOrderId)
	if err != nil {
		if err == model.ErrNotFound {
			return &types.UploadFilmPhotoResp{
				Code: 404,
				Msg:  "胶卷订单不存在",
			}, nil
		}
		return &types.UploadFilmPhotoResp{
			Code: 500,
			Msg:  "查询胶卷订单失败: " + err.Error(),
		}, nil
	}

	// 获取文件扩展名
	ext := filepath.Ext(fileName)
	if ext == "" {
		// 根据ContentType设置默认扩展名
		switch contentType {
		case "image/jpeg", "image/jpg":
			ext = ".jpeg"
		case "image/png":
			ext = ".png"
		case "image/gif":
			ext = ".gif"
		default:
			ext = ".jpg"
		}
	}

	// 使用毫秒级时间戳+纳秒作为文件名，确保唯一性
	newFileName := fmt.Sprintf("%d_%d%s", req.FilmOrderId, time.Now().UnixNano(), ext)
	logx.Infof("生成的新文件名: %s", newFileName)

	// 设置上传文件的保存目录
	uploadDir := "D:/graduation/FilmMall/service/film/api/uploads"

	// 检查uploadDir是绝对路径
	if !filepath.IsAbs(uploadDir) {
		logx.Errorf("上传目录不是绝对路径: %s", uploadDir)
		return &types.UploadFilmPhotoResp{
			Code: 500,
			Msg:  "上传目录配置错误",
		}, nil
	}

	// 确保上传目录存在
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		logx.Errorf("创建上传目录失败: %s, 错误: %v", uploadDir, err)
		return &types.UploadFilmPhotoResp{
			Code: 500,
			Msg:  "创建上传目录失败: " + err.Error(),
		}, nil
	} else {
		logx.Infof("上传目录已确认: %s", uploadDir)
	}

	// 计算文件大小
	fileSize := fileHeader.Size
	logx.Infof("上传文件大小: %d 字节", fileSize)

	// 构建目标文件的完整路径
	saveFilePath := filepath.Join(uploadDir, newFileName)
	logx.Infof("文件将保存到: %s", saveFilePath)

	// 直接创建目标文件
	destFile, err := os.Create(saveFilePath)
	if err != nil {
		logx.Errorf("创建目标文件失败: %s, 错误: %v", saveFilePath, err)
		return &types.UploadFilmPhotoResp{
			Code: 500,
			Msg:  "创建目标文件失败: " + err.Error(),
		}, nil
	}
	defer destFile.Close()

	// 复制内容到目标文件
	written, err := io.Copy(destFile, file)
	if err != nil {
		logx.Errorf("复制文件内容失败: %v", err)
		return &types.UploadFilmPhotoResp{
			Code: 500,
			Msg:  "复制文件内容失败: " + err.Error(),
		}, nil
	}
	logx.Infof("成功写入 %d 字节到文件", written)

	// 确保所有数据写入磁盘
	if err := destFile.Sync(); err != nil {
		logx.Errorf("同步文件到磁盘失败: %v", err)
		return &types.UploadFilmPhotoResp{
			Code: 500,
			Msg:  "同步文件到磁盘失败: " + err.Error(),
		}, nil
	} else {
		logx.Infof("文件已同步到磁盘")
	}

	// 关闭文件（虽然defer会做这个，但在此之前显式关闭更好）
	if err := destFile.Close(); err != nil {
		logx.Errorf("关闭文件失败: %v", err)
		return &types.UploadFilmPhotoResp{
			Code: 500,
			Msg:  "关闭文件失败: " + err.Error(),
		}, nil
	} else {
		logx.Infof("文件已关闭")
	}

	// 检查文件是否确实已保存
	if fileInfo, err := os.Stat(saveFilePath); err != nil {
		if os.IsNotExist(err) {
			logx.Errorf("保存后文件不存在: %s", saveFilePath)
		} else {
			logx.Errorf("检查文件状态失败: %v", err)
		}
		return &types.UploadFilmPhotoResp{
			Code: 500,
			Msg:  "文件保存失败",
		}, nil
	} else {
		logx.Infof("确认文件已保存: %s, 大小: %d 字节", saveFilePath, fileInfo.Size())
		if fileInfo.Size() == 0 {
			logx.Errorf("保存的文件大小为0")
			return &types.UploadFilmPhotoResp{
				Code: 500,
				Msg:  "保存的文件大小为0",
			}, nil
		}
		if fileInfo.Size() != written {
			logx.Errorf("文件大小不匹配: 预期 %d, 实际 %d", written, fileInfo.Size())
		}
	}

	// 构建文件URL
	fileUrl := fmt.Sprintf("/uploads/%s", newFileName)
	logx.Infof("文件URL: %s", fileUrl)

	// 获取排序值
	sort := req.Sort
	if sort == 0 {
		// 如果没有指定排序值，获取当前最大排序值并加1
		photos, err := l.svcCtx.FilmPhotoModel.FindByFilmOrderId(l.ctx, req.FilmOrderId)
		if err != nil && err != model.ErrNotFound {
			logx.Errorf("获取照片列表失败: %v", err)
		}

		if len(photos) > 0 {
			maxSort := int64(0)
			for _, photo := range photos {
				if photo.Sort > maxSort {
					maxSort = photo.Sort
				}
			}
			sort = maxSort + 1
		} else {
			sort = 1 // 第一张照片，排序值为1
		}
	}

	// 将图片路径保存到数据库
	filmPhotoModel := model.FilmPhoto{
		FilmOrderId: req.FilmOrderId,
		Url:         fileUrl,
		Sort:        sort,
	}

	_, err = l.svcCtx.FilmPhotoModel.Insert(l.ctx, &filmPhotoModel)
	if err != nil {
		logx.Errorf("保存图片记录到数据库失败: %v", err)
		// 数据库操作失败时，尝试删除已上传的文件
		if removeErr := os.Remove(saveFilePath); removeErr != nil {
			logx.Errorf("删除上传的文件失败: %v", removeErr)
		}
		return &types.UploadFilmPhotoResp{
			Code: 500,
			Msg:  "保存图片记录到数据库失败: " + err.Error(),
		}, nil
	}

	logx.Infof("图片上传成功: URL=%s", fileUrl)
	return &types.UploadFilmPhotoResp{
		Code: 0,
		Msg:  "上传成功",
		Data: types.UploadFilmPhotoData{
			Url: fileUrl,
		},
	}, nil
}

// getFormFile 从请求上下文中获取上传的文件
func (l *UploadPhotoLogic) getFormFile(ctx context.Context) (multipart.File, *multipart.FileHeader, error) {
	r := ctx.Value("request").(*http.Request)
	if r == nil {
		return nil, nil, fmt.Errorf("无法获取HTTP请求")
	}

	// 解析表单数据，最大内存10MB
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		return nil, nil, err
	}

	// 获取文件
	file, header, err := r.FormFile("file")
	if err != nil {
		return nil, nil, err
	}

	return file, header, nil
}

// validateFileType 校验文件类型
func (l *UploadPhotoLogic) validateFileType(fileHeader *multipart.FileHeader) error {
	// 获取文件名和扩展名
	ext := strings.ToLower(path.Ext(fileHeader.Filename))
	// 检查文件类型
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
	}

	if !allowedExts[ext] {
		return fmt.Errorf("不支持的文件类型: %s, 仅支持JPG、PNG、GIF格式", ext)
	}

	// 检查文件大小，限制为5MB
	if fileHeader.Size > 5<<20 {
		return fmt.Errorf("文件大小超过限制，最大允许5MB")
	}

	return nil
}
