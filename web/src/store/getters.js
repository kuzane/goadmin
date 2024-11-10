const getters = {
  sidebar: state => state.app.sidebar,
  size: state => state.app.size,
  device: state => state.app.device,
  visitedViews: state => state.tagsView.visitedViews,
  cachedViews: state => state.tagsView.cachedViews,
  token: state => state.user.token,
  avatar: state => state.user.avatar,
  name: state => state.user.name,
  phone: state => state.user.phone,
  email: state => state.user.email,
  description: state => state.user.description,
  nickname: sate => sate.user.nickname,
  roles: state => state.user.roles,
  permission_routes: state => state.permission.routes,
  errorLogs: state => state.errorLog.logs,
  hasPermission: state => (permission) => {
    return state.user.permissions.includes('*') || state.user.permissions.includes(permission)
  }

}
export default getters
