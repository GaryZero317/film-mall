-- 清空现有数据
TRUNCATE TABLE product_view_log;
TRUNCATE TABLE user_activity_log;
TRUNCATE TABLE product_sales_daily;
TRUNCATE TABLE category_sales_daily;

-- 设置时间范围：最近30天
SET @start_date = DATE_SUB(CURDATE(), INTERVAL 30 DAY);
SET @end_date = CURDATE();

-- 创建用户ID范围 (1-100)
-- 创建商品ID范围 (1-50)
-- 创建类别ID - 使用实际存在的类别ID
-- 1(135胶卷), 2(120胶卷), 3(拍立得相纸), 4(彩色负片), 5(彩色反转片), 6(黑白负片),
-- 7(120彩色负片), 8(120彩色反转片), 9(宝丽来相纸), 10(富士相纸), 11(电影卷),
-- 13(电影彩色负片), 14(电影彩色反转片), 15(电影黑白负片), 16(胶卷冲洗), 17(哈苏X5), 18(富士SP3000)

-- 1. 填充商品浏览记录 (product_view_log)
INSERT INTO product_view_log (user_id, product_id, view_time, ip, user_agent)
SELECT 
    FLOOR(1 + RAND() * 100) AS user_id,
    FLOOR(1 + RAND() * 50) AS product_id,
    DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 30) DAY) + INTERVAL FLOOR(RAND() * 24) HOUR + INTERVAL FLOOR(RAND() * 60) MINUTE + INTERVAL FLOOR(RAND() * 60) SECOND AS view_time,
    CONCAT('192.168.', FLOOR(RAND() * 255), '.', FLOOR(RAND() * 255)) AS ip,
    ELT(FLOOR(1 + RAND() * 3), 
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36',
        'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.1 Safari/605.1.15',
        'Mozilla/5.0 (iPhone; CPU iPhone OS 14_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/604.1'
    ) AS user_agent
FROM 
    information_schema.columns
LIMIT 1000;

-- 2. 填充用户活跃度记录 (user_activity_log)
-- 活动类型分布: 浏览(60%), 加购(25%), 下单(10%), 支付(5%)
INSERT INTO user_activity_log (user_id, activity_type, activity_time, related_id)
WITH activity_data AS (
    SELECT 
        FLOOR(1 + RAND() * 100) AS user_id,
        CASE 
            WHEN RAND() < 0.6 THEN 'view'
            WHEN RAND() < 0.85 THEN 'cart'
            WHEN RAND() < 0.95 THEN 'order'
            ELSE 'payment'
        END AS activity_type,
        DATE_ADD(@start_date, INTERVAL FLOOR(RAND() * 30) DAY) + INTERVAL FLOOR(RAND() * 24) HOUR + INTERVAL FLOOR(RAND() * 60) MINUTE + INTERVAL FLOOR(RAND() * 60) SECOND AS activity_time,
        FLOOR(1 + RAND() * 50) AS related_id
    FROM 
        information_schema.columns
    LIMIT 2000
)
SELECT 
    user_id, activity_type, activity_time, related_id
FROM 
    activity_data;

