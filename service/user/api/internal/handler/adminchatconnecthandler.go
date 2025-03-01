package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"mall/service/user/api/internal/svc"
	"mall/service/user/api/internal/ws"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func AdminChatConnectHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// 从URL参数获取管理员ID
		adminIdStr := r.URL.Query().Get("adminId")
		adminId, err := strconv.ParseInt(adminIdStr, 10, 64)
		if err != nil || adminId <= 0 {
			httpx.Error(w, errors.New("无效的管理员ID"))
			return
		}

		// 获取用户ID（可选，用于直接与特定用户聊天）
		userIdStr := r.URL.Query().Get("userId")
		userId, _ := strconv.ParseInt(userIdStr, 10, 64)

		// 升级HTTP连接为WebSocket
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		// 创建客户端
		client := &ws.Client{
			ID:        adminId,
			IsAdmin:   true,
			Conn:      conn,
			Send:      make(chan []byte, 256),
			Manager:   svcCtx.WSManager,
			LastPing:  time.Now(),
			PartnerID: userId,
		}

		// 注册客户端
		client.Manager.Register <- client

		// 启动处理
		go client.WritePump()
		go client.ReadPump()
	}
}
