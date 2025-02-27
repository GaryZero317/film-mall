-- 胶片冲洗照片表
CREATE TABLE `film_photo` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `film_order_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '冲洗订单ID',
  `url` varchar(255) NOT NULL DEFAULT '' COMMENT '照片URL',
  `sort` int unsigned NOT NULL DEFAULT '0' COMMENT '排序',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_film_order_id` (`film_order_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;
