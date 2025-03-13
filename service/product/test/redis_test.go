package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

// 测试 Redis 连接配置
const (
	redisAddr = "localhost:6379"
	redisPass = "123456"
)

// Redis 键前缀
const (
	productKeyPrefix      = "mall:product:gorm:" // 产品基本信息缓存前缀
	productCategoryPrefix = "product:category:"  // 产品分类列表前缀
	productHotKey         = "product:hot"        // 热门产品键
)

func TestProductRedisKeys(t *testing.T) {
	// 创建 Redis 客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPass,
		DB:       0,
	})
	defer rdb.Close()

	ctx := context.Background()

	// 1. 测试产品基本信息缓存
	t.Run("测试产品基本信息缓存", func(t *testing.T) {
		productId := 1001
		key := fmt.Sprintf("%s%d", productKeyPrefix, productId)

		// 设置测试数据
		testProduct := map[string]interface{}{
			"id":          float64(productId),
			"name":        "测试商品",
			"desc":        "这是一个测试商品",
			"stock":       float64(100),
			"amount":      float64(9900),
			"status":      float64(1),
			"category_id": float64(1),
		}

		// 写入缓存
		err := rdb.HMSet(ctx, key, testProduct).Err()
		assert.NoError(t, err, "写入产品缓存失败")

		// 设置过期时间
		err = rdb.Expire(ctx, key, 24*time.Hour).Err()
		assert.NoError(t, err, "设置过期时间失败")

		// 读取缓存
		result, err := rdb.HGetAll(ctx, key).Result()
		assert.NoError(t, err, "读取产品缓存失败")
		assert.NotEmpty(t, result, "产品缓存为空")

		// 清理测试数据
		err = rdb.Del(ctx, key).Err()
		assert.NoError(t, err, "清理产品缓存失败")
	})

	// 2. 测试产品分类列表缓存
	t.Run("测试产品分类列表缓存", func(t *testing.T) {
		categoryId := 101
		key := fmt.Sprintf("%s%d", productCategoryPrefix, categoryId)

		// 设置测试数据
		testProducts := []string{"1001", "1002", "1003"}
		err := rdb.LPush(ctx, key, testProducts).Err()
		assert.NoError(t, err, "写入分类商品列表失败")

		// 设置过期时间
		err = rdb.Expire(ctx, key, time.Hour).Err()
		assert.NoError(t, err, "设置过期时间失败")

		// 读取缓存
		result, err := rdb.LRange(ctx, key, 0, -1).Result()
		assert.NoError(t, err, "读取分类商品列表失败")
		assert.Equal(t, len(testProducts), len(result), "分类商品列表长度不匹配")

		// 清理测试数据
		err = rdb.Del(ctx, key).Err()
		assert.NoError(t, err, "清理分类商品列表失败")
	})

	// 3. 测试热门商品排行缓存
	t.Run("测试热门商品排行缓存", func(t *testing.T) {
		// 设置测试数据
		testHotProducts := []*redis.Z{
			{Score: 100, Member: "1001"},
			{Score: 80, Member: "1002"},
			{Score: 60, Member: "1003"},
		}

		err := rdb.ZAdd(ctx, productHotKey, testHotProducts...).Err()
		assert.NoError(t, err, "写入热门商品排行失败")

		// 设置过期时间
		err = rdb.Expire(ctx, productHotKey, time.Hour).Err()
		assert.NoError(t, err, "设置过期时间失败")

		// 读取缓存
		result, err := rdb.ZRevRangeWithScores(ctx, productHotKey, 0, -1).Result()
		assert.NoError(t, err, "读取热门商品排行失败")
		assert.Equal(t, len(testHotProducts), len(result), "热门商品排行长度不匹配")

		// 清理测试数据
		err = rdb.Del(ctx, productHotKey).Err()
		assert.NoError(t, err, "清理热门商品排行失败")
	})

	// 4. 测试缓存过期时间
	t.Run("测试缓存过期时间", func(t *testing.T) {
		productId := 1001
		key := fmt.Sprintf("%s%d", productKeyPrefix, productId)

		// 设置测试数据
		testProduct := map[string]interface{}{
			"id":   productId,
			"name": "测试商品",
		}

		err := rdb.HMSet(ctx, key, testProduct).Err()
		assert.NoError(t, err, "写入产品缓存失败")

		// 设置1秒过期
		err = rdb.Expire(ctx, key, time.Second).Err()
		assert.NoError(t, err, "设置过期时间失败")

		// 等待过期
		time.Sleep(2 * time.Second)

		// 验证是否过期
		exists, err := rdb.Exists(ctx, key).Result()
		assert.NoError(t, err, "检查键是否存在失败")
		assert.Equal(t, int64(0), exists, "缓存未按预期过期")
	})
}

func TestProductCacheOperations(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPass,
		DB:       0,
	})
	defer rdb.Close()

	ctx := context.Background()

	// 测试产品缓存的完整操作流程
	t.Run("测试产品缓存完整操作流程", func(t *testing.T) {
		productId := 2001
		key := fmt.Sprintf("%s%d", productKeyPrefix, productId)

		// 1. 创建产品缓存
		testProduct := map[string]interface{}{
			"id":          float64(productId),
			"name":        "新测试商品",
			"desc":        "这是一个新的测试商品",
			"stock":       float64(50),
			"amount":      float64(5900),
			"status":      float64(1),
			"category_id": float64(2),
		}

		err := rdb.HMSet(ctx, key, testProduct).Err()
		assert.NoError(t, err, "创建产品缓存失败")

		// 2. 读取并验证缓存内容
		result, err := rdb.HGetAll(ctx, key).Result()
		assert.NoError(t, err, "读取产品缓存失败")
		assert.Equal(t, testProduct["name"], result["name"], "产品名称不匹配")

		// 3. 更新缓存
		err = rdb.HSet(ctx, key, "stock", 45).Err()
		assert.NoError(t, err, "更新产品库存失败")

		// 4. 验证更新结果
		stock, err := rdb.HGet(ctx, key, "stock").Result()
		assert.NoError(t, err, "读取更新后的库存失败")
		assert.Equal(t, "45", stock, "更新后的库存不匹配")

		// 5. 删除缓存
		err = rdb.Del(ctx, key).Err()
		assert.NoError(t, err, "删除产品缓存失败")

		// 6. 验证删除结果
		exists, err := rdb.Exists(ctx, key).Result()
		assert.NoError(t, err, "检查缓存是否删除失败")
		assert.Equal(t, int64(0), exists, "缓存未被成功删除")
	})
}
