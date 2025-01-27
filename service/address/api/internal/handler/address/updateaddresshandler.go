package address

import (
	"net/http"

	"mall/service/address/api/internal/logic/address"
	"mall/service/address/api/internal/svc"
	"mall/service/address/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func UpdateAddressHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateAddressReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := address.NewUpdateAddressLogic(r.Context(), svcCtx)
		resp, err := l.UpdateAddress(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
