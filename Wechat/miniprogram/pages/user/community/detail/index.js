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
    if (!options || !options.id) {
      wx.showToast({
        title: '参数错误',
        icon: 'none'
      })
      setTimeout(() => {
        wx.navigateBack()
      }, 1500)
      return
    }

    const id = parseInt(options.id)
    if (isNaN(id)) {
      wx.showToast({
        title: '无效的作品ID',
        icon: 'none'
      })
      setTimeout(() => {
        wx.navigateBack()
      }, 1500)
      return
    }

    this.setData({ id })
    this.loadWorkDetail()
    this.loadComments()
  },

  // 加载作品详情
  async loadWorkDetail() {
    if (!this.data.id) return
    this.setData({ loading: true })
    
    try {
      const res = await getWorkDetail(this.data.id)
      console.log('作品详情原始响应:', res)
      
      if (res.code === 0 || res.code === 200) {  // 兼容两种成功状态码
        // 检查是否是作品所有者
        const userInfo = wx.getStorageSync('userInfo')
        
        // 获取work对象，兼容不同的返回结构
        let work = null;
        if (res.data && res.data.work) {
          work = res.data.work;
          console.log('从data.work获取作品数据');
        } else {
          work = res.data;
          console.log('直接从data获取作品数据');
        }
        
        // 确保images字段存在
        if (!work.images && res.data && res.data.images) {
          work.images = res.data.images;
          console.log('从data.images获取图片数据:', work.images.length);
        }
        
        if (!work.images) {
          work.images = [];
          console.warn('未找到图片数据');
        }
        
        // 获取服务器基础URL
        const baseUrl = getApp().globalData.baseUrl.community || "http://localhost:8008";
        
        // 处理图片URL，确保包含完整域名
        work.images = work.images.map(img => {
          if (img.url && img.url.startsWith('/')) {
            // 如果是相对路径，添加域名
            img.url = baseUrl + img.url;
          }
          return img;
        });
        
        // 处理封面图片URL
        if (work.cover_url && work.cover_url.startsWith('/')) {
          work.cover_url = baseUrl + work.cover_url;
        }
        
        // 确保作者信息存在
        if (!work.author && res.data && res.data.author) {
          work.author = res.data.author;
          console.log('从data.author获取作者数据');
        }
        
        if (!work.author) {
          work.author = {
            uid: work.uid,
            nickname: '用户' + work.uid,
            avatar: '/assets/images/default-avatar.png'
          };
          console.warn('未找到作者数据，使用默认值');
        }
        
        // 处理作者头像URL
        if (work.author && work.author.avatar && work.author.avatar.startsWith('/') && !work.author.avatar.startsWith('/assets/')) {
          work.author.avatar = baseUrl + work.author.avatar;
        }
        
        // 检查点赞状态
        if (res.data && res.data.like_status !== undefined) {
          work.like_status = res.data.like_status;
        }
        
        // 记录图片URL列表用于预览
        this.imageUrls = work.images.map(img => img.url);
        console.log('图片URL列表:', this.imageUrls);
        
        const isOwner = userInfo && userInfo.id === work.uid;
        
        this.setData({
          work,
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
        title: error.message || '网络错误，请重试',
        icon: 'none'
      })
    } finally {
      this.setData({ loading: false })
    }
  },

  // 加载评论
  async loadComments() {
    if (!this.data.id) return
    this.setData({ commentLoading: true })
    
    try {
      const res = await getComments(this.data.id)
      console.log('评论列表原始响应:', res)
      
      if (res.code === 0 || res.code === 200) {  // 兼容两种成功状态码
        // 更详细的数据检查
        let commentList = [];
        if (res.data && res.data.list && Array.isArray(res.data.list)) {
          commentList = res.data.list;
          console.log('找到评论列表数组:', commentList.length);
        } else if (res.data && Array.isArray(res.data)) {
          commentList = res.data;
          console.log('直接使用data数组:', commentList.length);
        } else {
          console.warn('未找到评论列表数组:', res.data);
          commentList = [];
        }

        // 获取服务器基础URL
        const baseUrl = getApp().globalData.baseUrl.community || "http://localhost:8008";

        // 确保每个评论对象都有必要的字段
        commentList = commentList.map(item => {
          // 确保comment字段存在
          const comment = item.comment || item;
          
          // 确保user字段存在
          const user = item.user || {
            nickname: '用户' + (comment.uid || 'Unknown'),
            avatar: '/assets/images/default-avatar.png'
          };
          
          // 处理用户头像URL
          if (user.avatar && user.avatar.startsWith('/') && !user.avatar.startsWith('/assets/')) {
            user.avatar = baseUrl + user.avatar;
          }
          
          // 处理回复评论
          let replies = item.replies || [];
          
          // 处理每个回复的用户头像URL
          replies = replies.map(reply => {
            if (reply.user && reply.user.avatar && reply.user.avatar.startsWith('/') && !reply.user.avatar.startsWith('/assets/')) {
              reply.user.avatar = baseUrl + reply.user.avatar;
            }
            return reply;
          });
          
          return {
            id: comment.id,
            content: comment.content,
            create_time: comment.create_time,
            user: user,
            replies: replies
          };
        });
        
        console.log('处理后的评论列表:', commentList);
        
        this.setData({
          comments: commentList
        });
      } else {
        wx.showToast({
          title: res.msg || '加载评论失败',
          icon: 'none'
        })
      }
    } catch (error) {
      console.error('加载评论失败:', error)
      wx.showToast({
        title: error.message || '加载评论失败',
        icon: 'none'
      })
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
      this.confirmDeleteWork()
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
    const current = e.currentTarget.dataset.src;
    if (!current) return;
    
    // 使用保存的图片URL列表
    const urls = this.imageUrls || [];
    if (urls.length === 0) {
      // 如果没有图片列表，只预览当前图片
      urls.push(current);
    }
    
    console.log('预览图片:', current, '图片列表:', urls);
    
    wx.previewImage({
      current, // 当前显示图片的http链接
      urls, // 需要预览的图片http链接列表
      fail: (err) => {
        console.error('预览图片失败:', err);
        wx.showToast({
          title: '预览图片失败',
          icon: 'none'
        });
      }
    });
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