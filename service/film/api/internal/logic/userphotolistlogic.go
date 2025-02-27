package logic

import (
	"context"
	"fmt"

	"mall/service/film/api/internal/svc"
	"mall/service/film/api/internal/types"
	"mall/service/film/model"

	"github.com/golang-jwt/jwt"
	"github.com/zeromicro/go-zero/core/logx"
)

type UserPhotoListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserPhotoListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserPhotoListLogic {
	return &UserPhotoListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserPhotoListLogic) UserPhotoList(req *types.FilmPhotoListReq) (resp *types.FilmPhotoListResp, err error) {
	l.Infof("用户请求获取订单 %d 的照片列表", req.FilmOrderId)

	// 从JWT中获取用户ID
	userId, err := l.getUserId()
	if err != nil {
		l.Errorf("获取用户ID失败: %v", err)
		return &types.FilmPhotoListResp{
			Code: 401,
			Msg:  "未授权操作: " + err.Error(),
		}, nil
	}

	// 验证订单是否存在
	filmOrder, err := l.svcCtx.FilmOrderModel.FindOne(l.ctx, req.FilmOrderId)
	if err != nil {
		if err == model.ErrNotFound {
			l.Errorf("订单 %d 不存在", req.FilmOrderId)
			return &types.FilmPhotoListResp{
				Code: 404,
				Msg:  fmt.Sprintf("订单 %d 不存在", req.FilmOrderId),
			}, nil
		}
		l.Errorf("查询订单信息失败: %v", err)
		return &types.FilmPhotoListResp{
			Code: 500,
			Msg:  "查询订单信息失败: " + err.Error(),
		}, nil
	}

	// 验证订单是否属于当前用户
	if filmOrder.Uid != userId {
		l.Errorf("用户 %d 无权查看订单 %d 的照片", userId, req.FilmOrderId)
		return &types.FilmPhotoListResp{
			Code: 403,
			Msg:  "您无权查看该订单的照片",
		}, nil
	}

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

	l.Infof("成功获取用户 %d 订单 %d 的照片列表，共 %d 张照片", userId, req.FilmOrderId, len(photoList))
	return &types.FilmPhotoListResp{
		Code: 0,
		Msg:  "成功",
		Data: types.FilmPhotoListData{
			List: photoList,
		},
	}, nil
}

// 从JWT中获取用户ID
func (l *UserPhotoListLogic) getUserId() (int64, error) {
	token, ok := l.ctx.Value("Authorization").(*jwt.Token)
	if !ok {
		return 0, fmt.Errorf("JWT令牌不存在")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("无法解析JWT令牌内容")
	}

	// 从claims中获取用户ID
	uidFloat, ok := claims["uid"].(float64)
	if !ok {
		return 0, fmt.Errorf("JWT令牌中缺少uid字段或格式错误")
	}

	return int64(uidFloat), nil
}
