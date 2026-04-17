<script setup>
import { computed, ref, watch } from 'vue'
import FindingDrawer from '../components/plan/FindingDrawer.vue'
import FindingsBoard from '../components/plan/FindingsBoard.vue'

const props = defineProps({
  phases: { type: Array, default: () => [] },
  findingViews: { type: Array, default: () => [] },
  categories: { type: Array, default: () => [] },
  selectedNamespace: { type: String, default: '' },
})

const severityFilter = ref('all')
const categoryFilter = ref('all')
const namespaceFilter = ref('all')
const selectedFindingId = ref('')

const namespaceOptions = computed(() => {
  const namespaces = new Set(
    props.findingViews
      .map((finding) => finding.namespace)
      .filter(Boolean),
  )

  return Array.from(namespaces).sort()
})

const filteredPhases = computed(() => {
  return props.phases
    .map((phase) => {
      const items = phase.items.filter((finding) => {
        const severityMatch = severityFilter.value === 'all' || finding.severity === severityFilter.value
        const categoryMatch = categoryFilter.value === 'all' || finding.category === categoryFilter.value
        const namespaceMatch = namespaceFilter.value === 'all' || finding.namespace === namespaceFilter.value
        return severityMatch && categoryMatch && namespaceMatch
      })

      const counts = {
        high: items.filter((finding) => finding.severity === 'High').length,
        medium: items.filter((finding) => finding.severity === 'Medium').length,
        low: items.filter((finding) => finding.severity === 'Low').length,
        info: items.filter((finding) => finding.severity === 'Info').length,
      }

      return { ...phase, items, counts }
    })
    .filter((phase) => phase.items.length > 0)
})

const selectedFinding = computed(() => {
  return props.findingViews.find((finding) => finding.id === selectedFindingId.value) || null
})

const filteredFindingCount = computed(() => {
  return filteredPhases.value.reduce((total, phase) => total + phase.items.length, 0)
})

watch(
  () => props.findingViews,
  (findings) => {
    if (!findings.length) {
      selectedFindingId.value = ''
      return
    }

    if (!findings.some((finding) => finding.id === selectedFindingId.value)) {
      selectedFindingId.value = findings[0].id
    }
  },
  { immediate: true },
)

watch(filteredPhases, (phases) => {
  const available = phases.flatMap((phase) => phase.items)
  if (!available.length) {
    selectedFindingId.value = ''
    return
  }

  if (!available.some((finding) => finding.id === selectedFindingId.value)) {
    selectedFindingId.value = available[0].id
  }
})

async function copyToClipboard(text) {
  try {
    await navigator.clipboard.writeText(text)
  } catch (error) {
    console.error('Failed to copy text:', error)
  }
}

function clearSelection() {
  selectedFindingId.value = ''
}
</script>

<template>
  <div class="page page--plan">
    <section class="content-card">
      <div class="section-heading">
        <div>
          <p class="section-heading__eyebrow">Action Plan</p>
          <h2 class="section-heading__title">Migration work queue</h2>
          <p class="muted-copy">
            {{ filteredFindingCount }} action{{ filteredFindingCount === 1 ? '' : 's' }}
            for {{ selectedNamespace || 'the selected control plane' }}.
          </p>
        </div>
      </div>

      <div class="plan-filters">
        <label class="picker">
          <span class="picker__label">Severity</span>
          <select v-model="severityFilter" class="picker__select">
            <option value="all">All severities</option>
            <option value="High">High</option>
            <option value="Medium">Medium</option>
            <option value="Low">Low</option>
            <option value="Info">Info</option>
          </select>
        </label>

        <label class="picker">
          <span class="picker__label">Category</span>
          <select v-model="categoryFilter" class="picker__select">
            <option value="all">All categories</option>
            <option v-for="category in categories" :key="category.id" :value="category.id">
              {{ category.label }}
            </option>
          </select>
        </label>

        <label class="picker">
          <span class="picker__label">Namespace</span>
          <select v-model="namespaceFilter" class="picker__select">
            <option value="all">All namespaces</option>
            <option v-for="namespace in namespaceOptions" :key="namespace" :value="namespace">
              {{ namespace }}
            </option>
          </select>
        </label>
      </div>
    </section>

    <div v-if="!filteredPhases.length" class="empty-inline">
      No action items match the current filters.
    </div>

    <section v-else class="plan-layout">
      <FindingsBoard
        :phases="filteredPhases"
        :selected-finding-id="selectedFindingId"
        @select="selectedFindingId = $event.id"
      />

      <FindingDrawer
        :finding="selectedFinding"
        @close="clearSelection"
        @copy="copyToClipboard"
      />
    </section>
  </div>
</template>
