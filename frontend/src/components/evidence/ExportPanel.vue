<script setup>
defineProps({
  summary: { type: Object, default: () => ({}) },
  selectedNamespace: { type: String, default: '' },
  lastScannedAt: { type: String, default: '' },
  counts: { type: Object, default: () => ({}) },
})

const emit = defineEmits(['export-json', 'export-markdown'])

function formatTimestamp(value) {
  if (!value) return 'Not scanned yet'
  return new Intl.DateTimeFormat(undefined, { dateStyle: 'medium', timeStyle: 'short' }).format(new Date(value))
}
</script>

<template>
  <section class="content-card evidence-export">
    <div class="section-heading">
      <div>
        <p class="section-heading__eyebrow">Evidence</p>
        <h2 class="section-heading__title">Export and metadata</h2>
      </div>
    </div>

    <div class="evidence-export__meta">
      <div class="evidence-export__meta-item">
        <span class="picker__label">Control Plane</span>
        <strong>{{ selectedNamespace || 'None selected' }}</strong>
      </div>
      <div class="evidence-export__meta-item">
        <span class="picker__label">Readiness</span>
        <strong>{{ summary.readiness?.status || 'unknown' }}</strong>
      </div>
      <div class="evidence-export__meta-item">
        <span class="picker__label">Last Scan</span>
        <strong>{{ formatTimestamp(lastScannedAt) }}</strong>
      </div>
      <div class="evidence-export__meta-item">
        <span class="picker__label">Counts</span>
        <strong>{{ counts.checks_total || 0 }} checks / {{ counts.findings_total || 0 }} findings</strong>
      </div>
    </div>

    <div class="evidence-export__actions">
      <button class="app-button" @click="emit('export-json')">Export JSON</button>
      <button class="app-button" @click="emit('export-markdown')">Export Markdown</button>
    </div>
  </section>
</template>
