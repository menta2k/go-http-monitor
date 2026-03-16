<script setup>
import { ref, watch, onMounted } from 'vue'
import { getMonitorStats, getMonitorTimeline, getMonitorHistory } from '../api/client'
import ResponseTimeChart from './ResponseTimeChart.vue'
import UptimeChart from './UptimeChart.vue'
import StatusCodeChart from './StatusCodeChart.vue'

const props = defineProps({
  monitorId: { type: Number, required: true },
})

const period = ref('1h')
const periods = [
  { value: '1h', label: '1H' },
  { value: '6h', label: '6H' },
  { value: '24h', label: '24H' },
  { value: '7d', label: '7D' },
  { value: '30d', label: '30D' },
]

const summary = ref(null)
const timeline = ref([])
const history = ref([])
const loading = ref(true)

async function fetchStats() {
  loading.value = true
  try {
    const [s, t, h] = await Promise.all([
      getMonitorStats(props.monitorId, period.value),
      getMonitorTimeline(props.monitorId, period.value, 60),
      getMonitorHistory(props.monitorId, 100),
    ])
    summary.value = s
    timeline.value = t
    history.value = h
  } catch {
    // ignore
  } finally {
    loading.value = false
  }
}

onMounted(fetchStats)
watch(period, fetchStats)

function fmtPct(v) {
  if (v == null) return '-'
  return v.toFixed(1) + '%'
}

function fmtMs(v) {
  if (v == null) return '-'
  return Math.round(v) + 'ms'
}
</script>

<template>
  <div class="space-y-4">
    <!-- Period selector -->
    <div class="flex items-center justify-between">
      <h2 class="text-sm font-semibold text-gray-900 dark:text-gray-100">Analytics</h2>
      <div class="flex items-center gap-1 bg-gray-100 dark:bg-gray-700 rounded-lg p-0.5">
        <button
          v-for="p in periods"
          :key="p.value"
          @click="period = p.value"
          :class="[
            'px-2.5 py-1 text-xs font-medium rounded-md transition-colors',
            period === p.value
              ? 'bg-white dark:bg-gray-600 text-gray-900 dark:text-gray-100 shadow-sm'
              : 'text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200'
          ]"
        >
          {{ p.label }}
        </button>
      </div>
    </div>

    <!-- Summary cards -->
    <div v-if="summary" class="grid grid-cols-2 sm:grid-cols-4 gap-3">
      <div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg px-4 py-3">
        <p class="text-xs text-gray-500 dark:text-gray-400">Uptime</p>
        <p class="text-xl font-semibold" :class="summary.uptime_pct >= 99 ? 'text-emerald-600 dark:text-emerald-400' : summary.uptime_pct >= 95 ? 'text-amber-600 dark:text-amber-400' : 'text-red-600 dark:text-red-400'">
          {{ fmtPct(summary.uptime_pct) }}
        </p>
      </div>
      <div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg px-4 py-3">
        <p class="text-xs text-gray-500 dark:text-gray-400">Avg Response</p>
        <p class="text-xl font-semibold text-gray-900 dark:text-gray-100">{{ fmtMs(summary.avg_response_ms) }}</p>
      </div>
      <div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg px-4 py-3">
        <p class="text-xs text-gray-500 dark:text-gray-400">P95 Response</p>
        <p class="text-xl font-semibold text-gray-900 dark:text-gray-100">{{ fmtMs(summary.p95_response_ms) }}</p>
      </div>
      <div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg px-4 py-3">
        <p class="text-xs text-gray-500 dark:text-gray-400">Checks</p>
        <p class="text-xl font-semibold text-gray-900 dark:text-gray-100">
          {{ summary.total_checks }}
          <span v-if="summary.failed_checks > 0" class="text-sm text-red-500 dark:text-red-400">({{ summary.failed_checks }} failed)</span>
        </p>
      </div>
    </div>

    <!-- Charts -->
    <div v-if="!loading" class="grid grid-cols-1 lg:grid-cols-2 gap-4">
      <!-- Response Time -->
      <div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-4">
        <h3 class="text-xs font-medium text-gray-500 dark:text-gray-400 mb-2">Response Time</h3>
        <ResponseTimeChart :timeline="timeline" />
      </div>

      <!-- Healthy vs Failed -->
      <div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-4">
        <h3 class="text-xs font-medium text-gray-500 dark:text-gray-400 mb-2">Health Status</h3>
        <UptimeChart :timeline="timeline" />
      </div>

      <!-- Status Code Distribution -->
      <div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-4 lg:col-span-2">
        <h3 class="text-xs font-medium text-gray-500 dark:text-gray-400 mb-2">HTTP Status Code Distribution</h3>
        <StatusCodeChart :history="history" />
      </div>
    </div>

    <div v-else class="text-center py-8 text-sm text-gray-500 dark:text-gray-400">
      Loading analytics...
    </div>
  </div>
</template>
