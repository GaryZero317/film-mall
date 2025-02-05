package utils

import (
	"fmt"
	"time"
)

// GetTimeRange 根据时间范围类型返回开始和结束时间
func GetTimeRange(rangeType string) (start, end time.Time) {
	now := time.Now()
	end = now

	switch rangeType {
	case "7days":
		start = now.AddDate(0, 0, -7)
	case "30days":
		start = now.AddDate(0, 0, -30)
	case "month":
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	case "quarter":
		quarter := (int(now.Month())-1)/3*3 + 1
		start = time.Date(now.Year(), time.Month(quarter), 1, 0, 0, 0, 0, now.Location())
	case "year":
		start = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
	default:
		start = now.AddDate(0, 0, -7) // 默认最近7天
	}

	return start, end
}

// GetDatesBetween 获取两个时间之间的所有日期
func GetDatesBetween(start, end time.Time) []string {
	var dates []string
	for d := start; d.Before(end) || d.Equal(end); d = d.AddDate(0, 0, 1) {
		dates = append(dates, d.Format("2006-01-02"))
	}
	return dates
}

// GetHoursList 获取小时列表 (0-23)
func GetHoursList() []string {
	hours := make([]string, 24)
	for i := 0; i < 24; i++ {
		hours[i] = fmt.Sprintf("%02d:00", i)
	}
	return hours
}

// GetWeekdays 获取星期列表
func GetWeekdays() []string {
	return []string{"周一", "周二", "周三", "周四", "周五", "周六", "周日"}
}
