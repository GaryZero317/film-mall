<template>
  <div class="dashboard-container">
    <!-- 时间维度选择 -->
    <el-row class="mb-4">
      <el-col :span="24">
        <el-card shadow="hover">
          <div class="filter-container">
            <el-select v-model="timeRange" @change="handleTimeRangeChange" class="mr-3">
              <el-option label="最近7天" value="7days" />
              <el-option label="最近30天" value="30days" />
              <el-option label="本月" value="month" />
              <el-option label="本季度" value="quarter" />
              <el-option label="本年" value="year" />
            </el-select>
            <el-button type="primary" @click="exportData">
              <el-icon><Download /></el-icon>
              导出统计数据
            </el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 数据统计卡片 -->
    <el-row :gutter="20" class="data-overview">
      <el-col :span="8">
        <el-card shadow="hover" class="data-card">
          <div class="card-header">
            <div class="icon-wrapper blue">
              <el-icon><ShoppingBag /></el-icon>
            </div>
            <div class="data-info">
              <div class="label">商品总数</div>
              <div class="value">{{ statistics.products }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover" class="data-card">
          <div class="card-header">
            <div class="icon-wrapper orange">
              <el-icon><List /></el-icon>
            </div>
            <div class="data-info">
              <div class="label">订单总数</div>
              <div class="value">{{ statistics.orders }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover" class="data-card">
          <div class="card-header">
            <div class="icon-wrapper green">
              <el-icon><Money /></el-icon>
            </div>
            <div class="data-info">
              <div class="label">支付总额</div>
              <div class="value">{{ statistics.payments }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 折线图 -->
    <el-row :gutter="20" class="chart-row">
      <el-col :span="24">
        <el-card shadow="hover">
          <div class="chart-container">
            <div ref="lineChartRef" style="width: 100%; height: 350px"></div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 雷达图、饼图、柱状图 -->
    <el-row :gutter="20" class="chart-row">
      <el-col :span="8">
        <el-card shadow="hover">
          <div class="chart-container">
            <div ref="radarChartRef" style="width: 100%; height: 300px"></div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover">
          <div class="chart-container">
            <div ref="pieChartRef" style="width: 100%; height: 300px"></div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover">
          <div class="chart-container">
            <div ref="barChartRef" style="width: 100%; height: 300px"></div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 热门商品分析 -->
    <el-row :gutter="20" class="chart-row">
      <el-col :span="12">
        <el-card shadow="hover">
          <div class="chart-container">
            <div class="chart-title">热门商品TOP10</div>
            <div ref="hotProductsChartRef" style="height: 400px"></div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card shadow="hover">
          <div class="chart-container">
            <div class="chart-title">商品类别销售分布</div>
            <div ref="categoryChartRef" style="height: 400px"></div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 用户行为分析 -->
    <el-row :gutter="20" class="chart-row">
      <el-col :span="12">
        <el-card shadow="hover">
          <div class="chart-container">
            <div class="chart-title">用户购买行为分析</div>
            <div ref="userBehaviorChartRef" style="height: 400px"></div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card shadow="hover">
          <div class="chart-container">
            <div class="chart-title">用户活跃度分析</div>
            <div ref="userActivityChartRef" style="height: 400px"></div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, reactive } from 'vue'
import { ShoppingBag, List, Money, Download } from '@element-plus/icons-vue'
import * as echarts from 'echarts'
import { getAdminProductList } from '@/api/product'
import { getAdminOrderList } from '@/api/order'
import { getAdminPaymentList } from '@/api/payment'
import { ElMessage } from 'element-plus'
import { getHotProducts } from '@/api/statistics'
import { getCategoryStats } from '@/api/statistics'
import { getUserBehavior } from '@/api/statistics'
import { getUserActivity } from '@/api/statistics'
import { exportStatistics } from '@/api/statistics'

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
let radarChart = null
let pieChart = null
let barChart = null

// 图表DOM引用
const lineChartRef = ref(null)
const radarChartRef = ref(null)
const pieChartRef = ref(null)
const barChartRef = ref(null)

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

// 初始化雷达图
const initRadarChart = () => {
  radarChart = echarts.init(radarChartRef.value)
  const option = {
    title: {
      text: '系统模块完成度',
      left: 'center',
      top: 20,
      textStyle: {
        fontSize: 14
      }
    },
    radar: {
      indicator: [
        { name: '商品管理', max: 100 },
        { name: '订单系统', max: 100 },
        { name: '用户中心', max: 100 },
        { name: '支付系统', max: 100 },
        { name: '数据统计', max: 100 }
      ]
    },
    series: [{
      type: 'radar',
      data: [
        {
          value: [95, 90, 85, 90, 70],
          name: '功能完成度',
          itemStyle: {
            color: '#4caf50'
          },
          areaStyle: {
            color: 'rgba(76, 175, 80, 0.2)'
          }
        }
      ]
    }]
  }
  radarChart.setOption(option)
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
        center: ['50%', '40%'],
        radius: ['50%', '70%'],
        avoidLabelOverlap: true,
        label: {
          show: true,
          position: 'outside',
          formatter: '{b}: {d}%',
          distanceToLabelLine: 5
        },
        labelLine: {
          length: 15,
          length2: 10,
          smooth: true
        },
        emphasis: {
          label: {
            show: true,
            fontSize: '16',
            fontWeight: 'bold'
          }
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

// 初始化柱状图
const initBarChart = async () => {
  try {
    // 获取最近7天的订单数据，按商品类型统计
    const endDate = new Date()
    const startDate = new Date()
    startDate.setDate(startDate.getDate() - 6)
    
    const res = await getAdminOrderList({
      startDate: startDate.toISOString().split('T')[0],
      endDate: endDate.toISOString().split('T')[0],
      page: 1,
      pageSize: 999999
    })

    // 统计每天不同状态的订单数量
    const dailyStats = {
      pending: Array(7).fill(0),    // 待支付
      paid: Array(7).fill(0),       // 已支付
      completed: Array(7).fill(0)   // 已完成
    }

    if (res?.data?.list) {
      res.data.list.forEach(order => {
        const createTime = order.create_time || order.createTime || order.created_at
        if (createTime) {
          const orderDate = new Date(createTime)
          const dayIndex = 6 - Math.floor((endDate - orderDate) / (1000 * 60 * 60 * 24))
          if (dayIndex >= 0 && dayIndex < 7) {
            if (order.status === 0) dailyStats.pending[dayIndex]++
            else if (order.status === 1) dailyStats.paid[dayIndex]++
            else if (order.status === 2) dailyStats.completed[dayIndex]++
          }
        }
      })
    }

    // 生成日期标签
    const dateLabels = Array(7).fill().map((_, i) => {
      const date = new Date(startDate)
      date.setDate(date.getDate() + i)
      return `${date.getMonth() + 1}/${date.getDate()}`
    })

    barChart = echarts.init(barChartRef.value)
    const option = {
      title: {
        text: '订单状态统计',
        left: 'center',
        top: 20,
        textStyle: {
          fontSize: 14
        }
      },
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          type: 'shadow'
        }
      },
      legend: {
        bottom: 0,
        data: ['待支付', '已支付', '已完成']
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '15%',
        top: '15%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        data: dateLabels
      },
      yAxis: {
        type: 'value'
      },
      series: [
        {
          name: '待支付',
          type: 'bar',
          stack: 'total',
          data: dailyStats.pending,
          itemStyle: {
            color: '#ff9800'
          }
        },
        {
          name: '已支付',
          type: 'bar',
          stack: 'total',
          data: dailyStats.paid,
          itemStyle: {
            color: '#4caf50'
          }
        },
        {
          name: '已完成',
          type: 'bar',
          stack: 'total',
          data: dailyStats.completed,
          itemStyle: {
            color: '#2196f3'
          }
        }
      ]
    }
    barChart.setOption(option)
  } catch (error) {
    console.error('获取订单统计数据失败:', error)
    // 发生错误时显示空数据
    const dateLabels = Array(7).fill().map((_, i) => {
      const date = new Date()
      date.setDate(date.getDate() - 6 + i)
      return `${date.getMonth() + 1}/${date.getDate()}`
    })
    
    barChart = echarts.init(barChartRef.value)
    barChart.setOption({
      title: {
        text: '订单状态统计',
        left: 'center',
        top: 20,
        textStyle: {
          fontSize: 14
        }
      },
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          type: 'shadow'
        }
      },
      legend: {
        bottom: 0,
        data: ['待支付', '已支付', '已完成']
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '15%',
        top: '15%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        data: dateLabels
      },
      series: [
        {
          name: '待支付',
          type: 'bar',
          stack: 'total',
          data: Array(7).fill(0),
          itemStyle: { color: '#ff9800' }
        },
        {
          name: '已支付',
          type: 'bar',
          stack: 'total',
          data: Array(7).fill(0),
          itemStyle: { color: '#4caf50' }
        },
        {
          name: '已完成',
          type: 'bar',
          stack: 'total',
          data: Array(7).fill(0),
          itemStyle: { color: '#2196f3' }
        }
      ]
    })
  }
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
        data: products
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

    categoryChart = echarts.init(categoryChartRef.value)
    const option = {
      tooltip: {
        trigger: 'item',
        formatter: '{b}: {c} ({d}%)'
      },
      legend: {
        orient: 'vertical',
        right: 10,
        top: 'center'
      },
      series: [
        {
          name: '类别销售',
          type: 'pie',
          radius: ['50%', '70%'],
          avoidLabelOverlap: false,
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
          data: categories.map((category, index) => ({
            name: category,
            value: sales[index]
          }))
        }
      ]
    }
    categoryChart.setOption(option)
  } catch (error) {
    console.error('获取类别统计数据失败:', error)
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
  radarChart?.resize()
  pieChart?.resize()
  barChart?.resize()
  hotProductsChart?.resize()
  categoryChart?.resize()
  userBehaviorChart?.resize()
  userActivityChart?.resize()
}

