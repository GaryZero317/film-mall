<template>
  <div class="dashboard-container">
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
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, reactive } from 'vue'
import { ShoppingBag, List, Money } from '@element-plus/icons-vue'
import * as echarts from 'echarts'
import { getAdminProductList } from '@/api/product'
import { getAdminOrderList } from '@/api/order'
import { getAdminPaymentList } from '@/api/payment'
import { ElMessage } from 'element-plus'

// 统计数据
const statistics = reactive({
  products: '0',  // 商品总数
  orders: '0',    // 订单总数
  payments: '0'   // 支付总额
})

// 订单数据
const orderData = ref([])

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

// 监听窗口大小变化
const handleResize = () => {
  lineChart?.resize()
  radarChart?.resize()
  pieChart?.resize()
  barChart?.resize()
}

onMounted(async () => {
  await fetchStatistics() // 获取统计数据
  initLineChart()
  await fetchOrderData()
  initRadarChart()
  initPieChart()
  initBarChart()
  window.addEventListener('resize', handleResize)
})

// 组件卸载时移除事件监听
onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  lineChart?.dispose()
  radarChart?.dispose()
  pieChart?.dispose()
  barChart?.dispose()
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
</style> 