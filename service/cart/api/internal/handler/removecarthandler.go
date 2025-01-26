package handler

import (
	"net/http"
	"strconv"

	"mall/service/cart/api/internal/logic"
	"mall/service/cart/api/internal/svc"
	"mall/service/cart/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func RemoveCartHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RemoveCartReq
		id, err := strconv.ParseInt(r.URL.Path[len("/api/cart/"):], 10, 64)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		req.Id = id

		l := logic.NewRemoveCartLogic(r.Context(), svcCtx)
		resp, err := l.RemoveCart(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
