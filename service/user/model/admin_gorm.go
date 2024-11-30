package model

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/syncx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// GormAdmin represents the admin model for GORM
type GormAdmin struct {
	ID         int64     `gorm:"primaryKey;column:id;autoIncrement"`
	Username   string    `gorm:"column:username;uniqueIndex;type:varchar(255)"`
	Password   string    `gorm:"column:password;type:varchar(255)"`
	Level      int       `gorm:"column:level;default:1"` // 0: 超级管理员, 1: 普通管理员
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime"`
	UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime"`
}

// TableName specifies the table name for GORM
func (GormAdmin) TableName() string {
	return "admin"
}

// GormAdminModel represents the GORM implementation of AdminModel
type GormAdminModel struct {
	db    *gorm.DB
	cache cache.Cache
}

// NewGormAdminModel creates a new instance of GormAdminModel
func NewGormAdminModel(sqlDB *sql.DB, c cache.CacheConf) (*GormAdminModel, error) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&GormAdmin{})
	if err != nil {
		return nil, err
	}

	// Create cache with required parameters
	cacheInstance := cache.New(c, syncx.NewSingleFlight(), cache.NewStat("admin"), nil)

	return &GormAdminModel{
		db:    db,
		cache: cacheInstance,
	}, nil
}

// Insert adds a new admin to the database
func (m *GormAdminModel) Insert(ctx context.Context, data *GormAdmin) error {
	return m.db.WithContext(ctx).Create(data).Error
}

// FindOne retrieves an admin by ID
func (m *GormAdminModel) FindOne(ctx context.Context, id int64) (*GormAdmin, error) {
	var admin GormAdmin
	err := m.db.WithContext(ctx).First(&admin, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &admin, nil
}

// FindOneByUsername retrieves an admin by username
func (m *GormAdminModel) FindOneByUsername(ctx context.Context, username string) (*GormAdmin, error) {
	var admin GormAdmin
	err := m.db.WithContext(ctx).Where("username = ?", username).First(&admin).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &admin, nil
}

// Update modifies an existing admin
func (m *GormAdminModel) Update(ctx context.Context, data *GormAdmin) error {
	return m.db.WithContext(ctx).Save(data).Error
}

// Delete removes an admin from the database
func (m *GormAdminModel) Delete(ctx context.Context, id int64) error {
	return m.db.WithContext(ctx).Delete(&GormAdmin{}, id).Error
}
