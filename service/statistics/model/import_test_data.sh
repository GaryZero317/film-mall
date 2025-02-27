#!/bin/bash
echo "正在导入胶片商城统计测试数据..."
echo "请确保MySQL服务已启动"

mysql -u root -p200317 mall < test_data.sql

if [ $? -eq 0 ]; then
    echo "测试数据导入成功！"
    echo "现在可以启动统计服务并在前端查看统计图表"
else
    echo "导入失败，请检查MySQL连接和权限"
fi

read -p "按任意键继续..." key 