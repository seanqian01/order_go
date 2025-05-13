import request from './request'

// 获取订单列表
export function getOrderList(params) {
  return request({
    url: '/api/orders',
    method: 'get',
    params
  })
}

// 获取订单详情
export function getOrderDetail(id) {
  return request({
    url: `/api/orders/${id}`,
    method: 'get'
  })
}