-- 3. 填充商品销售统计 (product_sales_daily)
-- 为每个商品在每一天生成随机销售数据
INSERT INTO product_sales_daily (product_id, category_id, sales_date, sales_count, sales_amount)
WITH dates AS (
    SELECT 
        DATE_ADD(@start_date, INTERVAL seq DAY) AS date
    FROM 
        (SELECT 0 AS seq UNION SELECT 1 UNION SELECT 2 UNION SELECT 3 UNION SELECT 4
         UNION SELECT 5 UNION SELECT 6 UNION SELECT 7 UNION SELECT 8 UNION SELECT 9
         UNION SELECT 10 UNION SELECT 11 UNION SELECT 12 UNION SELECT 13 UNION SELECT 14
         UNION SELECT 15 UNION SELECT 16 UNION SELECT 17 UNION SELECT 18 UNION SELECT 19
         UNION SELECT 20 UNION SELECT 21 UNION SELECT 22 UNION SELECT 23 UNION SELECT 24
         UNION SELECT 25 UNION SELECT 26 UNION SELECT 27 UNION SELECT 28 UNION SELECT 29
        ) AS seq
    WHERE DATE_ADD(@start_date, INTERVAL seq DAY) <= @end_date
),
products AS (
    SELECT 
        p AS product_id, 
        -- 根据实际商品ID分配类别
        -- 这里使用CASE WHEN进行分类，映射到实际类别
        CASE 
            WHEN p <= 10 THEN 4  -- 前10个商品为彩色负片(135)
            WHEN p <= 15 THEN 5  -- 11-15号商品为彩色反转片(135)
            WHEN p <= 20 THEN 6  -- 16-20号商品为黑白负片(135)
            WHEN p <= 25 THEN 7  -- 21-25号商品为彩色负片(120)
            WHEN p <= 30 THEN 8  -- 26-30号商品为彩色反转片(120)
            WHEN p <= 35 THEN 9  -- 31-35号商品为宝丽来相纸
            WHEN p <= 40 THEN 10 -- 36-40号商品为富士相纸
            WHEN p <= 43 THEN 13 -- 41-43号商品为电影彩色负片
            WHEN p <= 46 THEN 14 -- 44-46号商品为电影彩色反转片
            WHEN p <= 50 THEN 15 -- 47-50号商品为电影黑白负片
        END AS category_id
    FROM 
        (SELECT 1 AS p UNION SELECT 2 UNION SELECT 3 UNION SELECT 4 UNION SELECT 5
         UNION SELECT 6 UNION SELECT 7 UNION SELECT 8 UNION SELECT 9 UNION SELECT 10
         UNION SELECT 11 UNION SELECT 12 UNION SELECT 13 UNION SELECT 14 UNION SELECT 15
         UNION SELECT 16 UNION SELECT 17 UNION SELECT 18 UNION SELECT 19 UNION SELECT 20
         UNION SELECT 21 UNION SELECT 22 UNION SELECT 23 UNION SELECT 24 UNION SELECT 25
         UNION SELECT 26 UNION SELECT 27 UNION SELECT 28 UNION SELECT 29 UNION SELECT 30
         UNION SELECT 31 UNION SELECT 32 UNION SELECT 33 UNION SELECT 34 UNION SELECT 35
         UNION SELECT 36 UNION SELECT 37 UNION SELECT 38 UNION SELECT 39 UNION SELECT 40
         UNION SELECT 41 UNION SELECT 42 UNION SELECT 43 UNION SELECT 44 UNION SELECT 45
         UNION SELECT 46 UNION SELECT 47 UNION SELECT 48 UNION SELECT 49 UNION SELECT 50
        ) AS p
)
SELECT 
    p.product_id, 
    p.category_id,
    d.date,
    -- 胶片销量模式：周末比平日高30%，平均每天每种胶片销售5-30卷
    FLOOR(5 + RAND() * 25) * IF(DAYOFWEEK(d.date) IN (1, 7), 1.3, 1) AS sales_count,
    -- 根据实际类别设置价格区间
    FLOOR(5 + RAND() * 25) * IF(DAYOFWEEK(d.date) IN (1, 7), 1.3, 1) * 
    CASE p.category_id
        WHEN 4 THEN 50 + RAND() * 70   -- 彩色负片(135)：¥50-120
        WHEN 5 THEN 80 + RAND() * 70   -- 彩色反转片(135)：¥80-150
        WHEN 6 THEN 40 + RAND() * 60   -- 黑白负片(135)：¥40-100
        WHEN 7 THEN 60 + RAND() * 90   -- 彩色负片(120)：¥60-150
        WHEN 8 THEN 90 + RAND() * 80   -- 彩色反转片(120)：¥90-170
        WHEN 9 THEN 60 + RAND() * 140  -- 宝丽来相纸：¥60-200
        WHEN 10 THEN 50 + RAND() * 120 -- 富士相纸：¥50-170
        WHEN 13 THEN 150 + RAND() * 200 -- 电影彩色负片：¥150-350
        WHEN 14 THEN 180 + RAND() * 220 -- 电影彩色反转片：¥180-400
        WHEN 15 THEN 120 + RAND() * 180 -- 电影黑白负片：¥120-300
        ELSE 50 + RAND() * 100 -- 默认价格区间
    END AS sales_amount
