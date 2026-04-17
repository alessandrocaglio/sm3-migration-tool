<script setup>
defineProps({
  finding: { type: Object, required: true },
  selected: { type: Boolean, default: false },
})

const emit = defineEmits(['select'])

function severityTone(severity) {
  if (severity === 'High') return 'danger'
  if (severity === 'Medium') return 'warning'
  if (severity === 'Low') return 'neutral'
  return 'success'
}
</script>

<template>
  <article :class="['plan-card', { 'plan-card--selected': selected }]" @click="emit('select', finding)">
    <div class="plan-card__header">
      <span :class="['status-chip', `status-chip--${severityTone(finding.severity)}`]">
        {{ finding.severity }}
      </span>
      <span class="plan-card__category">{{ finding.category_label }}</span>
    </div>

    <h4 class="plan-card__title">{{ finding.title }}</h4>
    <p class="plan-card__message">{{ finding.message }}</p>

    <div class="plan-card__meta">
      <span>{{ finding.kind }}</span>
      <span class="plan-card__resource">{{ finding.namespace }}/{{ finding.resource_name }}</span>
    </div>
  </article>
</template>
