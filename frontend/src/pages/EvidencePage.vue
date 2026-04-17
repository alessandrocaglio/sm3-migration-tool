<script setup>
defineProps({
  checks: { type: Array, default: () => [] },
})
</script>

<template>
  <div class="page">
    <section class="content-card">
      <div class="section-heading">
        <div>
          <p class="section-heading__eyebrow">Evidence</p>
          <h2 class="section-heading__title">Validation results</h2>
        </div>
      </div>

      <div v-if="checks.length === 0" class="empty-inline">
        No checks available yet for the selected control plane.
      </div>

      <table v-else class="checks-table">
        <thead>
          <tr>
            <th>Status</th>
            <th>Validation Check</th>
            <th>Target Resource</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(check, index) in checks" :key="`${check.target}-${index}`">
            <td>
              <span :class="['pill', check.status === 'Failure' ? 'pill--danger' : 'pill--success']">
                {{ check.status }}
              </span>
            </td>
            <td>{{ check.title }}</td>
            <td class="checks-table__target">{{ check.target }}</td>
          </tr>
        </tbody>
      </table>
    </section>
  </div>
</template>
