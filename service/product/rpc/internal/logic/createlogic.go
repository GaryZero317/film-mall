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
	res, err := l.svcCtx.ProductModel.Insert(l.ctx, &model.Product{
		Name:       in.Name,
		Desc:       in.Desc,
		Stock:      in.Stock,
		Amount:     in.Amount,
		Status:     in.Status,
		CategoryId: in.CategoryId,
	})
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &product.CreateResponse{
		Id: id,
	}, nil
}
