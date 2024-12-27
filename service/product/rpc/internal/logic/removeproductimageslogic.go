package logic

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"mall/service/product/model"
	"mall/service/product/rpc/internal/svc"
	"mall/service/product/rpc/pb/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveProductImagesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRemoveProductImagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveProductImagesLogic {
	return &RemoveProductImagesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RemoveProductImagesLogic) RemoveProductImages(in *product.RemoveProductImagesRequest) (*product.RemoveProductImagesResponse, error) {
	l.Logger.Infof("开始删除商品图片: productId=%d, imageUrls=%v", in.ProductId, in.ImageUrls)

	// 先删除数据库记录
	productImage := &model.ProductImage{}
	if err := productImage.BatchDelete(l.ctx, l.svcCtx.DB, in.ProductId, in.ImageUrls); err != nil {
		l.Logger.Errorf("删除商品图片记录失败: %v", err)
		return nil, fmt.Errorf("删除商品图片记录失败: %v", err)
	}

	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		l.Logger.Errorf("获取当前目录失败: %v", err)
		// 数据库记录已删除，继续尝试删除文件
	}
	l.Logger.Infof("当前工作目录: %s", currentDir)

	// 删除物理文件
	for _, imageUrl := range in.ImageUrls {
		// 从URL中提取文件名，去掉"/uploads/"前缀
		filename := strings.TrimPrefix(imageUrl, "/uploads/")
		// 构建完整的文件路径
		filePath := filepath.Join(currentDir, "uploads", filename)
		l.Logger.Infof("处理文件: imageUrl=%s, filename=%s, filePath=%s", imageUrl, filename, filePath)

		// 检查文件是否存在
		if _, err := os.Stat(filePath); err != nil {
			if os.IsNotExist(err) {
				l.Logger.Errorf("文件不存在: %s", filePath)
			} else {
				l.Logger.Errorf("检查文件失败: %v", err)
			}
			continue
		}

		// 删除文件
		if err := os.Remove(filePath); err != nil {
			l.Logger.Errorf("删除文件失败: %v", err)
		} else {
			l.Logger.Infof("成功删除文件: %s", filePath)
		}
	}

	l.Logger.Info("删除商品图片完成")
	return &product.RemoveProductImagesResponse{}, nil
}
