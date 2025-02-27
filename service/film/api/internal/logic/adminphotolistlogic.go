package logic

import (
	"context"

	"mall/service/film/api/internal/svc"
	"mall/service/film/api/internal/types"
	"mall/service/film/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminPhotoListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminPhotoListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminPhotoListLogic {
	return &AdminPhotoListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminPhotoListLogic) AdminPhotoList(req *types.FilmPhotoListReq) (resp *types.FilmPhotoListResp, err error) {
	l.Infof("获取订单 %d 的照片列表", req.FilmOrderId)

	// 查询该订单的所有照片
	photos, err := l.svcCtx.FilmPhotoModel.FindByFilmOrderId(l.ctx, req.FilmOrderId)
	if err != nil && err != model.ErrNotFound {
		l.Errorf("查询照片列表失败: %v", err)
		return &types.FilmPhotoListResp{
			Code: 500,
			Msg:  "查询照片列表失败: " + err.Error(),
		}, nil
	}

	// 如果没有照片或者查询结果为空
	if photos == nil || len(photos) == 0 {
		l.Infof("订单 %d 没有照片", req.FilmOrderId)
		return &types.FilmPhotoListResp{
			Code: 0,
			Msg:  "成功",
			Data: types.FilmPhotoListData{
				List: []types.FilmPhoto{},
			},
		}, nil
	}

	// 构建返回数据
	photoList := make([]types.FilmPhoto, 0, len(photos))
	for _, photo := range photos {
		photoList = append(photoList, types.FilmPhoto{
			Id:          photo.Id,
			FilmOrderId: photo.FilmOrderId,
			Url:         photo.Url,
			Sort:        photo.Sort,
			CreateTime:  photo.CreateTime.Format("2006-01-02 15:04:05"),
		})
	}

	l.Infof("成功获取订单 %d 的照片列表，共 %d 张照片", req.FilmOrderId, len(photoList))
	return &types.FilmPhotoListResp{
		Code: 0,
		Msg:  "成功",
		Data: types.FilmPhotoListData{
			List: photoList,
		},
	}, nil
}
