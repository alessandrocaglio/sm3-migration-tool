<script setup>
defineProps({
  items: { type: Array, default: () => [] },
})
</script>

<template>
  <section class="content-card">
    <div class="section-heading">
      <div>
        <p class="section-heading__eyebrow">Namespace Impact</p>
        <h3 class="section-heading__title">Namespaces in scope</h3>
      </div>
    </div>

    <ul v-if="items.length" class="namespace-list">
      <li v-for="item in items" :key="item.namespace" class="namespace-list__item">
        <div>
          <strong>{{ item.namespace }}</strong>
          <p class="muted-copy">{{ item.issues }} issue{{ item.issues === 1 ? '' : 's' }} detected</p>
        </div>
        <span :class="['status-chip', item.issues > 0 ? 'status-chip--warning' : 'status-chip--success']">
          {{ item.issues > 0 ? 'Needs review' : 'Clear' }}
        </span>
      </li>
    </ul>
    <p v-else class="muted-copy">No mesh namespaces were reported.</p>
  </section>
</template>
