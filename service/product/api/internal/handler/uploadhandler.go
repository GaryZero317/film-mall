package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"mall/service/product/api/internal/svc"
	"mall/service/product/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func UploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 解析multipart表单
		err := r.ParseMultipartForm(10 << 20) // 10MB
		if err != nil {
			httpx.Error(w, fmt.Errorf("解析表单失败: %v", err))
			return
		}

		// 获取上传的文件
		file, header, err := r.FormFile("file")
		if err != nil {
			httpx.Error(w, fmt.Errorf("获取文件失败: %v", err))
			return
		}
		defer file.Close()

		// 检查文件类型
		ext := strings.ToLower(filepath.Ext(header.Filename))
		validTypes := []string{".jpg", ".jpeg", ".png", ".gif"}
		isValidType := false
		for _, validType := range validTypes {
			if ext == validType {
				isValidType = true
				break
			}
		}
		if !isValidType {
			httpx.Error(w, fmt.Errorf("不支持的文件类型"))
			return
		}

		// 创建上传目录
		uploadDir := "uploads"
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			httpx.Error(w, fmt.Errorf("创建目录失败: %v", err))
			return
		}

		// 生成文件名
		timestamp := time.Now().UnixNano()
		filename := fmt.Sprintf("%d%s", timestamp, ext)
		filePath := filepath.Join(uploadDir, filename)

		// 创建目标文件
		dst, err := os.Create(filePath)
		if err != nil {
			httpx.Error(w, fmt.Errorf("创建文件失败: %v", err))
			return
		}
		defer dst.Close()

		// 复制文件内容
		if _, err := io.Copy(dst, file); err != nil {
			httpx.Error(w, fmt.Errorf("保存文件失败: %v", err))
			return
		}

		// 返回文件URL
		fileURL := fmt.Sprintf("%s%s", svcCtx.Config.FileUpload.UrlPrefix, filename)
		resp := &types.UploadResponse{
			Url: fileURL,
		}

		httpx.OkJson(w, map[string]interface{}{
			"code": 0,
			"msg":  "success",
			"data": resp,
		})
	}
}
