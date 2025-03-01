-- 聊天消息表
CREATE TABLE IF NOT EXISTS `chat_message` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '消息ID',
  `user_id` bigint(20) NOT NULL COMMENT '用户ID',
  `admin_id` bigint(20) NOT NULL COMMENT '管理员ID',
  `direction` tinyint(4) NOT NULL DEFAULT '1' COMMENT '消息方向：1用户到管理员，2管理员到用户',
  `content` text COMMENT '消息内容',
  `type` tinyint(4) NOT NULL DEFAULT '1' COMMENT '消息类型：1文本，2图片',
  `read_status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '读取状态：0未读，1已读',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_admin` (`user_id`,`admin_id`),
  KEY `idx_admin_user` (`admin_id`,`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='聊天消息表'; 