/**
 * 格式化时间
 * @param {string|Date} time 时间字符串或Date对象
 * @param {string} format 格式化模板，默认为 'YYYY-MM-DD HH:mm:ss'
 * @returns {string} 格式化后的时间字符串
 */
export function formatTime(time, format = 'YYYY-MM-DD HH:mm:ss') {
  if (!time) return '';
  
  const date = typeof time === 'string' ? new Date(time) : time;
  
  // 检查日期是否有效
  if (isNaN(date.getTime())) {
    return '';
  }
  
  // 检查是否为Go语言的时间零值（0001-01-01）
  if (date.getFullYear() === 1 && date.getMonth() === 0 && date.getDate() === 1) {
    return '-';
  }
  
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  const hours = String(date.getHours()).padStart(2, '0');
  const minutes = String(date.getMinutes()).padStart(2, '0');
  const seconds = String(date.getSeconds()).padStart(2, '0');
  
  return format
    .replace('YYYY', year)
    .replace('MM', month)
    .replace('DD', day)
    .replace('HH', hours)
    .replace('mm', minutes)
    .replace('ss', seconds);
}
