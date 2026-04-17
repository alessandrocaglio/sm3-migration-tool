<script setup>
defineProps({
  rows: { type: Array, default: () => [] },
  selectedResourceId: { type: String, default: '' },
})

const emit = defineEmits(['select'])
</script>

<template>
  <section class="content-card inventory-table-card">
    <div class="section-heading">
      <div>
        <p class="section-heading__eyebrow">Inventory</p>
        <h2 class="section-heading__title">Discovered resources</h2>
      </div>
    </div>

    <div v-if="rows.length === 0" class="empty-inline">
      No resources match the current filters.
    </div>

    <div v-else class="inventory-table">
      <div class="inventory-table__head">
        <span>Name</span>
        <span>Kind</span>
        <span>Namespace</span>
      </div>

      <button
        v-for="row in rows"
        :key="row.id"
        :class="['inventory-table__row', { 'inventory-table__row--active': row.id === selectedResourceId }]"
        @click="emit('select', row)"
      >
        <span class="inventory-table__name">{{ row.name }}</span>
        <span>{{ row.kind }}</span>
        <span class="inventory-table__mono">{{ row.namespace }}</span>
      </button>
    </div>
  </section>
</template>
