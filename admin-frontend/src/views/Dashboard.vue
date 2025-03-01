<template>
  <div class="dashboard-container">
    <!-- 筛选与导出区域 -->
    <el-card shadow="hover" class="filter-card">
      <div class="filter-header">
        <h2 class="section-title">数据分析</h2>
        <div class="filter-controls">
          <el-select v-model="timeRange" @change="handleTimeRangeChange" class="time-select" size="large">
            <el-option label="最近7天" value="7days" />
            <el-option label="最近30天" value="30days" />
            <el-option label="本月" value="month" />
            <el-option label="本季度" value="quarter" />
            <el-option label="本年" value="year" />
          </el-select>
          <el-button type="primary" @click="exportData" class="export-btn" size="large">
            <el-icon><Download /></el-icon>
            导出统计数据
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- 数据统计卡片 -->
    <div class="stats-card-row">
      <el-card shadow="hover" class="stats-card sales-card">
        <div class="stats-card-inner">
          <div class="stats-icon">
            <el-icon><ShoppingBag /></el-icon>
          </div>
          <div class="stats-content">
            <div class="stats-title">商品总数</div>
            <div class="stats-value">{{ statistics.products }}</div>
            <div class="stats-trend">
              <span class="trend-label">较上期</span>
              <span class="trend-value positive">
                <el-icon><ArrowUp /></el-icon> 12.5%
              </span>
            </div>
          </div>
        </div>
      </el-card>
      
      <el-card shadow="hover" class="stats-card orders-card">
        <div class="stats-card-inner">
          <div class="stats-icon">
            <el-icon><List /></el-icon>
          </div>
          <div class="stats-content">
            <div class="stats-title">订单总数</div>
            <div class="stats-value">{{ statistics.orders }}</div>
            <div class="stats-trend">
              <span class="trend-label">较上期</span>
              <span class="trend-value positive">
                <el-icon><ArrowUp /></el-icon> 8.2%
              </span>
            </div>
          </div>
        </div>
      </el-card>
      
      <el-card shadow="hover" class="stats-card revenue-card">
        <div class="stats-card-inner">
          <div class="stats-icon">
            <el-icon><Money /></el-icon>
          </div>
          <div class="stats-content">
            <div class="stats-title">支付总额</div>
            <div class="stats-value">{{ statistics.payments }}</div>
            <div class="stats-trend">
              <span class="trend-label">较上期</span>
              <span class="trend-value positive">
                <el-icon><ArrowUp /></el-icon> 15.3%
              </span>
            </div>
          </div>
        </div>
      </el-card>
      
      <el-card shadow="hover" class="stats-card users-card">
        <div class="stats-card-inner">
          <div class="stats-icon">
            <el-icon><User /></el-icon>
          </div>
          <div class="stats-content">
            <div class="stats-title">用户总数</div>
            <div class="stats-value">4,382</div>
            <div class="stats-trend">
              <span class="trend-label">较上期</span>
              <span class="trend-value positive">
                <el-icon><ArrowUp /></el-icon> 5.7%
              </span>
            </div>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 折线图与环形图 -->
    <div class="chart-row">
      <el-card shadow="hover" class="chart-card trend-chart">
        <div class="chart-header">
          <h3 class="chart-title">销售趋势</h3>
          <div class="chart-actions">
            <el-radio-group v-model="trendView" size="small">
              <el-radio-button label="sales">销售额</el-radio-button>
              <el-radio-button label="orders">订单数</el-radio-button>
            </el-radio-group>
          </div>
        </div>
        <div class="chart-content">
          <div ref="lineChartRef" style="width: 100%; height: 350px"></div>
        </div>
      </el-card>
      
      <el-card shadow="hover" class="chart-card distribution-chart">
        <div class="chart-header">
          <h3 class="chart-title">编程占比</h3>
        </div>
        <div class="chart-content">
          <div ref="pieChartRef" style="width: 100%; height: 350px"></div>
        </div>
      </el-card>
    </div>

    <!-- 热门商品与类别分布 -->
    <div class="chart-row">
      <el-card shadow="hover" class="chart-card hot-products">
        <div class="chart-header">
          <h3 class="chart-title">热门商品TOP10</h3>
          <el-tooltip content="查看更多商品数据" placement="top">
            <el-button link type="primary" @click="$router.push('/products')">
              更多 <el-icon><ArrowRight /></el-icon>
            </el-button>
          </el-tooltip>
        </div>
        <div class="chart-content">
          <div ref="hotProductsChartRef" style="height: 400px"></div>
        </div>
      </el-card>
      
      <el-card shadow="hover" class="chart-card category-chart">
        <div class="chart-header">
          <h3 class="chart-title">商品类别销售分布</h3>
          <el-tooltip content="查看类别详情" placement="top">
            <el-button link type="primary">
              详情 <el-icon><ArrowRight /></el-icon>
            </el-button>
          </el-tooltip>
        </div>
        <div class="chart-content">
          <div ref="categoryChartRef" style="height: 400px"></div>
        </div>
      </el-card>
    </div>

    <!-- 用户行为与活跃度 -->
    <div class="chart-row">
      <el-card shadow="hover" class="chart-card behavior-chart">
        <div class="chart-header">
          <h3 class="chart-title">用户购买行为分析</h3>
        </div>
        <div class="chart-content">
          <div ref="userBehaviorChartRef" style="height: 400px"></div>
        </div>
      </el-card>
      
      <el-card shadow="hover" class="chart-card activity-chart">
        <div class="chart-header">
          <h3 class="chart-title">用户活跃时段分析</h3>
        </div>
        <div class="chart-content">
          <div ref="userActivityChartRef" style="height: 400px"></div>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, reactive, watch } from 'vue'
