import { getWorkList, getUserWorks } from '../../../api/community'

Page({
  data: {
    activeTab: 0, // 0: 社区, 1: 我的作品
    communityWorks: [],
    myWorks: [],
    loading: false,
    page: 1,
    pageSize: 10,
    hasMore: true,
    refreshing: false,
    baseUrl: getApp().globalData.baseUrl?.community || 'http://localhost:8008' // 使用全局配置的社区服务地址
  },

  onLoad() {
    // 使用全局配置的社区服务地址
    const baseUrl = getApp().globalData.baseUrl?.community || 'http://localhost:8008'
    this.setData({ baseUrl })
    console.log('社区服务地址:', baseUrl)
    this.loadCommunityWorks()
  },

  // 切换标签
  onTabChange(e) {
    const index = parseInt(e.currentTarget.dataset.index)
    this.setData({ activeTab: index })
    
    if (index === 0 && this.data.communityWorks.length === 0) {
      this.loadCommunityWorks()
    } else if (index === 1 && this.data.myWorks.length === 0) {
      this.loadMyWorks()
    }
  },

  // 加载社区作品
  async loadCommunityWorks(refresh = false) {
    if (this.data.loading) return
    
    if (refresh) {
      this.setData({
        page: 1,
        hasMore: true,
        communityWorks: []
      })
    }
    
    if (!this.data.hasMore && !refresh) return
    
    this.setData({ loading: true })
    
    try {
      const params = {
        page: this.data.page,
        pageSize: this.data.pageSize
      }
      
      console.log('请求参数:', params)
      const res = await getWorkList(params)
      console.log('API原始响应:', res)
      
      if (!res) {
        throw new Error('API响应为空')
      }
      
      if (res.code === 0) {
        if (!res.data || !res.data.list) {
          console.error('API响应数据格式错误:', res)
          throw new Error('返回数据格式错误')
        }
        
        // 处理嵌套的work对象并添加baseUrl
        const list = res.data.list.map(item => {
          const work = { ...item.work }
          // 处理封面图URL
          if (work.cover_url && !work.cover_url.startsWith('http')) {
            // 确保URL格式正确
            let url = work.cover_url
            if (!url.startsWith('/')) {
              url = '/' + url
            }
            work.cover_url = this.data.baseUrl + url
            console.log('处理后的图片URL:', work.cover_url)
          }
          return {
            ...work,
            author: item.author
          }
        })
        
        const finalList = this.data.page === 1 ? list : [...this.data.communityWorks, ...list]
        console.log('处理后的列表数据:', finalList)
        
        if (finalList && finalList.length > 0) {
          console.log('第一项数据:', finalList[0])
          if (!finalList[0].id || !finalList[0].title || !finalList[0].cover_url) {
            console.warn('数据字段缺失:', {
              id: finalList[0].id,
              title: finalList[0].title,
              cover_url: finalList[0].cover_url
            })
          }
        } else {
          console.log('列表为空')
        }
        
        this.setData({
          communityWorks: finalList,
          hasMore: finalList.length < res.data.total,
          page: this.data.page + 1
        })
      } else {
        console.error('API返回错误:', res.msg)
        wx.showToast({
          title: res.msg || '加载失败',
          icon: 'none'
        })
      }
    } catch (error) {
      console.error('加载社区作品失败:', error)
      wx.showToast({
        title: error.message || '网络错误，请重试',
        icon: 'none'
      })
    } finally {
      this.setData({ 
        loading: false,
        refreshing: false
      })
    }
  },

  // 加载我的作品
  async loadMyWorks() {
    if (this.data.loading) return
    
    this.setData({ loading: true })
    console.log('开始加载我的作品')
    
    try {
      const res = await getUserWorks()
      console.log('加载我的作品API响应:', res)
      
      if (res.code === 0) {
        // 验证数据格式
        if (!res.data) {
          console.error('API响应数据为空')
          throw new Error('返回数据为空')
        }

        // 确保 res.data 是数组
        const worksList = Array.isArray(res.data) ? res.data : res.data.list || []
        console.log('解析后的作品列表:', worksList)

        // 处理作品数据，添加完整的图片URL
        const works = worksList.map(work => {
          // 确保work对象存在
          if (!work) return null

          // 如果work是嵌套的，提取work对象
          const workData = work.work || work

          if (workData.cover_url && !workData.cover_url.startsWith('http')) {
            let url = workData.cover_url
            if (!url.startsWith('/')) {
              url = '/' + url
            }
            workData.cover_url = this.data.baseUrl + url
          }
          return workData
        }).filter(work => work !== null) // 过滤掉无效的数据
        
        console.log('处理后的我的作品数据:', works)
        this.setData({
          myWorks: works
        })
      } else {
        console.error('加载我的作品失败:', res.msg)
        wx.showToast({
          title: res.msg || '加载失败',
          icon: 'none'
        })
      }
    } catch (error) {
      console.error('加载我的作品失败:', error)
      wx.showToast({
        title: error.message || '网络错误，请重试',
        icon: 'none'
      })
    } finally {
      this.setData({ loading: false })
    }
  },

  // 刷新
  onRefresh() {
    this.setData({ refreshing: true })
    if (this.data.activeTab === 0) {
      this.loadCommunityWorks(true)
    } else {
      this.loadMyWorks()
    }
  },

  // 加载更多
  onLoadMore() {
    if (this.data.activeTab === 0) {
      this.loadCommunityWorks()
    }
  },

  // 创建新作品
  onCreateWork() {
    console.log('点击发布按钮')
    wx.navigateTo({
      url: '/pages/user/community/create/index',
      fail: (err) => {
        console.error('跳转失败:', err)
        wx.showToast({
          title: '跳转失败',
          icon: 'none'
        })
      }
    })
  },

  // 查看作品详情
  onViewWorkDetail(e) {
    const id = e.currentTarget.dataset.id
    if (!id) {
      wx.showToast({
        title: '无效的作品ID',
        icon: 'none'
      })
      return
    }
    wx.navigateTo({
      url: `/pages/user/community/detail/index?id=${id}`
    })
  },

  // 下拉刷新
  onPullDownRefresh() {
    this.onRefresh()
    wx.stopPullDownRefresh()
  },

  // 上拉加载更多
  onReachBottom() {
    this.onLoadMore()
  }
}) 