package logic

import (
	"context"
	"mall/service/product/api/internal/svc"
	"mall/service/product/api/internal/types"
	"mall/service/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(req *types.CreateRequest) (resp *types.CreateResponse, err error) {
	l.Logger.Infof("接收到创建商品请求: %+v", req)
	l.Logger.Infof("图片URLs: %v", req.ImageUrls)

	// 调用 RPC 创建商品
	res, err := l.svcCtx.ProductRpc.Create(l.ctx, &product.CreateRequest{
		Name:       req.Name,
		Desc:       req.Desc,
		Stock:      req.Stock,
		Amount:     req.Amount,
		Status:     req.Status,
		ImageUrls:  req.ImageUrls,
		CategoryId: req.CategoryId,
	})

	if err != nil {
		l.Logger.Errorf("创建商品失败: %v", err)
		return nil, err
	}

	l.Logger.Infof("商品创建成功，ID: %d", res.Id)
	return &types.CreateResponse{
		Code: 0,
		Msg:  "success",
		Data: types.CreateData{
			Id: res.Id,
		},
	}, nil
}
