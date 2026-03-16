<script setup>
import { computed } from 'vue'

const props = defineProps({
  result: { type: Object, default: null },
  monitor: { type: Object, required: true },
})

const status = computed(() => {
  if (!props.result) return { label: 'Pending', color: 'gray' }
  if (props.result.error) return { label: 'Error', color: 'red' }
  if (props.result.status_code !== props.monitor.expected_status) return { label: 'Status Mismatch', color: 'yellow' }
  if (props.result.body_matched === false) return { label: 'Body Mismatch', color: 'yellow' }
  return { label: 'Healthy', color: 'green' }
})

const colorClasses = computed(() => {
  const map = {
    gray: 'bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-300',
    red: 'bg-red-100 text-red-700 dark:bg-red-900/40 dark:text-red-400',
    yellow: 'bg-amber-100 text-amber-700 dark:bg-amber-900/40 dark:text-amber-400',
    green: 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-400',
  }
  return map[status.value.color]
})
</script>

<template>
  <span :class="['inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium', colorClasses]">
    {{ status.label }}
  </span>
</template>
