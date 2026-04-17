<script setup>
defineProps({
  routes: { type: Array, default: () => [] },
  currentRoute: { type: String, required: true },
})

const emit = defineEmits(['navigate'])
</script>

<template>
  <aside class="side-nav">
    <div class="side-nav__brand">
      <div class="side-nav__logo">SM</div>
      <div>
        <p class="side-nav__title">Mesh Migration</p>
        <p class="side-nav__subtitle">Service Mesh 2 to 3</p>
      </div>
    </div>

    <nav class="side-nav__links" aria-label="Primary">
      <button
        v-for="route in routes"
        :key="route.key"
        :class="['side-nav__link', { 'side-nav__link--active': route.key === currentRoute }]"
        @click="$emit('navigate', route.key)"
      >
        {{ route.label }}
      </button>
    </nav>
  </aside>
</template>
