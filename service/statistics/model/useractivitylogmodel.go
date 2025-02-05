package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	userActivityLogFieldNames          = builder.RawFieldNames(&UserActivityLog{})
	userActivityLogRows                = strings.Join(userActivityLogFieldNames, ",")
	userActivityLogRowsExpectAutoSet   = strings.Join(stringx.Remove(userActivityLogFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	userActivityLogRowsWithPlaceHolder = strings.Join(stringx.Remove(userActivityLogFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"
)

type (
	// UserActivityLogModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserActivityLogModel.
	UserActivityLogModel interface {
		Insert(ctx context.Context, data *UserActivityLog) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*UserActivityLog, error)
		Update(ctx context.Context, data *UserActivityLog) error
		Delete(ctx context.Context, id int64) error
		FindUserBehaviors(ctx context.Context, start, end time.Time) ([]*UserActivityLog, error)
		FindUserActivities(ctx context.Context, start, end time.Time) ([]*UserActivityLog, error)
	}

	defaultUserActivityLogModel struct {
		conn  sqlx.SqlConn
		table string
	}

	UserActivityLog struct {
		Id           int64         `db:"id"`
		UserId       int64         `db:"user_id"`       // 用户ID
		ActivityType string        `db:"activity_type"` // 活动类型：view(浏览),cart(加购),order(下单),payment(支付)
		ActivityTime time.Time     `db:"activity_time"` // 活动时间
		RelatedId    sql.NullInt64 `db:"related_id"`    // 关联ID(商品ID/订单ID等)
	}
)

func NewUserActivityLogModel(conn sqlx.SqlConn) UserActivityLogModel {
	return &defaultUserActivityLogModel{
		conn:  conn,
		table: "`user_activity_log`",
	}
}

func (m *defaultUserActivityLogModel) Insert(ctx context.Context, data *UserActivityLog) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?)", m.table, userActivityLogRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.UserId, data.ActivityType, data.ActivityTime, data.RelatedId)
	return ret, err
}

func (m *defaultUserActivityLogModel) FindOne(ctx context.Context, id int64) (*UserActivityLog, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", userActivityLogRows, m.table)
	var resp UserActivityLog
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

func (m *defaultUserActivityLogModel) Update(ctx context.Context, data *UserActivityLog) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, userActivityLogRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.UserId, data.ActivityType, data.ActivityTime, data.RelatedId, data.Id)
	return err
}

func (m *defaultUserActivityLogModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

// FindUserBehaviors 查询指定时间范围内的用户行为数据
func (m *defaultUserActivityLogModel) FindUserBehaviors(ctx context.Context, start, end time.Time) ([]*UserActivityLog, error) {
	query := fmt.Sprintf("select %s from %s where activity_time between ? and ?", userActivityLogRows, m.table)
	var resp []*UserActivityLog
	err := m.conn.QueryRowsCtx(ctx, &resp, query, start.Format("2006-01-02 15:04:05"), end.Format("2006-01-02 15:04:05"))
	return resp, err
}

// FindUserActivities 查询指定时间范围内的用户活跃度数据
func (m *defaultUserActivityLogModel) FindUserActivities(ctx context.Context, start, end time.Time) ([]*UserActivityLog, error) {
	query := fmt.Sprintf("select %s from %s where activity_time between ? and ?", userActivityLogRows, m.table)
	var resp []*UserActivityLog
	err := m.conn.QueryRowsCtx(ctx, &resp, query, start.Format("2006-01-02 15:04:05"), end.Format("2006-01-02 15:04:05"))
	return resp, err
}
