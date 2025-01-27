package address

import (
	"net/http"

	"mall/service/address/api/internal/logic/address"
	"mall/service/address/api/internal/svc"
	"mall/service/address/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func DeleteAddressHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteAddressReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := address.NewDeleteAddressLogic(r.Context(), svcCtx)
		err := l.DeleteAddress(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
