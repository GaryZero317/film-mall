-- 作品表
CREATE TABLE IF NOT EXISTS `works` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '作品ID',
  `uid` bigint(20) NOT NULL COMMENT '用户ID',
  `title` varchar(255) NOT NULL COMMENT '作品标题',
  `description` text COMMENT '作品描述',
  `cover_url` varchar(255) DEFAULT NULL COMMENT '封面图URL',
  `film_type` varchar(50) DEFAULT NULL COMMENT '胶片类型',
  `film_brand` varchar(50) DEFAULT NULL COMMENT '胶片品牌',
  `camera` varchar(100) DEFAULT NULL COMMENT '相机型号',
  `lens` varchar(100) DEFAULT NULL COMMENT '镜头型号',
  `exif_info` text COMMENT 'EXIF信息(JSON格式)',
  `view_count` bigint(20) NOT NULL DEFAULT '0' COMMENT '浏览次数',
  `like_count` bigint(20) NOT NULL DEFAULT '0' COMMENT '点赞数',
  `comment_count` bigint(20) NOT NULL DEFAULT '0' COMMENT '评论数',
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '状态:0草稿,1已发布,2已删除',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `update_time` datetime NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_uid` (`uid`),
  KEY `idx_status` (`status`),
  KEY `idx_film_type` (`film_type`),
  KEY `idx_film_brand` (`film_brand`),
  KEY `idx_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='作品表';

-- 作品图片表
CREATE TABLE IF NOT EXISTS `work_images` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '图片ID',
  `work_id` bigint(20) NOT NULL COMMENT '作品ID',
  `url` varchar(255) NOT NULL COMMENT '图片URL',
  `sort` int(11) NOT NULL DEFAULT '0' COMMENT '排序',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_work_id` (`work_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='作品图片表';

-- 点赞表
CREATE TABLE IF NOT EXISTS `likes` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '点赞ID',
  `uid` bigint(20) NOT NULL COMMENT '用户ID',
  `work_id` bigint(20) NOT NULL COMMENT '作品ID',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_uid_work_id` (`uid`,`work_id`),
  KEY `idx_work_id` (`work_id`),
  KEY `idx_uid` (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='点赞表';

-- 评论表
CREATE TABLE IF NOT EXISTS `comments` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '评论ID',
  `uid` bigint(20) NOT NULL COMMENT '用户ID',
  `work_id` bigint(20) NOT NULL COMMENT '作品ID',
  `content` text NOT NULL COMMENT '评论内容',
  `reply_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '回复的评论ID(为0表示顶级评论)',
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '状态:0正常,1已删除',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_work_id` (`work_id`),
  KEY `idx_uid` (`uid`),
  KEY `idx_reply_id` (`reply_id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='评论表'; 