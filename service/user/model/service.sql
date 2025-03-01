CREATE TABLE `customer_service`
(
    `id`          bigint unsigned     NOT NULL AUTO_INCREMENT,
    `user_id`     bigint unsigned     NOT NULL COMMENT '用户ID',
    `title`       varchar(255)        NOT NULL DEFAULT '' COMMENT '问题标题',
    `content`     text                NOT NULL COMMENT '问题内容',
    `type`        tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '问题类型：1-订单问题，2-产品咨询，3-售后服务，4-其他',
    `status`      tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '状态：1-待处理，2-处理中，3-已解决',
    `reply`       text                COMMENT '客服回复',
    `reply_time`  timestamp           NULL     DEFAULT NULL COMMENT '回复时间',
    `contact_way` varchar(255)        NOT NULL DEFAULT '' COMMENT '联系方式',
    `create_time` timestamp           NULL     DEFAULT CURRENT_TIMESTAMP,
    `update_time` timestamp           NULL     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

CREATE TABLE `faq`
(
    `id`          bigint unsigned     NOT NULL AUTO_INCREMENT,
    `question`    varchar(255)        NOT NULL DEFAULT '' COMMENT '问题',
    `answer`      text                NOT NULL COMMENT '答案',
    `category`    tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '分类：1-订单相关，2-产品相关，3-配送相关，4-其他',
    `priority`    int                 NOT NULL DEFAULT '0' COMMENT '优先级（排序用）',
    `create_time` timestamp           NULL     DEFAULT CURRENT_TIMESTAMP,
    `update_time` timestamp           NULL     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_category` (`category`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4; 