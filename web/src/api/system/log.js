import request from '@/utils/request'

export function listLog(query) {
  return request({
    url: '/sys/logs',
    method: 'get',
    params: query
  })
}

export function deleteLog(ids) {
  return request({
    url: '/sys/logs/' + ids,
    method: 'delete'
  })
}

export function emptyLog() {
  return request({
    url: '/sys/logs',
    method: 'delete'
  })
}
