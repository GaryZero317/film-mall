package logic

import (
	"context"

	"mall/service/product/api/internal/svc"
	"mall/service/product/api/internal/types"
	"mall/service/product/rpc/pb/product"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.UpdateRequest) (resp *types.UpdateResponse, err error) {
	// 开启事务
	err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		// 更新商品基本信息
		_, err = l.svcCtx.ProductRpc.Update(l.ctx, &product.UpdateRequest{
			Id:     req.Id,
			Name:   req.Name,
			Desc:   req.Desc,
			Stock:  req.Stock,
			Amount: req.Amount,
			Status: req.Status,
		})
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &types.UpdateResponse{}, nil
}
