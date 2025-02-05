<template>
  <div class="product-list-container">
    <div class="action-bar">
      <el-button type="primary" @click="handleAdd">
        <el-icon><Plus /></el-icon>添加商品
      </el-button>
    </div>

    <!-- 商品列表 -->
    <el-table :data="productList" style="width: 100%" v-loading="loading">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="商品名称" width="200" />
      <el-table-column prop="category_name" label="商品分类" width="150">
        <template #default="scope">
          {{ getCategoryName(scope.row.category_id) }}
        </template>
      </el-table-column>
      <el-table-column label="商品图片" width="120">
        <template #default="scope">
          <el-image
            v-if="scope.row.mainImage"
            :src="scope.row.mainImage"
            :preview-src-list="scope.row.images"
            fit="cover"
            class="product-image"
            style="width: 80px; height: 80px"
          />
          <el-empty v-else description="暂无图片" :image-size="40" />
        </template>
      </el-table-column>
      <el-table-column prop="desc" label="商品描述" show-overflow-tooltip />
      <el-table-column prop="stock" label="库存" width="100">
        <template #default="scope">
          <el-tag :type="scope.row.stock > 10 ? 'success' : 'danger'">
            {{ scope.row.stock }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="amount" label="价格" width="120">
        <template #default="scope">
          ¥{{ (scope.row.amount / 100).toFixed(2) }}
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template #default="scope">
          <el-tag :type="scope.row.status === 1 ? 'success' : 'info'">
            {{ scope.row.status === 1 ? '上架' : '下架' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="scope">
          <el-button type="primary" size="small" @click="handleEdit(scope.row)">
            编辑
          </el-button>
          <el-button type="danger" size="small" @click="handleDelete(scope.row)">
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 添加/编辑对话框 -->
    <el-dialog
      :title="dialogType === 'add' ? '添加商品' : '编辑商品'"
      v-model="dialogVisible"
      width="600px">
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px">
        <el-form-item label="商品名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入商品名称" />
        </el-form-item>
        <el-form-item label="商品分类" prop="category_id">
          <el-cascader
            v-model="form.category_id"
            :options="categoryOptions"
            :props="{
              value: 'id',
              label: 'name',
              children: 'children',
              checkStrictly: false,
              emitPath: false
            }"
            placeholder="请选择商品分类"
          />
        </el-form-item>
        <el-form-item label="商品描述" prop="desc">
          <el-input
            v-model="form.desc"
            type="textarea"
            :rows="3"
            placeholder="请输入商品描述" />
        </el-form-item>
        <el-form-item label="库存数量" prop="stock">
          <el-input-number
            v-model="form.stock"
            :min="0"
            :max="9999"
            controls-position="right" />
        </el-form-item>
        <el-form-item label="商品价格" prop="amount">
          <el-input-number
            v-model="form.amount"
            :min="0"
            :precision="2"
            :step="0.01"
            controls-position="right">
            <template #prefix>¥</template>
          </el-input-number>
        </el-form-item>
        <el-form-item label="商品状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio :value="1">上架</el-radio>
            <el-radio :value="0">下架</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="商品图片" prop="images">
          <el-upload
            class="product-image-uploader"
            :action="'/api/upload'"
            :headers="{
              Authorization: 'Bearer ' + userStore.token
            }"
            :show-file-list="true"
            :on-success="handleUploadSuccess"
            :on-error="handleUploadError"
            :before-upload="beforeUpload"
            name="file"
            accept="image/*"
            list-type="picture-card"
            multiple>
            <el-icon><Plus /></el-icon>
            <template #tip>
              <div class="el-upload__tip">
                只能上传jpg/png/gif文件，且不超过2MB
              </div>
            </template>
          </el-upload>
          <div v-if="form.images && form.images.length > 0" class="image-list">
            <div v-for="(image, index) in form.images" :key="index" class="image-item">
              <el-image
                :src="image"
                fit="cover"
                class="preview-image"
                :preview-src-list="[image]"
              />
              <div class="image-actions">
                <el-button
                  type="primary"
                  size="small"
                  :disabled="form.mainImage === image"
                  @click="setMainImage(image)">
                  设为主图
                </el-button>
                <el-button
                  type="danger"
                  size="small"
                  @click="removeImage(index)">
                  删除
                </el-button>
              </div>
            </div>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmit" :loading="submitting">
            确定
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { 
  createProduct, 
  updateProduct, 
  removeProduct, 
  getAdminProductList, 
  setMainImage as setProductMainImage,
  uploadImage,
  addProductImages,
  removeProductImages
} from '../../api/product'
import { useUserStore } from '../../stores/user'

const userStore = useUserStore()
const loading = ref(false)
const submitting = ref(false)
const dialogVisible = ref(false)
const dialogType = ref('add')
const productList = ref([])

// 表单相关
const formRef = ref(null)
const form = ref({
  name: '',
  desc: '',
  stock: 0,
  amount: 0,
  status: 1,
  images: [],
  mainImage: '',
  category_id: null
})

const rules = {
  name: [
    { required: true, message: '请输入商品名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  desc: [
    { required: true, message: '请输入商品描述', trigger: 'blur' }
  ],
  stock: [
    { required: true, message: '请输入库存数量', trigger: 'blur' }
  ],
  amount: [
    { required: true, message: '请输入商品价格', trigger: 'blur' }
  ],
  status: [
    { required: true, message: '请选择商品状态', trigger: 'change' }
  ],
  category_id: [
    { required: true, message: '请选择商品分类', trigger: 'change' }
  ]
}

// 分类选项
const categoryOptions = [
  {
    id: 1,
    name: '135胶卷',
    children: [
      { id: 4, name: '彩色负片' },
      { id: 5, name: '彩色反转片' },
      { id: 6, name: '黑白负片' }
    ]
  },
  {
    id: 2,
    name: '120胶卷',
    children: [
      { id: 7, name: '彩色负片' },
      { id: 8, name: '彩色反转片' }
    ]
  },
  {
    id: 3,
    name: '拍立得相纸',
    children: [
      { id: 9, name: '宝丽来相纸' },
      { id: 10, name: '富士相纸' }
    ]
  },
  {
    id: 11,
    name: '电影卷',
    children: [
      { id: 13, name: '彩色负片' },
      { id: 14, name: '彩色反转片' },
      { id: 15, name: '黑白负片' }
    ]
  }
]

// 获取分类名称
const getCategoryName = (categoryId) => {
  if (!categoryId) return '未分类'
  
  for (const mainCategory of categoryOptions) {
    for (const subCategory of mainCategory.children || []) {
      if (subCategory.id === categoryId) {
        return `${mainCategory.name} > ${subCategory.name}`
      }
    }
  }
  return '未分类'
}

// 获取商品列表
const fetchProductList = async () => {
  loading.value = true
  try {
    const res = await getAdminProductList({
      page: 1,
      pageSize: 100
    })
    console.log('商品列表完整响应:', res)
    console.log('商品列表data字段:', res.data)
    if (res.code === 0 && res.data) {
      console.log('商品列表总数:', res.data.total)
      console.log('商品列表详细数据:', JSON.stringify(res.data.list, null, 2))
      productList.value = res.data.list || []
      // 打印第一个商品的数据结构
      if (productList.value.length > 0) {
        console.log('第一个商品数据示例:', productList.value[0])
      }
    } else {
      ElMessage.error(res.msg || '获取商品列表失败')
    }
  } catch (error) {
    console.error('获取商品列表失败:', error)
    ElMessage.error('获取商品列表失败')
  } finally {
    loading.value = false
  }
}

// 添加商品
const handleAdd = () => {
  dialogType.value = 'add'
  form.value = {
    name: '',
    desc: '',
    stock: 0,
    amount: 0,
    status: 1,
    images: [],
    mainImage: '',
    category_id: null
  }
  dialogVisible.value = true
}

// 编辑商品
const handleEdit = (row) => {
  console.log('编辑商品:', row)
  dialogType.value = 'edit'
  form.value = {
    id: row.id,
    name: row.name,
    desc: row.desc,
    stock: row.stock,
    amount: row.amount / 100,
    status: row.status,
    images: row.images || [],
    mainImage: row.mainImage || '',
    category_id: row.category_id || null
  }
  dialogVisible.value = true
}

// 删除商品
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除该商品吗？', '提示', {
      type: 'warning'
    })
    
    await removeProduct({ id: row.id })
    ElMessage.success('删除成功')
    fetchProductList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除商品失败:', error)
      ElMessage.error('删除失败')
    }
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) {
    console.error('表单引用为空')
    return
  }
  
  try {
    console.log('开始验证表单')
    await formRef.value.validate()
    submitting.value = true
    
    const submitData = {
      ...form.value,
      amount: Math.round(form.value.amount * 100), // 转换为分
      imageUrls: form.value.images || [],
      mainImage: form.value.mainImage || '',
      category_id: parseInt(form.value.category_id) || 0  // 确保category_id是整数
    }
    console.log('提交的数据:', submitData)
    
    let res
    if (dialogType.value === 'add') {
      console.log('执行添加操作')
      res = await createProduct(submitData)
      console.log('创建商品响应:', res)
    } else {
      console.log('执行更新操作')
      res = await updateProduct(submitData)
      console.log('更新商品响应:', res)
      
      // 只有在明确成功的情况下才继续处理
      if (res && res.code === 0) {
        try {
          const product = productList.value.find(p => p.id === form.value.id)
          if (product && product.mainImage !== form.value.mainImage) {
            await setProductMainImage({
              productId: form.value.id,
              imageUrl: form.value.mainImage
            })
          }
        } catch (error) {
          console.error('设置主图失败:', error)
          ElMessage.warning('商品更新成功，但设置主图失败')
        }
      }
    }
    console.log('服务器响应:', res)

    // 更新响应处理逻辑
    if (res && (res.code === 0 || Object.keys(res).length === 0)) {
      ElMessage.success(dialogType.value === 'add' ? '添加成功' : '更新成功')
      dialogVisible.value = false
      fetchProductList()
    } else {
      throw new Error(res?.msg || '操作失败')
    }
  } catch (error) {
    console.error('提交失败，详细错误:', error)
    if (error.response) {
      console.error('服务器响应:', error.response)
      ElMessage.error(error.response.data || '服务器错误')
    } else if (error.request) {
      console.error('请求错误:', error.request)
      ElMessage.error('网络请求失败，请检查网络连接')
    } else {
      console.error('其他错误:', error)
      ElMessage.error(error?.message || '提交失败')
    }
  } finally {
    submitting.value = false
  }
}