onMounted(async () => {
  await fetchStatistics()
  initLineChart()
  await fetchOrderData()
  initRadarChart()
  initPieChart()
  initBarChart()
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
  radarChart?.dispose()
  pieChart?.dispose()
  barChart?.dispose()
  hotProductsChart?.dispose()
  categoryChart?.dispose()
  userBehaviorChart?.dispose()
  userActivityChart?.dispose()
})
</script>

<style scoped>
.dashboard-container {
  padding: 20px;
}

.data-overview {
  margin-bottom: 20px;
}

.data-card .card-header {
  display: flex;
  align-items: center;
}

.data-card .icon-wrapper {
  width: 48px;
  height: 48px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 16px;
}

.data-card .icon-wrapper .el-icon {
  font-size: 24px;
  color: #fff;
}

.data-card .icon-wrapper.blue {
  background-color: #2196f3;
}

.data-card .icon-wrapper.green {
  background-color: #4caf50;
}

.data-card .icon-wrapper.orange {
  background-color: #ff9800;
}

.data-card .icon-wrapper.purple {
  background-color: #9c27b0;
}

.data-card .data-info .label {
  font-size: 14px;
  color: #909399;
  margin-bottom: 4px;
}

.data-card .data-info .value {
  font-size: 24px;
  font-weight: bold;
  color: #303133;
}

.chart-row {
  margin-bottom: 20px;
}

.chart-row:last-child {
  margin-bottom: 0;
}

.chart-container {
  padding: 20px;
}

/* 路由过渡动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.filter-container {
  display: flex;
  align-items: center;
  padding: 10px;
}

.mr-3 {
  margin-right: 12px;
}

.mb-4 {
  margin-bottom: 16px;
}

.chart-title {
  font-size: 16px;
  color: #303133;
  margin-bottom: 20px;
  text-align: center;
  font-weight: bold;
}

.el-card {
  margin-bottom: 20px;
}

.el-card:last-child {
  margin-bottom: 0;
}

.chart-container {
  position: relative;
  padding: 20px;
  height: 100%;
}

/* 响应式布局 */
@media (max-width: 1200px) {
  .el-col {
    width: 100% !important;
  }
}

/* 图表加载状态 */
.chart-loading {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(255, 255, 255, 0.9);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1;
}
</style> 