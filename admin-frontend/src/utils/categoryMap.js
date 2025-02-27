/**
 * 商品类别ID到名称的映射
 */
const categoryMap = {
  "1": "135胶卷",
  "2": "120胶卷",
  "3": "拍立得相纸",
  "4": "彩色负片(135)",
  "5": "彩色反转片(135)",
  "6": "黑白负片(135)",
  "7": "彩色负片(120)",
  "8": "彩色反转片(120)",
  "9": "宝丽来相纸",
  "10": "富士相纸",
  "11": "电影卷",
  "13": "电影彩色负片",
  "14": "电影彩色反转片",
  "15": "电影黑白负片",
  "16": "胶卷冲洗",
  "17": "哈苏X5",
  "18": "富士SP3000"
};

/**
 * 根据类别ID获取类别名称
 * @param {string|number} id 类别ID
 * @returns {string} 类别名称
 */
export function getCategoryName(id) {
  return categoryMap[id] || `类别${id}`;
}

/**
 * 转换类别ID数组为名称数组
 * @param {Array<string|number>} categoryIds 类别ID数组
 * @returns {Array<string>} 类别名称数组
 */
export function convertCategoryIdsToNames(categoryIds) {
  return categoryIds.map(id => getCategoryName(id));
}

export default categoryMap; 