import { ShoppingBag, List, Money, Download, User, ArrowUp, ArrowDown, ArrowRight } from '@element-plus/icons-vue'
import * as echarts from 'echarts'
import { getAdminProductList } from '@/api/product'
import { getAdminOrderList } from '@/api/order'
import { getAdminPaymentList } from '@/api/payment'
import { ElMessage } from 'element-plus'
import { getHotProducts, getCategoryStats, getUserBehavior, getUserActivity, exportStatistics } from '@/api/statistics'
import { getCategoryName } from '@/utils/categoryMap'

// 统计数据
const statistics = reactive({
  products: '0',  // 商品总数
  orders: '0',    // 订单总数
  payments: '0'   // 支付总额
})

// 订单数据
const orderData = ref([])

// 时间范围选择
const timeRange = ref('7days')
const trendView = ref('sales')

// 新增图表引用
const hotProductsChartRef = ref(null)
const categoryChartRef = ref(null)
const userBehaviorChartRef = ref(null)
const userActivityChartRef = ref(null)

// 新增图表实例
let hotProductsChart = null
let categoryChart = null
let userBehaviorChart = null
let userActivityChart = null

// 生成模拟数据
const generateMockData = () => {
  const today = new Date()
  const data = []
  for (let i = 6; i >= 0; i--) {
    const date = new Date(today)
    date.setDate(date.getDate() - i)
    data.push({
      date: date.toLocaleDateString('zh-CN', { month: '2-digit', day: '2-digit' }),
      count: Math.floor(Math.random() * 100) + 50 // 生成50-150之间的随机数
    })
  }
  return data
}

// 获取订单数据
const fetchOrderData = async () => {
  try {
    // 获取最近7天的订单数据
    const endDate = new Date()
    const startDate = new Date()
    startDate.setDate(startDate.getDate() - 6)
    
    const res = await getAdminOrderList({
      startDate: startDate.toISOString().split('T')[0],
      endDate: endDate.toISOString().split('T')[0],
      page: 1,
      pageSize: 999999
    })

    if (res?.data?.list) {
      // 按日期分组统计订单数量
      const orderStats = {}
      res.data.list.forEach(order => {
        // 检查订单创建时间是否存在
        const createTime = order.create_time || order.createTime || order.created_at
        if (createTime) {
          const date = createTime.split(' ')[0]
          orderStats[date] = (orderStats[date] || 0) + 1
        }
      })

      // 构造最近7天的数据数组
      const data = []
      for (let i = 6; i >= 0; i--) {
        const date = new Date()
        date.setDate(date.getDate() - i)
        const dateStr = date.toISOString().split('T')[0]
        const count = orderStats[dateStr] || 0
        data.push({
          date: dateStr.slice(5).replace('-', '/'), // 格式化为 MM/DD
          count: count
        })
      }
      orderData.value = data
    } else {
      // 如果没有数据，使用空数组
      const data = []
      for (let i = 6; i >= 0; i--) {
        const date = new Date()
        date.setDate(date.getDate() - i)
        data.push({
          date: date.toISOString().split('T')[0].slice(5).replace('-', '/'),
          count: 0
        })
      }
      orderData.value = data
    }
    updateLineChart()
  } catch (error) {
    console.error('获取订单数据失败:', error)
    // 发生错误时显示空数据
    const data = []
    for (let i = 6; i >= 0; i--) {
      const date = new Date()
      date.setDate(date.getDate() - i)
      data.push({
        date: date.toISOString().split('T')[0].slice(5).replace('-', '/'),
        count: 0
      })
    }
    orderData.value = data
    updateLineChart()
  }
}

