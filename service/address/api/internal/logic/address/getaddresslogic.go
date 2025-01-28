package address

import (
	"context"
	"encoding/json"
	"mall/service/address/api/internal/svc"
	"mall/service/address/api/internal/types"
	"mall/service/address/model"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetAddressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAddressLogic {
	return &GetAddressLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAddressLogic) GetAddress(req *types.GetAddressReq) (resp *types.GetAddressResp, err error) {
	// 从JWT中获取用户ID并转换类型
	uidNumber := l.ctx.Value("uid").(json.Number)
	userId, err := uidNumber.Int64()
	if err != nil {
		return nil, errors.New("无效的用户ID")
	}

	// 从数据库获取地址信息
	address, err := l.svcCtx.AddressModel.FindOne(l.ctx, req.Id)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, errors.New("地址不存在")
		}
		return nil, err
	}

	// 验证地址所属用户
	if address.UserId != userId {
		return nil, errors.New("无权访问该地址信息")
	}

	return &types.GetAddressResp{
		Address: types.Address{
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
		},
	}, nil
}
