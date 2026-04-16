<script setup>
import { ref, onMounted, onBeforeUnmount, watch } from 'vue'
import * as monaco from 'monaco-editor/esm/vs/editor/editor.api'

const props = defineProps({
  modelValue: { type: String, default: '' },
  language: { type: String, default: 'yaml' },
})

const editorContainer = ref(null)
let editor = null

onMounted(() => {
  editor = monaco.editor.create(editorContainer.value, {
    value: props.modelValue,
    language: props.language,
    theme: 'vs-dark',
    readOnly: true,
    automaticLayout: true,
    minimap: { enabled: false },
    fontSize: 12,
    scrollBeyondLastLine: false,
    lineNumbers: 'on',
    renderWhitespace: 'selection',
    fontFamily: "'Fira Code', 'Monaco', 'Menlo', 'Ubuntu Mono', 'Consolas', monospace",
  })
})

onBeforeUnmount(() => {
  if (editor) {
    editor.dispose()
  }
})

watch(() => props.modelValue, (newVal) => {
  if (editor && newVal !== editor.getValue()) {
    editor.setValue(newVal)
  }
})

watch(() => props.language, (newLang) => {
  if (editor) {
    monaco.editor.setModelLanguage(editor.getModel(), newLang)
  }
})
</script>

<template>
  <div ref="editorContainer" class="w-full h-full min-h-[400px]"></div>
</template>

<style scoped>
/* Monaco Editor usually needs a defined height from parent */
</style>
