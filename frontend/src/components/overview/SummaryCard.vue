<script setup>
defineProps({
  label: { type: String, required: true },
  value: { type: [String, Number], required: true },
  hint: { type: String, default: '' },
  tone: { type: String, default: 'neutral' },
})
</script>

<template>
  <article :class="['summary-card', `summary-card--${tone}`]">
    <span class="summary-card__label">{{ label }}</span>
    <strong class="summary-card__value">{{ value }}</strong>
    <span v-if="hint" class="summary-card__hint">{{ hint }}</span>
  </article>
</template>
