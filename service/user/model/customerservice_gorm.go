package model

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// GormCustomerService 表示GORM的客服服务模型
type GormCustomerService struct {
	ID         int64      `gorm:"primaryKey;column:id;autoIncrement"`
	UserId     int64      `gorm:"column:user_id"`
	Title      string     `gorm:"column:title;type:varchar(255)"`
	Content    string     `gorm:"column:content;type:text"`
	Type       int64      `gorm:"column:type"`
	Status     int64      `gorm:"column:status"`
	Reply      string     `gorm:"column:reply;type:text"`
	ReplyTime  *time.Time `gorm:"column:reply_time"`
	ContactWay string     `gorm:"column:contact_way;type:varchar(255)"`
	CreateTime time.Time  `gorm:"column:create_time;autoCreateTime"`
	UpdateTime time.Time  `gorm:"column:update_time;autoUpdateTime"`
}

// TableName 指定GORM的表名
func (GormCustomerService) TableName() string {
	return "customer_service"
}

// GormCustomerServiceModel 表示CustomerServiceModel的GORM实现
type GormCustomerServiceModel struct {
	db    *gorm.DB
	table string
}

// NewGormCustomerServiceModel 创建一个新的GormCustomerServiceModel实例
func NewGormCustomerServiceModel(sqlDB *sql.DB) (*GormCustomerServiceModel, error) {
	// 创建GORM数据库连接
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 返回模型实例
	return &GormCustomerServiceModel{
		db:    gormDB,
		table: "customer_service",
	}, nil
}

// Insert 插入一条新记录
func (m *GormCustomerServiceModel) Insert(ctx context.Context, data *CustomerService) (sql.Result, error) {
	// 转换为GORM模型
	gormData := &GormCustomerService{
		UserId:     data.UserId,
		Title:      data.Title,
		Content:    data.Content,
		Type:       data.Type,
		Status:     data.Status,
		Reply:      data.Reply,
		ReplyTime:  data.ReplyTime,
		ContactWay: data.ContactWay,
		CreateTime: data.CreateTime,
		UpdateTime: data.UpdateTime,
	}

	// 使用GORM插入
	result := m.db.WithContext(ctx).Create(gormData)
	if result.Error != nil {
		return nil, result.Error
	}

	// 更新ID
	data.Id = gormData.ID

	// 创建一个符合sql.Result接口的结果对象
	return &sqlResult{
		lastInsertId: gormData.ID,
		rowsAffected: result.RowsAffected,
	}, nil
}

// FindOne 根据ID查找记录
func (m *GormCustomerServiceModel) FindOne(ctx context.Context, id int64) (*CustomerService, error) {
	var gormData GormCustomerService
	result := m.db.WithContext(ctx).Where("id = ?", id).First(&gormData)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, result.Error
	}

	// 转换为普通模型
	return &CustomerService{
		Id:         gormData.ID,
		UserId:     gormData.UserId,
		Title:      gormData.Title,
		Content:    gormData.Content,
		Type:       gormData.Type,
		Status:     gormData.Status,
		Reply:      gormData.Reply,
		ReplyTime:  gormData.ReplyTime,
		ContactWay: gormData.ContactWay,
		CreateTime: gormData.CreateTime,
		UpdateTime: gormData.UpdateTime,
	}, nil
}

// Update 更新记录
func (m *GormCustomerServiceModel) Update(ctx context.Context, data *CustomerService) error {
	// 转换为GORM模型
	gormData := &GormCustomerService{
		ID:         data.Id,
		UserId:     data.UserId,
		Title:      data.Title,
		Content:    data.Content,
		Type:       data.Type,
		Status:     data.Status,
		Reply:      data.Reply,
		ReplyTime:  data.ReplyTime,
		ContactWay: data.ContactWay,
		UpdateTime: time.Now(), // 自动更新时间
	}

	// 使用GORM更新
	result := m.db.WithContext(ctx).Model(&GormCustomerService{}).Where("id = ?", data.Id).Updates(gormData)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// Delete 删除记录
func (m *GormCustomerServiceModel) Delete(ctx context.Context, id int64) error {
	result := m.db.WithContext(ctx).Delete(&GormCustomerService{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// FindByUserId 通过用户ID查找记录
func (m *GormCustomerServiceModel) FindByUserId(ctx context.Context, userId int64, page, pageSize int64, status int64) ([]*CustomerService, int64, error) {
	var gormDataList []GormCustomerService
	var count int64
	query := m.db.WithContext(ctx).Model(&GormCustomerService{}).Where("user_id = ?", userId)

	// 添加状态筛选
	if status > 0 {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	if err := query.Order("create_time DESC").
		Limit(int(pageSize)).
		Offset(int((page - 1) * pageSize)).
		Find(&gormDataList).Error; err != nil {
		return nil, 0, err
	}

	// 转换为普通模型列表
	var result []*CustomerService
	for _, gormData := range gormDataList {
		result = append(result, &CustomerService{
			Id:         gormData.ID,
			UserId:     gormData.UserId,
			Title:      gormData.Title,
			Content:    gormData.Content,
			Type:       gormData.Type,
			Status:     gormData.Status,
			Reply:      gormData.Reply,
			ReplyTime:  gormData.ReplyTime,
			ContactWay: gormData.ContactWay,
			CreateTime: gormData.CreateTime,
			UpdateTime: gormData.UpdateTime,
		})
	}

	return result, count, nil
}

// FindAll 获取所有客服问题（管理员）
func (m *GormCustomerServiceModel) FindAll(ctx context.Context, page, pageSize int64, status int64, serviceType int64) ([]*CustomerService, int64, error) {
	var gormDataList []GormCustomerService
	var count int64
	query := m.db.WithContext(ctx).Model(&GormCustomerService{})

	// 添加筛选条件
	if status > 0 {
		query = query.Where("status = ?", status)
	}
	if serviceType > 0 {
		query = query.Where("type = ?", serviceType)
	}

	// 获取总数
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	if err := query.Order("create_time DESC").
		Limit(int(pageSize)).
		Offset(int((page - 1) * pageSize)).
		Find(&gormDataList).Error; err != nil {
		return nil, 0, err
	}

	// 转换为普通模型列表
	var result []*CustomerService
	for _, gormData := range gormDataList {
		result = append(result, &CustomerService{
			Id:         gormData.ID,
			UserId:     gormData.UserId,
			Title:      gormData.Title,
			Content:    gormData.Content,
			Type:       gormData.Type,
			Status:     gormData.Status,
			Reply:      gormData.Reply,
			ReplyTime:  gormData.ReplyTime,
			ContactWay: gormData.ContactWay,
			CreateTime: gormData.CreateTime,
			UpdateTime: gormData.UpdateTime,
		})
	}

	return result, count, nil
}

// sqlResult 实现sql.Result接口
type sqlResult struct {
	lastInsertId int64
	rowsAffected int64
}

func (r *sqlResult) LastInsertId() (int64, error) {
	return r.lastInsertId, nil
}

func (r *sqlResult) RowsAffected() (int64, error) {
	return r.rowsAffected, nil
}
