<script setup>
import { computed, ref, watch } from 'vue'
import yaml from 'js-yaml'
import CodeExplorer from '../components/CodeExplorer.vue'

const props = defineProps({
  resources: { type: Object, default: () => ({}) },
  meshNamespaces: { type: Array, default: () => [] },
  selectedNamespace: { type: String, default: '' },
})

const selectedResourceNamespace = ref('')
const selectedResourceType = ref('')
const selectedResource = ref(null)
const selectedFormat = ref('yaml')

const typeOrder = [
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
]

watch(
  () => props.selectedNamespace,
  (namespace) => {
    selectedResourceNamespace.value = namespace || props.meshNamespaces[0] || ''
  },
  { immediate: true },
)

const filteredResourcesByType = computed(() => {
  if (!selectedResourceNamespace.value) return {}

  const result = {}

  for (const type of typeOrder) {
    const items = props.resources[type] || []
    const filtered = items.filter((resource) => resource.namespace === selectedResourceNamespace.value)
    if (filtered.length > 0) {
      result[type] = filtered
    }
  }

  return result
})

const formattedResource = computed(() => {
  if (!selectedResource.value) return ''

  if (selectedFormat.value === 'yaml') {
    return yaml.dump(selectedResource.value, { indent: 2, lineWidth: -1 })
  }

  return JSON.stringify(selectedResource.value, null, 2)
})

function selectResource(type, resource) {
  selectedResourceType.value = type
  selectedResource.value = resource
}

async function copyToClipboard(text) {
  try {
    await navigator.clipboard.writeText(text)
  } catch (error) {
    console.error('Failed to copy text:', error)
  }
}
</script>

<template>
  <div class="page page--inventory">
    <section class="inventory-layout">
      <aside class="content-card inventory-sidebar">
        <div class="section-heading">
          <div>
            <p class="section-heading__eyebrow">Inventory</p>
            <h2 class="section-heading__title">Mesh resources</h2>
          </div>
        </div>

        <label class="picker">
          <span class="picker__label">Namespace Filter</span>
          <select v-model="selectedResourceNamespace" class="picker__select">
            <option v-for="namespace in meshNamespaces" :key="namespace" :value="namespace">
              {{ namespace }}
            </option>
          </select>
        </label>

        <div class="resource-groups">
          <div v-for="(group, type) in filteredResourcesByType" :key="type" class="resource-group">
            <p class="resource-group__title">{{ type }}</p>
            <button
              v-for="resource in group"
              :key="`${type}/${resource.name}`"
              :class="['resource-group__item', { 'resource-group__item--active': selectedResource?.name === resource.name && selectedResourceType === type }]"
              @click="selectResource(type, resource)"
            >
              {{ resource.name }}
            </button>
          </div>
        </div>
      </aside>

      <section class="content-card inventory-detail">
        <div v-if="selectedResource" class="inventory-detail__header">
          <div>
            <h3 class="section-heading__title">{{ selectedResource.name }}</h3>
            <p class="muted-copy">{{ selectedResource.namespace }} / {{ selectedResourceType }}</p>
          </div>
          <div class="inventory-detail__actions">
            <div class="format-toggle">
              <button :class="['format-toggle__button', { 'format-toggle__button--active': selectedFormat === 'yaml' }]" @click="selectedFormat = 'yaml'">
                YAML
              </button>
              <button :class="['format-toggle__button', { 'format-toggle__button--active': selectedFormat === 'json' }]" @click="selectedFormat = 'json'">
                JSON
              </button>
            </div>
            <button class="app-button" @click="copyToClipboard(formattedResource)">Copy</button>
          </div>
        </div>

        <div v-if="selectedResource" class="inventory-detail__body">
          <CodeExplorer :model-value="formattedResource" :language="selectedFormat" />
        </div>
        <div v-else class="empty-inline">
          Select a resource from the inventory list to inspect it.
        </div>
      </section>
    </section>
  </div>
</template>
