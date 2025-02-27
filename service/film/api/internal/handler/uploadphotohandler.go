package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"mall/service/film/api/internal/svc"
	"mall/service/film/api/internal/types"
	"mall/service/film/model"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UploadPhotoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UploadFilmPhotoReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		// 解析multipart表单数据
		err := r.ParseMultipartForm(10 << 20) // 10MB
		if err != nil {
			httpx.OkJson(w, &types.UploadFilmPhotoResp{
				Code: 400,
				Msg:  "解析表单失败: " + err.Error(),
			})
			return
		}

		// 获取上传的文件
		file, header, err := r.FormFile("file")
		if err != nil {
			httpx.OkJson(w, &types.UploadFilmPhotoResp{
				Code: 400,
				Msg:  "获取上传文件失败: " + err.Error(),
			})
			return
		}
		defer file.Close()

		// 校验文件类型
		ext := strings.ToLower(filepath.Ext(header.Filename))
		// 检查文件类型
		allowedExts := map[string]bool{
			".jpg":  true,
			".jpeg": true,
			".png":  true,
			".gif":  true,
		}

		if !allowedExts[ext] {
			httpx.OkJson(w, &types.UploadFilmPhotoResp{
				Code: 400,
				Msg:  fmt.Sprintf("不支持的文件类型: %s, 仅支持JPG、PNG、GIF格式", ext),
			})
			return
		}

		// 检查文件大小，限制为5MB
		if header.Size > 5<<20 {
			httpx.OkJson(w, &types.UploadFilmPhotoResp{
				Code: 400,
				Msg:  "文件大小超过限制，最大允许5MB",
			})
			return
		}

		// 验证冲洗订单是否存在
		ctx := r.Context()
		filmOrder, err := svcCtx.FilmOrderModel.FindOne(ctx, req.FilmOrderId)
		if err != nil {
			httpx.OkJson(w, &types.UploadFilmPhotoResp{
				Code: 404,
				Msg:  fmt.Sprintf("订单 %d 不存在", req.FilmOrderId),
			})
			return
		}

		// 创建保存目录
		uploadDir := "D:/graduation/FilmMall/service/film/api/uploads"
		if err = os.MkdirAll(uploadDir, 0755); err != nil {
			httpx.OkJson(w, &types.UploadFilmPhotoResp{
				Code: 500,
				Msg:  "服务器错误: 创建上传目录失败",
			})
			return
		}

		// 获取当前工作目录，用于调试
		workDir, _ := os.Getwd()
		logx.Infof("[上传处理] 当前工作目录: %s", workDir)
		logx.Infof("[上传处理] 上传目录设置为: %s", uploadDir)

		// 记录上传文件信息
		logx.Infof("[上传处理] 上传文件: %s, 大小: %d 字节, 类型: %s",
			header.Filename, header.Size, header.Header.Get("Content-Type"))

		// 生成文件名和保存路径
		// 使用时间戳和订单号组成唯一文件名
		fileName := fmt.Sprintf("%d_%d%s", filmOrder.Id, time.Now().UnixNano(), ext)
		savePath := filepath.Join(uploadDir, fileName)
		logx.Infof("[上传处理] 生成的文件名: %s", fileName)
		logx.Infof("[上传处理] 完整保存路径: %s", savePath)

		// 保存文件
		destFile, err := os.Create(savePath)
		if err != nil {
			httpx.OkJson(w, &types.UploadFilmPhotoResp{
				Code: 500,
				Msg:  "服务器错误: 创建文件失败",
			})
			return
		}
		defer destFile.Close()

		if _, err = io.Copy(destFile, file); err != nil {
			httpx.OkJson(w, &types.UploadFilmPhotoResp{
				Code: 500,
				Msg:  "服务器错误: 保存文件失败",
			})
			return
		}

		// 强制刷新到磁盘
		if err = destFile.Sync(); err != nil {
			logx.Errorf("[上传处理] 同步文件到磁盘失败: %v", err)
		}

		// 关闭文件
		if err = destFile.Close(); err != nil {
			logx.Errorf("[上传处理] 关闭文件失败: %v", err)
		}

		// 验证文件是否已保存
		if fileInfo, err := os.Stat(savePath); err != nil {
			logx.Errorf("[上传处理] 保存后无法验证文件: %v", err)
		} else {
			logx.Infof("[上传处理] 文件已保存: %s, 大小: %d 字节", savePath, fileInfo.Size())
		}

		// 构建文件URL (使用前缀/uploads/，不带目录部分)
		fileUrl := fmt.Sprintf("/uploads/%s", fileName)
		logx.Infof("[上传处理] 构建的文件URL: %s", fileUrl)

		// 获取排序值
		sort := req.Sort
		if sort == 0 {
			// 如果没有指定排序值，获取当前最大排序值并加1
			photos, err := svcCtx.FilmPhotoModel.FindByFilmOrderId(ctx, req.FilmOrderId)
			if err != nil && err != model.ErrNotFound {
				// 获取照片列表失败，但这不是致命错误，可以继续
				sort = 1
			} else if len(photos) > 0 {
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

		// 保存到数据库
		photo := &model.FilmPhoto{
			FilmOrderId: req.FilmOrderId,
			Url:         fileUrl,
			Sort:        sort,
		}

		_, err = svcCtx.FilmPhotoModel.Insert(ctx, photo)
		if err != nil {
			// 删除已保存的文件
			_ = os.Remove(savePath)
			httpx.OkJson(w, &types.UploadFilmPhotoResp{
				Code: 500,
				Msg:  "服务器错误: 保存照片到数据库失败",
			})
			return
		}

		// 返回成功响应
		httpx.OkJson(w, &types.UploadFilmPhotoResp{
			Code: 0,
			Msg:  "上传成功",
			Data: types.UploadFilmPhotoData{
				Url: fileUrl,
			},
		})
	}
}

// 添加小写版本函数
func uploadPhotoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return UploadPhotoHandler(svcCtx)
}
