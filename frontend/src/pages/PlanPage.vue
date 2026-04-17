<script setup>
import { computed, ref } from 'vue'

const props = defineProps({
  findings: { type: Array, default: () => [] },
})

const expandedRows = ref(new Set())

const sortedFindings = computed(() => {
  const severityOrder = { High: 0, Medium: 1, Low: 2, Info: 3 }

  return [...props.findings].sort((left, right) => {
    return (severityOrder[left.severity] ?? 99) - (severityOrder[right.severity] ?? 99)
  })
})

function toggleRow(index) {
  if (expandedRows.value.has(index)) {
    expandedRows.value.delete(index)
  } else {
    expandedRows.value.add(index)
  }
}

async function copyToClipboard(text) {
  try {
    await navigator.clipboard.writeText(text)
  } catch (error) {
    console.error('Failed to copy text:', error)
  }
}
</script>

<template>
  <div class="page">
    <section class="content-card">
      <div class="section-heading">
        <div>
          <p class="section-heading__eyebrow">Action Plan</p>
          <h2 class="section-heading__title">Remediations</h2>
        </div>
      </div>

      <div v-if="sortedFindings.length === 0" class="empty-inline">
        No blockers were found for the selected control plane.
      </div>

      <ul v-else class="finding-list finding-list--expandable">
        <li v-for="(finding, index) in sortedFindings" :key="`${finding.namespace}/${finding.resource_name}/${index}`" class="finding-list__item">
          <button class="finding-list__toggle" @click="toggleRow(index)">
            <div class="finding-list__header">
              <span :class="['pill', finding.severity === 'High' ? 'pill--danger' : 'pill--neutral']">
                {{ finding.severity }}
              </span>
              <strong>{{ finding.kind }}</strong>
              <span class="finding-list__resource">{{ finding.namespace }}/{{ finding.resource_name }}</span>
            </div>
            <p class="finding-list__message">{{ finding.message }}</p>
          </button>

          <div v-if="expandedRows.has(index)" class="finding-detail">
            <p class="finding-detail__description">{{ finding.remediation.description }}</p>

            <div v-if="finding.remediation.commands?.length" class="finding-detail__commands">
              <div v-for="command in finding.remediation.commands" :key="command" class="command-block">
                <pre>{{ command }}</pre>
                <button class="app-button" @click="copyToClipboard(command)">Copy</button>
              </div>
            </div>

            <div v-if="finding.remediation.docs_links?.length" class="finding-detail__links">
              <a v-for="link in finding.remediation.docs_links" :key="link" :href="link" target="_blank" rel="noreferrer">
                {{ link }}
              </a>
            </div>
          </div>
        </li>
      </ul>
    </section>
  </div>
</template>
