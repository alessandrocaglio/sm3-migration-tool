<script setup>
defineProps({
  findings: { type: Array, default: () => [] },
})
</script>

<template>
  <section class="content-card">
    <div class="section-heading">
      <div>
        <p class="section-heading__eyebrow">Top Findings</p>
        <h3 class="section-heading__title">Immediate blockers</h3>
      </div>
    </div>

    <ul v-if="findings.length" class="blocker-list">
      <li
        v-for="finding in findings"
        :key="`${finding.namespace}/${finding.resource_name}/${finding.message}`"
        class="blocker-list__item"
      >
        <div class="blocker-list__header">
          <span :class="['status-chip', finding.severity === 'High' ? 'status-chip--danger' : 'status-chip--warning']">
            {{ finding.severity }}
          </span>
          <strong>{{ finding.kind }}</strong>
        </div>
        <p class="blocker-list__message">{{ finding.message }}</p>
        <p class="blocker-list__resource">{{ finding.namespace }}/{{ finding.resource_name }}</p>
      </li>
    </ul>
    <p v-else class="muted-copy">No blockers available for the selected control plane.</p>
  </section>
</template>
