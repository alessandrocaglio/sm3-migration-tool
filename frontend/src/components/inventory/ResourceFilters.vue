<script setup>
defineProps({
  filters: { type: Object, required: true },
  options: { type: Object, required: true },
})

const emit = defineEmits(['update:search', 'update:namespace', 'update:kind', 'update:relevance'])
</script>

<template>
  <section class="content-card">
    <div class="section-heading">
      <div>
        <p class="section-heading__eyebrow">Inventory Filters</p>
        <h2 class="section-heading__title">Find resources quickly</h2>
      </div>
    </div>

    <div class="inventory-filters">
      <label class="picker">
        <span class="picker__label">Search</span>
        <input
          class="picker__input"
          type="search"
          :value="filters.search"
          placeholder="Search by name, namespace, or kind"
          @input="emit('update:search', $event.target.value)"
        />
      </label>

      <label class="picker">
        <span class="picker__label">Namespace</span>
        <select class="picker__select" :value="filters.namespace" @change="emit('update:namespace', $event.target.value)">
          <option value="all">All namespaces</option>
          <option v-for="namespace in options.namespaces" :key="namespace" :value="namespace">{{ namespace }}</option>
        </select>
      </label>

      <label class="picker">
        <span class="picker__label">Kind</span>
        <select class="picker__select" :value="filters.kind" @change="emit('update:kind', $event.target.value)">
          <option value="all">All kinds</option>
          <option v-for="kind in options.kinds" :key="kind" :value="kind">{{ kind }}</option>
        </select>
      </label>

      <label class="picker">
        <span class="picker__label">Migration Area</span>
        <select class="picker__select" :value="filters.relevance" @change="emit('update:relevance', $event.target.value)">
          <option value="all">All areas</option>
          <option v-for="item in options.relevance" :key="item" :value="item">{{ item }}</option>
        </select>
      </label>
    </div>
  </section>
</template>
