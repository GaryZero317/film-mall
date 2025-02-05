-- 添加category_id字段
ALTER TABLE `product`
ADD COLUMN `category_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '分类ID' AFTER `status`;

-- 添加category_id索引
ALTER TABLE `product`
ADD INDEX `idx_category_id` (`category_id`);

-- 更新135胶卷商品的分类
-- 柯达彩色负片
UPDATE product SET category_id = 4 WHERE id IN (13,14,15,16,17,18,19);  -- Gold200, Portra160/400等
-- 富士彩色负片
UPDATE product SET category_id = 5 WHERE id IN (21,22);  -- C200, C400
-- 伊尔福黑白负片
UPDATE product SET category_id = 6 WHERE id IN (23,24,25,26,27,28);  -- Delta系列, HP5等

-- 更新120胶卷商品的分类
-- 彩色负片
UPDATE product SET category_id = 7 WHERE id IN (29,30,33);  -- PRO160NS, Portra400, Ektar100
-- 彩色反转片
UPDATE product SET category_id = 8 WHERE id IN (31,32);  -- RVP100, Provia100F

-- 更新拍立得相纸的分类
-- 宝丽来相纸
UPDATE product SET category_id = 9 WHERE id IN (34,35);  -- itype系列
-- 富士相纸
UPDATE product SET category_id = 10 WHERE id IN (36,37,38,39,40,41);  -- instax系列 