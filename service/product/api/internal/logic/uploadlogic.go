package logic

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"time"

	"mall/service/product/api/internal/svc"
	"mall/service/product/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadLogic {
	return &UploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadLogic) Upload(file multipart.File, handler *multipart.FileHeader) (*types.UploadResponse, error) {
	l.Logger.Infof("开始处理文件上传: %s, 大小: %d bytes", handler.Filename, handler.Size)

	// 检查文件类型
	ext := path.Ext(handler.Filename)
	l.Logger.Infof("文件扩展名: %s", ext)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
		l.Logger.Errorf("不支持的文件类型: %s", ext)
		return nil, fmt.Errorf("不支持的文件类型: %s", ext)
	}

	// 检查文件大小（限制为2MB）
	if handler.Size > 2*1024*1024 {
		l.Logger.Errorf("文件大小超过限制: %d bytes", handler.Size)
		return nil, fmt.Errorf("文件大小超过限制")
	}

	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		l.Logger.Errorf("获取当前目录失败: %v", err)
		return nil, fmt.Errorf("获取当前目录失败: %v", err)
	}
	l.Logger.Infof("当前工作目录: %s", currentDir)

	// 创建上传目录
	uploadDir := filepath.Join(currentDir, "uploads")
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		l.Logger.Errorf("创建上传目录失败: %v", err)
		return nil, fmt.Errorf("创建上传目录失败: %v", err)
	}
	l.Logger.Infof("上传目录: %s", uploadDir)

	// 生成文件名
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filepath := filepath.Join(uploadDir, filename)
	l.Logger.Infof("目标文件路径: %s", filepath)

	// 创建目标文件
	dst, err := os.Create(filepath)
	if err != nil {
		l.Logger.Errorf("创建文件失败: %v", err)
		return nil, fmt.Errorf("创建文件失败: %v", err)
	}
	defer dst.Close()

	// 复制文件内容
	written, err := io.Copy(dst, file)
	if err != nil {
		l.Logger.Errorf("保存文件失败: %v", err)
		// 如果保存失败，尝试删除已创建的文件
		os.Remove(filepath)
		return nil, fmt.Errorf("保存文件失败: %v", err)
	}
	l.Logger.Infof("文件保存成功，写入 %d bytes", written)

	// 设置文件权限
	if err := os.Chmod(filepath, 0644); err != nil {
		l.Logger.Errorf("设置文件权限失败: %v", err)
		return nil, fmt.Errorf("设置文件权限失败: %v", err)
	}

	// 返回文件URL
	url := fmt.Sprintf("/uploads/%s", filename)
	l.Logger.Infof("文件上传成功，URL: %s", url)
	return &types.UploadResponse{
		Url: url,
	}, nil
}
