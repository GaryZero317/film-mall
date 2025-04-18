# film-mall
基于GoZero的在线胶卷售卖平台

## 项目简介
Film Mall 是一个专注于胶卷售卖的在线商城平台，采用微服务架构设计，提供完整的购物体验和社区互动功能。

## 技术栈
- 后端框架：Go-Zero (微服务框架)
- 数据库：MySQL 8.0
- 前端框架：Vue (管理后台)
- 微信小程序：用户端
- API文档：内置Swagger文档

## 主要功能模块

### 1. 用户系统
- 用户注册/登录
- 个人信息管理
- 收货地址管理
- 用户行为日志

### 2. 商品系统
- 商品分类管理
- 商品信息管理
- 商品图片管理
- 商品浏览记录

### 3. 购物车系统
- 购物车管理
- 商品数量调整
- 商品选择状态管理

### 4. 订单系统
- 订单创建和管理
- 支付集成
- 订单状态追踪
- 胶卷回收服务

### 5. 社区系统
- 作品展示
- 点赞功能
- 评论互动
- 用户作品管理

### 6. 客服系统
- 在线客服聊天
- 常见问题解答(FAQ)
- 问题工单处理

### 7. 数据统计
- 商品销售统计
- 类别销售分析
- 用户行为分析

## 项目结构
```
film-mall/
├── admin-frontend/    # 管理后台前端
├── app/              # 微信小程序
├── service/          # 后端微服务
│   ├── user/        # 用户服务
│   ├── product/     # 商品服务
│   ├── order/       # 订单服务
│   ├── cart/        # 购物车服务
│   ├── pay/         # 支付服务
│   ├── community/   # 社区服务
│   ├── film/        # 胶卷服务
│   ├── statistics/  # 统计服务
│   └── address/     # 地址服务
├── common/          # 公共组件
└── docs/           # 项目文档
```

## 开发环境要求
- Go 1.16+
- MySQL 8.0+
- Node.js 14+
- 微信开发者工具

## 快速开始
1. 克隆项目
```bash
git clone https://github.com/your-username/film-mall.git
```

2. 初始化数据库
```bash
mysql -u root -p < mall.sql
```

3. 启动服务
```bash
# 在各个服务目录下执行
go run user.go
```

4. 启动管理后台
```bash
cd admin-frontend
npm install
npm start
```

## API 文档
API文档使用Swagger生成，启动服务后访问：
```
http://localhost:8888/swagger/index.html
```

## 测试
项目提供了完整的Postman测试集合，位于：
```
gomall.postman_collection.json
```

## 贡献指南
欢迎提交Issue和Pull Request。

## 许可证
MIT License
