import { computed, reactive, readonly } from 'vue'
import { fetchNamespaces, fetchScan } from '../services/api'
import {
  buildNamespaceImpact,
  buildRecentFindings,
  buildResourceCounts,
  buildSeverityDistribution,
  buildSummary,
  buildTopResourceTypes,
} from '../services/mappers'

const state = reactive({
  allowedNamespaces: [],
  selectedNamespace: '',
  scanStatus: 'idle',
  error: null,
  data: {
    checks: [],
    findings: [],
    mesh_namespaces: [],
    resources: {},
    scanned_namespace: '',
  },
  initialized: false,
  lastScannedAt: '',
})

const summary = computed(() => buildSummary(state.data))
const resourceCounts = computed(() => buildResourceCounts(state.data.resources))
const recentFindings = computed(() => buildRecentFindings(state.data))
const severityDistribution = computed(() => buildSeverityDistribution(state.data))
const namespaceImpact = computed(() => buildNamespaceImpact(state.data))
const topResourceTypes = computed(() => buildTopResourceTypes(state.data.resources))

async function loadNamespaces() {
  const payload = await fetchNamespaces()
  state.allowedNamespaces = payload.namespaces || []

  if (!state.selectedNamespace && state.allowedNamespaces.length === 1) {
    state.selectedNamespace = state.allowedNamespaces[0]
  }
}

async function runScan(namespace = state.selectedNamespace) {
  if (!namespace) {
    state.data = {
      checks: [],
      findings: [],
      mesh_namespaces: [],
      resources: {},
      scanned_namespace: '',
    }
    state.error = null
    state.scanStatus = 'idle'
    return
  }

  state.scanStatus = 'loading'
  state.error = null

  try {
    const payload = await fetchScan(namespace)
    state.data = payload
    state.selectedNamespace = payload.scanned_namespace || namespace
    state.scanStatus = 'ready'
    state.lastScannedAt = new Date().toISOString()
  } catch (error) {
    state.error = error.message
    state.scanStatus = 'error'
    state.data = {
      checks: [],
      findings: [],
      mesh_namespaces: [],
      resources: {},
      scanned_namespace: '',
    }
  }
}

async function initialize() {
  if (state.initialized) return

  state.scanStatus = 'loading'

  try {
    await loadNamespaces()
    state.initialized = true

    if (state.selectedNamespace) {
      await runScan(state.selectedNamespace)
    } else {
      state.scanStatus = 'idle'
    }
  } catch (error) {
    state.error = error.message
    state.scanStatus = 'error'
    state.initialized = true
  }
}

async function selectNamespace(namespace) {
  state.selectedNamespace = namespace
  await runScan(namespace)
}

async function refresh() {
  await runScan(state.selectedNamespace)
}

export function useScanStore() {
  return {
    state: readonly(state),
    summary,
    resourceCounts,
    recentFindings,
    severityDistribution,
    namespaceImpact,
    topResourceTypes,
    initialize,
    selectNamespace,
    refresh,
  }
}
