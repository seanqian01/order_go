/**
 * 格式化时间
 * @param {string|Date} time 时间字符串或Date对象
 * @param {string} format 格式化模板，默认为 'YYYY-MM-DD HH:mm:ss'
 * @returns {string} 格式化后的时间字符串
 */
export function formatTime(time, format = 'YYYY-MM-DD HH:mm:ss') {
  if (!time) return '';
  
  // 将字符串转换为Date对象
  const date = typeof time === 'string' ? new Date(time) : time;
  
  // 检查日期是否有效
  if (isNaN(date.getTime())) {
    return '';
  }
  
  // 检查是否为Go语言的时间零值（0001-01-01）
  if (date.getFullYear() === 1 && date.getMonth() === 0 && date.getDate() === 1) {
    return '-';
  }
  
  // 转换为上海时区（UTC+8）
  const shanghaiDate = new Date(date.getTime() + 8 * 60 * 60 * 1000);
  
  // 使用toISOString()获取ISO格式的时间字符串，然后取前19位（YYYY-MM-DDTHH:mm:ss）
  const isoString = shanghaiDate.toISOString().slice(0, 19);
  const year = isoString.slice(0, 4);
  const month = isoString.slice(5, 7);
  const day = isoString.slice(8, 10);
  const hours = isoString.slice(11, 13);
  const minutes = isoString.slice(14, 16);
  const seconds = isoString.slice(17, 19);
  
  return format
    .replace('YYYY', year)
    .replace('MM', month)
    .replace('DD', day)
    .replace('HH', hours)
    .replace('mm', minutes)
    .replace('ss', seconds);
}
