-- 插入胶卷产品相关常见问题 (类别1)
INSERT INTO `faq` (`question`, `answer`, `category`, `priority`, `create_time`, `update_time`) VALUES 
('如何选择适合我相机的胶卷？', '您可以在产品页面查看胶卷规格参数，确认ISO感光度和胶卷类型是否与您的相机兼容。您也可以根据拍摄场景选择合适的胶卷，如日光型、阴天型等。', 1, 100, NOW(), NOW()),
('彩色胶卷和黑白胶卷有什么区别？', '彩色胶卷能够记录自然色彩，适合风景、人像等多彩场景；黑白胶卷只记录光线明暗变化，没有色彩信息，适合注重光影和形式的摄影，且冲洗处理相对简单。', 1, 90, NOW(), NOW()),
('不同ISO值的胶卷有什么不同？', 'ISO值表示胶卷的感光度。低ISO胶卷(如ISO 100)颗粒细腻但需要充足光线；高ISO胶卷(如ISO 800)在弱光环境下也能使用，但颗粒感较强。选择时应根据拍摄环境的光线条件决定。', 1, 80, NOW(), NOW()),
('胶卷的保质期是多久？', '未开封的胶卷一般在冷藏条件下可保存2-3年。已过期胶卷仍可使用但可能出现色偏、对比度降低等效果，有些摄影师会特意使用过期胶卷创造特殊效果。', 1, 70, NOW(), NOW()),
('胶卷需要特殊存储条件吗？', '为延长胶卷保质期，建议将未使用的胶卷存放在冰箱(不是冷冻室)，使用前应取出至少2小时回温，避免结露损坏胶片。已曝光未冲洗的胶卷应尽快冲洗或妥善保存。', 1, 60, NOW(), NOW());

-- 插入账户相关常见问题 (类别2)
INSERT INTO `faq` (`question`, `answer`, `category`, `priority`, `create_time`, `update_time`) VALUES 
('如何注册新账户？', '点击网站右上角的"注册"按钮，填写必要信息如手机号、密码等完成注册。我们也支持第三方账号快速注册。', 2, 100, NOW(), NOW()),
('忘记密码怎么办？', '点击登录页面的"忘记密码"链接，通过验证手机号或邮箱接收验证码后重设密码。', 2, 90, NOW(), NOW()),
('如何修改个人资料？', '登录后，点击右上角头像进入"个人中心"，在"基本资料"选项中可以修改昵称、头像等信息。', 2, 80, NOW(), NOW()),
('账号安全问题如何设置？', '在"个人中心"的"账号安全"选项中，您可以设置登录密码、支付密码和绑定手机号等安全信息。', 2, 70, NOW(), NOW()),
('如何注销账户？', '请联系客服申请账户注销，注销前请确保账户内没有余额和未完成的订单。', 2, 60, NOW(), NOW());

-- 插入支付与配送相关常见问题 (类别3)
INSERT INTO `faq` (`question`, `answer`, `category`, `priority`, `create_time`, `update_time`) VALUES 
('支持哪些支付方式？', '我们支持微信支付、支付宝、银联卡支付以及平台余额支付等多种支付方式。', 3, 100, NOW(), NOW()),
('订单发货时间是多久？', '正常情况下，我们会在付款确认后的1-2个工作日内发货。特殊商品或定制商品可能需要更长时间，具体可参考商品详情页说明。', 3, 90, NOW(), NOW()),
('如何申请退换货？', '在"订单中心"找到相应订单，点击"申请退款/退货"按钮，填写申请原因并上传相关凭证后提交申请。未拆封的商品支持7天无理由退换。', 3, 80, NOW(), NOW()),
('退款多久能到账？', '退款审核通过后，一般1-7个工作日内退回原支付账户，具体到账时间以银行或第三方支付平台为准。', 3, 70, NOW(), NOW()),
('如何查询物流信息？', '在"订单中心"找到相应订单，点击"查看物流"按钮即可查看物流跟踪信息。您也可以通过短信中的物流单号在快递公司官网查询。', 3, 60, NOW(), NOW());

-- 插入胶卷冲洗与使用相关常见问题 (类别4)
INSERT INTO `faq` (`question`, `answer`, `category`, `priority`, `create_time`, `update_time`) VALUES 
('你们提供冲洗服务吗？', '是的，我们提供专业的胶卷冲洗和扫描服务。您可以将已曝光的胶卷寄送至我们的冲洗中心，或在部分城市使用上门取件服务。', 4, 100, NOW(), NOW()),
('冲洗服务的价格和周期是多少？', '标准冲洗服务(含高分辨率扫描)的价格为每卷60-80元不等，具体取决于胶卷类型。正常处理周期为3-5个工作日，加急服务可在24小时内完成。', 4, 90, NOW(), NOW()),
('如何正确装载胶卷到相机中？', '请在暗处或避光环境下操作，打开相机后盖，将胶卷装入左侧仓，将胶片前端拉出并插入右侧卷轴，关闭后盖后适度摇动并按下快门进行试拍，确认胶片正常前进。', 4, 80, NOW(), NOW()),
('拍摄完成后如何保存胶卷？', '拍摄完毕后，按下相机底部的回片按钮，将胶片完全回卷入胶盒。取出胶卷后应放入原塑料盒中，避免光线和潮湿，并尽快送去冲洗。已曝光未冲洗的胶卷应避免X光检查和高温环境。', 4, 70, NOW(), NOW()),
('我可以自己冲洗胶卷吗？', '是的，我们销售胶卷冲洗套件和相关化学药剂，适合自己在家冲洗。黑白胶卷冲洗相对简单，彩色胶卷冲洗则需要更精确的温度控制和专业设备。我们网站上有详细的教程可供参考。', 4, 60, NOW(), NOW()); 