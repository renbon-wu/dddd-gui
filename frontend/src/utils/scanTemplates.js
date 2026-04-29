export const scanTemplates = [
  {
    id: 'internal',
    name: '内网扫描',
    description: '扫描内网资产，使用默认端口',
    icon: 'Connection',
    config: {
      ports: 'Top1000',
      scanType: 'tcp',
      enablePoc: true,
      enableFingerprint: true
    }
  },
  {
    id: 'external',
    name: '外网扫描',
    description: '扫描外网资产，使用全端口',
    icon: 'Promotion',
    config: {
      ports: '1-65535',
      scanType: 'syn',
      enablePoc: true,
      enableFingerprint: true
    }
  },
  {
    id: 'fingerprint-only',
    name: '仅指纹识别',
    description: '仅识别资产指纹，不进行漏洞扫描',
    icon: 'Search',
    config: {
      ports: 'Top1000',
      scanType: 'tcp',
      enablePoc: false,
      enableFingerprint: true
    }
  },
  {
    id: 'full-scan',
    name: '完整漏洞扫描',
    description: '全面扫描，包括指纹识别和漏洞检测',
    icon: 'Warning',
    config: {
      ports: 'Top2000',
      scanType: 'tcp',
      enablePoc: true,
      enableFingerprint: true,
      deepScan: true
    }
  },
  {
    id: 'subdomain',
    name: '子域名枚举',
    description: '枚举子域名并扫描',
    icon: 'Link',
    config: {
      ports: 'Top100',
      scanType: 'tcp',
      enablePoc: true,
      enableFingerprint: true,
      subdomain: true
    }
  }
]

// 获取所有模板
export const getAllTemplates = () => {
  return [...scanTemplates]
}

// 根据ID获取模板
export const getTemplateById = (id) => {
  return scanTemplates.find(t => t.id === id)
}

// 过滤模板
export const filterTemplates = (keyword) => {
  if (!keyword) return getAllTemplates()
  const lowerKeyword = keyword.toLowerCase()
  return scanTemplates.filter(t => 
    t.name.toLowerCase().includes(lowerKeyword) ||
    t.description.toLowerCase().includes(lowerKeyword)
  )
}
