<script setup>
import { ref, computed, onMounted } from 'vue'
import yaml from 'js-yaml'
import CodeExplorer from './components/CodeExplorer.vue'

const findings = ref([])
const allChecks = ref([])
const meshNamespaces = ref([])
const rawResources = ref({})
const loading = ref(true)
const scanning = ref(false)
const error = ref(null)

const allowedNamespaces = ref([])
const selectedNamespace = ref('')
const selectedResourceNamespace = ref('') // Namespace filter within the Resources tab
const activeTab = ref('assessment') // assessment, remediations, resources, namespaces

const expandedRows = ref(new Set())
const selectedResourceType = ref('')
const selectedResource = ref(null)
const selectedFormat = ref('yaml') // yaml, json

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
    
    if (allowedNamespaces.value.length === 1 && !selectedNamespace.value) {
      selectedNamespace.value = allowedNamespaces.value[0]
      fetchScan()
    }
  } catch (e) {
    console.error(e)
  }
}

const fetchScan = async () => {
  if (!selectedNamespace.value) {
    findings.value = []
    allChecks.value = []
    loading.value = false
    return
  }

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
    allChecks.value = data.checks || []
    meshNamespaces.value = data.mesh_namespaces || []
    rawResources.value = data.resources || {}

    // Force-reset the sub-filter to the newly selected control plane namespace
    selectedResourceNamespace.value = selectedNamespace.value
  } catch (e) {
    error.value = e.message
    findings.value = []
  } finally {
    loading.value = false
  }
}

const triggerScan = async () => {
  if (!selectedNamespace.value) return
  scanning.value = true
  await fetchScan()
  scanning.value = false
}

const onNamespaceChange = () => {
  findings.value = []
  selectedResourceNamespace.value = '' // Reset sub-filter
  fetchScan()
}

// Filtered resources for the Resources Tab
const filteredResourcesByType = computed(() => {
  if (!rawResources.value || !selectedResourceNamespace.value) return {}
  
  const result = {}
  const keys = [
    'sm2_control_planes', 'smmr', 'smms', 
    'virtual_services', 'gateways', 'service_entries', 
    'authorization_policies', 'peer_authentications', 'sidecars', 'telemetries',
    'routes'
  ]
  
  keys.forEach(key => {
    const items = rawResources.value[key] || []
    const filtered = items.filter(r => r.namespace === selectedResourceNamespace.value)
    if (filtered.length > 0) {
      result[key] = filtered
    }
  })
  
  return result
})

const formattedResource = computed(() => {
  if (!selectedResource.value) return ''
  if (selectedFormat.value === 'yaml') {
    return yaml.dump(selectedResource.value, { indent: 2, lineWidth: -1 })
  }
  return JSON.stringify(selectedResource.value, null, 2)
})

const selectResource = (type, res) => {
  selectedResourceType.value = type
  selectedResource.value = res
}

onMounted(async () => {
  loading.value = true
  await fetchNamespaces()
  if (!selectedNamespace.value) {
    loading.value = false
  }
})

