import yaml from 'js-yaml'

const RESOURCE_CONFIG = {
  sm2_control_planes: { label: 'SM2 Control Planes', relevance: 'control-plane' },
  sm3_control_planes: { label: 'SM3 Control Planes', relevance: 'control-plane' },
  namespaces: { label: 'Namespaces', relevance: 'supporting' },
  virtual_services: { label: 'Virtual Services', relevance: 'traffic' },
  gateways: { label: 'Gateways', relevance: 'traffic' },
  service_entries: { label: 'Service Entries', relevance: 'traffic' },
  routes: { label: 'Routes', relevance: 'traffic' },
  smmr: { label: 'Member Rolls', relevance: 'membership' },
  smms: { label: 'Members', relevance: 'membership' },
  authorization_policies: { label: 'Authorization Policies', relevance: 'security' },
  peer_authentications: { label: 'Peer Authentications', relevance: 'security' },
  sidecars: { label: 'Sidecars', relevance: 'security' },
  telemetries: { label: 'Telemetries', relevance: 'observability' },
}

const RESOURCE_ORDER = [
  'sm2_control_planes',
  'sm3_control_planes',
  'namespaces',
  'virtual_services',
  'gateways',
  'service_entries',
  'routes',
  'smmr',
  'smms',
  'authorization_policies',
  'peer_authentications',
  'sidecars',
  'telemetries',
]

export function buildInventoryRows(resources = {}) {
  const rows = []

  for (const bucket of RESOURCE_ORDER) {
    const items = Array.isArray(resources[bucket]) ? resources[bucket] : []
    const config = RESOURCE_CONFIG[bucket] || { label: bucket, relevance: 'supporting' }

    for (const item of items) {
      rows.push({
        id: `${bucket}/${item.namespace || '_cluster'}/${item.name}`,
        bucket,
        relevance: config.relevance,
        kind: item.kind || config.label,
        name: item.name,
        namespace: item.namespace || 'cluster-scoped',
        spec: item.spec ?? null,
        status: item.status ?? null,
        raw: item,
      })
    }
  }

  return rows.sort((left, right) => {
    if (left.bucket !== right.bucket) return RESOURCE_ORDER.indexOf(left.bucket) - RESOURCE_ORDER.indexOf(right.bucket)
    if (left.namespace !== right.namespace) return left.namespace.localeCompare(right.namespace)
    return left.name.localeCompare(right.name)
  })
}

export function buildInventoryFilterOptions(rows) {
  return {
    kinds: uniqueSorted(rows.map((row) => row.kind)),
    namespaces: uniqueSorted(rows.map((row) => row.namespace)),
    relevance: uniqueSorted(rows.map((row) => row.relevance)),
  }
}

function uniqueSorted(values) {
  return Array.from(new Set(values.filter(Boolean))).sort()
}

export function formatResourceDocument(value, format) {
  if (value == null) return ''
  if (format === 'yaml') {
    return yaml.dump(value, { indent: 2, lineWidth: -1 })
  }
  return JSON.stringify(value, null, 2)
}

export function formatLabel(value) {
  return value
    .split(/[-_]/g)
    .filter(Boolean)
    .map((part) => part.charAt(0).toUpperCase() + part.slice(1))
    .join(' ')
}
