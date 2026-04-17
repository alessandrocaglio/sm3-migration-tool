<script setup>
import { computed, ref } from 'vue'
import CodeExplorer from '../CodeExplorer.vue'
import { formatResourceDocument } from '../../utils/resources'

const props = defineProps({
  resource: { type: Object, default: null },
})

const selectedFormat = ref('yaml')

const formattedRaw = computed(() => {
  if (!props.resource) return ''
  return formatResourceDocument(props.resource.raw, selectedFormat.value)
})

async function copyToClipboard(text) {
  try {
    await navigator.clipboard.writeText(text)
  } catch (error) {
    console.error('Failed to copy text:', error)
  }
}
</script>

<template>
  <section class="content-card inventory-detail-panel">
    <div v-if="resource" class="inventory-detail-panel__content">
      <div class="inventory-detail-panel__header">
        <div>
          <p class="section-heading__eyebrow">Resource Detail</p>
          <h2 class="section-heading__title">{{ resource.name }}</h2>
          <p class="muted-copy">{{ resource.kind }} · {{ resource.namespace }}</p>
        </div>

        <div class="inventory-detail-panel__actions">
          <div class="format-toggle">
            <button :class="['format-toggle__button', { 'format-toggle__button--active': selectedFormat === 'yaml' }]" @click="selectedFormat = 'yaml'">
              YAML
            </button>
            <button :class="['format-toggle__button', { 'format-toggle__button--active': selectedFormat === 'json' }]" @click="selectedFormat = 'json'">
              JSON
            </button>
          </div>
          <button
            class="app-button"
            @click="copyToClipboard(formattedRaw)"
          >
            Copy
          </button>
        </div>
      </div>

      <div class="inventory-detail-panel__viewer">
        <CodeExplorer :model-value="formattedRaw || 'No data available.'" :language="selectedFormat" />
      </div>
    </div>

    <div v-else class="empty-inline">
      Select a resource to inspect its manifest.
    </div>
  </section>
</template>