FROM 
    products p
CROSS JOIN 
    dates d;

-- 4. 填充类别销售统计 (category_sales_daily)
-- 基于商品销售统计聚合计算
INSERT INTO category_sales_daily (category_id, sales_date, sales_count, sales_amount)
SELECT 
    category_id,
    sales_date,
    SUM(sales_count) AS sales_count,
    SUM(sales_amount) AS sales_amount
FROM 
    product_sales_daily
GROUP BY 
    category_id, sales_date;

-- 添加一些热门商品偏好
-- 热门胶片：富士C200、柯达金200、乐凯100等(与实际热门商品对应)
UPDATE product_sales_daily 
SET 
    sales_count = sales_count * 2.5,
    sales_amount = sales_amount * 2.5
WHERE 
    product_id IN (3, 8, 13, 21, 31) -- 热门胶片ID
    AND sales_date >= DATE_SUB(CURDATE(), INTERVAL 7 DAY);

-- 假期销售高峰（如某个节日）
UPDATE product_sales_daily 
SET 
    sales_count = sales_count * 3,
    sales_amount = sales_amount * 3
WHERE 
    sales_date IN (
        DATE_SUB(CURDATE(), INTERVAL 12 DAY),
        DATE_SUB(CURDATE(), INTERVAL 11 DAY),
        DATE_SUB(CURDATE(), INTERVAL 10 DAY)
    );

-- 添加一些用户购买模式
-- 摄影爱好者和专业摄影师购买模式
INSERT INTO user_activity_log (user_id, activity_type, activity_time, related_id)
SELECT 
    user_id,
    'payment',
    DATE_ADD(CURDATE(), INTERVAL - FLOOR(RAND() * 30) DAY) + INTERVAL FLOOR(RAND() * 24) HOUR,
    -- 摄影爱好者更倾向于购买专业胶片
    ELT(FLOOR(1 + RAND() * 5), 8, 15, 20, 42, 45)
FROM (
    SELECT 
        FLOOR(1 + RAND() * 10) AS user_id
    FROM 
        information_schema.columns
    LIMIT 100
) AS active_users;

-- 新手用户购买偏好（倾向于购买入门级胶片和一次性相机）
INSERT INTO user_activity_log (user_id, activity_type, activity_time, related_id)
SELECT 
    user_id + 50, -- 从51-60的用户ID，区分于上面的专业用户
    'payment',
    DATE_ADD(CURDATE(), INTERVAL - FLOOR(RAND() * 30) DAY) + INTERVAL FLOOR(RAND() * 24) HOUR,
    -- 新手用户更倾向于购买入门级产品
    ELT(FLOOR(1 + RAND() * 5), 1, 3, 13, 31, 37)
FROM (
    SELECT 
        FLOOR(1 + RAND() * 10) AS user_id
    FROM 
        information_schema.columns
    LIMIT 50
) AS novice_users;

-- 更新类别销售统计以匹配修改后的商品销售数据
TRUNCATE TABLE category_sales_daily;
INSERT INTO category_sales_daily (category_id, sales_date, sales_count, sales_amount)
SELECT 
    category_id,
    sales_date,
    SUM(sales_count) AS sales_count,
    SUM(sales_amount) AS sales_amount
FROM 
    product_sales_daily
GROUP BY 
    category_id, sales_date; 