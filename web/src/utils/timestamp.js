
export function formatTimestamp(timestamp) { // 转化时间戳
  const date = new Date(timestamp * 1000) // 转换为毫秒
  const options = { year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit', second: '2-digit', hour12: false }
  return date.toLocaleString('zh-CN', options).replace(/\//g, '-') // 格式化并替换斜杠为短横线
}
