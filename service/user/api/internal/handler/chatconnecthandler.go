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

func ChatConnectHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// 从URL参数获取用户ID
		userIdStr := r.URL.Query().Get("userId")
		userId, err := strconv.ParseInt(userIdStr, 10, 64)
		if err != nil || userId <= 0 {
			httpx.Error(w, errors.New("无效的用户ID"))
			return
		}

		// 获取管理员ID
		adminIdStr := r.URL.Query().Get("adminId")
		adminId, _ := strconv.ParseInt(adminIdStr, 10, 64)
		if adminId <= 0 {
			adminId = 1 // 默认管理员
		}

		// 升级HTTP连接为WebSocket
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		// 创建客户端
		client := &ws.Client{
			ID:        userId,
			IsAdmin:   false,
			Conn:      conn,
			Send:      make(chan []byte, 256),
			Manager:   svcCtx.WSManager,
			LastPing:  time.Now(),
			PartnerID: adminId,
		}

		// 注册客户端
		client.Manager.Register <- client

		// 启动处理
		go client.WritePump()
		go client.ReadPump()
	}
}
