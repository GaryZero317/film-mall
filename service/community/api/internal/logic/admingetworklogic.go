package logic

import (
	"context"

	"mall/common/ctxdata"
	"mall/service/community/api/internal/svc"
	"mall/service/community/api/internal/types"
	"mall/service/community/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminGetWorkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminGetWorkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminGetWorkLogic {
	return &AdminGetWorkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminGetWorkLogic) AdminGetWork(req *types.GetWorkReq) (resp *types.GetWorkResp, err error) {
	resp = &types.GetWorkResp{
		Code: 0,
		Msg:  "获取成功",
	}

	// 获取管理员ID
	adminId, ok := ctxdata.GetUserIdFromCtx(l.ctx)
	if !ok || adminId <= 0 {
		resp.Code = 401
		resp.Msg = "管理员未登录或登录已过期"
		return resp, nil
	}

	// 获取作品信息（管理员可以获取任何状态的作品）
	work, err := l.svcCtx.WorkModel.FindOne(l.ctx, req.Id)
	if err != nil {
		if err == model.ErrNotFound {
			resp.Code = 404
			resp.Msg = "作品不存在"
		} else {
			resp.Code = 500
			resp.Msg = "获取作品失败: " + err.Error()
		}
		return resp, nil
	}

	// 获取作品图片
	images, err := l.svcCtx.WorkImageModel.FindByWorkId(l.ctx, req.Id)
	if err != nil {
		// 忽略图片获取错误，可能没有图片
		images = []*model.WorkImage{}
	}

	// 转换为API类型
	workData := types.Work{
		Id:           work.Id,
		Uid:          work.Uid,
		Title:        work.Title,
		Description:  work.Description,
		CoverUrl:     work.CoverUrl,
		FilmType:     work.FilmType,
		FilmBrand:    work.FilmBrand,
		Camera:       work.Camera,
		Lens:         work.Lens,
		ExifInfo:     work.ExifInfo,
		ViewCount:    work.ViewCount,
		LikeCount:    work.LikeCount,
		CommentCount: work.CommentCount,
		Status:       work.Status,
		CreateTime:   work.CreateTime.Format("2006-01-02 15:04:05"),
		UpdateTime:   work.UpdateTime.Format("2006-01-02 15:04:05"),
	}

	// 转换图片数据
	imageList := make([]types.WorkImage, 0, len(images))
	for _, img := range images {
		imageList = append(imageList, types.WorkImage{
			Id:         img.Id,
			WorkId:     img.WorkId,
			Url:        img.Url,
			Sort:       img.Sort,
			CreateTime: img.CreateTime.Format("2006-01-02 15:04:05"),
		})
	}

	// 声明一个默认作者信息
	authorInfo := types.UserSimple{
		Uid:      work.Uid,
		Nickname: "未知用户",
		Avatar:   "",
	}

	// 获取作者信息
	author, err := l.svcCtx.UserModel.FindOne(l.ctx, work.Uid)
	if err == nil {
		// 只在获取成功时更新作者信息
		authorInfo = types.UserSimple{
			Uid:      author.Id,
			Nickname: author.Nickname,
			Avatar:   author.Avatar,
		}
	}

	// 管理员查看不需要更新点赞状态
	resp.Data = types.GetWorkData{
		Work:       workData,
		Images:     imageList,
		LikeStatus: false,
		Author:     authorInfo,
	}

	l.Logger.Infof("管理员[%d]查看了作品[%d]", adminId, req.Id)
	return resp, nil
}
