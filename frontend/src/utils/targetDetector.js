// IP地址检测
const ipRegex = /^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/

// CIDR网段检测
const cidrRegex = /^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\/([0-9]|[1-2][0-9]|3[0-2])$/

// 域名检测
const domainRegex = /^[a-zA-Z0-9][a-zA-Z0-9-]{1,61}[a-zA-Z0-9]\.[a-zA-Z]{2,}$/

// URL检测
const urlRegex = /^https?:\/\/[^\s]+$/

// Hunter查询检测
const hunterQueryRegex = /^ip="|^icp\.name=|^icp\.domain=/

// Fofa查询检测
const fofaQueryRegex = /^domain="|^ip="|^host=/

// 检测目标类型
export const detectTargetType = (target) => {
  if (!target || target.trim() === '') {
    return { type: 'unknown', value: target, label: '未知类型' }
  }
  const trimmedTarget = target.trim()
  
  if (cidrRegex.test(trimmedTarget)) {
    return { type: 'cidr', value: trimmedTarget, label: 'CIDR网段' }
  } else if (ipRegex.test(trimmedTarget)) {
    return { type: 'ip', value: trimmedTarget, label: 'IP地址' }
  } else if (urlRegex.test(trimmedTarget)) {
    return { type: 'url', value: trimmedTarget, label: 'URL地址' }
  } else if (domainRegex.test(trimmedTarget)) {
    return { type: 'domain', value: trimmedTarget, label: '域名' }
  } else if (hunterQueryRegex.test(trimmedTarget)) {
    return { type: 'hunter', value: trimmedTarget, label: 'Hunter查询' }
  } else if (fofaQueryRegex.test(trimmedTarget)) {
    return { type: 'fofa', value: trimmedTarget, label: 'Fofa查询' }
  }
  return { type: 'unknown', value: trimmedTarget, label: '未知类型' }
}

// 获取类型对应的标签颜色
export const getTypeTagType = (type) => {
  switch (type) {
    case 'ip':
      return 'success'
    case 'cidr':
      return 'warning'
    case 'url':
      return 'primary'
    case 'domain':
      return 'info'
    case 'hunter':
      return 'danger'
    case 'fofa':
      return 'danger'
    default:
      return ''
  }
}

// 批量检测目标类型
export const detectTargetsTypes = (targets) => {
  return targets.map(t => detectTargetType(t))
}
