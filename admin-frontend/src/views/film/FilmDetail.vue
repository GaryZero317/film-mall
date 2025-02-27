<template>
  <div class="container">
    <div class="page-header">
      <h2>胶片冲洗订单详情</h2>
      <el-button @click="$router.push('/film/list')">返回列表</el-button>
    </div>

    <el-card v-loading="loading" class="info-card">
      <template #header>
        <div class="card-header">
          <span>订单基本信息</span>
        </div>
      </template>

      <el-descriptions :column="2" border>
        <el-descriptions-item label="订单ID">{{ orderInfo.id }}</el-descriptions-item>
        <el-descriptions-item label="订单号">{{ orderInfo.foid }}</el-descriptions-item>
        <el-descriptions-item label="用户ID">{{ orderInfo.uid }}</el-descriptions-item>
        <el-descriptions-item label="地址ID">{{ orderInfo.address_id }}</el-descriptions-item>
        <el-descriptions-item label="订单状态">
          <el-tag :type="getStatusTagType(orderInfo.status)">
            {{ orderInfo.status_desc }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="回寄底片">
          <el-tag :type="orderInfo.return_film ? 'success' : 'info'">
            {{ orderInfo.return_film ? '是' : '否' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="总价">
          {{ orderInfo.total_price ? (orderInfo.total_price / 100).toFixed(2) + ' 元' : '0.00 元' }}
        </el-descriptions-item>
        <el-descriptions-item label="运费">
          {{ orderInfo.shipping_fee ? (orderInfo.shipping_fee / 100).toFixed(2) + ' 元' : '0.00 元' }}
        </el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ orderInfo.create_time }}</el-descriptions-item>
        <el-descriptions-item label="更新时间">{{ orderInfo.update_time }}</el-descriptions-item>
        <el-descriptions-item label="备注" :span="2">{{ orderInfo.remark || '无' }}</el-descriptions-item>
      </el-descriptions>
    </el-card>

    <el-card class="info-card">
      <template #header>
        <div class="card-header">
          <span>已冲洗照片</span>
          <div>
            <el-button type="primary" size="small" @click="handleBatchUpload" :disabled="orderInfo.status < 1">
              批量上传照片
            </el-button>
            <input
              type="file"
              ref="batchUploadInput"
              multiple
              accept="image/*"
              style="display: none"
              @change="onBatchFileSelected"
            />
          </div>
        </div>
      </template>

      <div v-if="photoLoading" class="photo-loading">
        <el-skeleton :rows="5" animated />
      </div>
      <div v-else-if="photos.length === 0" class="no-photos">
        <el-empty description="暂无照片，请上传" />
      </div>
      <div v-else class="photo-grid">
        <div v-for="(photo, index) in photos" :key="index" class="photo-item">
          <el-image 
            :src="getFullPhotoUrl(photo.url)" 
            fit="cover" 
            :preview-src-list="getPhotoUrlList()"
            :initial-index="index"
            @error="handleImageLoadError"
          >
            <template #error>
              <div class="image-error">
                <el-icon><Picture /></el-icon>
                <span>图片加载失败</span>
              </div>
            </template>
          </el-image>
          <div class="photo-actions">
            <el-button 
              type="danger" 
              size="small" 
              icon="Delete" 
              circle
              @click="removePhoto(photo.id)"
            />
          </div>
        </div>
      </div>

      <div v-if="uploadingCount > 0" class="upload-progress">
        <el-progress 
          :percentage="uploadProgress" 
          :format="progressFormat"
          status="success"
        />
      </div>
    </el-card>

    <el-card class="info-card">
      <template #header>
        <div class="card-header">
          <span>订单项信息</span>
        </div>
      </template>

      <el-table :data="orderInfo.items || []" border style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="film_type" label="胶片类型" />
        <el-table-column prop="film_brand" label="胶片品牌" />
        <el-table-column prop="size" label="尺寸规格" />
        <el-table-column prop="quantity" label="数量" width="80" />
        <el-table-column label="单价" width="120">
          <template #default="scope">
            {{ (scope.row.price / 100).toFixed(2) }} 元
          </template>
        </el-table-column>
        <el-table-column label="总价" width="120">
          <template #default="scope">
            {{ (scope.row.amount / 100).toFixed(2) }} 元
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" />
      </el-table>
    </el-card>

    <el-card class="info-card">
      <template #header>
        <div class="card-header">
          <span>订单操作</span>
        </div>
      </template>

      <el-form :model="form" label-width="120px">
        <el-form-item label="订单状态">
          <el-select v-model="form.status">
            <el-option v-for="item in statusOptions" :key="item.value" :label="item.label" :value="item.value"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="回寄底片">
          <el-switch v-model="form.return_film" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.remark" type="textarea" :rows="3"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleUpdate">更新订单</el-button>
          <el-button v-if="orderInfo.status === 3" type="danger" @click="handleDelete">删除订单</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Picture } from '@element-plus/icons-vue'
