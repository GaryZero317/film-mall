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
    // 模拟数据
    const mockData = [
      { date: '12-01', count: 120 },
      { date: '12-02', count: 132 },
      { date: '12-03', count: 101 },
      { date: '12-04', count: 134 },
      { date: '12-05', count: 90 },
      { date: '12-06', count: 230 },
      { date: '12-07', count: 210 }
    ]
    
    orderData.value = mockData
    updateLineChart()
  } catch (error) {
    console.error('获取订单数据失败:', error)
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
    statistics.products = productsRes.total.toLocaleString()

    // 获取订单总数
    const ordersRes = await getAdminOrderList({
      page: 1,
      pageSize: 1
    })
    statistics.orders = ordersRes.total.toLocaleString()

    // 获取支付列表计算总额
    const paymentsRes = await getAdminPaymentList({
      page: 1,
      pageSize: 999999 // 获取所有支付记录
    })
    const totalPayments = paymentsRes.list.reduce((sum, payment) => {
      return sum + (payment.status === 1 ? payment.amount : 0) // 只计算支付成功的订单
    }, 0)
    statistics.payments = `¥${totalPayments.toLocaleString()}`
  } catch (error) {
    console.error('获取统计数据失败:', error)
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
    radar: {
      indicator: [
        { name: '管理', max: 100 },
        { name: '市场', max: 100 },
        { name: '开发', max: 100 },
        { name: '支持', max: 100 },
        { name: '技术', max: 100 }
      ]
    },
    series: [{
      type: 'radar',
      data: [
        {
          value: [80, 70, 90, 85, 75],
          name: '预算分配',
          itemStyle: {
            color: '#4caf50'
          },
          areaStyle: {
            color: 'rgba(76, 175, 80, 0.2)'
          }
        },
        {
          value: [70, 75, 85, 80, 80],
          name: '实际开销',
          itemStyle: {
            color: '#2196f3'
          },
          areaStyle: {
            color: 'rgba(33, 150, 243, 0.2)'
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
      trigger: 'item'
    },
    legend: {
      orient: 'vertical',
      right: 10,
      top: 'center'
    },
    series: [
      {
        type: 'pie',
        radius: ['50%', '70%'],
        avoidLabelOverlap: false,
        label: {
          show: false
        },
        data: [
          { value: 35, name: '工业', itemStyle: { color: '#4caf50' } },
          { value: 25, name: '技术', itemStyle: { color: '#9c27b0' } },
          { value: 20, name: '金融', itemStyle: { color: '#2196f3' } },
          { value: 15, name: '预测', itemStyle: { color: '#ff9800' } },
          { value: 5, name: '其他', itemStyle: { color: '#f44336' } }
        ]
      }
    ]
  }
  pieChart.setOption(option)
}

// 初始化柱��图
const initBarChart = () => {
  barChart = echarts.init(barChartRef.value)
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
      type: 'category',
      data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日']
    },
    yAxis: {
      type: 'value'
    },
    series: [
      {
        name: '访问量',
        type: 'bar',
        stack: 'total',
        data: [320, 302, 301, 334, 390, 330, 320],
        itemStyle: {
          color: '#4caf50'
        }
      },
      {
        name: '下载量',
        type: 'bar',
        stack: 'total',
        data: [120, 132, 101, 134, 90, 230, 210],
        itemStyle: {
          color: '#2196f3'
        }
      },
      {
        name: '购买量',
        type: 'bar',
        stack: 'total',
        data: [220, 182, 191, 234, 290, 330, 310],
        itemStyle: {
          color: '#ff9800'
        }
      }
    ]
  }
  barChart.setOption(option)
}

// 监听窗口大小变化
const handleResize = () => {
  lineChart?.resize()
  radarChart?.resize()
  pieChart?.resize()
  barChart?.resize()
}

onMounted(async () => {
  await fetchStatistics() // ��取统计数据
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