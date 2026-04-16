<script setup>
import { computed, onMounted, watch } from 'vue'
import Prism from 'prismjs'
import 'prismjs/themes/prism-tomorrow.css' // Premium dark theme
import 'prismjs/components/prism-yaml'
import 'prismjs/components/prism-json'

const props = defineProps({
  modelValue: { type: String, default: '' },
  language: { type: String, default: 'yaml' },
})

const highlightedCode = computed(() => {
  const lang = props.language === 'json' ? 'json' : 'yaml'
  return Prism.highlight(props.modelValue, Prism.languages[lang], lang)
})

onMounted(() => {
  Prism.highlightAll()
})

watch(() => props.modelValue, () => {
  // Prism highlight stays reactive via computed highlightedCode
})
</script>

<template>
  <div class="relative group h-full overflow-hidden flex flex-col bg-[#1d1f21]">
    <pre class="flex-1 overflow-auto p-6 text-[13px] leading-relaxed font-mono custom-scrollbar m-0"><code 
        v-html="highlightedCode" 
        :class="`language-${language}`"
      ></code></pre>
  </div>
</template>

<style>
/* Override Prism default backgrounds for deeper integration */
pre[class*="language-"] {
  background: transparent !important;
  margin: 0 !important;
  border: none !important;
  box-shadow: none !important;
}

code[class*="language-"] {
  text-shadow: none !important;
  font-family: 'Fira Code', 'Monaco', 'Menlo', 'Ubuntu Mono', 'Consolas', monospace !important;
}

.custom-scrollbar::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: #0a0a0a;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #333;
  border-radius: 4px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: #444;
}
</style>
