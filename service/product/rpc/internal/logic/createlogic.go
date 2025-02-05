package logic

import (
	"context"
	"mall/service/product/model"
	"mall/service/product/rpc/internal/svc"
	"mall/service/product/rpc/pb/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateLogic) Create(in *product.CreateRequest) (*product.CreateResponse, error) {
	l.Logger.Infof("接收到创建商品请求: %+v", in)
	l.Logger.Infof("图片URLs: %v", in.ImageUrls)

	// 创建商品
	newProduct := &model.Product{
		Name:       in.Name,
		Desc:       in.Desc,
		Stock:      in.Stock,
		Amount:     in.Amount,
		Status:     in.Status,
		CategoryId: in.CategoryId,
		Images:     in.ImageUrls, // 设置图片列表
		MainImage:  "",           // 主图将在保存图片时设置
	}

	// 如果有图片，设置第一张为主图
	if len(in.ImageUrls) > 0 {
		newProduct.MainImage = in.ImageUrls[0]
		l.Logger.Infof("设置主图: %s", newProduct.MainImage)
	}

	l.Logger.Infof("准备保存商品数据: %+v", newProduct)
	res, err := l.svcCtx.ProductModel.Insert(l.ctx, newProduct)
	if err != nil {
		l.Logger.Errorf("创建商品失败: %v", err)
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		l.Logger.Errorf("获取插入ID失败: %v", err)
		return nil, err
	}

	l.Logger.Infof("商品创建成功，ID: %d", id)
	return &product.CreateResponse{
		Id: id,
	}, nil
}
