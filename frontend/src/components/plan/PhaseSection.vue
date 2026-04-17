<script setup>
import FindingCard from './FindingCard.vue'

defineProps({
  phase: { type: Object, required: true },
  selectedFindingId: { type: String, default: '' },
})

const emit = defineEmits(['select'])
</script>

<template>
  <section class="phase-section">
    <header class="phase-section__header">
      <div>
        <p class="section-heading__eyebrow">{{ phase.title }}</p>
        <h3 class="phase-section__title">{{ phase.items.length }} action{{ phase.items.length === 1 ? '' : 's' }}</h3>
      </div>
      <div class="phase-section__counts">
        <span v-if="phase.counts.high" class="status-chip status-chip--danger">{{ phase.counts.high }} High</span>
        <span v-if="phase.counts.medium" class="status-chip status-chip--warning">{{ phase.counts.medium }} Medium</span>
        <span v-if="phase.counts.low" class="status-chip status-chip--neutral">{{ phase.counts.low }} Low</span>
      </div>
    </header>

    <div class="phase-section__items">
      <FindingCard
        v-for="finding in phase.items"
        :key="finding.id"
        :finding="finding"
        :selected="finding.id === selectedFindingId"
        @select="emit('select', $event)"
      />
    </div>
  </section>
</template>
