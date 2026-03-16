<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { getMonitorStatus } from '../api/client'
import StatusBadge from './StatusBadge.vue'

const props = defineProps({
  monitor: { type: Object, required: true },
})

const emit = defineEmits(['delete'])

const latestResult = ref(null)
let interval = null

async function fetchStatus() {
  try {
    latestResult.value = await getMonitorStatus(props.monitor.id)
  } catch {
    // No results yet
  }
}

onMounted(() => {
  fetchStatus()
  interval = setInterval(fetchStatus, 15000)
})

onUnmounted(() => {
  clearInterval(interval)
})

function formatTime(ms) {
  if (!ms) return '-'
  return `${ms}ms`
}

function formatDate(dateStr) {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString()
}
</script>

<template>
  <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-5 hover:border-gray-300 dark:hover:border-gray-600 transition-colors">
    <div class="flex items-start justify-between gap-4">
      <div class="min-w-0 flex-1">
        <div class="flex items-center gap-2 mb-1">
          <StatusBadge :result="latestResult" :monitor="monitor" />
          <span class="text-xs text-gray-400 dark:text-gray-500">every {{ monitor.interval_seconds }}s</span>
        </div>
        <router-link
          :to="`/monitors/${monitor.id}`"
          class="text-sm font-medium text-gray-900 dark:text-gray-100 hover:text-indigo-600 dark:hover:text-indigo-400 transition-colors truncate block"
        >
          {{ monitor.url }}
        </router-link>
        <div class="mt-2 flex items-center gap-4 text-xs text-gray-500 dark:text-gray-400">
          <span v-if="latestResult">
            Status: <span class="font-medium text-gray-700 dark:text-gray-300">{{ latestResult.status_code }}</span>
            (expect {{ monitor.expected_status }})
          </span>
          <span v-if="latestResult">
            Response: <span class="font-medium text-gray-700 dark:text-gray-300">{{ formatTime(latestResult.response_time_ms) }}</span>
          </span>
          <span v-if="monitor.body_contains" class="truncate max-w-48">
            Match: "<span class="font-medium text-gray-700 dark:text-gray-300">{{ monitor.body_contains }}</span>"
          </span>
        </div>
        <div v-if="latestResult" class="mt-1 text-xs text-gray-400 dark:text-gray-500">
          Last check: {{ formatDate(latestResult.checked_at) }}
        </div>
      </div>
      <div class="flex items-center gap-1 shrink-0">
        <router-link
          :to="`/monitors/${monitor.id}/edit`"
          class="p-1.5 text-gray-400 dark:text-gray-500 hover:text-indigo-600 dark:hover:text-indigo-400 transition-colors rounded"
          title="Edit"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
          </svg>
        </router-link>
        <button
          @click="emit('delete', monitor.id)"
          class="p-1.5 text-gray-400 dark:text-gray-500 hover:text-red-600 dark:hover:text-red-400 transition-colors rounded"
          title="Delete"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
          </svg>
        </button>
      </div>
    </div>
  </div>
</template>
