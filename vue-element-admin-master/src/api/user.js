import request from '@/utils/request'

export function login(data) {
  return request({
    url: 'http://localhost:8000/api/admin/login',
    method: 'post',
    data: {
      username: data.username,
      password: data.password
    }
  })
}

export function getInfo(token) {
  return request({
    url: 'http://localhost:8000/api/admin/info',
    method: 'post',
    headers: {
      'Authorization': `Bearer ${token}`
    }
  }).then(response => {
    if (!Object.prototype.hasOwnProperty.call(response.data, 'level')) {
      throw new Error('getInfo: response must include a level field')
    }
    return response
  })
}

export function logout() {
  return request({
    url: '/vue-element-admin/user/logout',
    method: 'post'
  })
}
