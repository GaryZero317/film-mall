<template>
  <div class="container">
    <div class="page-header">
      <h2>社区作品管理</h2>
    </div>
    
    <div class="toolbar">
      <el-form :inline="true" :model="queryParams" class="demo-form-inline">
        <el-form-item label="用户ID">
          <el-input v-model="queryParams.uid" placeholder="用户ID" clearable></el-input>
        </el-form-item>
        <el-form-item label="作品状态">
          <el-select v-model="queryParams.status" placeholder="全部状态" clearable>
            <el-option v-for="item in statusOptions" :key="item.value" :label="item.label" :value="item.value"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="胶片类型">
          <el-input v-model="queryParams.film_type" placeholder="胶片类型" clearable></el-input>
        </el-form-item>
        <el-form-item label="胶片品牌">
          <el-input v-model="queryParams.film_brand" placeholder="胶片品牌" clearable></el-input>
        </el-form-item>
        <el-form-item label="关键词">
          <el-input v-model="queryParams.keyword" placeholder="标题/描述关键词" clearable></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleQuery">搜索</el-button>
          <el-button @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <el-table v-loading="loading" :data="workList" border style="width: 100%">
      <el-table-column prop="id" label="作品ID" width="80" />
      <el-table-column prop="uid" label="用户ID" width="80" />
      <el-table-column prop="title" label="作品标题" width="150" />
      <el-table-column prop="description" label="描述" width="200" show-overflow-tooltip />
      <el-table-column label="封面图" width="120" align="center">
        <template #default="scope">
          <el-image 
            v-if="scope.row.cover_url" 
            :src="scope.row.cover_url" 
            :preview-src-list="[scope.row.cover_url]"
            fit="contain"
            style="width: 80px; height: 80px">
          </el-image>
          <span v-else>无封面</span>
        </template>
      </el-table-column>
      <el-table-column prop="film_type" label="胶片类型" width="100" />
      <el-table-column prop="film_brand" label="胶片品牌" width="100" />
      <el-table-column label="互动数据" width="180">
        <template #default="scope">
          <div>浏览：{{ scope.row.view_count }}</div>
          <div>点赞：{{ scope.row.like_count }}</div>
          <div>评论：{{ scope.row.comment_count }}</div>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="100">
        <template #default="scope">
          <el-tag :type="getStatusTagType(scope.row.status)">
            {{ getStatusText(scope.row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="create_time" label="创建时间" width="180" />
      <el-table-column label="操作" width="200">
        <template #default="scope">
          <el-button size="small" type="primary" @click="viewWork(scope.row)">查看</el-button>
          <el-button size="small" type="warning" @click="editWork(scope.row)">编辑</el-button>
          <el-button size="small" type="danger" @click="handleDelete(scope.row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <div class="pagination-container">
      <el-pagination
        :current-page="queryParams.page"
        :page-size="queryParams.page_size"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        :total="total"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>

    <!-- 编辑作品对话框 -->
    <el-dialog v-model="dialogVisible" title="编辑社区作品" width="500px">
      <el-form ref="formRef" :model="form" label-width="120px">
        <el-form-item label="作品ID">
          <span>{{ form.id }}</span>
        </el-form-item>
        <el-form-item label="用户ID">
          <span>{{ form.uid }}</span>
        </el-form-item>
        <el-form-item label="作品标题">
          <el-input v-model="form.title"></el-input>
        </el-form-item>
        <el-form-item label="作品描述">
          <el-input v-model="form.description" type="textarea" rows="3"></el-input>
        </el-form-item>
        <el-form-item label="作品状态">
          <el-select v-model="form.status">
            <el-option v-for="item in statusOptions" :key="item.value" :label="item.label" :value="item.value"></el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitForm">确认</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useRouter } from 'vue-router'
import { getWorkList, getWorkDetail, updateWork, deleteWork } from '../../api/community'

const router = useRouter()
const loading = ref(false)
const workList = ref([])
const total = ref(0)
const dialogVisible = ref(false)
const formRef = ref(null)

// 查询条件
const queryParams = reactive({
  page: 1,
  page_size: 10,
  uid: '',
  status: '',
  film_type: '',
  film_brand: '',
  keyword: ''
})

// 表单数据
const form = reactive({
  id: '',
  uid: '',
  title: '',
  description: '',
  status: 0
})

// 状态选项
const statusOptions = [
  { value: 0, label: '草稿' },
  { value: 1, label: '已发布' },
  { value: 2, label: '已删除' }
]

// 获取状态标签类型
const getStatusTagType = (status) => {
  if (status === 0) return 'info'
  if (status === 1) return 'success'
  if (status === 2) return 'danger'
  return 'info'
}

// 获取状态文本
const getStatusText = (status) => {
  if (status === 0) return '草稿'
  if (status === 1) return '已发布'
  if (status === 2) return '已删除'
  return '未知'
}

// 查询作品列表
const handleQuery = () => {
  loading.value = true
  getWorkList(queryParams).then(res => {
    // 处理数据结构，将work对象中的数据展开到顶层
    const formattedList = res.data.list.map(item => {
      return {
        ...item.work,
        author: item.author
      }
    })
    workList.value = formattedList
    total.value = res.data.total
  }).catch(err => {
    console.error('获取作品列表失败', err)
    ElMessage.error('获取作品列表失败')
  }).finally(() => {
    loading.value = false
  })
}

// 重置查询条件
const resetQuery = () => {
  Object.keys(queryParams).forEach(key => {
    if (key !== 'page' && key !== 'page_size') {
      queryParams[key] = ''
    }
  })
  queryParams.page = 1
  handleQuery()
}

// 页码变化
const handleSizeChange = (newSize) => {
  queryParams.page_size = newSize
  handleQuery()
}

const handleCurrentChange = (newPage) => {
  queryParams.page = newPage
  handleQuery()
}

// 查看作品详情
const viewWork = (row) => {
  if (row && row.id) {
    console.log('跳转到作品详情页面，ID:', row.id)
    router.push({
      name: 'CommunityWorkDetail',
      params: { id: row.id }
    })
  } else {
    ElMessage.error('作品ID无效')
  }
}

// 编辑作品
const editWork = (row) => {
  // 填充表单数据
  Object.keys(form).forEach(key => {
    if (key in row) {
      form[key] = row[key]
    }
  })
  dialogVisible.value = true
}

// 提交表单
const submitForm = () => {
  updateWork(form.id, {
    title: form.title,
    description: form.description,
    status: form.status
  }).then(() => {
    ElMessage.success('更新作品成功')
    dialogVisible.value = false
    handleQuery()
  }).catch(err => {
    console.error('更新作品失败', err)
    ElMessage.error('更新作品失败')
  })
}

// 删除作品
const handleDelete = (row) => {
  ElMessageBox.confirm(
    '确定要删除该作品吗？删除后不可恢复。',
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(() => {
    deleteWork(row.id).then(() => {
      ElMessage.success('删除成功')
      handleQuery()
    }).catch(err => {
      console.error('删除作品失败', err)
      ElMessage.error('删除作品失败')
    })
  }).catch(() => {
    // 取消删除
  })
}

// 页面加载时获取数据
onMounted(() => {
  handleQuery()
})
</script>

<style scoped>
.container {
  padding: 20px;
}

.page-header {
  margin-bottom: 20px;
}

.toolbar {
  margin-bottom: 20px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style> 