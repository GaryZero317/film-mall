package address

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"mall/service/address/api/internal/logic/address"
	"mall/service/address/api/internal/svc"
	"mall/service/address/api/internal/types"
)

func AddAddressHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddAddressReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := address.NewAddAddressLogic(r.Context(), svcCtx)
		resp, err := l.AddAddress(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
