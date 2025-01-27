package address

import (
	"context"
	"encoding/json"
	"errors"

	"mall/service/address/api/internal/svc"
	"mall/service/address/api/internal/types"
	"mall/service/address/api/internal/utils"
	"mall/service/address/model"

	"github.com/zeromicro/go-zero/core/logx"
)

func boolToInt64(b bool) int64 {
	if b {
		return 1
	}
	return 0
}

type AddAddressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddAddressLogic {
	return &AddAddressLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddAddressLogic) AddAddress(req *types.AddAddressReq) (resp *types.AddAddressResp, err error) {
	uidAny := l.ctx.Value("uid")
	if uidAny == nil {
		return nil, errors.New("未登录")
	}

	var userId int64
	if jsonNumber, ok := uidAny.(json.Number); ok {
		userId, err = jsonNumber.Int64()
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("无效的用户ID")
	}

	// 如果设置为默认地址，需要将其他地址设置为非默认
	if req.IsDefault {
		err = l.svcCtx.AddressModel.UpdateDefaultByUserId(l.ctx, userId, false)
		if err != nil {
			return nil, err
		}
	}

	result, err := l.svcCtx.AddressModel.Insert(l.ctx, &model.Address{
		UserId:     userId,
		Name:       req.Name,
		Phone:      req.Phone,
		Province:   req.Province,
		City:       req.City,
		District:   req.District,
		DetailAddr: req.DetailAddr,
		IsDefault:  utils.BoolToInt64(req.IsDefault),
	})
	if err != nil {
		return nil, err
	}

	newId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &types.AddAddressResp{
		Id: newId,
	}, nil
}
