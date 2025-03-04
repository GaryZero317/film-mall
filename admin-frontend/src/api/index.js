// 直接API调用助手
import { adminService } from './request'
import * as communityApi from './community'

// 直接调用API，返回原始响应
export function callDirectApi(url, method, data) {
  return adminService({
    url,
    method: method || 'post',
    data
  })
}

export default {
  callDirectApi,
  community: communityApi
} 