package handler

import (
	"net/http"

	"mall/service/cart/api/internal/logic"
	"mall/service/cart/api/internal/svc"
	"mall/service/cart/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func AddCartHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logx.Infof("DEBUG AddCartHandler - Request headers: %+v", r.Header)
		logx.Infof("DEBUG AddCartHandler - Authorization header: %s", r.Header.Get("Authorization"))

		var req types.AddCartReq
		if err := httpx.Parse(r, &req); err != nil {
			logx.Errorf("DEBUG AddCartHandler - Parse request error: %v", err)
			httpx.Error(w, err)
			return
		}

		logx.Infof("DEBUG AddCartHandler - Request body: %+v", req)

		l := logic.NewAddCartLogic(r.Context(), svcCtx)
		resp, err := l.AddCart(&req)
		if err != nil {
			logx.Errorf("DEBUG AddCartHandler - AddCart error: %v", err)
			httpx.Error(w, err)
		} else {
			logx.Infof("DEBUG AddCartHandler - Response: %+v", resp)
			httpx.OkJson(w, resp)
		}
	}
}
