<script setup>
import { computed, ref } from 'vue'
import yaml from 'js-yaml'
import ChecksTable from '../components/evidence/ChecksTable.vue'
import ExportPanel from '../components/evidence/ExportPanel.vue'
import CodeExplorer from '../components/CodeExplorer.vue'
import { exportFindingsAsMarkdown, exportScanAsJSON } from '../utils/reportExport'

const props = defineProps({
  scanData: { type: Object, default: () => ({}) },
  selectedNamespace: { type: String, default: '' },
  lastScannedAt: { type: String, default: '' },
})

const statusFilter = ref('all')
const checkSearch = ref('')
const rawFormat = ref('json')

const checks = computed(() => props.scanData.checks || [])
const summary = computed(() => props.scanData.summary || {})
const counts = computed(() => summary.value.counts || {})

const filteredChecks = computed(() => {
  const search = checkSearch.value.trim().toLowerCase()
  return checks.value.filter((check) => {
    const statusMatch = statusFilter.value === 'all' || check.status === statusFilter.value
    const searchMatch =
      !search ||
      check.title.toLowerCase().includes(search) ||
      check.target.toLowerCase().includes(search) ||
      (check.checker || '').toLowerCase().includes(search)
    return statusMatch && searchMatch
  })
})

const rawPayload = computed(() => {
  if (rawFormat.value === 'json') return JSON.stringify(props.scanData, null, 2)
  return yaml.dump(props.scanData || {}, { indent: 2, lineWidth: -1 })
})

function onExportJSON() {
  exportScanAsJSON(props.scanData, props.selectedNamespace)
}

function onExportMarkdown() {
  exportFindingsAsMarkdown(props.scanData, props.selectedNamespace, props.lastScannedAt)
}
</script>

<template>
  <div class="page page--evidence">
    <ExportPanel
      :summary="summary"
      :selected-namespace="selectedNamespace"
      :last-scanned-at="lastScannedAt"
      :counts="counts"
      @export-json="onExportJSON"
      @export-markdown="onExportMarkdown"
    />

    <section class="content-card">
      <div class="section-heading">
        <div>
          <p class="section-heading__eyebrow">Filter</p>
          <h2 class="section-heading__title">Check filters</h2>
        </div>
      </div>

      <div class="evidence-filters">
        <label class="picker">
          <span class="picker__label">Status</span>
          <select v-model="statusFilter" class="picker__select">
            <option value="all">All statuses</option>
            <option value="Failure">Failure</option>
            <option value="Success">Success</option>
          </select>
        </label>
        <label class="picker">
          <span class="picker__label">Search</span>
          <input v-model="checkSearch" class="picker__input" type="search" placeholder="Search check title, target, checker" />
        </label>
      </div>
    </section>

    <ChecksTable :checks="filteredChecks" />

    <section class="content-card">
      <div class="section-heading">
        <div>
          <p class="section-heading__eyebrow">Raw Payload</p>
          <h2 class="section-heading__title">API response</h2>
        </div>
        <div class="format-toggle">
          <button :class="['format-toggle__button', { 'format-toggle__button--active': rawFormat === 'json' }]" @click="rawFormat = 'json'">
            JSON
          </button>
          <button :class="['format-toggle__button', { 'format-toggle__button--active': rawFormat === 'yaml' }]" @click="rawFormat = 'yaml'">
            YAML
          </button>
        </div>
      </div>

      <div class="evidence-raw">
        <CodeExplorer :model-value="rawPayload" :language="rawFormat === 'json' ? 'json' : 'yaml'" />
      </div>
    </section>
  </div>
</template>
