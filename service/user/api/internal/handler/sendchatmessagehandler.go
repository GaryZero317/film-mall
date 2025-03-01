package handler

import (
	"encoding/json"
	"net/http"

	"mall/service/user/api/internal/logic"
	"mall/service/user/api/internal/svc"
	"mall/service/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func SendChatMessageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 记录请求头信息
		logx.Infof("请求URI: %s", r.RequestURI)
		logx.Infof("Authorization: %s", r.Header.Get("Authorization"))

		var req types.SendChatMessageRequest
		if err := httpx.Parse(r, &req); err != nil {
			logx.Errorf("解析请求体失败: %v", err)
			httpx.Error(w, err)
			return
		}

		// 记录请求体信息
		reqBody, _ := json.Marshal(req)
		logx.Infof("请求体内容: %s", string(reqBody))

		l := logic.NewSendChatMessageLogic(r.Context(), svcCtx)
		resp, err := l.SendChatMessage(&req)
		if err != nil {
			logx.Errorf("发送消息失败: %v", err)
			httpx.Error(w, err)
		} else {
			logx.Info("发送消息成功")
			httpx.OkJson(w, resp)
		}
	}
}
