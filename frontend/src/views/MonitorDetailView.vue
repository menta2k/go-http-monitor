<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getMonitor, getMonitorStatus, getMonitorHistory, deleteMonitor } from '../api/client'
import StatusBadge from '../components/StatusBadge.vue'
import NotificationList from '../components/NotificationList.vue'
import StatsPanel from '../components/StatsPanel.vue'

const route = useRoute()
const router = useRouter()
const id = Number(route.params.id)

const monitor = ref(null)
const latest = ref(null)
const history = ref([])
const historyTotal = ref(0)
const historyLimit = 20
const historyOffset = ref(0)
const loading = ref(true)
const error = ref('')
let interval = null

const currentPage = computed(() => Math.floor(historyOffset.value / historyLimit) + 1)
const totalPages = computed(() => Math.max(1, Math.ceil(historyTotal.value / historyLimit)))
const hasPrev = computed(() => historyOffset.value > 0)
const hasNext = computed(() => historyOffset.value + historyLimit < historyTotal.value)

async function fetchData() {
  try {
    const [m, h] = await Promise.all([
      getMonitor(id),
      getMonitorHistory(id, historyLimit, historyOffset.value),
    ])
    monitor.value = m
    history.value = h.results
    historyTotal.value = h.total
    if (h.results.length > 0 && historyOffset.value === 0) {
      latest.value = h.results[0]
    }
  } catch (e) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}

async function refreshStatus() {
  try {
    const s = await getMonitorStatus(id)
    latest.value = s
    // Only refresh the table if on page 1
    if (historyOffset.value === 0) {
      const h = await getMonitorHistory(id, historyLimit, 0)
      history.value = h.results
      historyTotal.value = h.total
    }
  } catch {
    // ignore
  }
}

async function fetchHistory() {
  try {
    const h = await getMonitorHistory(id, historyLimit, historyOffset.value)
    history.value = h.results
    historyTotal.value = h.total
  } catch {
    // ignore
  }
}

function prevPage() {
  if (!hasPrev.value) return
  historyOffset.value = Math.max(0, historyOffset.value - historyLimit)
  fetchHistory()
}

function nextPage() {
  if (!hasNext.value) return
  historyOffset.value += historyLimit
  fetchHistory()
}

function goToPage(page) {
  historyOffset.value = (page - 1) * historyLimit
  fetchHistory()
}

onMounted(() => {
  fetchData()
  interval = setInterval(refreshStatus, 15000)
})

onUnmounted(() => {
  clearInterval(interval)
})

async function handleDelete() {
  if (!confirm('Delete this monitor and all its check history?')) return
  try {
    await deleteMonitor(id)
    router.push('/')
  } catch (e) {
    alert(e.message)
  }
}

function formatDate(dateStr) {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString()
}

function statusClass(result) {
  if (result.error) return 'text-red-600 dark:text-red-400'
  if (monitor.value && result.status_code !== monitor.value.expected_status) return 'text-amber-600 dark:text-amber-400'
  if (result.body_matched === false) return 'text-amber-600 dark:text-amber-400'
  return 'text-emerald-600 dark:text-emerald-400'
}

// Generate visible page numbers (max 5 around current)
const visiblePages = computed(() => {
  const total = totalPages.value
  const current = currentPage.value
  if (total <= 7) {
    return Array.from({ length: total }, (_, i) => i + 1)
  }
  const pages = new Set([1, total])
  for (let i = Math.max(2, current - 1); i <= Math.min(total - 1, current + 1); i++) {
    pages.add(i)
  }
  const sorted = [...pages].sort((a, b) => a - b)
  const result = []
  for (let i = 0; i < sorted.length; i++) {
    if (i > 0 && sorted[i] - sorted[i - 1] > 1) {
      result.push(null) // gap marker
    }
    result.push(sorted[i])
  }
  return result
})
</script>

