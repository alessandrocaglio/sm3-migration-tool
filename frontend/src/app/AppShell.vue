<script setup>
import { computed, onMounted } from 'vue'
import AppHeader from '../components/layout/AppHeader.vue'
import SideNav from '../components/layout/SideNav.vue'
import EmptyState from '../components/scan/EmptyState.vue'
import ErrorState from '../components/scan/ErrorState.vue'
import ScanStatusBar from '../components/scan/ScanStatusBar.vue'
import OverviewPage from '../pages/OverviewPage.vue'
import PlanPage from '../pages/PlanPage.vue'
import InventoryPage from '../pages/InventoryPage.vue'
import EvidencePage from '../pages/EvidencePage.vue'
import { useRouter } from './router'
import { useScanStore } from '../stores/scanStore'

const { routes, currentRoute, navigate } = useRouter()
const { state, summary, recentFindings, severityDistribution, namespaceImpact, topResourceTypes, initialize, selectNamespace, refresh } = useScanStore()

const busy = computed(() => state.scanStatus === 'loading')

const currentPage = computed(() => {
  switch (currentRoute.value) {
    case 'plan':
      return PlanPage
    case 'inventory':
      return InventoryPage
    case 'evidence':
      return EvidencePage
    default:
      return OverviewPage
  }
})

const pageProps = computed(() => {
  switch (currentRoute.value) {
    case 'plan':
      return {
        phases: state.data.phases || [],
        findingViews: state.data.finding_views || [],
        categories: state.data.categories || [],
        selectedNamespace: state.selectedNamespace,
      }
    case 'inventory':
      return {
        resources: state.data.resources || {},
        meshNamespaces: state.data.mesh_namespaces || [],
        selectedNamespace: state.selectedNamespace,
      }
    case 'evidence':
      return { checks: state.data.checks || [] }
    default:
      return {
        summary: summary.value,
        recentFindings: recentFindings.value,
        severityDistribution: severityDistribution.value,
        namespaceImpact: namespaceImpact.value,
        topResourceTypes: topResourceTypes.value,
      }
  }
})

onMounted(() => {
  initialize()
})
</script>

<template>
  <div class="workspace">
    <SideNav :routes="routes" :current-route="currentRoute" @navigate="navigate" />

    <div class="workspace__main">
      <AppHeader
        :namespaces="state.allowedNamespaces"
        :selected-namespace="state.selectedNamespace"
        :busy="busy"
        @select-namespace="selectNamespace"
        @refresh="refresh"
      />

      <main class="workspace__content">
        <ScanStatusBar
          :status="state.scanStatus"
          :last-scanned-at="state.lastScannedAt"
          :selected-namespace="state.selectedNamespace"
        />

        <ErrorState v-if="state.error" :message="state.error" />
        <EmptyState
          v-else-if="!state.selectedNamespace && state.scanStatus !== 'loading'"
          title="Select a control plane"
          description="Choose a control plane namespace to load migration findings and discovered resources."
        />
        <EmptyState
          v-else-if="state.scanStatus === 'loading' && !(state.data.checks || []).length"
          title="Loading assessment"
          description="Scanning the cluster and collecting migration data."
        />
        <component :is="currentPage" v-else v-bind="pageProps" />
      </main>
    </div>
  </div>
</template>
