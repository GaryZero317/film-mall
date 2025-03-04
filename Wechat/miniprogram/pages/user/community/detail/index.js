import { getWorkDetail, likeWork, getComments, addComment, deleteComment, deleteWork } from '../../../../api/community'

Page({
  data: {
    id: null,
    work: null,
    comments: [],
    loading: false,
    commentLoading: false,
    commentContent: '',
    isOwner: false,
    showActionSheet: false,
    actions: [
      { name: '编辑', color: '#1989fa' },
      { name: '删除', color: '#ee0a24' }
    ]
  },

  onLoad(options) {
    const id = options.id
    if (id) {
      this.setData({ id })
      this.loadWorkDetail()
      this.loadComments()
    } else {
      wx.showToast({
        title: '参数错误',
        icon: 'none'
      })
      setTimeout(() => {
        wx.navigateBack()
      }, 1500)
    }
  },

  // 加载作品详情
  async loadWorkDetail() {
    this.setData({ loading: true })
    
    try {
      const res = await getWorkDetail(this.data.id)
      console.log('作品详情:', res)
      
      if (res.code === 0) {
        // 检查是否是作品所有者
        const userInfo = wx.getStorageSync('userInfo')
        const isOwner = userInfo && userInfo.id === res.data.work.uid
        
        this.setData({
          work: res.data.work,
          isOwner
        })
      } else {
        wx.showToast({
          title: res.msg || '加载失败',
          icon: 'none'
        })
      }
    } catch (error) {
      console.error('加载作品详情失败:', error)
      wx.showToast({
        title: '网络错误，请重试',
        icon: 'none'
      })
    } finally {
      this.setData({ loading: false })
    }
  },

  // 加载评论
  async loadComments() {
    this.setData({ commentLoading: true })
    
    try {
      const res = await getComments(this.data.id)
      console.log('评论列表:', res)
      
      if (res.code === 0) {
        this.setData({
          comments: res.data
        })
      } else {
        wx.showToast({
          title: res.msg || '加载评论失败',
          icon: 'none'
        })
      }
    } catch (error) {
      console.error('加载评论失败:', error)
    } finally {
      this.setData({ commentLoading: false })
    }
  },

  // 点赞/取消点赞
  async onLike() {
    if (!this.data.work) return
    
    try {
      const isLiked = this.data.work.like_status
      const res = await likeWork(this.data.id, !isLiked)
      
      if (res.code === 0) {
        const work = { ...this.data.work }
        work.like_status = !isLiked
        work.like_count = isLiked ? work.like_count - 1 : work.like_count + 1
        
        this.setData({ work })
        
        wx.showToast({
          title: isLiked ? '已取消点赞' : '点赞成功',
          icon: 'success'
        })
      } else {
        wx.showToast({
          title: res.msg || '操作失败',
          icon: 'none'
        })
      }
    } catch (error) {
      console.error('点赞操作失败:', error)
      wx.showToast({
        title: '网络错误，请重试',
        icon: 'none'
      })
    }
  },

  // 提交评论
  async onSubmitComment() {
    if (!this.data.commentContent.trim()) {
      wx.showToast({
        title: '请输入评论内容',
        icon: 'none'
      })
      return
    }
    
    try {
      const res = await addComment(this.data.id, this.data.commentContent)
      
      if (res.code === 0) {
        this.setData({ commentContent: '' })
        this.loadComments()
        
        // 更新评论数
        const work = { ...this.data.work }
        work.comment_count += 1
        this.setData({ work })
        
        wx.showToast({
          title: '评论成功',
          icon: 'success'
        })
      } else {
        wx.showToast({
          title: res.msg || '评论失败',
          icon: 'none'
        })
      }
    } catch (error) {
      console.error('提交评论失败:', error)
      wx.showToast({
        title: '网络错误，请重试',
        icon: 'none'
      })
    }
  },

  // 删除评论
  async onDeleteComment(e) {
    const commentId = e.currentTarget.dataset.id
    
    wx.showModal({
      title: '提示',
      content: '确定要删除这条评论吗？',
      success: async (res) => {
        if (res.confirm) {
          try {
            const res = await deleteComment(commentId)
            
            if (res.code === 0) {
              this.loadComments()
              
              // 更新评论数
              const work = { ...this.data.work }
              work.comment_count -= 1
              this.setData({ work })
              
              wx.showToast({
                title: '删除成功',
                icon: 'success'
              })
            } else {
              wx.showToast({
                title: res.msg || '删除失败',
                icon: 'none'
              })
            }
          } catch (error) {
            console.error('删除评论失败:', error)
            wx.showToast({
              title: '网络错误，请重试',
              icon: 'none'
            })
          }
        }
      }
    })
  },

  // 显示操作菜单
  onShowActions() {
    this.setData({ showActionSheet: true })
  },

  // 关闭操作菜单
  onCloseActions() {
    this.setData({ showActionSheet: false })
  },

  // 处理操作菜单选择
  onSelectAction(e) {
    const index = e.detail.index
    
    if (index === 0) {
      // 编辑
      wx.navigateTo({
        url: `/pages/user/community/create/index?id=${this.data.id}`
      })
    } else if (index === 1) {
      // 删除
      this.onDeleteWork()
    }
    
    this.setData({ showActionSheet: false })
  },

  // 删除作品
  async onDeleteWork() {
    wx.showModal({
      title: '提示',
      content: '确定要删除这个作品吗？删除后无法恢复。',
      success: async (res) => {
        if (res.confirm) {
          try {
            const res = await deleteWork(this.data.id)
            
            if (res.code === 0) {
              wx.showToast({
                title: '删除成功',
                icon: 'success'
              })
              
              setTimeout(() => {
                wx.navigateBack()
              }, 1500)
            } else {
              wx.showToast({
                title: res.msg || '删除失败',
                icon: 'none'
              })
            }
          } catch (error) {
            console.error('删除作品失败:', error)
            wx.showToast({
              title: '网络错误，请重试',
              icon: 'none'
            })
          }
        }
      }
    })
  },

  // 输入评论内容
  onInputComment(e) {
    this.setData({
      commentContent: e.detail.value
    })
  },

  // 预览图片
  onPreviewImage(e) {
    const current = e.currentTarget.dataset.src
    const urls = this.data.work.images.map(img => img.url)
    
    wx.previewImage({
      current,
      urls
    })
  },

  // 返回
  onBack() {
    wx.navigateBack()
  },

  // 分享
  onShareAppMessage() {
    const work = this.data.work
    if (!work) return {}
    
    return {
      title: work.title,
      path: `/pages/user/community/detail/index?id=${this.data.id}`,
      imageUrl: work.cover_url
    }
  }
}) 