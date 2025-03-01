package ws

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

// 消息类型和方向常量
const (
	TextMessage  = 1
	ImageMessage = 2

	UserToAdmin = 1
	AdminToUser = 2
)

// Message 聊天消息结构
type Message struct {
	ID        int64  `json:"id,omitempty"`
	UserId    int64  `json:"userId"`
	AdminId   int64  `json:"adminId"`
	Direction int    `json:"direction"`
	Content   string `json:"content"`
	Type      int    `json:"type,omitempty"`
	Timestamp int64  `json:"timestamp"`
}

// Client WebSocket客户端
type Client struct {
	ID        int64
	IsAdmin   bool
	Conn      *websocket.Conn
	Send      chan []byte
	Manager   *Manager
	mu        sync.Mutex
	LastPing  time.Time
	PartnerID int64
}

// Manager WebSocket连接管理器
type Manager struct {
	Users      map[int64]*Client
	Admins     map[int64]*Client
	Broadcast  chan *Message
	Register   chan *Client
	Unregister chan *Client
	mu         sync.Mutex
	// 回调函数
	StoreMessage       func(message *Message) (int64, error)
	GetUnreadCount     func(userId, adminId int64) (int, error)
	UpdateLastActivity func(userId, adminId int64) error
}

// NewManager 创建管理器
func NewManager() *Manager {
	logx.Info("创建新的WebSocket管理器")
	manager := &Manager{
		Users:      make(map[int64]*Client),
		Admins:     make(map[int64]*Client),
		Broadcast:  make(chan *Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
	logx.Info("WebSocket管理器创建成功")
	return manager
}

// Start 启动管理器
func (m *Manager) Start() {
	logx.Info("WebSocket管理器启动")
	for {
		select {
		case client := <-m.Register:
			logx.Infof("接收到客户端注册请求: ID=%d, IsAdmin=%v", client.ID, client.IsAdmin)
			m.registerClient(client)
		case client := <-m.Unregister:
			logx.Infof("接收到客户端注销请求: ID=%d, IsAdmin=%v", client.ID, client.IsAdmin)
			m.unregisterClient(client)
		case message := <-m.Broadcast:
			logx.Infof("接收到广播消息: ID=%d, UserId=%d, AdminId=%d, Direction=%d, Content=%s",
				message.ID, message.UserId, message.AdminId, message.Direction, message.Content)
			m.broadcastMessage(message)
		}
	}
}

// 注册客户端
func (m *Manager) registerClient(client *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if client.IsAdmin {
		m.Admins[client.ID] = client
		logx.Infof("管理员 %d 已连接", client.ID)
	} else {
		m.Users[client.ID] = client
		logx.Infof("用户 %d 已连接", client.ID)
	}

	logx.Infof("当前连接数: 用户=%d, 管理员=%d", len(m.Users), len(m.Admins))
}

// 注销客户端
func (m *Manager) unregisterClient(client *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if client.IsAdmin {
		if _, ok := m.Admins[client.ID]; ok {
			delete(m.Admins, client.ID)
			close(client.Send)
			logx.Infof("管理员 %d 已断开连接", client.ID)
		} else {
			logx.Errorf("尝试注销不存在的管理员: %d", client.ID)
		}
	} else {
		if _, ok := m.Users[client.ID]; ok {
			delete(m.Users, client.ID)
			close(client.Send)
			logx.Infof("用户 %d 已断开连接", client.ID)
		} else {
			logx.Errorf("尝试注销不存在的用户: %d", client.ID)
		}
	}

	logx.Infof("当前连接数: 用户=%d, 管理员=%d", len(m.Users), len(m.Admins))
}

// 广播消息
func (m *Manager) broadcastMessage(message *Message) {
	logx.Info("开始处理广播消息")

	// 存储消息
	if m.StoreMessage != nil {
		logx.Info("尝试存储消息")
		messageID, err := m.StoreMessage(message)
		if err != nil {
			logx.Errorf("存储消息失败: %v", err)
		} else {
			message.ID = messageID
			logx.Infof("消息存储成功, ID: %d", messageID)
		}
	} else {
		logx.Error("StoreMessage回调未设置，消息将不会被持久化")
	}

	// 更新会话最后活动时间
	if m.UpdateLastActivity != nil {
		logx.Infof("尝试更新会话活动时间: userId=%d, adminId=%d", message.UserId, message.AdminId)
		if err := m.UpdateLastActivity(message.UserId, message.AdminId); err != nil {
			logx.Errorf("更新会话活动时间失败: %v", err)
		} else {
			logx.Info("更新会话活动时间成功")
		}
	} else {
		logx.Error("UpdateLastActivity回调未设置，会话活动时间将不会更新")
	}

	// 转换为JSON
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		logx.Errorf("消息序列化失败: %v", err)
		return
	}
	logx.Infof("消息序列化成功: %s", string(jsonMessage))

	// 根据消息方向发送
	m.mu.Lock()
	defer m.mu.Unlock()

	if message.Direction == UserToAdmin {
		// 发送给管理员
		logx.Infof("消息方向: 用户->管理员, 目标管理员ID: %d", message.AdminId)
		if admin, ok := m.Admins[message.AdminId]; ok {
			logx.Infof("找到目标管理员 %d，尝试发送消息", message.AdminId)
			select {
			case admin.Send <- jsonMessage:
				logx.Infof("消息成功发送至管理员 %d", message.AdminId)
			default:
				logx.Errorf("管理员 %d 发送通道已满，关闭连接", message.AdminId)
				m.unregisterClient(admin)
			}
		} else {
			logx.Infof("管理员 %d 不在线，消息已存储", message.AdminId)
		}
	} else if message.Direction == AdminToUser {
		// 发送给用户
		logx.Infof("消息方向: 管理员->用户, 目标用户ID: %d", message.UserId)
		if user, ok := m.Users[message.UserId]; ok {
			logx.Infof("找到目标用户 %d，尝试发送消息", message.UserId)
			select {
			case user.Send <- jsonMessage:
				logx.Infof("消息成功发送至用户 %d", message.UserId)
			default:
				logx.Errorf("用户 %d 发送通道已满，关闭连接", message.UserId)
				m.unregisterClient(user)
			}
		} else {
			logx.Infof("用户 %d 不在线，消息已存储", message.UserId)
		}
	} else {
		logx.Errorf("无效的消息方向: %d", message.Direction)
	}

	logx.Info("广播消息处理完成")
}
