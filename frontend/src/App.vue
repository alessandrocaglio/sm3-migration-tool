<script setup>
import { ref, onMounted } from 'vue'

const findings = ref([])
const loading = ref(true)
const scanning = ref(false)
const error = ref(null)

const allowedNamespaces = ref([])
const selectedNamespace = ref('')

const expandedRows = ref(new Set())

const toggleRow = (index) => {
  if (expandedRows.value.has(index)) {
    expandedRows.value.delete(index)
  } else {
    expandedRows.value.add(index)
  }
}

const copyToClipboard = async (text) => {
  try {
    await navigator.clipboard.writeText(text)
  } catch (err) {
    console.error('Failed to copy text: ', err)
  }
}

const fetchNamespaces = async () => {
  try {
    const res = await fetch('/api/namespaces')
    if (!res.ok) throw new Error('Failed to fetch namespaces')
    const data = await res.json()
    allowedNamespaces.value = data.namespaces || []
    if (allowedNamespaces.value.length > 0 && !selectedNamespace.value) {
      selectedNamespace.value = allowedNamespaces.value[0]
    }
  } catch (e) {
    console.error(e)
  }
}

const fetchScan = async () => {
  loading.value = true
  error.value = null
  try {
    const res = await fetch(`/api/scan?namespace=${selectedNamespace.value}`)
    if (!res.ok) {
      const errData = await res.json().catch(() => ({}))
      throw new Error(errData.error || `HTTP error! status: ${res.status}`)
    }
    const data = await res.json()
    findings.value = data.findings || []
  } catch (e) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}

const triggerScan = async () => {
  scanning.value = true
  await fetchScan()
  scanning.value = false
}

const onNamespaceChange = () => {
  triggerScan()
}

onMounted(async () => {
  await fetchNamespaces()
  fetchScan()
})
</script>

