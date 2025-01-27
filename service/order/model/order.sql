CREATE TABLE `order`
(
    `id`          bigint unsigned     NOT NULL AUTO_INCREMENT,
    `uid`         bigint unsigned     NOT NULL DEFAULT '0' COMMENT '用户ID',
    `amount`      int(10) unsigned    NOT NULL DEFAULT '0' COMMENT '订单金额',
    `status`      tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '订单状态',
    `create_time` timestamp           NULL     DEFAULT CURRENT_TIMESTAMP,
    `update_time` timestamp           NULL     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_uid` (`uid`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;
