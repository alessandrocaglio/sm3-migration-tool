<script setup>
defineProps({
  namespaces: { type: Array, default: () => [] },
  selectedNamespace: { type: String, default: '' },
  disabled: { type: Boolean, default: false },
})

const emit = defineEmits(['select'])
</script>

<template>
  <label class="picker">
    <span class="picker__label">Control Plane</span>
    <select
      class="picker__select"
      :value="selectedNamespace"
      :disabled="disabled"
      @change="emit('select', $event.target.value)"
    >
      <option disabled value="">Select a namespace</option>
      <option v-for="namespace in namespaces" :key="namespace" :value="namespace">
        {{ namespace }}
      </option>
    </select>
  </label>
</template>
