const RESOURCE_KEYS = [
  'sm2_control_planes',
  'sm3_control_planes',
  'smmr',
  'smms',
  'virtual_services',
  'gateways',
  'service_entries',
  'authorization_policies',
  'peer_authentications',
  'sidecars',
  'telemetries',
  'routes',
  'namespaces',
]

export function buildSummary(data) {
  const checks = data.checks || []
  const findings = data.findings || []
  const meshNamespaces = data.mesh_namespaces || []

  const passedChecks = checks.filter((check) => check.status === 'Success').length
  const failedChecks = checks.filter((check) => check.status === 'Failure').length
  const highSeverityFindings = findings.filter((finding) => finding.severity === 'High').length

  return {
    selectedNamespace: data.scanned_namespace || '',
    passedChecks,
    failedChecks,
    totalChecks: checks.length,
    findingsCount: findings.length,
    highSeverityFindings,
    meshNamespaceCount: meshNamespaces.length,
    resourceCount: buildResourceCounts(data.resources).total,
  }
}

export function buildResourceCounts(resources = {}) {
  const byType = {}
  let total = 0

  for (const key of RESOURCE_KEYS) {
    const count = Array.isArray(resources[key]) ? resources[key].length : 0
    byType[key] = count
    total += count
  }

  return { byType, total }
}

export function buildRecentFindings(data) {
  return (data.findings || []).slice(0, 5)
}
