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
    refreshing: false
  },

  onLoad() {
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
      
      const res = await getWorkList(params)
      console.log('加载社区作品结果:', res)
      
      if (res.code === 0) {
        const list = this.data.page === 1 ? res.data.list : [...this.data.communityWorks, ...res.data.list]
        this.setData({
          communityWorks: list,
          hasMore: list.length < res.data.total,
          page: this.data.page + 1
        })
      } else {
        wx.showToast({
          title: res.msg || '加载失败',
          icon: 'none'
        })
      }
    } catch (error) {
      console.error('加载社区作品失败:', error)
      wx.showToast({
        title: '网络错误，请重试',
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
    
    try {
      const res = await getUserWorks()
      console.log('加载我的作品结果:', res)
      
      if (res.code === 0) {
        this.setData({
          myWorks: res.data
        })
      } else {
        wx.showToast({
          title: res.msg || '加载失败',
          icon: 'none'
        })
      }
    } catch (error) {
      console.error('加载我的作品失败:', error)
      wx.showToast({
        title: '网络错误，请重试',
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