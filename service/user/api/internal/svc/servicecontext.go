package svc

import (
	"context"
	"database/sql"
	"mall/service/user/api/internal/config"
	"mall/service/user/api/internal/ws"
	"mall/service/user/model"
	"mall/service/user/rpc/pb/user"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config               config.Config
	UserRpc              user.UserClient
	UserModel            model.UserModel
	AdminModel           *model.GormAdminModel
	CustomerServiceModel model.CustomerServiceModel
	FaqModel             model.FaqModel
	ChatMessageModel     model.ChatMessageModel
	WSManager            *ws.Manager
}

func NewServiceContext(c config.Config) *ServiceContext {
	logx.Info("开始初始化服务上下文")

	conn := sqlx.NewMysql(c.Mysql.DataSource)
	sqlDB, err := sql.Open("mysql", c.Mysql.DataSource)
	if err != nil {
		logx.Errorf("打开MySQL连接失败: %v", err)
		panic(err)
	}
	logx.Info("MySQL连接初始化成功")

	adminModel, err := model.NewGormAdminModel(sqlDB, c.CacheRedis)
	if err != nil {
		logx.Errorf("初始化Admin模型失败: %v", err)
		panic(err)
	}
	logx.Info("Admin模型初始化成功")

	// 创建WebSocket管理器
	logx.Info("开始创建WebSocket管理器")
	wsManager := ws.NewManager()
	logx.Info("WebSocket管理器创建成功")

	// 创建Chat消息模型
	logx.Info("开始创建Chat消息模型")
	chatMessageModel := model.NewChatMessageModel(conn)
	logx.Info("Chat消息模型创建成功")

	// 设置消息存储回调函数
	logx.Info("开始设置消息存储回调函数")
	wsManager.StoreMessage = func(message *ws.Message) (int64, error) {
		logx.Infof("尝试存储消息: userId=%d, adminId=%d, content=%s, direction=%d",
			message.UserId, message.AdminId, message.Content, message.Direction)

		// 确定发送者类型
		var senderType int64
		var senderId int64
		if message.Direction == ws.UserToAdmin {
			senderType = model.SenderTypeUser
			senderId = message.UserId
		} else {
			senderType = model.SenderTypeAdmin
			senderId = message.AdminId
		}

		// 实现消息存储逻辑 - 使用新的字段名称
		chatMessage := &model.ChatMessage{
			SessionId:  message.UserId,  // 使用用户ID作为会话ID
			SenderType: senderType,      // 根据方向确定发送者类型
			SenderId:   senderId,        // 根据方向确定发送者ID
			Content:    message.Content, // 消息内容
			ReadStatus: 0,               // 初始为未读
			CreateTime: time.Now(),      // 当前时间
		}

		logx.Infof("构建消息记录: %+v", chatMessage)
		result, err := chatMessageModel.Insert(context.Background(), chatMessage)
		if err != nil {
			logx.Errorf("消息插入数据库失败: %v", err)
			return 0, err
		}

		lastID, err := result.LastInsertId()
		if err != nil {
			logx.Errorf("获取LastInsertId失败: %v", err)
			return 0, err
		}

		logx.Infof("消息存储成功，ID: %d", lastID)
		return lastID, nil
	}
	logx.Info("消息存储回调函数设置成功")

	// 启动WebSocket管理器
	logx.Info("启动WebSocket管理器")
	go wsManager.Start()
	logx.Info("WebSocket管理器已在后台启动")

	logx.Info("服务上下文初始化完成")
	return &ServiceContext{
		Config:               c,
		UserRpc:              user.NewUserClient(zrpc.MustNewClient(c.UserRpc).Conn()),
		UserModel:            model.NewUserModel(conn, c.CacheRedis),
		AdminModel:           adminModel,
		CustomerServiceModel: model.NewCustomerServiceModel(conn),
		FaqModel:             model.NewFaqModel(conn),
		ChatMessageModel:     chatMessageModel,
		WSManager:            wsManager,
	}
}

func (s *ServiceContext) GetChatHistory(ctx context.Context, userId, adminId int64, page, pageSize int) ([]*model.ChatMessage, int64, error) {
	logx.Infof("获取聊天历史: userId=%d, adminId=%d, page=%d, pageSize=%d",
		userId, adminId, page, pageSize)
	messages, total, err := s.ChatMessageModel.FindByUserAndAdmin(ctx, userId, adminId, page, pageSize)
	if err != nil {
		logx.Errorf("获取聊天历史失败: %v", err)
		return nil, 0, err
	}
	logx.Infof("成功获取聊天历史，消息数: %d, 总记录数: %d", len(messages), total)
	return messages, total, nil
}