// 图片上传相关方法
const handleUploadSuccess = async (response, uploadFile) => {
  if (response.code === 0 && response.data) {
    const imageUrl = response.data.url || response.data
    console.log('上传成功，图片URL:', imageUrl)
    form.value.images = [...(form.value.images || []), imageUrl]
    
    // 如果是编辑模式，调用添加商品图片接口
    if (dialogType.value === 'edit' && form.value.id) {
      try {
        console.log('调用添加商品图片接口:', {
          productId: form.value.id,
          imageUrls: [imageUrl]
        })
        const res = await addProductImages({
          productId: form.value.id,
          imageUrls: [imageUrl]
        })
        console.log('添加商品图片响应:', res)
        
        if (!res || Object.keys(res).length === 0 || res.code === 0) {
          // 如果还没有主图，设置为主图
          if (!form.value.mainImage) {
            await setProductMainImage({
              productId: form.value.id,
              imageUrl: imageUrl
            })
            form.value.mainImage = imageUrl
          }
          ElMessage.success('上传成功')
          fetchProductList() // 刷新列表
        } else {
          ElMessage.error(res.msg || '添加商品图片失败')
        }
      } catch (error) {
        console.error('添加商品图片失败:', error)
        ElMessage.error('添加商品图片失败')
      }
    } else {
      // 如果是新增模式，直接设置主图
      if (!form.value.mainImage) {
        form.value.mainImage = imageUrl
      }
      ElMessage.success('上传成功')
    }
  } else {
    ElMessage.error(response.msg || '上传失败')
  }
}