// 获取统计数据
const fetchStatistics = async () => {
  try {
    // 获取商品总数
    const productsRes = await getAdminProductList({
      page: 1,
      pageSize: 1
    })
    statistics.products = productsRes?.data?.total?.toLocaleString() || '0'
      
    // 获取订单总数和总金额
    const ordersRes = await getAdminOrderList({
      page: 1,
      pageSize: 999999, // 获取所有订单
      status: 1 // 已支付的订单
    })
    
    // 计算订单总数
    statistics.orders = ordersRes?.data?.total?.toLocaleString() || '0'

    // 计算总金额 - 只统计已支付的订单
    let totalAmount = 0
    if (ordersRes?.data?.list) {
      totalAmount = ordersRes.data.list.reduce((sum, order) => {
        // 只统计已支付的订单
        if (order.status === 1) {
          return sum + (order.total_price || 0) + (order.shipping_fee || 0)
        }
        return sum
      }, 0)
    }
    statistics.payments = `¥${(totalAmount / 100).toFixed(2)}`

  } catch (error) {
    console.error('获取统计数据失败:', error)
    ElMessage.error('获取统计数据失败')
  }
}

// 图表实例
let lineChart = null
let pieChart = null

// 图表DOM引用
const lineChartRef = ref(null)
const pieChartRef = ref(null)

// 初始化折线图
const initLineChart = () => {
  if (!lineChartRef.value) return
  
  lineChart = echarts.init(lineChartRef.value)
  const option = {
    tooltip: {
      trigger: 'axis'
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: orderData.value?.length ? orderData.value.map(item => item.date) : []
    },
    yAxis: {
      type: 'value'
    },
    series: [
      {
        name: '订单量',
        type: 'line',
        data: orderData.value?.length ? orderData.value.map(item => item.count) : [],
        smooth: true,
        lineStyle: {
          color: '#2196f3'
        },
        itemStyle: {
          color: '#2196f3'
        }
      }
    ]
  }
  lineChart.setOption(option)
}

// 更新折线图数据
const updateLineChart = () => {
  if (!lineChart || !orderData.value?.length) return
  
  lineChart.setOption({
    xAxis: {
      data: orderData.value.map(item => item.date)
    },
    series: [
      {
        data: orderData.value.map(item => item.count)
      }
    ]
  })
}

// 初始化饼图
const initPieChart = () => {
  pieChart = echarts.init(pieChartRef.value)
  const option = {
    tooltip: {
      trigger: 'item',
      formatter: '{b}: {d}%'
    },
    legend: {
      orient: 'horizontal',
      bottom: 0,
      left: 'center',
      itemWidth: 25,
      itemHeight: 14,
      itemGap: 30,
      textStyle: {
        fontSize: 12,
        padding: [3, 0, 0, 0]
      }
    },
    series: [
      {
        name: '语言占比',
        type: 'pie',
        center: ['50%', '45%'],
        radius: ['45%', '65%'],
        avoidLabelOverlap: true,
        label: {
          show: false, // 关闭饼图上的标签
          position: 'outside'
        },
        emphasis: {
          label: {
            show: true,
            fontSize: '16',
            fontWeight: 'bold'
          }
        },
        labelLine: {
          show: false
        },
        data: [
          { value: 50, name: 'Go', itemStyle: { color: '#00ADD8' } },        // Go的标准色
          { value: 30, name: 'JavaScript', itemStyle: { color: '#F7DF1E' } }, // JS的标准色
          { value: 17, name: 'Vue', itemStyle: { color: '#4FC08D' } },       // Vue的标准色
          { value: 3, name: 'Other', itemStyle: { color: '#E34F26' } }   // Other
        ]
      }
    ]
  }
  pieChart.setOption(option)
}

// 处理时间范围变化
const handleTimeRangeChange = async () => {
  await Promise.all([
    fetchStatistics(),
    fetchOrderData(),
    initHotProductsChart(),
    initCategoryChart(),
    initUserBehaviorChart(),
    initUserActivityChart()
  ])
}

