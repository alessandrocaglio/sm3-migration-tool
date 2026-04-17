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

const RESOURCE_LABELS = {
  sm2_control_planes: 'SM2 Control Planes',
  sm3_control_planes: 'SM3 Control Planes',
  smmr: 'Member Rolls',
  smms: 'Members',
  virtual_services: 'Virtual Services',
  gateways: 'Gateways',
  service_entries: 'Service Entries',
  authorization_policies: 'Authorization Policies',
  peer_authentications: 'Peer Authentications',
  sidecars: 'Sidecars',
  telemetries: 'Telemetries',
  routes: 'Routes',
  namespaces: 'Namespaces',
}

const severityWeights = {
  High: 6,
  Medium: 3,
  Low: 1,
  Info: 0,
}

export function buildSummary(data) {
  const checks = data.checks || []
  const findings = data.findings || []
  const meshNamespaces = data.mesh_namespaces || []

  const passedChecks = checks.filter((check) => check.status === 'Success').length
  const failedChecks = checks.filter((check) => check.status === 'Failure').length
  const highSeverityFindings = findings.filter((finding) => finding.severity === 'High').length
  const mediumSeverityFindings = findings.filter((finding) => finding.severity === 'Medium').length
  const lowSeverityFindings = findings.filter((finding) => finding.severity === 'Low').length
  const readinessScore = buildReadinessScore(checks, findings)
  const readinessState = describeReadiness(readinessScore, highSeverityFindings, failedChecks)
  const resourceCounts = buildResourceCounts(data.resources)

  return {
    selectedNamespace: data.scanned_namespace || '',
    passedChecks,
    failedChecks,
    totalChecks: checks.length,
    findingsCount: findings.length,
    highSeverityFindings,
    mediumSeverityFindings,
    lowSeverityFindings,
    meshNamespaceCount: meshNamespaces.length,
    resourceCount: resourceCounts.total,
    readinessScore,
    readinessState,
    readinessTone: readinessStateToTone(readinessState),
    nextStep: buildNextStep(findings, failedChecks),
  }
}

function buildReadinessScore(checks = [], findings = []) {
  if (!checks.length) return 0

  const weightedPenalty = findings.reduce((total, finding) => {
    return total + (severityWeights[finding.severity] ?? 0)
  }, 0)

  const maxPenalty = Math.max(checks.length * 6, 1)
  const penaltyRatio = Math.min(weightedPenalty / maxPenalty, 1)
  const score = Math.round(100 - penaltyRatio * 100)

  return Math.max(0, Math.min(score, 100))
}

function describeReadiness(score, highSeverityFindings, failedChecks) {
  if (highSeverityFindings > 0 || score < 45) return 'Not ready'
  if (failedChecks > 0 || score < 75) return 'Needs attention'
  return 'Ready to proceed'
}

function readinessStateToTone(state) {
  if (state === 'Not ready') return 'danger'
  if (state === 'Needs attention') return 'warning'
  return 'success'
}

function buildNextStep(findings, failedChecks) {
  if (!failedChecks) return 'Review the action plan, then validate migration sequencing and export the assessment.'

  const highSeverity = findings.filter((finding) => finding.severity === 'High')
  if (highSeverity.length > 0) {
    return `Resolve ${highSeverity.length} high-severity blocker${highSeverity.length === 1 ? '' : 's'} before planning workload migration.`
  }

  return 'Triage remaining failures and convert them into an ordered migration action plan.'
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
  return [...(data.findings || [])]
    .sort((left, right) => (severityWeights[right.severity] ?? 0) - (severityWeights[left.severity] ?? 0))
    .slice(0, 5)
}

export function buildSeverityDistribution(data) {
  const findings = data.findings || []

  return [
    { key: 'High', label: 'High', count: findings.filter((finding) => finding.severity === 'High').length, tone: 'danger' },
    { key: 'Medium', label: 'Medium', count: findings.filter((finding) => finding.severity === 'Medium').length, tone: 'warning' },
    { key: 'Low', label: 'Low', count: findings.filter((finding) => finding.severity === 'Low').length, tone: 'neutral' },
  ]
}

export function buildNamespaceImpact(data) {
  const meshNamespaces = data.mesh_namespaces || []
  const findings = data.findings || []

  const issueCounts = new Map()
  for (const finding of findings) {
    issueCounts.set(finding.namespace, (issueCounts.get(finding.namespace) || 0) + 1)
  }

  return meshNamespaces
    .map((namespace) => ({
      namespace,
      issues: issueCounts.get(namespace) || 0,
    }))
    .sort((left, right) => right.issues - left.issues || left.namespace.localeCompare(right.namespace))
}

export function buildTopResourceTypes(resources = {}) {
  return Object.entries(buildResourceCounts(resources).byType)
    .filter(([, count]) => count > 0)
    .sort((left, right) => right[1] - left[1])
    .slice(0, 4)
    .map(([key, count]) => ({
      key,
      label: RESOURCE_LABELS[key] || key,
      count,
    }))
}
