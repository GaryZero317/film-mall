package ctxdata

import (
	"context"
	"encoding/json"
	"strconv"
)

// 从上下文获取用户ID，支持多种类型
func GetUserIdFromCtx(ctx context.Context) (int64, bool) {
	// 先尝试直接获取int64类型
	if uid, ok := ctx.Value("uid").(int64); ok {
		return uid, true
	}

	// 尝试从json.Number类型转换
	if jsonNum, ok := ctx.Value("uid").(json.Number); ok {
		uidInt, err := jsonNum.Int64()
		if err == nil {
			return uidInt, true
		}
	}

	// 尝试从float64转换
	if fuid, ok := ctx.Value("uid").(float64); ok {
		return int64(fuid), true
	}

	// 尝试从string转换
	if suid, ok := ctx.Value("uid").(string); ok {
		uidInt, err := strconv.ParseInt(suid, 10, 64)
		if err == nil {
			return uidInt, true
		}
	}

	return 0, false
}