<template>
  <div class="min-h-screen bg-true-black text-gray-200 font-inter">
    <!-- Header -->
    <header class="border-b border-gray-800 bg-dark-charcoal/50 backdrop-blur sticky top-0 z-50">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex justify-between items-center">
        <div class="flex items-center space-x-3">
          <div class="w-8 h-8 rounded-md bg-redhat flex items-center justify-center">
            <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
            </svg>
          </div>
          <div>
            <h1 class="text-xl font-bold text-white tracking-tight">Service Mesh 3 Migration</h1>
            <p class="text-xs text-gray-400">Readiness Assessment</p>
          </div>
        </div>
        <div class="flex items-center space-x-4">
          <div class="flex items-center space-x-2">
            <label for="namespace-select" class="text-xs font-semibold text-gray-500 uppercase tracking-widest">Namespace</label>
            <select 
              id="namespace-select"
              v-model="selectedNamespace" 
              @change="onNamespaceChange"
              class="bg-gray-800 border border-gray-700 text-white text-sm rounded-md focus:ring-redhat focus:border-redhat block p-2 transition-all outline-none"
            >
              <option disabled value="">Select a namespace</option>
              <option v-for="ns in allowedNamespaces" :key="ns" :value="ns">{{ ns }}</option>
            </select>
          </div>
          <button 
            @click="triggerScan" 
            :disabled="scanning || loading"
            class="bg-redhat hover:bg-red-700 text-white px-4 py-2 rounded-md font-medium text-sm transition-all duration-200 shadow-[0_0_15px_rgba(238,0,0,0.3)] hover:shadow-[0_0_25px_rgba(238,0,0,0.5)] disabled:opacity-50"
          >
            {{ scanning ? 'Scanning...' : 'Rescan' }}
          </button>
        </div>
      </div>
    </header>

    <!-- Main Content -->
    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      
      <!-- Metrics overview -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <div class="bg-dark-charcoal rounded-xl p-6 border border-gray-800 flex items-center space-x-4">
           <div class="p-3 bg-redhat/20 rounded-lg text-redhat">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"></path></svg>
           </div>
           <div>
            <p class="text-sm text-gray-400 font-medium">Issues Found</p>
            <p class="text-3xl font-bold text-white">{{ findings.length }}</p>
           </div>
        </div>
      </div>

      <div v-if="error" class="mb-6 p-4 rounded-lg bg-red-900/40 border border-red-500/50 text-red-200">
        <strong class="font-bold">Error connecting to discovery engine: </strong>
        <span>{{ error }}</span>
      </div>

      <!-- Findings List -->
      <div v-if="loading && !findings.length" class="flex justify-center py-20">
        <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-redhat"></div>
      </div>

      <div v-else-if="!loading && findings.length === 0" class="text-center py-20 bg-dark-charcoal rounded-xl border border-gray-800">
        <svg class="w-16 h-16 mx-auto text-green-500 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
        <h3 class="text-xl font-medium text-white mb-2">Cluster is Ready</h3>
        <p class="text-gray-400">No migration blockers found. You are ready to upgrade to OSSM 3.0.</p>
      </div>

      <div v-else class="bg-dark-charcoal rounded-xl border border-gray-800 overflow-hidden">
        <div class="px-6 py-4 border-b border-gray-800 bg-gray-900/50">
          <h2 class="text-lg font-semibold text-white">Migration Findings</h2>
        </div>
        <ul class="divide-y divide-gray-800">
          <li v-for="(finding, index) in findings" :key="index" class="hover:bg-gray-800/30 transition-colors duration-150">
            <div class="px-6 py-4 cursor-pointer" @click="toggleRow(index)">
              <div class="flex items-center justify-between">
                <div class="flex items-center space-x-3">
                  <span :class="{
                    'bg-red-900/50 text-red-400 border border-red-800': finding.severity === 'High',
                    'bg-yellow-900/50 text-yellow-400 border border-yellow-800': finding.severity === 'Medium',
                    'bg-blue-900/50 text-blue-400 border border-blue-800': finding.severity === 'Low'
                  }" class="px-2.5 py-1 rounded text-xs font-semibold uppercase tracking-wider">
                    {{ finding.severity }}
                  </span>
                  <div>
                    <h3 class="text-sm font-medium text-white">{{ finding.message }}</h3>
                    <p class="text-xs text-gray-500 mt-1">
                      <span class="font-mono text-gray-400">{{ finding.kind }}</span> &middot; {{ finding.namespace }}/{{ finding.resource_name }}
                    </p>
                  </div>
                </div>
                <div>
                  <svg class="w-5 h-5 text-gray-500 transform transition-transform duration-200" :class="{'rotate-180': expandedRows.has(index)}" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path></svg>
                </div>
              </div>
            </div>

            <!-- Expandable Remediation Guide -->
            <div v-show="expandedRows.has(index)" class="px-6 pb-6 pt-2 bg-gray-900/30 border-t border-gray-800/50">
              <div class="pl-[76px]">
                <p class="text-sm text-gray-300 mb-4">{{ finding.remediation.description }}</p>
                
                <!-- YAML block -->
                <div v-if="finding.remediation.yaml" class="mb-4">
                  <h4 class="text-xs font-semibold text-gray-500 uppercase tracking-widest mb-2">YAML Template</h4>
                  <div class="relative group">
                    <pre class="bg-black/50 p-4 rounded-lg overflow-x-auto text-sm font-mono text-green-400 border border-gray-700">{{ finding.remediation.yaml }}</pre>
                    <button @click.stop="copyToClipboard(finding.remediation.yaml)" class="absolute top-2 right-2 p-1.5 rounded-md bg-gray-800 text-gray-400 opacity-0 group-hover:opacity-100 transition-opacity hover:text-white hover:bg-gray-700">
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"></path></svg>
                    </button>
                  </div>
                </div>

                <!-- Commands block -->
                <div v-if="finding.remediation.commands && finding.remediation.commands.length > 0" class="mb-4">
                  <h4 class="text-xs font-semibold text-gray-500 uppercase tracking-widest mb-2">Execution Scripts</h4>
                  <div v-for="(cmd, cIdx) in finding.remediation.commands" :key="cIdx" class="relative group mb-2 last:mb-0">
                    <pre class="bg-black/50 p-3 rounded-lg overflow-x-auto text-sm font-mono text-blue-300 border border-gray-700 flex select-all"><code><span class="text-gray-600 select-none">$ </span>{{ cmd }}</code></pre>
                    <button @click.stop="copyToClipboard(cmd)" class="absolute top-2 right-2 p-1.5 rounded-md bg-gray-800 text-gray-400 opacity-0 group-hover:opacity-100 transition-opacity hover:text-white hover:bg-gray-700">
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"></path></svg>
                    </button>
                  </div>
                </div>

                <!-- External Docs -->
                <div v-if="finding.remediation.docs_links && finding.remediation.docs_links.length > 0" class="mt-4 flex space-x-2">
                  <a v-for="(link, lIdx) in finding.remediation.docs_links" :key="lIdx" :href="link" target="_blank" rel="noopener noreferrer" class="inline-flex items-center px-3 py-1.5 rounded bg-gray-800 hover:bg-gray-700 text-xs font-medium text-blue-400 transition-colors">
                    <svg class="w-3.5 h-3.5 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"></path></svg>
                    View Official Docs
                  </a>
                </div>

              </div>
            </div>
          </li>
        </ul>
      </div>

    </main>
  </div>
</template>