// 初始化热门商品图表
const initHotProductsChart = async () => {
  if (!hotProductsChartRef.value) return
  
  try {
    const res = await getHotProducts({ timeRange: timeRange.value })
    const { products, sales } = res.data

    // 产品名称已经通过API返回，无需额外转换

    hotProductsChart = echarts.init(hotProductsChartRef.value)
    const option = {
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          type: 'shadow'
        }
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true
      },
      xAxis: {
        type: 'value'
      },
      yAxis: {
        type: 'category',
        data: products, // 使用产品名称
        axisLabel: {
          formatter: function(value) {
            // 如果产品名称太长，进行截断
            if (value.length > 15) {
              return value.substring(0, 12) + '...';
            }
            return value;
          }
        }
      },
      series: [
        {
          name: '销售量',
          type: 'bar',
          data: sales,
          itemStyle: {
            color: '#2196f3'
          }
        }
      ]
    }
    hotProductsChart.setOption(option)
  } catch (error) {
    console.error('获取热门商品数据失败:', error)
  }
}

// 初始化商品类别图表
const initCategoryChart = async () => {
  if (!categoryChartRef.value) return
  
  try {
    const res = await getCategoryStats({ timeRange: timeRange.value })
    const { categories, sales } = res.data

    // 将类别ID转换为类别名称
    const categoryNames = categories.map(id => getCategoryName(id))
    
    // 构建饼图数据
    const pieData = categories.map((id, index) => ({
      name: getCategoryName(id),
      value: sales[index]
    }))

    categoryChart = echarts.init(categoryChartRef.value)
    const option = {
      tooltip: {
        trigger: 'item',
        formatter: '{b}: {c} ({d}%)'
      },
      legend: {
        type: 'scroll',        // 启用滚动
        orient: 'vertical',    // 垂直排列
        right: 10,
        top: 20,               // 调整顶部位置
        bottom: 20,            // 设置底部边界以允许滚动
        itemGap: 10,           // 图例项之间的间距
        formatter: function(name) {
          // 如果名称太长则截断
          if (name.length > 10) {
            return name.substring(0, 8) + '...';
          }
          return name;
        },
        textStyle: {
          fontSize: 12         // 调整文本大小
        },
        data: categoryNames    // 使用类别名称数组
      },
      series: [
        {
          name: '类别销售',
          type: 'pie',
          radius: ['40%', '60%'],  // 调小饼图的半径
          center: ['30%', '50%'],  // 将饼图向左移动更多，为图例腾出更多空间
          avoidLabelOverlap: true, // 防止标签重叠
          itemStyle: {
            borderRadius: 10,
            borderColor: '#fff',
            borderWidth: 2
          },
          label: {
            show: false,
            position: 'center'
          },
          emphasis: {
            label: {
              show: true,
              fontSize: '20',
              fontWeight: 'bold'
            }
          },
          labelLine: {
            show: false
          },
          data: pieData // 使用已转换的数据
        }
      ]
    }
    categoryChart.setOption(option)
  } catch (error) {
    console.error('获取类别销售数据失败:', error)
  }
}

// 初始化用户行为图表
const initUserBehaviorChart = async () => {
  if (!userBehaviorChartRef.value) return
  
  try {
    const res = await getUserBehavior({ timeRange: timeRange.value })
    const { dates, views, carts, orders } = res.data

    userBehaviorChart = echarts.init(userBehaviorChartRef.value)
    const option = {
      tooltip: {
        trigger: 'axis'
      },
      legend: {
        data: ['浏览量', '加购量', '下单量']
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        boundaryGap: false,
        data: dates
      },
      yAxis: {
        type: 'value'
      },
      series: [
        {
          name: '浏览量',
          type: 'line',
          data: views,
          smooth: true,
          itemStyle: { color: '#2196f3' }
        },
        {
          name: '加购量',
          type: 'line',
          data: carts,
          smooth: true,
          itemStyle: { color: '#ff9800' }
        },
        {
          name: '下单量',
          type: 'line',
          data: orders,
          smooth: true,
          itemStyle: { color: '#4caf50' }
        }
      ]
    }
    userBehaviorChart.setOption(option)
  } catch (error) {
    console.error('获取用户行为数据失败:', error)
  }
}

// 初始化用户活跃度图表
const initUserActivityChart = async () => {
  if (!userActivityChartRef.value) return
  
  try {
    const res = await getUserActivity({ timeRange: timeRange.value })
    const { hours, activity } = res.data

    userActivityChart = echarts.init(userActivityChartRef.value)
    const option = {
      tooltip: {
        position: 'top'
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        data: hours,
        splitArea: {
          show: true
        }
      },
      yAxis: {
        type: 'category',
        data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日'],
        splitArea: {
          show: true
        }
      },
      visualMap: {
        min: 0,
        max: 10,
        calculable: true,
        orient: 'horizontal',
        left: 'center',
        bottom: '0%'
      },
      series: [{
        name: '活跃度',
        type: 'heatmap',
        data: activity,
        label: {
          show: true
        },
        emphasis: {
          itemStyle: {
            shadowBlur: 10,
            shadowColor: 'rgba(0, 0, 0, 0.5)'
          }
        }
      }]
    }
    userActivityChart.setOption(option)
  } catch (error) {
    console.error('获取用户活跃度数据失败:', error)
  }
}

