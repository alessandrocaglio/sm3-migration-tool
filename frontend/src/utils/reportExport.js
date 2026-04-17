function download(filename, content, type) {
  const blob = new Blob([content], { type })
  const url = URL.createObjectURL(blob)
  const anchor = document.createElement('a')
  anchor.href = url
  anchor.download = filename
  anchor.click()
  URL.revokeObjectURL(url)
}

export function exportScanAsJSON(scanData, selectedNamespace) {
  const date = new Date().toISOString().slice(0, 19).replace(/:/g, '-')
  const name = selectedNamespace || 'cluster'
  download(`sm3-evidence-${name}-${date}.json`, JSON.stringify(scanData, null, 2), 'application/json')
}

export function exportFindingsAsMarkdown(scanData, selectedNamespace, lastScannedAt) {
  const findings = scanData.findings || []
  const summary = scanData.summary || {}
  const counts = summary.counts || {}
  const namespace = selectedNamespace || 'n/a'
  const scannedAt = lastScannedAt || new Date().toISOString()

  let markdown = '# Service Mesh 3 Evidence Report\n\n'
  markdown += `- Control Plane: \`${namespace}\`\n`
  markdown += `- Scanned At: \`${scannedAt}\`\n`
  markdown += `- Readiness: \`${summary.readiness?.status || 'unknown'}\` (${summary.readiness?.score ?? 0}/100)\n`
  markdown += `- Checks: \`${counts.checks_total || 0}\`\n`
  markdown += `- Findings: \`${counts.findings_total || 0}\`\n\n`

  markdown += '## Findings\n\n'
  markdown += '| Severity | Category | Resource | Message |\n'
  markdown += '|---|---|---|---|\n'

  const views = scanData.finding_views || []
  for (const finding of views) {
    markdown += `| ${finding.severity} | ${finding.category_label} | ${finding.namespace}/${finding.resource_name} | ${finding.message} |\n`
  }

  if (!views.length) {
    markdown += '| - | - | - | No findings |\n'
  }

  const date = new Date().toISOString().slice(0, 19).replace(/:/g, '-')
  const name = selectedNamespace || 'cluster'
  download(`sm3-evidence-${name}-${date}.md`, markdown, 'text/markdown')
}
