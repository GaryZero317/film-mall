import request from '../utils/request'

export function getAddressList() {
  return request({
    url: '/api/address/list',
    method: 'GET'
  })
}

export function getDefaultAddress() {
  return request({
    url: '/api/address/default',
    method: 'GET'
  })
}

export function createAddress(data) {
  return request({
    url: '/api/address/add',
    method: 'POST',
    data
  })
}

export function updateAddress(id, data) {
  return request({
    url: '/api/address/update',
    method: 'POST',
    data: {
      id,
      ...data
    }
  })
}

export function deleteAddress(id) {
  return request({
    url: '/api/address/delete',
    method: 'POST',
    data: {
      id
    }
  })
}

export function setDefaultAddress(id) {
  return request({
    url: '/api/address/setDefault',
    method: 'POST',
    data: {
      id
    }
  })
}

export function getAddressDetail(id) {
  return request({
    url: `/api/address/${id}`,
    method: 'GET'
  })
} 