import { getFilmOrderDetail, updateFilmOrder, deleteFilmOrder } from '@/api/film'
import { uploadFilmPhoto, getFilmPhotos, deleteFilmPhoto } from '@/api/film'

const route = useRoute()
const router = useRouter()
const orderId = ref(route.params.id)
const loading = ref(false)
const orderInfo = ref({})
const batchUploadInput = ref(null)

// 表单数据
const form = reactive({
  status: 0,
  return_film: false,
  remark: ''
})

// 订单状态选项
const statusOptions = [
  { value: 0, label: '待付款' },
  { value: 1, label: '冲洗处理中' },
  { value: 2, label: '待收货' },
  { value: 3, label: '已完成' }
]

// 照片相关
const photos = ref([])
const photoLoading = ref(false)
const uploadingCount = ref(0)
const uploadedCount = ref(0)
const uploadProgress = ref(0)

// 获取订单详情
const getOrderDetail = async () => {
  if (!orderId.value) return
  
  loading.value = true
  try {
    const res = await getFilmOrderDetail(orderId.value)
    if (res.code === 0 && res.data) {
      orderInfo.value = res.data
      form.status = res.data.status
      form.return_film = res.data.return_film
      form.remark = res.data.remark || ''
      // 加载照片
      loadPhotos()
    } else {
      ElMessage.error(res.msg || '获取订单详情失败')
    }
  } catch (error) {
    console.error('获取胶片冲洗订单详情失败:', error)
    ElMessage.error('获取订单详情失败')
  } finally {
    loading.value = false
  }
}

// 更新订单
const handleUpdate = async () => {
  try {
    const res = await updateFilmOrder(orderInfo.value.id, {
      status: form.status,
      return_film: form.return_film,
      remark: form.remark
    })
    
    if (res.code === 0) {
      ElMessage.success('更新订单成功')
      getOrderDetail()
    } else {
      ElMessage.error(res.msg || '更新订单失败')
    }
  } catch (error) {
    console.error('更新订单出错:', error)
    ElMessage.error('更新订单出错')
  }
}

// 删除订单
const handleDelete = () => {
  ElMessageBox.confirm('确定要删除该订单吗？', '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      const res = await deleteFilmOrder(orderInfo.value.id)
      if (res.code === 0) {
        ElMessage.success('删除成功')
        router.push('/film/list')
      } else {
        ElMessage.error(res.msg || '删除失败')
      }
    } catch (error) {
      console.error('删除订单出错:', error)
      ElMessage.error('删除订单出错')
    }
  }).catch(() => {
    // 取消删除
  })
}

// 获取状态标签类型
const getStatusTagType = (status) => {
  const statusMap = {
    0: 'info',
    1: 'warning',
    2: 'primary',
    3: 'success'
  }
  return statusMap[status] || 'info'
}

// 加载照片
const loadPhotos = async () => {
  if (!orderId.value) return
  
  photoLoading.value = true
  try {
    const res = await getFilmPhotos(orderId.value)
    if (res.code === 0) {
      // 检查API返回的数据结构
      if (res.data && res.data.list) {
        photos.value = res.data.list // 处理包含list字段的情况
      } else {
        photos.value = res.data || []
      }
    } else {
      ElMessage.error(res.msg || '获取照片列表失败')
    }
  } catch (error) {
    console.error('获取照片列表失败:', error)
    ElMessage.error('获取照片列表失败')
  } finally {
    photoLoading.value = false
  }
}

// 批量上传
const handleBatchUpload = () => {
  if (orderInfo.value.status < 1) {
    ElMessage.warning('订单状态为冲洗处理中才能上传照片')
    return
  }
  
  batchUploadInput.value.click()
}

