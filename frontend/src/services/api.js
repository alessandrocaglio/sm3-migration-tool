async function getJSON(url, options = {}) {
  const response = await fetch(url, options)
  if (!response.ok) {
    const errorPayload = await response.json().catch(() => ({}))
    throw new Error(errorPayload.error || `Request failed with status ${response.status}`)
  }

  return response.json()
}

export function fetchNamespaces() {
  return getJSON('/api/namespaces')
}

export function fetchScan(namespace) {
  const query = namespace ? `?namespace=${encodeURIComponent(namespace)}` : ''
  return getJSON(`/api/scan${query}`)
}
