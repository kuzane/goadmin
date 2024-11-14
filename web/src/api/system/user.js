import request from '@/utils/request'

export function login(data) {
  console.log("Base API URL:", process.env.VUE_APP_BASE_API);
  return request({
    url: '/login',
    method: 'post',
    data
  })
}

export function logout() {
  return request({
    url: '/logout',
    method: 'post'
  })
}

export function getInfo() {
  return request({
    url: '/userinfo',
    method: 'get'
  })
}

export function addUser(data) {
  return request({
    url: '/sys/users',
    method: 'post',
    data
  })
}

export function deleteUser(ids) {
  return request({
    url: '/sys/users/' + ids,
    method: 'delete'
  })
}

export function updateUser(id, data) {
  return request({
    url: `/sys/users/${id}`,
    method: 'patch',
    data
  })
}

export function listUser(query) {
  return request({
    url: '/sys/users',
    method: 'get',
    params: query
  })
}

export function detailUser(id) {
  return request({
    url: `/sys/users/${id}`,
    method: 'get'
  })
}

// 通过邮箱重置密码
export function forgotPwd(data) {
  return request({
    url: '/pwd/forgot',
    method: 'post',
    data: data
  })
}

// 获取验证吗
export function getEmailCaptcha(data) {
  return request({
    url: '/email/captcha',
    method: 'post',
    data: data
  })
}

// 修改密码
export function changePassword(data) {
  return request({
    url: '/password',
    method: 'put',
    data: data
  })
}

// 修改个人信息
export function changeProfile(data) {
  return request({
    url: '/profile',
    method: 'put',
    data: data
  })
}
