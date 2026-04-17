<script setup>
import ControlPlanePicker from '../scan/ControlPlanePicker.vue'

defineProps({
  namespaces: { type: Array, default: () => [] },
  selectedNamespace: { type: String, default: '' },
  busy: { type: Boolean, default: false },
})

const emit = defineEmits(['select-namespace', 'refresh'])
</script>

<template>
  <header class="app-header">
    <div>
      <p class="app-header__eyebrow">Migration Workspace</p>
      <h1 class="app-header__title">Service Mesh 3 Readiness Assistant</h1>
    </div>

    <div class="app-header__controls">
      <ControlPlanePicker
        :namespaces="namespaces"
        :selected-namespace="selectedNamespace"
        :disabled="busy"
        @select="$emit('select-namespace', $event)"
      />
      <button class="app-button app-button--primary" :disabled="busy" @click="$emit('refresh')">
        {{ busy ? 'Scanning…' : 'Rescan' }}
      </button>
    </div>
  </header>
</template>
