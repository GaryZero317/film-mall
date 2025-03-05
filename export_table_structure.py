import pandas as pd
from sqlalchemy import create_engine, text

# 数据库配置
DB_TYPE = 'mysql'  # 数据库类型
DB_DRIVER = 'pymysql'  # 驱动
DB_USER = 'root'  # 数据库用户名
DB_PASSWORD = '200317'  # 数据库密码
DB_HOST = 'localhost'  # 数据库主机
DB_PORT = '3306'  # 数据库端口
DB_NAME = 'mall'  # 数据库名称

# 创建数据库连接
engine = create_engine(f'{DB_TYPE}+{DB_DRIVER}://{DB_USER}:{DB_PASSWORD}@{DB_HOST}:{DB_PORT}/{DB_NAME}')

# 指定要导出的表名
table_name = input("请输入要导出的表名: ")

# 获取指定表的表结构
with engine.connect() as connection:
    structure = connection.execute(text(f"DESCRIBE `{table_name}`;"))
    
    # 准备表结构数据
    table_structure = []
    for row in structure:
        table_structure.append({
            '列名': row[0],
            '数据类型': row[1],
            '是否为空': row[2],
            '主键': row[3],
            '自增': row[4],
            '默认值': row[5],
            '备注': ''  # 可以根据需要添加备注
        })

# 转换为 DataFrame
df = pd.DataFrame(table_structure)

# 导出到 Excel 文件
output_file = f'{table_name}_structure.xlsx'
df.to_excel(output_file, index=False)
print(f'表结构已导出到 {output_file}') 