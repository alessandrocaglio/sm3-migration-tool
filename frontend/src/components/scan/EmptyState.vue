<script setup>
defineProps({
  title: { type: String, required: true },
  description: { type: String, required: true },
})
</script>

<template>
  <section class="state-card">
    <h2 class="state-card__title">{{ title }}</h2>
    <p class="state-card__description">{{ description }}</p>
  </section>
</template>
