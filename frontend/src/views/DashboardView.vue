<script setup>
import { onMounted } from 'vue'
import { useMonitorStore } from '../stores/monitors'
import MonitorCard from '../components/MonitorCard.vue'

const store = useMonitorStore()

onMounted(() => {
  store.fetchMonitors()
})

async function handleDelete(id) {
  if (!confirm('Delete this monitor and all its check history?')) return
  try {
    await store.remove(id)
  } catch (e) {
    alert(e.message)
  }
}
</script>

<template>
  <div class="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
    <div class="flex items-center justify-between mb-6">
      <div>
        <h1 class="text-lg font-semibold text-gray-900 dark:text-gray-100">Monitors</h1>
        <p class="text-sm text-gray-500 dark:text-gray-400 mt-0.5">{{ store.monitors.length }} URL{{ store.monitors.length !== 1 ? 's' : '' }} being monitored</p>
      </div>
      <router-link
        to="/monitors/new"
        class="inline-flex items-center gap-1.5 bg-indigo-600 text-white text-sm font-medium rounded-lg px-4 py-2 hover:bg-indigo-700 transition-colors"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        Add Monitor
      </router-link>
    </div>

    <div v-if="store.loading" class="text-center py-12 text-sm text-gray-500 dark:text-gray-400">
      Loading monitors...
    </div>

    <div v-else-if="store.error" class="bg-red-50 dark:bg-red-900/30 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-400 text-sm rounded-lg px-4 py-3">
      {{ store.error }}
    </div>

    <div v-else-if="store.monitors.length === 0" class="text-center py-16">
      <svg class="mx-auto w-10 h-10 text-gray-300 dark:text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 10h.01M15 10h.01M9.5 15.5c.83.83 2.17.83 3 0" />
      </svg>
      <p class="mt-3 text-sm text-gray-500 dark:text-gray-400">No monitors yet</p>
      <router-link
        to="/monitors/new"
        class="mt-2 inline-block text-sm text-indigo-600 dark:text-indigo-400 hover:text-indigo-700 dark:hover:text-indigo-300"
      >
        Create your first monitor
      </router-link>
    </div>

    <div v-else class="space-y-3">
      <MonitorCard
        v-for="m in store.monitors"
        :key="m.id"
        :monitor="m"
        @delete="handleDelete"
      />
    </div>
  </div>
</template>
