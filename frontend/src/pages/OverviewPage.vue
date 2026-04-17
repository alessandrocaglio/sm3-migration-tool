<script setup>
import NamespaceImpactList from '../components/overview/NamespaceImpactList.vue'
import ReadinessHero from '../components/overview/ReadinessHero.vue'
import RiskDistribution from '../components/overview/RiskDistribution.vue'
import SummaryCard from '../components/overview/SummaryCard.vue'
import TopBlockers from '../components/overview/TopBlockers.vue'

defineProps({
  summary: { type: Object, required: true },
  recentFindings: { type: Array, default: () => [] },
  severityDistribution: { type: Array, default: () => [] },
  namespaceImpact: { type: Array, default: () => [] },
  topResourceTypes: { type: Array, default: () => [] },
})
</script>

<template>
  <div class="page page--overview">
    <ReadinessHero :summary="summary" />

    <section class="summary-grid">
      <SummaryCard label="Passed Checks" :value="summary.passedChecks" hint="Validations already compatible" tone="success" />
      <SummaryCard label="Failed Checks" :value="summary.failedChecks" hint="Checks requiring remediation" tone="danger" />
      <SummaryCard label="Mesh Namespaces" :value="summary.meshNamespaceCount" hint="Namespaces currently in scope" />
      <SummaryCard label="Discovered Resources" :value="summary.resourceCount" hint="Resources loaded into this assessment" />
    </section>

    <section class="overview-grid">
      <TopBlockers :findings="recentFindings" />
      <RiskDistribution :distribution="severityDistribution" :top-resource-types="topResourceTypes" />
    </section>

    <NamespaceImpactList :items="namespaceImpact" />
  </div>
</template>
