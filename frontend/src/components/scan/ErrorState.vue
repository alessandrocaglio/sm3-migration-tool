<script setup>
defineProps({
  message: { type: String, required: true },
})
</script>

<template>
  <section class="state-card state-card--error" role="alert">
    <h2 class="state-card__title">Scan failed</h2>
    <p class="state-card__description">{{ message }}</p>
  </section>
</template>
