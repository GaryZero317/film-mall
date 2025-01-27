package address

import (
	"context"
	"encoding/json"
	"mall/service/address/api/internal/svc"
	"mall/service/address/api/internal/types"

	"errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAddressListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAddressListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAddressListLogic {
	return &GetAddressListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAddressListLogic) GetAddressList() (resp *types.GetAddressListResp, err error) {
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

	addresses, err := l.svcCtx.AddressModel.FindByUserId(l.ctx, userId)
	if err != nil {
		return nil, err
	}

	var list []types.Address
	for _, address := range addresses {
		list = append(list, types.Address{
			Id:         address.Id,
			UserId:     address.UserId,
			Name:       address.Name,
			Phone:      address.Phone,
			Province:   address.Province,
			City:       address.City,
			District:   address.District,
			DetailAddr: address.DetailAddr,
			IsDefault:  address.IsDefault == 1,
			CreateTime: address.CreateTime.Unix(),
			UpdateTime: address.UpdateTime.Unix(),
		})
	}

	return &types.GetAddressListResp{
		List: list,
	}, nil
}