<template>
  <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
    <div class="mb-6">
      <router-link to="/" class="text-sm text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200 transition-colors">
        &larr; Back to monitors
      </router-link>
    </div>

    <div v-if="loading" class="text-center py-12 text-sm text-gray-500 dark:text-gray-400">Loading...</div>
    <div v-else-if="error" class="bg-red-50 dark:bg-red-900/30 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-400 text-sm rounded-lg px-4 py-3">{{ error }}</div>

    <template v-else-if="monitor">
      <!-- Header -->
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-6 mb-6">
        <div class="flex items-start justify-between">
          <div>
            <div class="flex items-center gap-2 mb-2">
              <StatusBadge :result="latest" :monitor="monitor" />
              <span class="text-xs text-gray-400 dark:text-gray-500">checks every {{ monitor.interval_seconds }}s</span>
            </div>
            <h1 class="text-lg font-semibold text-gray-900 dark:text-gray-100 break-all">{{ monitor.url }}</h1>
            <div class="mt-3 grid grid-cols-2 sm:grid-cols-4 gap-4 text-sm">
              <div>
                <span class="text-gray-500 dark:text-gray-400">Expected Status</span>
                <p class="font-medium text-gray-900 dark:text-gray-100">{{ monitor.expected_status }}</p>
              </div>
              <div>
                <span class="text-gray-500 dark:text-gray-400">Body Match</span>
                <p class="font-medium text-gray-900 dark:text-gray-100">{{ monitor.body_contains || '-' }}</p>
              </div>
              <div v-if="latest">
                <span class="text-gray-500 dark:text-gray-400">Last Status</span>
                <p class="font-medium text-gray-900 dark:text-gray-100">{{ latest.status_code || '-' }}</p>
              </div>
              <div v-if="latest">
                <span class="text-gray-500 dark:text-gray-400">Response Time</span>
                <p class="font-medium text-gray-900 dark:text-gray-100">{{ latest.response_time_ms }}ms</p>
              </div>
            </div>
          </div>
          <div class="flex items-center gap-1 shrink-0 ml-4">
            <router-link
              :to="`/monitors/${monitor.id}/edit`"
              class="inline-flex items-center gap-1.5 text-sm text-gray-600 dark:text-gray-300 hover:text-indigo-600 dark:hover:text-indigo-400 border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-1.5 transition-colors"
            >
              Edit
            </router-link>
            <button
              @click="handleDelete"
              class="inline-flex items-center gap-1.5 text-sm text-gray-600 dark:text-gray-300 hover:text-red-600 dark:hover:text-red-400 border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-1.5 transition-colors"
            >
              Delete
            </button>
          </div>
        </div>
      </div>

      <!-- Analytics -->
      <StatsPanel :monitor-id="id" class="mb-6" />

      <!-- Notifications -->
      <NotificationList :monitor-id="id" class="mb-6" />

      <!-- History -->
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700">
        <div class="px-6 py-4 border-b border-gray-200 dark:border-gray-700 flex items-center justify-between">
          <h2 class="text-sm font-semibold text-gray-900 dark:text-gray-100">Check History</h2>
          <span class="text-xs text-gray-400 dark:text-gray-500">{{ historyTotal }} total</span>
        </div>
        <div v-if="history.length === 0" class="px-6 py-8 text-center text-sm text-gray-500 dark:text-gray-400">
          No checks recorded yet. The first check will run shortly.
        </div>
        <template v-else>
          <div class="overflow-x-auto">
            <table class="w-full text-sm">
              <thead>
                <tr class="text-left text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                  <th class="px-6 py-3 font-medium">Time</th>
                  <th class="px-6 py-3 font-medium">Status</th>
                  <th class="px-6 py-3 font-medium">Response</th>
                  <th class="px-6 py-3 font-medium">Body Match</th>
                  <th class="px-6 py-3 font-medium">Error</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-100 dark:divide-gray-700">
                <tr v-for="r in history" :key="r.id" class="hover:bg-gray-50 dark:hover:bg-gray-700/50">
                  <td class="px-6 py-3 text-gray-600 dark:text-gray-300 whitespace-nowrap">{{ formatDate(r.checked_at) }}</td>
                  <td class="px-6 py-3 whitespace-nowrap">
                    <span :class="['font-medium', statusClass(r)]">{{ r.status_code || '-' }}</span>
                  </td>
                  <td class="px-6 py-3 text-gray-600 dark:text-gray-300 whitespace-nowrap">{{ r.response_time_ms }}ms</td>
                  <td class="px-6 py-3 whitespace-nowrap">
                    <span v-if="r.body_matched === true" class="text-emerald-600 dark:text-emerald-400">Yes</span>
                    <span v-else-if="r.body_matched === false" class="text-amber-600 dark:text-amber-400">No</span>
                    <span v-else class="text-gray-400 dark:text-gray-500">-</span>
                  </td>
                  <td class="px-6 py-3 text-red-600 dark:text-red-400 max-w-xs truncate">{{ r.error || '-' }}</td>
                </tr>
              </tbody>
            </table>
          </div>

          <!-- Pagination -->
          <div v-if="totalPages > 1" class="px-6 py-3 border-t border-gray-200 dark:border-gray-700 flex items-center justify-between">
            <p class="text-xs text-gray-500 dark:text-gray-400">
              Showing {{ historyOffset + 1 }}–{{ Math.min(historyOffset + historyLimit, historyTotal) }} of {{ historyTotal }}
            </p>
            <div class="flex items-center gap-1">
              <button
                @click="prevPage"
                :disabled="!hasPrev"
                class="px-2 py-1 text-xs rounded border border-gray-300 dark:border-gray-600 text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 disabled:opacity-30 disabled:cursor-not-allowed transition-colors"
              >
                Prev
              </button>
              <template v-for="(p, i) in visiblePages" :key="i">
                <span v-if="p === null" class="px-1 text-xs text-gray-400">...</span>
                <button
                  v-else
                  @click="goToPage(p)"
                  :class="[
                    'w-7 h-7 text-xs rounded transition-colors',
                    p === currentPage
                      ? 'bg-indigo-600 text-white'
                      : 'text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700'
                  ]"
                >
                  {{ p }}
                </button>
              </template>
              <button
                @click="nextPage"
                :disabled="!hasNext"
                class="px-2 py-1 text-xs rounded border border-gray-300 dark:border-gray-600 text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 disabled:opacity-30 disabled:cursor-not-allowed transition-colors"
              >
                Next
              </button>
            </div>
          </div>
        </template>
      </div>
    </template>
  </div>
</template>
