-- 商品浏览记录表
CREATE TABLE IF NOT EXISTS `product_view_log` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `user_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '用户ID',
    `product_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '商品ID',
    `view_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '浏览时间',
    `ip` varchar(50) DEFAULT NULL COMMENT '访问IP',
    `user_agent` varchar(500) DEFAULT NULL COMMENT '用户代理',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_product_id` (`product_id`),
    KEY `idx_view_time` (`view_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品浏览记录表';

-- 用户活跃度记录表
CREATE TABLE IF NOT EXISTS `user_activity_log` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `user_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '用户ID',
    `activity_type` varchar(50) NOT NULL DEFAULT '' COMMENT '活动类型：view(浏览),cart(加购),order(下单),payment(支付)',
    `activity_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '活动时间',
    `related_id` bigint(20) DEFAULT NULL COMMENT '关联ID(商品ID/订单ID等)',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_activity_time` (`activity_time`),
    KEY `idx_activity_type` (`activity_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户活跃度记录表';

-- 商品销售统计表（按天）
CREATE TABLE IF NOT EXISTS `product_sales_daily` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `product_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '商品ID',
    `category_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '类别ID',
    `sales_date` date NOT NULL COMMENT '销售日期',
    `sales_count` int(11) NOT NULL DEFAULT 0 COMMENT '销售数量',
    `sales_amount` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '销售金额',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_product_date` (`product_id`,`sales_date`),
    KEY `idx_category_id` (`category_id`),
    KEY `idx_sales_date` (`sales_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品销售统计表（按天）';

-- 商品类别销售统计表（按天）
CREATE TABLE IF NOT EXISTS `category_sales_daily` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `category_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '类别ID',
    `sales_date` date NOT NULL COMMENT '销售日期',
    `sales_count` int(11) NOT NULL DEFAULT 0 COMMENT '销售数量',
    `sales_amount` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '销售金额',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_category_date` (`category_id`,`sales_date`),
    KEY `idx_sales_date` (`sales_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品类别销售统计表（按天）'; 