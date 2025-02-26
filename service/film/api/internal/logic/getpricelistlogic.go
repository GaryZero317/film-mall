package logic

import (
	"context"

	"mall/service/film/api/internal/svc"
	"mall/service/film/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPriceListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPriceListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPriceListLogic {
	return &GetPriceListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPriceListLogic) GetPriceList() (resp *types.FilmPriceListResp, err error) {
	l.Logger.Infof("获取胶片冲洗价格列表")

	// 构建预设的胶片价格列表数据
	priceList := &types.FilmPriceList{
		Types: []types.FilmType{
			{Id: 1, Name: "黑白胶片"},
			{Id: 2, Name: "彩色胶片"},
			{Id: 3, Name: "正片"},
			{Id: 4, Name: "负片"},
		},
		Brands: []types.FilmBrand{
			{Id: 1, Name: "柯达 Kodak"},
			{Id: 2, Name: "富士 Fujifilm"},
			{Id: 3, Name: "乐凯 Lucky"},
			{Id: 4, Name: "伊尔福 ILFORD"},
			{Id: 5, Name: "其他品牌"},
		},
		Sizes: []types.FilmSize{
			{Id: 1, Name: "135胶片"},
			{Id: 2, Name: "120胶片"},
		},
	}

	return &types.FilmPriceListResp{
		Code: 0,
		Msg:  "success",
		Data: *priceList,
	}, nil
}
