package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ WorkModel = (*customWorkModel)(nil)

type (
	// WorkModel 作品模型接口
	WorkModel interface {
		Insert(ctx context.Context, data *Work) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*Work, error)
		FindOneByUid(ctx context.Context, uid, id int64) (*Work, error)
		Update(ctx context.Context, data *Work) error
		Delete(ctx context.Context, id int64) error
		List(ctx context.Context, uid int64, keyword, filmType, filmBrand string, status, page, pageSize int64) ([]*Work, int64, error)
		IncrViewCount(ctx context.Context, id int64) error
		IncrLikeCount(ctx context.Context, id int64, incr int) error
		IncrCommentCount(ctx context.Context, id int64, incr int) error
	}

	// 数据库作品表映射结构体
	Work struct {
		Id           int64     `db:"id"`            // 作品ID
		Uid          int64     `db:"uid"`           // 用户ID
		Title        string    `db:"title"`         // 作品标题
		Description  string    `db:"description"`   // 作品描述
		CoverUrl     string    `db:"cover_url"`     // 封面图URL
		FilmType     string    `db:"film_type"`     // 胶片类型
		FilmBrand    string    `db:"film_brand"`    // 胶片品牌
		Camera       string    `db:"camera"`        // 相机型号
		Lens         string    `db:"lens"`          // 镜头型号
		ExifInfo     string    `db:"exif_info"`     // EXIF信息(JSON格式)
		ViewCount    int64     `db:"view_count"`    // 浏览次数
		LikeCount    int64     `db:"like_count"`    // 点赞数
		CommentCount int64     `db:"comment_count"` // 评论数
		Status       int64     `db:"status"`        // 状态:0草稿,1已发布,2已删除
		CreateTime   time.Time `db:"create_time"`   // 创建时间
		UpdateTime   time.Time `db:"update_time"`   // 更新时间
	}

	customWorkModel struct {
		conn  sqlx.SqlConn
		table string
	}
)

// NewWorkModel 创建作品模型
func NewWorkModel(conn sqlx.SqlConn) WorkModel {
	return &customWorkModel{
		conn:  conn,
		table: "works",
	}
}

// Insert 插入作品
func (m *customWorkModel) Insert(ctx context.Context, data *Work) (sql.Result, error) {
	data.CreateTime = time.Now()
	data.UpdateTime = time.Now()
	query := fmt.Sprintf("INSERT INTO %s (uid, title, description, cover_url, film_type, film_brand, camera, lens, exif_info, view_count, like_count, comment_count, status, create_time, update_time) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table)
	return m.conn.ExecCtx(ctx, query, data.Uid, data.Title, data.Description, data.CoverUrl, data.FilmType,
		data.FilmBrand, data.Camera, data.Lens, data.ExifInfo, data.ViewCount, data.LikeCount,
		data.CommentCount, data.Status, data.CreateTime, data.UpdateTime)
}

// FindOne 查询单个作品
func (m *customWorkModel) FindOne(ctx context.Context, id int64) (*Work, error) {
	var resp Work
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = ? LIMIT 1", m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// FindOneByUid 查询指定用户的作品
func (m *customWorkModel) FindOneByUid(ctx context.Context, uid, id int64) (*Work, error) {
	var resp Work
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = ? AND uid = ? LIMIT 1", m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, id, uid)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// Update 更新作品
func (m *customWorkModel) Update(ctx context.Context, data *Work) error {
	data.UpdateTime = time.Now()
	query := fmt.Sprintf("UPDATE %s SET title = ?, description = ?, cover_url = ?, film_type = ?, film_brand = ?, camera = ?, lens = ?, exif_info = ?, status = ?, update_time = ? WHERE id = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, data.Title, data.Description, data.CoverUrl, data.FilmType,
		data.FilmBrand, data.Camera, data.Lens, data.ExifInfo, data.Status, data.UpdateTime, data.Id)
	return err
}

// Delete 删除作品
func (m *customWorkModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("UPDATE %s SET status = 2, update_time = ? WHERE id = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, time.Now(), id)
	return err
}

// List 查询作品列表
func (m *customWorkModel) List(ctx context.Context, uid int64, keyword, filmType, filmBrand string, status, page, pageSize int64) ([]*Work, int64, error) {
	conditions := []string{}
	args := []interface{}{}

	if uid > 0 {
		conditions = append(conditions, "uid = ?")
		args = append(args, uid)
	}

	if keyword != "" {
		conditions = append(conditions, "(title LIKE ? OR description LIKE ?)")
		args = append(args, "%"+keyword+"%", "%"+keyword+"%")
	}

	if filmType != "" {
		conditions = append(conditions, "film_type = ?")
		args = append(args, filmType)
	}

	if filmBrand != "" {
		conditions = append(conditions, "film_brand = ?")
		args = append(args, filmBrand)
	}

	if status >= 0 { // status值有效
		conditions = append(conditions, "status = ?")
		args = append(args, status)
	} else {
		// 默认查询状态为已发布的作品
		conditions = append(conditions, "status = 1")
	}

	conditionStr := ""
	if len(conditions) > 0 {
		conditionStr = "WHERE " + strings.Join(conditions, " AND ")
	}

	// 查询总数
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s %s", m.table, conditionStr)
	var count int64
	err := m.conn.QueryRowCtx(ctx, &count, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// 查询列表
	offset := (page - 1) * pageSize
	query := fmt.Sprintf("SELECT * FROM %s %s ORDER BY create_time DESC LIMIT ?, ?",
		m.table, conditionStr)
	args = append(args, offset, pageSize)

	var resp []*Work
	err = m.conn.QueryRowsCtx(ctx, &resp, query, args...)
	if err != nil {
		return nil, 0, err
	}

	return resp, count, nil
}

// IncrViewCount 增加浏览次数
func (m *customWorkModel) IncrViewCount(ctx context.Context, id int64) error {
	query := fmt.Sprintf("UPDATE %s SET view_count = view_count + 1, update_time = ? WHERE id = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, time.Now(), id)
	return err
}

// IncrLikeCount 增加或减少点赞数量
func (m *customWorkModel) IncrLikeCount(ctx context.Context, id int64, incr int) error {
	var query string
	if incr > 0 {
		query = fmt.Sprintf("UPDATE %s SET like_count = like_count + 1, update_time = ? WHERE id = ?", m.table)
	} else {
		query = fmt.Sprintf("UPDATE %s SET like_count = GREATEST(0, like_count - 1), update_time = ? WHERE id = ?", m.table)
	}
	_, err := m.conn.ExecCtx(ctx, query, time.Now(), id)
	return err
}

// IncrCommentCount 增加或减少评论数量
func (m *customWorkModel) IncrCommentCount(ctx context.Context, id int64, incr int) error {
	var query string
	if incr > 0 {
		query = fmt.Sprintf("UPDATE %s SET comment_count = comment_count + 1, update_time = ? WHERE id = ?", m.table)
	} else {
		query = fmt.Sprintf("UPDATE %s SET comment_count = GREATEST(0, comment_count - 1), update_time = ? WHERE id = ?", m.table)
	}
	_, err := m.conn.ExecCtx(ctx, query, time.Now(), id)
	return err
}
