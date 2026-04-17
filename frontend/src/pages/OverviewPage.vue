<script setup>
defineProps({
  summary: { type: Object, required: true },
  recentFindings: { type: Array, default: () => [] },
})
</script>

<template>
  <div class="page page--overview">
    <section class="hero-card">
      <div>
        <p class="hero-card__eyebrow">Overview</p>
        <h2 class="hero-card__title">
          {{ summary.selectedNamespace || 'No control plane selected' }}
        </h2>
        <p class="hero-card__description">
          Review readiness, high-severity issues, and discovered scope before moving into the action plan.
        </p>
      </div>

      <div class="hero-card__metric">
        <span class="hero-card__metric-label">High severity blockers</span>
        <strong class="hero-card__metric-value">{{ summary.highSeverityFindings }}</strong>
      </div>
    </section>

    <section class="summary-grid">
      <article class="summary-card">
        <span class="summary-card__label">Passed Checks</span>
        <strong class="summary-card__value">{{ summary.passedChecks }}</strong>
      </article>
      <article class="summary-card">
        <span class="summary-card__label">Failed Checks</span>
        <strong class="summary-card__value">{{ summary.failedChecks }}</strong>
      </article>
      <article class="summary-card">
        <span class="summary-card__label">Mesh Namespaces</span>
        <strong class="summary-card__value">{{ summary.meshNamespaceCount }}</strong>
      </article>
      <article class="summary-card">
        <span class="summary-card__label">Discovered Resources</span>
        <strong class="summary-card__value">{{ summary.resourceCount }}</strong>
      </article>
    </section>

    <section class="content-card">
      <div class="section-heading">
        <div>
          <p class="section-heading__eyebrow">Top Findings</p>
          <h3 class="section-heading__title">Recent blockers</h3>
        </div>
      </div>

      <ul v-if="recentFindings.length" class="finding-list">
        <li v-for="finding in recentFindings" :key="`${finding.namespace}/${finding.resource_name}/${finding.message}`" class="finding-list__item">
          <div class="finding-list__header">
            <span :class="['pill', finding.severity === 'High' ? 'pill--danger' : 'pill--neutral']">
              {{ finding.severity }}
            </span>
            <strong>{{ finding.kind }}</strong>
            <span class="finding-list__resource">{{ finding.namespace }}/{{ finding.resource_name }}</span>
          </div>
          <p class="finding-list__message">{{ finding.message }}</p>
        </li>
      </ul>

      <p v-else class="muted-copy">No findings available for the selected control plane.</p>
    </section>
  </div>
</template>