const handleUploadError = (error) => {
  console.error('上传失败:', error)
  ElMessage.error('上传失败')
}

const beforeUpload = (file) => {
  const isImage = file.type.startsWith('image/')
  const isLt2M = file.size / 1024 / 1024 < 2

  if (!isImage) {
    ElMessage.error('只能上传图片文件!')
    return false
  }
  if (!isLt2M) {
    ElMessage.error('图片大小不能超过 2MB!')
    return false
  }
  return true
}

const setMainImage = async (image) => {
  try {
    if (dialogType.value === 'edit' && form.value.id) {
      const params = {
        productId: form.value.id,
        imageUrl: image
      }
      console.log('设置主图请求参数:', params)
      const res = await setProductMainImage(params)
      // 空对象响应也视为成功
      if (!res || Object.keys(res).length === 0 || res.code === 0) {
        form.value.mainImage = image
        // 更新商品列表中的主图
        const product = productList.value.find(p => p.id === form.value.id)
        if (product) {
          product.mainImage = image
        }
        ElMessage.success('设置主图成功')
      } else {
        ElMessage.error(res.msg || '设置主图失败')
      }
    } else {
      form.value.mainImage = image
    }
  } catch (error) {
    console.error('设置主图失败:', error)
    ElMessage.error('设置主图失败')
  }
}

