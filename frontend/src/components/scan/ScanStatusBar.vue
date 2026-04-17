<script setup>
defineProps({
  status: { type: String, default: 'idle' },
  lastScannedAt: { type: String, default: '' },
  selectedNamespace: { type: String, default: '' },
})

function formatTimestamp(value) {
  if (!value) return 'Not scanned yet'

  return new Intl.DateTimeFormat(undefined, {
    dateStyle: 'medium',
    timeStyle: 'short',
  }).format(new Date(value))
}
</script>

<template>
  <section class="status-bar">
    <div class="status-bar__item">
      <span class="status-bar__label">Status</span>
      <strong class="status-bar__value">{{ status }}</strong>
    </div>
    <div class="status-bar__item">
      <span class="status-bar__label">Selected Control Plane</span>
      <strong class="status-bar__value">{{ selectedNamespace || 'None selected' }}</strong>
    </div>
    <div class="status-bar__item">
      <span class="status-bar__label">Last Scan</span>
      <strong class="status-bar__value">{{ formatTimestamp(lastScannedAt) }}</strong>
    </div>
  </section>
</template>
