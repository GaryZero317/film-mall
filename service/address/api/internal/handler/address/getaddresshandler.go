package address

import (
	"net/http"

	"mall/service/address/api/internal/logic/address"
	"mall/service/address/api/internal/svc"
	"mall/service/address/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetAddressHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetAddressReq
		if err := httpx.ParsePath(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := address.NewGetAddressLogic(r.Context(), svcCtx)
		resp, err := l.GetAddress(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
