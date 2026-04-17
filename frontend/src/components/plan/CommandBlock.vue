<script setup>
defineProps({
  command: { type: String, required: true },
})

const emit = defineEmits(['copy'])
</script>

<template>
  <div class="command-block">
    <pre>{{ command }}</pre>
    <button class="app-button" @click="emit('copy', command)">Copy</button>
  </div>
</template>
