<script setup>
import { computed, ref, watch } from 'vue'
import ResourceDetailPanel from '../components/inventory/ResourceDetailPanel.vue'
import ResourceFilters from '../components/inventory/ResourceFilters.vue'
import ResourceTable from '../components/inventory/ResourceTable.vue'
import { buildInventoryFilterOptions, buildInventoryRows, formatLabel } from '../utils/resources'

const props = defineProps({
  resources: { type: Object, default: () => ({}) },
  meshNamespaces: { type: Array, default: () => [] },
  selectedNamespace: { type: String, default: '' },
})

const filters = ref({
  search: '',
  namespace: 'all',
  kind: 'all',
  relevance: 'all',
})

const selectedResourceId = ref('')

const inventoryRows = computed(() => buildInventoryRows(props.resources))
const rawFilterOptions = computed(() => buildInventoryFilterOptions(inventoryRows.value))
const filterOptions = computed(() => ({
  ...rawFilterOptions.value,
  relevance: rawFilterOptions.value.relevance.map((value) => formatLabel(value)),
}))

const relevanceLookup = computed(() => {
  return rawFilterOptions.value.relevance.reduce((acc, value) => {
    acc[formatLabel(value)] = value
    return acc
  }, {})
})

watch(
  () => props.selectedNamespace,
  (namespace) => {
    filters.value.namespace = namespace || 'all'
  },
  { immediate: true },
)

const filteredRows = computed(() => {
  const relevanceValue = relevanceLookup.value[filters.value.relevance] || filters.value.relevance

  return inventoryRows.value.filter((row) => {
    const search = filters.value.search.trim().toLowerCase()
    const matchesSearch =
      !search ||
      row.name.toLowerCase().includes(search) ||
      row.namespace.toLowerCase().includes(search) ||
      row.kind.toLowerCase().includes(search)

    const matchesNamespace = filters.value.namespace === 'all' || row.namespace === filters.value.namespace
    const matchesKind = filters.value.kind === 'all' || row.kind === filters.value.kind
    const matchesRelevance = relevanceValue === 'all' || row.relevance === relevanceValue

    return matchesSearch && matchesNamespace && matchesKind && matchesRelevance
  })
})

const selectedResource = computed(() => {
  return filteredRows.value.find((row) => row.id === selectedResourceId.value) || null
})

watch(
  filteredRows,
  (rows) => {
    if (!rows.length) {
      selectedResourceId.value = ''
      return
    }

    if (!rows.some((row) => row.id === selectedResourceId.value)) {
      selectedResourceId.value = rows[0].id
    }
  },
  { immediate: true },
)
</script>

<template>
  <div class="page page--inventory">
    <ResourceFilters
      :filters="filters"
      :options="filterOptions"
      @update:search="filters.search = $event"
      @update:namespace="filters.namespace = $event"
      @update:kind="filters.kind = $event"
      @update:relevance="filters.relevance = $event"
    />

    <section class="inventory-workspace">
      <ResourceTable
        :rows="filteredRows"
        :selected-resource-id="selectedResourceId"
        @select="selectedResourceId = $event.id"
      />

      <ResourceDetailPanel :resource="selectedResource" />
    </section>
  </div>
</template>