// 文件选择处理
const onBatchFileSelected = async (event) => {
  const files = event.target.files
  if (!files || files.length === 0) return
  
  // 验证文件类型和大小
  const validFiles = []
  for (let i = 0; i < files.length; i++) {
    const file = files[i]
    // 检查是否是图片
    if (!file.type.startsWith('image/')) {
      ElMessage.warning(`${file.name} 不是有效的图片文件`)
      continue
    }
    
    // 检查大小限制 (5MB)
    if (file.size > 5 * 1024 * 1024) {
      ElMessage.warning(`${file.name} 大小超过5MB限制`)
      continue
    }
    
    validFiles.push(file)
  }
  
  if (validFiles.length === 0) {
    ElMessage.error('没有有效的图片文件可上传')
    return
  }
  
  // 开始上传
  uploadingCount.value = validFiles.length
  uploadedCount.value = 0
  uploadProgress.value = 0
  
  // 并发上传
  const uploadPromises = validFiles.map(file => {
    const formData = new FormData()
    formData.append('file', file)
    formData.append('film_order_id', orderId.value)
    
    return uploadFilmPhoto(formData)
      .then(res => {
        if (res.code === 0) {
          uploadedCount.value++
          uploadProgress.value = Math.floor((uploadedCount.value / uploadingCount.value) * 100)
          return res.data
        } else {
          throw new Error(res.msg || '上传失败')
        }
      })
  })
  
  try {
    await Promise.all(uploadPromises)
    ElMessage.success('照片上传成功')
    // 重新加载照片列表
    loadPhotos()
  } catch (error) {
    console.error('上传照片失败:', error)
    ElMessage.error('部分照片上传失败')
  } finally {
    // 重置上传状态
    setTimeout(() => {
      uploadingCount.value = 0
      uploadProgress.value = 0
    }, 2000)
    // 重置文件输入框
    batchUploadInput.value.value = ''
  }
}

// 删除照片
const removePhoto = async (photoId) => {
  try {
    await ElMessageBox.confirm('确定要删除这张照片吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    const res = await deleteFilmPhoto(photoId)
    if (res.code === 0) {
      ElMessage.success('删除成功')
      // 从列表中移除
      photos.value = photos.value.filter(photo => photo.id !== photoId)
    } else {
      ElMessage.error(res.msg || '删除失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除照片失败:', error)
      ElMessage.error('删除照片失败')
    }
  }
}

// 获取所有照片URL列表，用于预览大图
const getPhotoUrlList = () => {
  // 确保 photos.value 是一个数组
  if (!photos.value) return []
  // 检查是否是包含 list 的对象（API返回格式）
  if (photos.value.list && Array.isArray(photos.value.list)) {
    return photos.value.list.map(photo => getFullPhotoUrl(photo.url))
  }
  // 直接是数组的情况
  if (Array.isArray(photos.value)) {
    return photos.value.map(photo => getFullPhotoUrl(photo.url))
  }
  // 不是数组，返回空数组
  return []
}

// 获取完整图片URL
const getFullPhotoUrl = (url) => {
  if (!url) return ''
  
  // 已经是完整URL的情况
  if (url.startsWith('http')) return url
  
  // 明确指定后端API地址
  // 使用film服务的正确端口8007
  const backendUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8007'
  
  // 确保url指向正确的路径格式（即指向服务器上的实际文件位置）
  let fileUrl = url
  // 如果url不是以'/uploads/'开头，则添加
  if (!url.startsWith('/uploads/')) {
    fileUrl = '/uploads/' + url.replace(/^\//, '')
  }
  
  console.log('构建图片URL:', backendUrl + fileUrl)
  return backendUrl + fileUrl
}

// 格式化进度条显示文本
const progressFormat = (percentage) => {
  return `${uploadedCount.value}/${uploadingCount.value} (${percentage}%)`
}

// 处理图片加载错误
const handleImageLoadError = (e) => {
  console.error('图片加载失败:', e)
}

onMounted(() => {
  getOrderDetail()
})
</script>

<style scoped>
.container {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.info-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.photo-loading {
  padding: 20px;
}

.no-photos {
  padding: 20px;
}

.photo-grid {
  display: flex;
  flex-wrap: wrap;
}

.photo-item {
  width: 200px;
  height: 200px;
  margin: 10px;
  position: relative;
}

.photo-item img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.photo-actions {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  opacity: 0;
  transition: opacity 0.3s;
}

.photo-item:hover .photo-actions {
  opacity: 1;
}

.upload-progress {
  margin-top: 20px;
}

.image-error {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  padding: 10px;
}
</style> 