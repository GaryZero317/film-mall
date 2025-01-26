package handler

import (
	"net/http"

	"mall/service/cart/api/internal/logic"
	"mall/service/cart/api/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func RemoveCartHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewRemoveCartLogic(r.Context(), svcCtx)
		resp, err := l.RemoveCart(r.PathValue("id"))
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
