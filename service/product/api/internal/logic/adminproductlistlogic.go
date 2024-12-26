package logic

import (
	"context"

	"mall/service/product/api/internal/svc"
	"mall/service/product/api/internal/types"
	"mall/service/product/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminProductListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminProductListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminProductListLogic {
	return &AdminProductListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminProductListLogic) AdminProductList(req *types.AdminProductListRequest) (resp *types.AdminProductListResponse, err error) {
	// 获取分页参数
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 10
	}

	// 调用model层获取数据
	products, total, err := l.svcCtx.ProductModel.FindPageListByPage(l.ctx, page, pageSize)
	if err != nil {
		logx.Errorf("获取商品列表失败: %v", err)
		return nil, err
	}

	// 转换数据
	var list []types.Product
	for _, p := range products {
		// 获取商品图片
		var images []*model.ProductImage
		if err := l.svcCtx.DB.WithContext(l.ctx).
			Where("product_id = ?", p.Id).
			Order("is_main DESC, sort_order ASC").
			Find(&images).Error; err != nil {
			logx.Errorf("获取商品[%d]图片失败: %v", p.Id, err)
			continue
		}

		// 提取图片URL列表和主图
		var imageUrls []string
		var mainImage string
		for _, img := range images {
			imageUrls = append(imageUrls, img.ImageUrl)
			if img.IsMain {
				mainImage = img.ImageUrl
			}
		}

		// 如果没有主图，但有其他图片，则使用第一张图片作为主图
		if mainImage == "" && len(imageUrls) > 0 {
			mainImage = imageUrls[0]
		}

		// 添加到列表
		list = append(list, types.Product{
			Id:        p.Id,
			Name:      p.Name,
			Desc:      p.Desc,
			Stock:     p.Stock,
			Amount:    p.Amount,
			Status:    p.Status,
			Images:    imageUrls,
			MainImage: mainImage,
		})
	}

	logx.Infof("获取商品列表成功: total=%d, list=%d", total, len(list))
	resp = &types.AdminProductListResponse{
		Code: 0,
		Msg:  "success",
		Data: &types.AdminProductListData{
			Total: total,
			List:  list,
		},
	}
	logx.Infof("返回数据: %+v", resp)
	return resp, nil
}
