<script setup>
defineProps({
  summary: { type: Object, required: true },
})

function scoreCaption(score) {
  if (score >= 80) return 'Most compatibility checks are in good shape.'
  if (score >= 50) return 'The environment is partly prepared but still needs review.'
  return 'Migration blockers need attention before execution planning.'
}
</script>

<template>
  <section class="readiness-hero">
    <div class="readiness-hero__copy">
      <p class="section-heading__eyebrow">Overview</p>
      <h2 class="readiness-hero__title">
        {{ summary.selectedNamespace || 'No control plane selected' }}
      </h2>
      <p class="readiness-hero__description">{{ summary.nextStep }}</p>

      <div class="readiness-hero__status-row">
        <span :class="['status-chip', `status-chip--${summary.readinessTone}`]">
          {{ summary.readinessState }}
        </span>
        <span class="readiness-hero__caption">{{ scoreCaption(summary.readinessScore) }}</span>
      </div>
    </div>

    <div class="readiness-hero__score">
      <span class="readiness-hero__score-label">Readiness score</span>
      <strong class="readiness-hero__score-value">{{ summary.readinessScore }}</strong>
      <span class="readiness-hero__score-scale">out of 100</span>
    </div>
  </section>
</template>