// Metrics
const successCount = computed(() => allChecks.value.filter(c => c.status === 'Success').length)
const failureCount = computed(() => allChecks.value.filter(c => c.status === 'Failure').length)
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
            <p class="text-xs text-gray-400">Readiness Assistant</p>
          </div>
        </div>
        <div class="flex items-center space-x-4">
          <div class="flex items-center space-x-2">
            <label for="namespace-select" class="text-xs font-semibold text-gray-500 uppercase tracking-widest">Control Plane</label>
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
      <div v-if="selectedNamespace" class="grid grid-cols-1 md:grid-cols-4 gap-4 mb-8">
        <div class="bg-dark-charcoal rounded-xl p-4 border border-gray-800 flex items-center space-x-4">
           <div class="p-2.5 bg-green-500/20 rounded-lg text-green-500">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path></svg>
           </div>
           <div>
            <p class="text-[10px] text-gray-500 font-bold uppercase tracking-wider">Passed Checks</p>
            <p class="text-2xl font-bold text-white">{{ successCount }}</p>
           </div>
        </div>
        <div class="bg-dark-charcoal rounded-xl p-4 border border-gray-800 flex items-center space-x-4">
           <div class="p-2.5 bg-redhat/20 rounded-lg text-redhat">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"></path></svg>
           </div>
           <div>
            <p class="text-[10px] text-gray-500 font-bold uppercase tracking-wider">Issues Found</p>
            <p class="text-2xl font-bold text-white">{{ failureCount }}</p>
           </div>
        </div>
        <div class="bg-dark-charcoal rounded-xl p-4 border border-gray-800 flex items-center space-x-4">
           <div class="p-2.5 bg-blue-500/20 rounded-lg text-blue-500">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"></path></svg>
           </div>
           <div>
            <p class="text-[10px] text-gray-500 font-bold uppercase tracking-wider">Mesh Namespaces</p>
            <p class="text-2xl font-bold text-white">{{ meshNamespaces.length }}</p>
           </div>
        </div>
        <div class="bg-dark-charcoal rounded-xl p-4 border border-gray-800 flex items-center space-x-4">
           <div class="p-2.5 bg-purple-500/20 rounded-lg text-purple-500">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4m0 5c0 2.21-3.582 4-8 4s-8-1.79-8-4"></path></svg>
           </div>
           <div>
            <p class="text-[10px] text-gray-500 font-bold uppercase tracking-wider">Total Resources</p>
            <p class="text-2xl font-bold text-white">{{ findings.length + successCount }}</p>
           </div>
        </div>
      </div>

      <!-- Navigation Tabs -->
      <div v-if="selectedNamespace" class="flex space-x-1 bg-gray-900 p-1 rounded-lg mb-6 max-w-2xl border border-gray-800">
        <button 
          v-for="tab in ['assessment', 'remediations', 'resources', 'namespaces']" 
          :key="tab"
          @click="activeTab = tab"
          :class="[
            'flex-1 py-2 px-4 text-xs font-semibold rounded-md transition-all uppercase tracking-widest',
            activeTab === tab ? 'bg-gray-800 text-white shadow-sm' : 'text-gray-500 hover:text-gray-300'
          ]"
        >
          {{ tab }}
        </button>
      </div>

      <div v-if="error" class="mb-6 p-4 rounded-lg bg-red-900/40 border border-red-500/50 text-red-200">
        <strong class="font-bold">Error: </strong>
        <span>{{ error }}</span>
      </div>

      <!-- Content Sections -->
      <div v-if="loading" class="flex justify-center py-20">
        <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-redhat"></div>
      </div>

      <div v-else-if="!selectedNamespace" class="text-center py-20 bg-dark-charcoal/30 rounded-xl border border-dashed border-gray-800">
        <svg class="w-16 h-16 mx-auto text-gray-600 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 21h7a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v11m0 5l4.879-4.879m0 0a3 3 0 104.243-4.242 3 3 0 00-4.243 4.242z"></path></svg>
        <h3 class="text-xl font-medium text-white mb-2">Ready for Assessment</h3>
        <p class="text-gray-400">Please select a control plane to begin the migration audit.</p>
      </div>

      <!-- ASESSMENT TAB -->
      <div v-else-if="activeTab === 'assessment'" class="space-y-4">
        <div class="bg-dark-charcoal rounded-xl border border-gray-800 overflow-hidden">
          <table class="w-full text-left text-sm">
            <thead class="bg-gray-900/50 text-gray-400 uppercase text-[10px] tracking-widest border-b border-gray-800">
              <tr>
                <th class="px-6 py-3 font-semibold">Status</th>
                <th class="px-6 py-3 font-semibold">Validation Check</th>
                <th class="px-6 py-3 font-semibold">Target Resource</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-800">
              <tr v-for="(check, idx) in allChecks" :key="idx" class="hover:bg-gray-800/30">
                <td class="px-6 py-4">
                  <span v-if="check.status === 'Success'" class="text-green-500 bg-green-500/10 px-2 py-1 rounded text-[10px] font-bold">✅ PASSED</span>
                  <span v-else class="text-red-500 bg-red-500/10 px-2 py-1 rounded text-[10px] font-bold">❌ FAILED</span>
                </td>
                <td class="px-6 py-4 font-medium text-white">{{ check.title }}</td>
                <td class="px-6 py-4 text-gray-500 font-mono text-xs">{{ check.target }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- REMEDIATIONS TAB -->
      <div v-else-if="activeTab === 'remediations'">
        <div v-if="findings.length === 0" class="text-center py-20 bg-dark-charcoal rounded-xl border border-gray-800">
          <svg class="w-16 h-16 mx-auto text-green-500 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
          <h3 class="text-xl font-medium text-white mb-2">No Blockers Found</h3>
          <p class="text-gray-400">All compatibility checks passed for <span class="text-white font-mono">{{ selectedNamespace }}</span>.</p>
        </div>
        <div v-else class="bg-dark-charcoal rounded-xl border border-gray-800 overflow-hidden">
          <ul class="divide-y divide-gray-800">
            <li v-for="(finding, index) in findings" :key="index" class="hover:bg-gray-800/30 transition-colors">
              <div class="px-6 py-4 cursor-pointer" @click="toggleRow(index)">
                <div class="flex items-center justify-between">
                  <div class="flex items-center space-x-3">
                    <span :class="finding.severity === 'High' ? 'bg-red-900/40 text-red-500' : 'bg-yellow-900/40 text-yellow-500'" class="px-2.5 py-1 rounded text-[10px] font-bold uppercase tracking-wider">
                      {{ finding.severity }}
                    </span>
                    <div>
                      <h3 class="text-sm font-medium text-white">{{ finding.message }}</h3>
                      <p class="text-xs text-gray-500 mt-1">{{ finding.kind }} &middot; {{ finding.namespace }}/{{ finding.resource_name }}</p>
                    </div>
                  </div>
                  <svg class="w-5 h-5 text-gray-600 transition-transform" :class="{'rotate-180': expandedRows.has(index)}" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path></svg>
                </div>
              </div>
              <div v-show="expandedRows.has(index)" class="px-6 pb-6 pt-2 bg-black/20 border-t border-gray-800/50">
                <div class="pl-[68px]">
                  <p class="text-xs text-gray-400 mb-4 leading-relaxed">{{ finding.remediation.description }}</p>
                  <div v-if="finding.remediation.commands?.length" class="space-y-2">
                    <div v-for="(cmd, cIdx) in finding.remediation.commands" :key="cIdx" class="relative group">
                      <pre class="bg-black/40 p-3 rounded-lg text-xs font-mono text-blue-300 border border-gray-800 select-all overflow-x-auto">{{ cmd }}</pre>
                      <button @click.stop="copyToClipboard(cmd)" class="absolute top-2 right-2 opacity-0 group-hover:opacity-100 p-1.5 bg-gray-800 rounded text-gray-400 hover:text-white transition-opacity"><svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"></path></svg></button>
                    </div>
                  </div>
                  <a v-for="link in finding.remediation.docs_links" :key="link" :href="link" target="_blank" class="inline-block mt-4 text-[10px] font-bold text-redhat uppercase tracking-widest hover:underline">View Migration Guide &rarr;</a>
                </div>
              </div>
            </li>
          </ul>
        </div>
      </div>

      <!-- Resources Tab -->
      <div v-else-if="activeTab === 'resources'" class="grid grid-cols-1 lg:grid-cols-12 gap-8">
        <div class="lg:col-span-4 space-y-6">
          <!-- Namespace Selector -->
          <div class="bg-dark-charcoal rounded-xl p-4 border border-gray-800">
            <label class="block text-[10px] font-bold text-gray-500 uppercase tracking-widest mb-2">Namespace Filter</label>
            <select v-model="selectedResourceNamespace" class="w-full bg-black/40 border-gray-700 rounded-lg text-sm text-gray-200 focus:ring-redhat focus:border-redhat px-3 py-2.5 outline-none transition-all">
               <option v-for="ns in meshNamespaces" :key="ns" :value="ns">{{ ns }}</option>
            </select>
          </div>

          <div class="bg-dark-charcoal rounded-xl border border-gray-800 overflow-hidden">
            <div class="px-5 py-4 bg-gray-900/50 border-b border-gray-800">
              <h3 class="text-xs font-bold text-white uppercase tracking-wider">Mesh Resources</h3>
            </div>
            <div class="p-2 max-h-[600px] overflow-y-auto custom-scrollbar">
              <div v-for="(resources, type) in filteredResourcesByType" :key="type" class="mb-4 last:mb-0">
                <div class="px-3 py-1.5 text-[10px] font-black text-gray-500 uppercase tracking-tighter">{{ type }}</div>
                <div class="space-y-1 mt-1">
                  <button 
                    v-for="r in resources" 
                    :key="r.name"
                    @click="selectResource(type, r)"
                    :class="selectedResource?.name === r.name && selectedResourceType === type ? 'bg-redhat/10 text-redhat border-redhat/20' : 'text-gray-400 hover:bg-gray-800/50 border-transparent'"
                    class="w-full text-left px-3 py-2 rounded-lg text-xs font-medium border transition-all truncate"
                  >
                    {{ r.name }}
                  </button>
                </div>
              </div>
              <div v-if="Object.keys(filteredResourcesByType).length === 0" class="p-6 text-center text-gray-600 text-xs italic">
                No resources found in this namespace
              </div>
            </div>
          </div>
        </div>

        <div class="lg:col-span-8">
          <div v-if="selectedResource" class="bg-dark-charcoal rounded-xl border border-gray-800 overflow-hidden h-full flex flex-col">
            <div class="px-6 py-4 bg-gray-900/50 border-b border-gray-800 flex justify-between items-center">
              <div>
                <h3 class="text-sm font-bold text-white">{{ selectedResource.name }}</h3>
                <p class="text-[10px] text-gray-500 font-mono">{{ selectedResource.namespace }} / {{ selectedResourceType }}</p>
              </div>
              <div class="flex items-center space-x-4">
                <button @click="copyToClipboard(formattedResource)" class="text-[10px] bg-redhat/20 text-redhat hover:bg-redhat/30 border border-redhat/30 font-bold px-3 py-1.5 rounded-md transition-colors uppercase tracking-widest">Copy to Clipboard</button>
                
                <!-- Format Toggle -->
                <div class="flex bg-black/40 p-0.5 rounded-lg border border-gray-700">
                  <button 
                    @click="selectedFormat = 'yaml'"
                    :class="selectedFormat === 'yaml' ? 'bg-gray-700 text-white shadow-sm' : 'text-gray-500 hover:text-gray-400'"
                    class="px-2.5 py-1 text-[10px] font-bold uppercase rounded-md transition-all"
                  >YAML</button>
                  <button 
                    @click="selectedFormat = 'json'"
                    :class="selectedFormat === 'json' ? 'bg-gray-700 text-white shadow-sm' : 'text-gray-500 hover:text-gray-400'"
                    class="px-2.5 py-1 text-[10px] font-bold uppercase rounded-md transition-all"
                  >JSON</button>
                </div>
              </div>
            </div>
            <div class="flex-1 overflow-hidden">
               <CodeExplorer :modelValue="formattedResource" :language="selectedFormat" />
            </div>
          </div>
          <div v-else class="bg-dark-charcoal rounded-xl border border-gray-800 border-dashed h-64 flex items-center justify-center text-gray-600 font-medium">
            Select a resource from the left to view its definition
          </div>
        </div>
      </div>

      <!-- NAMESPACES TAB -->
      <div v-else-if="activeTab === 'namespaces'">
        <div class="bg-dark-charcoal rounded-xl border border-gray-800 overflow-hidden">
          <div class="px-6 py-4 bg-gray-900/50 border-b border-gray-800">
            <h3 class="text-sm font-bold text-white tracking-tight uppercase">Namespaces in Mesh</h3>
            <p class="text-xs text-gray-500 mt-1">Discovered via ServiceMeshMemberRoll and Control Plane targeting</p>
          </div>
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 p-6">
            <div v-for="ns in meshNamespaces" :key="ns" class="flex items-center space-x-3 p-3 bg-black/20 rounded-lg border border-gray-800 hover:border-gray-700 transition-colors">
              <div class="w-2 h-2 rounded-full bg-blue-500"></div>
              <span class="text-sm font-medium text-gray-300">{{ ns }}</span>
            </div>
          </div>
        </div>
      </div>

    </main>
  </div>
</template>

<style>
@import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap');
body { 
  margin: 0; 
  background-color: #000;
  -webkit-font-smoothing: antialiased;
}

::-webkit-scrollbar {
  width: 8px;
}
::-webkit-scrollbar-track {
  background: #0a0a0a;
}
::-webkit-scrollbar-thumb {
  background: #222;
  border-radius: 4px;
}
::-webkit-scrollbar-thumb:hover {
  background: #333;
}
</style>
