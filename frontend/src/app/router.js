import { computed, ref } from 'vue'

const routes = [
  { key: 'overview', label: 'Overview', hash: '#/overview' },
  { key: 'plan', label: 'Action Plan', hash: '#/plan' },
  { key: 'inventory', label: 'Inventory', hash: '#/inventory' },
  { key: 'evidence', label: 'Evidence', hash: '#/evidence' },
]

const routeMap = new Map(routes.map((route) => [route.key, route]))
const currentRoute = ref(resolveRoute(window.location.hash))

function resolveRoute(hash) {
  const normalized = hash.replace(/^#\//, '').trim() || 'overview'
  return routeMap.get(normalized)?.key || 'overview'
}

function ensureHash() {
  const route = routeMap.get(currentRoute.value) || routes[0]
  if (window.location.hash !== route.hash) {
    window.location.hash = route.hash
  }
}

function handleHashChange() {
  currentRoute.value = resolveRoute(window.location.hash)
}

window.addEventListener('hashchange', handleHashChange)
if (!window.location.hash) {
  ensureHash()
}

export function useRouter() {
  const activeRoute = computed(() => routeMap.get(currentRoute.value) || routes[0])

  function navigate(key) {
    const route = routeMap.get(key)
    if (!route) return
    currentRoute.value = route.key
    ensureHash()
  }

  return {
    routes,
    currentRoute,
    activeRoute,
    navigate,
  }
}
