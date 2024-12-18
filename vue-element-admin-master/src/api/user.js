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
    url: '/vue-element-admin/user/info',
    method: 'get',
    params: { token }
  })
}

export function logout() {
  return request({
    url: '/vue-element-admin/user/logout',
    method: 'post'
  })
}
