package address

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"mall/service/address/api/internal/logic/address"
	"mall/service/address/api/internal/svc"
)

func GetAddressListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := address.NewGetAddressListLogic(r.Context(), svcCtx)
		resp, err := l.GetAddressList()
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