const removeImage = async (index) => {
  const image = form.value.images[index]
  
  // 如果是编辑模式，调用删除图片接口
  if (dialogType.value === 'edit' && form.value.id) {
    try {
      console.log('调用删除商品图片接口:', {
        productId: form.value.id,
        imageUrls: [image]
      })
      const res = await removeProductImages({
        productId: form.value.id,
        imageUrls: [image]
      })
      console.log('删除商品图片响应:', res)
      
      if (!res || Object.keys(res).length === 0 || res.code === 0) {
        // 如果删除的是主图，清空主图
        if (form.value.mainImage === image) {
          form.value.mainImage = form.value.images[0] || ''
          if (form.value.mainImage) {
            await setProductMainImage({
              productId: form.value.id,
              imageUrl: form.value.mainImage
            })
          }
        }
        // 从当前表单中移除图片
        form.value.images.splice(index, 1)
        // 更新商品列表中的图片
        const product = productList.value.find(p => p.id === form.value.id)
        if (product) {
          product.images = [...form.value.images]
          product.mainImage = form.value.mainImage
        }
        ElMessage.success('删除成功')
      } else {
        ElMessage.error(res.msg || '删除图片失败')
      }
    } catch (error) {
      console.error('删除商品图片失败:', error)
      ElMessage.error('删除图片失败')
    }
  } else {
    // 如果是新增模式，直接从数组中移除
    form.value.images.splice(index, 1)
    if (form.value.mainImage === image) {
      form.value.mainImage = form.value.images[0] || ''
    }
  }
}

onMounted(() => {
  fetchProductList()
})
</script>

<style scoped>
.product-list-container {
  padding: 20px;
}

.action-bar {
  margin-bottom: 20px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

:deep(.el-input-number .el-input__wrapper) {
  width: 200px;
}

.product-image-uploader {
  margin-bottom: 20px;
}

.image-list {
  display: flex;
  flex-wrap: wrap;
  gap: 20px;
}

.image-item {
  position: relative;
  width: 200px;
}

.preview-image {
  width: 100%;
  height: 200px;
  border-radius: 4px;
  object-fit: cover;
}

.image-actions {
  margin-top: 8px;
  display: flex;
  gap: 8px;
  justify-content: center;
}

.product-image {
  border-radius: 4px;
  cursor: pointer;
}
</style> 