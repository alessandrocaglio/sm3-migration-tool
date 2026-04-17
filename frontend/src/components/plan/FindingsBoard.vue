<script setup>
import PhaseSection from './PhaseSection.vue'

defineProps({
  phases: { type: Array, default: () => [] },
  selectedFindingId: { type: String, default: '' },
})

const emit = defineEmits(['select'])
</script>

<template>
  <div class="findings-board">
    <PhaseSection
      v-for="phase in phases"
      :key="phase.id"
      :phase="phase"
      :selected-finding-id="selectedFindingId"
      @select="emit('select', $event)"
    />
  </div>
</template>
