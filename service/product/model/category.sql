-- 商品分类表
CREATE TABLE `product_category` (
    `id`          bigint unsigned NOT NULL AUTO_INCREMENT,
    `name`        varchar(50)     NOT NULL COMMENT '分类名称',
    `parent_id`   bigint unsigned NOT NULL DEFAULT '0' COMMENT '父分类ID',
    `level`       int unsigned    NOT NULL DEFAULT '1' COMMENT '分类层级',
    `sort_order`  int unsigned    NOT NULL DEFAULT '0' COMMENT '排序',
    `status`      tinyint        NOT NULL DEFAULT '1' COMMENT '状态：0-禁用，1-启用',
    `create_time` timestamp      NULL     DEFAULT CURRENT_TIMESTAMP,
    `update_time` timestamp      NULL     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_parent_id` (`parent_id`),
    KEY `idx_level` (`level`),
    KEY `idx_sort_order` (`sort_order`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT='商品分类表';

-- 插入一级分类
INSERT INTO `product_category` (`name`, `parent_id`, `level`, `sort_order`) VALUES
('135胶卷', 0, 1, 1),
('120胶卷', 0, 1, 2),
('拍立得相纸', 0, 1, 3);

-- 插入135胶卷的子分类
INSERT INTO `product_category` (`name`, `parent_id`, `level`, `sort_order`) VALUES
('柯达彩色负片', 1, 2, 1),    -- 适用于：Gold200, Portra160/400, ColorPlus200, UltraMax400, ProImage100, Ektar100
('富士彩色负片', 1, 2, 2),    -- 适用于：C200, C400
('伊尔福黑白负片', 1, 2, 3);  -- 适用于：Delta系列, HP5 Plus, Pan400

-- 插入120胶卷的子分类
INSERT INTO `product_category` (`name`, `parent_id`, `level`, `sort_order`) VALUES
('彩色负片', 2, 2, 1),        -- 适用于：Portra400, PRO160NS, Ektar100
('彩色反转片', 2, 2, 2);      -- 适用于：RVP100, Provia100F

-- 插入拍立得相纸的子分类
INSERT INTO `product_category` (`name`, `parent_id`, `level`, `sort_order`) VALUES
('宝丽来相纸', 3, 2, 1),      -- 适用于：itype系列
('富士相纸', 3, 2, 2);        -- 适用于：instax系列 