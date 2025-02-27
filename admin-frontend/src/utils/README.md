# 工具函数说明

## categoryMap.js - 类别映射工具

该工具用于将后端返回的数字类别ID转换为可读的类别名称，优化前端显示效果。

### 功能

- `categoryMap` 对象: 包含所有类别ID到类别名称的映射关系
- `getCategoryName(id)`: 根据类别ID获取对应的类别名称，如果ID不存在则返回默认格式"类别{id}"
- `convertCategoryIdsToNames(categoryIds)`: 批量转换类别ID数组为名称数组

### 使用示例

```javascript
import { getCategoryName } from '@/utils/categoryMap';

// 单个ID转换
const categoryName = getCategoryName('1'); // 返回 "135胶卷"

// 处理API返回的类别数据
const categories = [1, 2, 3];
const categoryNames = categories.map(id => getCategoryName(id));
// 结果: ["135胶卷", "120胶卷", "拍立得相纸"]
```

### 在Dashboard中的应用

类别映射工具在Dashboard.vue中被用于以下场景：

1. 类别销售饼图：将后端返回的类别ID转换为可读名称
2. 热门商品图表：优化超长产品名称的显示

使用此工具使得前端展示更加友好，无需修改后端API和数据库结构。 