<script setup>
defineProps({
  distribution: { type: Array, default: () => [] },
  topResourceTypes: { type: Array, default: () => [] },
})
</script>

<template>
  <section class="content-card">
    <div class="section-heading">
      <div>
        <p class="section-heading__eyebrow">Risk Shape</p>
        <h3 class="section-heading__title">Issue distribution</h3>
      </div>
    </div>

    <div class="risk-distribution">
      <div class="risk-distribution__severity">
        <div v-for="item in distribution" :key="item.key" class="risk-row">
          <div class="risk-row__title">
            <span :class="['status-dot', `status-dot--${item.tone}`]"></span>
            <span>{{ item.label }}</span>
          </div>
          <strong class="risk-row__count">{{ item.count }}</strong>
        </div>
      </div>

      <div class="risk-distribution__resources">
        <p class="risk-distribution__label">Most represented resource groups</p>
        <ul v-if="topResourceTypes.length" class="compact-list">
          <li v-for="item in topResourceTypes" :key="item.key">
            <span>{{ item.label }}</span>
            <strong>{{ item.count }}</strong>
          </li>
        </ul>
        <p v-else class="muted-copy">No resource counts available yet.</p>
      </div>
    </div>
  </section>
</template>
