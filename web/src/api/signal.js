import request from './request'

// 获取信号列表
export function getSignalList(params) {
  return request({
    url: '/api/signals',
    method: 'get',
    params
  })
}

// 获取信号详情
export function getSignalDetail(id) {
  return request({
    url: `/api/signals/${id}`,
    method: 'get'
  })
}