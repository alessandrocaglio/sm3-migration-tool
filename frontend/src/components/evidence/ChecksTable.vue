<script setup>
defineProps({
  checks: { type: Array, default: () => [] },
})
</script>

<template>
  <section class="content-card">
    <div class="section-heading">
      <div>
        <p class="section-heading__eyebrow">Checks</p>
        <h2 class="section-heading__title">Validation results</h2>
      </div>
    </div>

    <div v-if="checks.length === 0" class="empty-inline">
      No checks match the current filter set.
    </div>

    <table v-else class="checks-table">
      <thead>
        <tr>
          <th>Status</th>
          <th>Validation Check</th>
          <th>Target Resource</th>
          <th>Checker</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(check, index) in checks" :key="`${check.target}-${check.title}-${index}`">
          <td>
            <span :class="['pill', check.status === 'Failure' ? 'pill--danger' : 'pill--success']">
              {{ check.status }}
            </span>
          </td>
          <td>{{ check.title }}</td>
          <td class="checks-table__target">{{ check.target }}</td>
          <td>{{ check.checker || 'n/a' }}</td>
        </tr>
      </tbody>
    </table>
  </section>
</template>
