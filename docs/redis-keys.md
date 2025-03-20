# Redis 键值文档

## 概述
本文档列出了系统中所有使用的 Redis 键，包括其格式、用途和过期时间。

## Redis 配置信息
- 主机地址：localhost:6379
- 密码：123456
- 类型：node

## 用户服务 (User Service)
mall:product:gorm:
### 用户相关键
| 键模式 | 数据类型 | 说明 | 示例 |
|-------|---------|------|------|
| `user:{userId}` | string | 用户基本信息缓存 | `user:10001` |
| `user:mobile:{mobile}` | string | 用户手机号索引 | `user:mobile:13800138000` |

### 管理员相关键
| 键模式 | 数据类型 | 说明 | 示例 |
|-------|---------|------|------|
| `admin:{adminId}` | string | 管理员信息缓存 | `admin:1001` |

## 购物车服务 (Cart Service)

### 购物车相关键
| 键模式 | 数据类型 | 说明 | 示例 |
|-------|---------|------|------|
| `cart:{userId}` | hash | 用户购物车商品列表 | `cart:10001` |
| `cart:count:{userId}` | string | 用户购物车商品总数 | `cart:count:10001` |

## 产品服务 (Product Service)

### 产品相关键
| 键模式 | 数据类型 | 说明 | 示例 |
|-------|---------|------|------|
| `product:{productId}` | string | 产品详情缓存 | `product:2001` |
| `product:category:{categoryId}` | list | 分类商品列表 | `product:category:101` |
| `product:hot` | zset | 热门商品排行 | `product:hot` |

## 订单服务 (Order Service)

### 订单相关键
| 键模式 | 数据类型 | 说明 | 示例 |
|-------|---------|------|------|
| `order:{orderId}` | string | 订单详情缓存 | `order:3001` |
| `order:user:{userId}` | list | 用户订单列表 | `order:user:10001` |
| `order_item:{orderItemId}` | string | 订单项详情缓存 | `order_item:4001` |

## 支付服务 (Pay Service)

### 支付相关键
| 键模式 | 数据类型 | 说明 | 示例 |
|-------|---------|------|------|
| `pay:{payId}` | string | 支付记录缓存 | `pay:5001` |
| `pay:order:{orderId}` | string | 订单支付状态 | `pay:order:3001` |

## 地址服务 (Address Service)

### 地址相关键
| 键模式 | 数据类型 | 说明 | 示例 |
|-------|---------|------|------|
| `address:{addressId}` | string | 地址详情缓存 | `address:6001` |
| `address:user:{userId}` | list | 用户地址列表 | `address:user:10001` |

## 缓存统计键
| 键模式 | 数据类型 | 说明 | 示例 |
|-------|---------|------|------|
| `cache:user:total` | string | 用户缓存统计 | `cache:user:total` |
| `cache:product:total` | string | 产品缓存统计 | `cache:product:total` |
| `cache:order:total` | string | 订单缓存统计 | `cache:order:total` |
| `cache:pay:total` | string | 支付缓存统计 | `cache:pay:total` |

## 缓存策略说明

### 过期时间
- 基础信息缓存：24小时
- 列表类缓存：1小时
- 计数器缓存：5分钟
- 状态类缓存：30分钟

### 缓存更新策略
1. 写操作时自动删除相关缓存
2. 读操作时如缓存不存在则自动重建
3. 使用 Single Flight 机制防止缓存击穿
4. 热点数据不过期，通过更新操作来保证一致性

### 缓存穿透防护
1. 对空值进行缓存，设置较短的过期时间
2. 使用布隆过滤器过滤非法请求

### 缓存击穿防护
1. 使用 go-zero 框架的 Single Flight 机制
2. 热点数据永不过期

### 缓存雪崩防护
1. 设置随机过期时间
2. 多级缓存
3. 熔断降级机制

## 注意事项
1. 键名使用冒号 `:` 作为命名空间分隔符
2. 所有键都应该设置合理的过期时间
3. 确保缓存与数据库的一致性
4. 监控缓存命中率和性能指标 