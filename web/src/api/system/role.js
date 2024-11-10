import request from '@/utils/request'

export function listRole(query) {
  return request({
    url: '/sys/roles',
    method: 'get',
    params: query
  })
}

export function addRole(data) {
  return request({
    url: '/sys/roles',
    method: 'post',
    data
  })
}

export function updateRole(id, data) {
  return request({
    url: `/sys/roles/${id}`,
    method: 'patch',
    data
  })
}

export function deleteRole(ids) {
  return request({
    url: '/sys/roles/' + ids,
    method: 'delete'
  })
}

export function detailRole(id) {
  return request({
    url: `/sys/roles/${id}`,
    method: 'get'
  })
}

export function treeEndpoint() {
  return request({
    url: '/sys/roles/apis',
    method: 'get'
  })
}
