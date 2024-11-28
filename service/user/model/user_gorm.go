package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/syncx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// GormUser represents the user model for GORM
type GormUser struct {
	ID         int64     `gorm:"primaryKey;column:id;autoIncrement"`
	Name       string    `gorm:"column:name;type:varchar(255)"`
	Gender     int64     `gorm:"column:gender"`
	Mobile     string    `gorm:"column:mobile;uniqueIndex;type:varchar(255)"`
	Password   string    `gorm:"column:password;type:varchar(255)"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime"`
	UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime"`
}

// TableName specifies the table name for GORM
func (GormUser) TableName() string {
	return "user"
}

// GormUserModel represents the GORM implementation of UserModel
type GormUserModel struct {
	db    *gorm.DB
	cache cache.Cache
}

var (
	ErrDuplicateEntry = errors.New("sql: duplicate entry")
)

// NewGormUserModel creates a new instance of GormUserModel
func NewGormUserModel(sqlDB *sql.DB, c cache.CacheConf) (*GormUserModel, error) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize GORM: %v", err)
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&GormUser{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate schema: %v", err)
	}

	cacheInstance := cache.New(c, syncx.NewSingleFlight(), cache.NewStat("user"), nil)
	return &GormUserModel{
		db:    db,
		cache: cacheInstance,
	}, nil
}

// Insert adds a new user to the database
func (m *GormUserModel) Insert(ctx context.Context, data *User) (sql.Result, error) {
	// 先检查手机号是否已存在
	var count int64
	err := m.db.WithContext(ctx).Model(&GormUser{}).Where("mobile = ?", data.Mobile).Count(&count).Error
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, ErrDuplicateEntry
	}

	gormUser := &GormUser{
		Name:     data.Name,
		Gender:   data.Gender,
		Mobile:   data.Mobile,
		Password: data.Password,
	}

	err = m.db.WithContext(ctx).Create(gormUser).Error
	if err != nil {
		return nil, err
	}

	// Update the original user with the new ID and timestamps
	data.Id = gormUser.ID
	data.CreateTime = gormUser.CreateTime
	data.UpdateTime = gormUser.UpdateTime

	return &lastInsertIDResult{id: gormUser.ID}, nil
}

// lastInsertIDResult implements sql.Result interface
type lastInsertIDResult struct {
	id int64
}

func (r *lastInsertIDResult) LastInsertId() (int64, error) {
	return r.id, nil
}

func (r *lastInsertIDResult) RowsAffected() (int64, error) {
	return 1, nil
}

// FindOne retrieves a user by ID
func (m *GormUserModel) FindOne(ctx context.Context, id int64) (*User, error) {
	var gormUser GormUser
	err := m.db.WithContext(ctx).First(&gormUser, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, sqlx.ErrNotFound
		}
		return nil, err
	}

	user := &User{
		Id:         gormUser.ID,
		Name:       gormUser.Name,
		Gender:     gormUser.Gender,
		Mobile:     gormUser.Mobile,
		Password:   gormUser.Password,
		CreateTime: gormUser.CreateTime,
		UpdateTime: gormUser.UpdateTime,
	}

	return user, nil
}

// FindOneByMobile retrieves a user by mobile number
func (m *GormUserModel) FindOneByMobile(ctx context.Context, mobile string) (*User, error) {
	var gormUser GormUser
	err := m.db.WithContext(ctx).Where("mobile = ?", mobile).First(&gormUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, sqlx.ErrNotFound
		}
		return nil, err
	}

	fmt.Printf("DEBUG FindOneByMobile - Found user: ID=%d, Name=%s, Mobile=%s\n", 
		gormUser.ID, gormUser.Name, gormUser.Mobile)

	user := &User{
		Id:         gormUser.ID,
		Name:       gormUser.Name,
		Gender:     gormUser.Gender,
		Mobile:     gormUser.Mobile,
		Password:   gormUser.Password,
		CreateTime: gormUser.CreateTime,
		UpdateTime: gormUser.UpdateTime,
	}

	fmt.Printf("DEBUG FindOneByMobile - Converted user: ID=%d, Name=%s, Mobile=%s\n", 
		user.Id, user.Name, user.Mobile)

	return user, nil
}

// Update modifies an existing user
func (m *GormUserModel) Update(ctx context.Context, data *User) error {
	gormUser := &GormUser{
		ID:       data.Id,
		Name:     data.Name,
		Gender:   data.Gender,
		Mobile:   data.Mobile,
		Password: data.Password,
	}

	err := m.db.WithContext(ctx).Save(gormUser).Error
	if err != nil {
		return err
	}

	data.UpdateTime = gormUser.UpdateTime

	// Delete cache
	m.cache.Del(fmt.Sprintf("%s%v", cacheUserIdPrefix, data.Id))
	m.cache.Del(fmt.Sprintf("%s%v", cacheUserMobilePrefix, data.Mobile))

	return nil
}

// Delete removes a user from the database
func (m *GormUserModel) Delete(ctx context.Context, id int64) error {
	// Get the user first for cache deletion
	user, err := m.FindOne(ctx, id)
	if err != nil {
		return err
	}

	err = m.db.WithContext(ctx).Delete(&GormUser{}, id).Error
	if err != nil {
		return err
	}

	// Delete cache
	m.cache.Del(fmt.Sprintf("%s%v", cacheUserIdPrefix, id))
	m.cache.Del(fmt.Sprintf("%s%v", cacheUserMobilePrefix, user.Mobile))

	return nil
}