// 导出统计数据
const exportData = async () => {
  try {
    const res = await exportStatistics({ timeRange: timeRange.value })
    const blob = new Blob([res.data], { type: 'application/vnd.ms-excel' })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', `统计数据_${timeRange.value}.xlsx`)
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
  } catch (error) {
    console.error('导出数据失败:', error)
    ElMessage.error('导出数据失败')
  }
}

// 监听窗口大小变化
const handleResize = () => {
  lineChart?.resize()
  pieChart?.resize()
  hotProductsChart?.resize()
  categoryChart?.resize()
  userBehaviorChart?.resize()
  userActivityChart?.resize()
}

onMounted(async () => {
  await fetchStatistics()
  initLineChart()
  await fetchOrderData()
  initPieChart()
  initHotProductsChart()
  initCategoryChart()
  initUserBehaviorChart()
  initUserActivityChart()
  window.addEventListener('resize', handleResize)
})

// 组件卸载时移除事件监听
onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  lineChart?.dispose()
  pieChart?.dispose()
  hotProductsChart?.dispose()
  categoryChart?.dispose()
  userBehaviorChart?.dispose()
  userActivityChart?.dispose()
})
</script>

<style scoped>
.dashboard-container {
  padding: 24px;
}

/* 筛选区域 */
.filter-card {
  margin-bottom: 24px;
  border-radius: 8px;
  overflow: hidden;
}

.filter-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.section-title {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: #262626;
}

.filter-controls {
  display: flex;
  align-items: center;
  gap: 16px;
}

.time-select {
  width: 150px;
}

.export-btn {
  display: flex;
  align-items: center;
  gap: 6px;
}

/* 统计卡片 */
.stats-card-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 24px;
  margin-bottom: 24px;
}

.stats-card {
  border-radius: 8px;
  overflow: hidden;
  transition: transform 0.3s, box-shadow 0.3s;
  height: 100%;
}

.stats-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.1);
}

.stats-card-inner {
  display: flex;
  padding: 4px;
  height: 100%;
}

.stats-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 60px;
  height: 60px;
  border-radius: 12px;
  margin-right: 16px;
  font-size: 24px;
  color: white;
}

.sales-card .stats-icon {
  background: linear-gradient(135deg, #1890ff, #096dd9);
}

.orders-card .stats-icon {
  background: linear-gradient(135deg, #ff9800, #ed6c02);
}

.revenue-card .stats-icon {
  background: linear-gradient(135deg, #52c41a, #389e0d);
}

.users-card .stats-icon {
  background: linear-gradient(135deg, #722ed1, #531dab);
}

.stats-content {
  display: flex;
  flex-direction: column;
  justify-content: center;
  flex: 1;
}

.stats-title {
  font-size: 14px;
  color: #8c8c8c;
  margin-bottom: 4px;
}

.stats-value {
  font-size: 24px;
  font-weight: 600;
  color: #262626;
  margin-bottom: 8px;
}

.stats-trend {
  display: flex;
  align-items: center;
  font-size: 12px;
}

.trend-label {
  color: #8c8c8c;
  margin-right: 8px;
}

.trend-value {
  display: flex;
  align-items: center;
  font-weight: 500;
}

.trend-value.positive {
  color: #52c41a;
}

.trend-value.negative {
  color: #f5222d;
}

/* 图表卡片 */
.chart-row {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 24px;
  margin-bottom: 24px;
}

.chart-row:last-child {
  margin-bottom: 0;
}

.chart-card {
  border-radius: 8px;
  overflow: hidden;
  height: 100%;
}

.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid #f0f0f0;
}

.chart-title {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #262626;
}

.chart-content {
  padding: 16px;
}

/* 响应式 */
@media (max-width: 1400px) {
  .stats-card-row {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .chart-row {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .stats-card-row {
    grid-template-columns: 1fr;
  }
  
  .filter-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }
  
  .filter-controls {
    width: 100%;
    flex-direction: column;
    align-items: flex-start;
  }
  
  .time-select {
    width: 100%;
  }
  
  .export-btn {
    width: 100%;
  }
}
</style> 