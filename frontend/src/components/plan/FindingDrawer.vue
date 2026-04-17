<script setup>
import CommandBlock from './CommandBlock.vue'

defineProps({
  finding: { type: Object, default: null },
})

const emit = defineEmits(['close', 'copy'])

function severityTone(severity) {
  if (severity === 'High') return 'danger'
  if (severity === 'Medium') return 'warning'
  if (severity === 'Low') return 'neutral'
  return 'success'
}
</script>

<template>
  <aside class="finding-drawer" :class="{ 'finding-drawer--empty': !finding }">
    <div v-if="finding" class="finding-drawer__content">
      <div class="finding-drawer__header">
        <div>
          <div class="finding-drawer__badges">
            <span :class="['status-chip', `status-chip--${severityTone(finding.severity)}`]">
              {{ finding.severity }}
            </span>
            <span class="pill pill--neutral">{{ finding.phase_label }}</span>
          </div>
          <h3 class="finding-drawer__title">{{ finding.title }}</h3>
          <p class="finding-drawer__subtitle">{{ finding.kind }} · {{ finding.namespace }}/{{ finding.resource_name }}</p>
        </div>
        <button class="app-button" @click="emit('close')">Close</button>
      </div>

      <section class="finding-drawer__section">
        <p class="section-heading__eyebrow">Why It Matters</p>
        <p class="finding-drawer__copy">{{ finding.why_it_matters }}</p>
      </section>

      <section class="finding-drawer__section">
        <p class="section-heading__eyebrow">Detected Issue</p>
        <p class="finding-drawer__copy">{{ finding.message }}</p>
      </section>

      <section class="finding-drawer__section">
        <p class="section-heading__eyebrow">Remediation</p>
        <p class="finding-drawer__copy">{{ finding.remediation.description || 'No remediation text provided.' }}</p>
      </section>

      <section v-if="finding.remediation.commands?.length" class="finding-drawer__section">
        <p class="section-heading__eyebrow">Commands</p>
        <div class="finding-drawer__commands">
          <CommandBlock
            v-for="command in finding.remediation.commands"
            :key="command"
            :command="command"
            @copy="emit('copy', $event)"
          />
        </div>
      </section>

      <section v-if="finding.remediation.yaml" class="finding-drawer__section">
        <p class="section-heading__eyebrow">YAML Example</p>
        <pre class="finding-drawer__code">{{ finding.remediation.yaml }}</pre>
      </section>

      <section v-if="finding.remediation.docs_links?.length" class="finding-drawer__section">
        <p class="section-heading__eyebrow">References</p>
        <div class="finding-drawer__links">
          <a v-for="link in finding.remediation.docs_links" :key="link" :href="link" target="_blank" rel="noreferrer">
            {{ link }}
          </a>
        </div>
      </section>
    </div>

    <div v-else class="finding-drawer__empty">
      <p class="section-heading__eyebrow">Finding Detail</p>
      <h3 class="finding-drawer__title">Select an action</h3>
      <p class="finding-drawer__copy">
        Choose an item from the action plan to inspect its migration rationale, commands, and documentation.
      </p>
    </div>
  </aside>
</template>
