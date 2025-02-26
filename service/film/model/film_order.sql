-- 胶片冲洗订单表
CREATE TABLE `film_order`
(
    `id`           bigint unsigned     NOT NULL AUTO_INCREMENT,
    `foid`         varchar(64)         NOT NULL DEFAULT '' COMMENT '订单号',
    `uid`          bigint unsigned     NOT NULL DEFAULT '0' COMMENT '用户ID',
    `address_id`   bigint unsigned     NOT NULL DEFAULT '0' COMMENT '收货地址ID',
    `return_film`  tinyint(1)          NOT NULL DEFAULT '0' COMMENT '是否回寄底片',
    `total_price`  int(10) unsigned    NOT NULL DEFAULT '0' COMMENT '订单总价(分)',
    `shipping_fee` int(10) unsigned    NOT NULL DEFAULT '0' COMMENT '运费(分)',
    `status`       tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '订单状态:0待付款,1冲洗处理中,2待收货,3已完成',
    `remark`       varchar(255)        NOT NULL DEFAULT '' COMMENT '订单备注',
    `create_time`  timestamp           NULL     DEFAULT CURRENT_TIMESTAMP,
    `update_time`  timestamp           NULL     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_foid` (`foid`),
    KEY `idx_uid` (`uid`),
    KEY `idx_status` (`status`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

-- 胶片冲洗订单项表
CREATE TABLE `film_order_item`
(
    `id`            bigint unsigned  NOT NULL AUTO_INCREMENT,
    `film_order_id` bigint unsigned  NOT NULL DEFAULT '0' COMMENT '冲洗订单ID',
    `film_type`     varchar(64)      NOT NULL DEFAULT '' COMMENT '胶片类型',
    `film_brand`    varchar(64)      NOT NULL DEFAULT '' COMMENT '胶片品牌',
    `size`          varchar(64)      NOT NULL DEFAULT '' COMMENT '尺寸规格',
    `quantity`      int(10) unsigned NOT NULL DEFAULT '0' COMMENT '数量',
    `price`         int(10) unsigned NOT NULL DEFAULT '0' COMMENT '单价(分)',
    `amount`        int(10) unsigned NOT NULL DEFAULT '0' COMMENT '总价(分)',
    `remark`        varchar(255)     NOT NULL DEFAULT '' COMMENT '备注',
    `create_time`   timestamp        NULL     DEFAULT CURRENT_TIMESTAMP,
    `update_time`   timestamp        NULL     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_film_order_id` (`film_order_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4; 