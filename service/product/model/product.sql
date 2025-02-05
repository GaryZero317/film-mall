CREATE TABLE `product`
(
    `id`          bigint unsigned     NOT NULL AUTO_INCREMENT,
    `name`        varchar(255)        NOT NULL DEFAULT '' COMMENT '产品名称',
    `desc`        varchar(255)        NOT NULL DEFAULT '' COMMENT '产品描述',
    `stock`       int(10) unsigned    NOT NULL DEFAULT '0' COMMENT '产品库存',
    `amount`      int(10) unsigned    NOT NULL DEFAULT '0' COMMENT '产品金额',
    `status`      tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '产品状态',
    `category_id` bigint unsigned     NOT NULL DEFAULT '0' COMMENT '分类ID',
    `create_time` timestamp           NULL     DEFAULT CURRENT_TIMESTAMP,
    `update_time` timestamp           NULL     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_category_id` (`category_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;
