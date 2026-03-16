<script setup>
import { computed } from 'vue'
import { use } from 'echarts/core'
import { PieChart } from 'echarts/charts'
import { TooltipComponent, LegendComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import VChart from 'vue-echarts'
import { useThemeStore } from '../stores/theme'

use([PieChart, TooltipComponent, LegendComponent, CanvasRenderer])

const props = defineProps({
  history: { type: Array, default: () => [] },
})

const theme = useThemeStore()

const statusColors = {
  '2xx': '#10b981',
  '3xx': '#6366f1',
  '4xx': '#f59e0b',
  '5xx': '#ef4444',
  'Error': '#991b1b',
}

const distribution = computed(() => {
  const counts = {}
  for (const r of props.history) {
    let bucket
    if (r.error && r.error !== '-' && r.error !== '') {
      bucket = 'Error'
    } else if (r.status_code >= 200 && r.status_code < 300) {
      bucket = '2xx'
    } else if (r.status_code >= 300 && r.status_code < 400) {
      bucket = '3xx'
    } else if (r.status_code >= 400 && r.status_code < 500) {
      bucket = '4xx'
    } else if (r.status_code >= 500) {
      bucket = '5xx'
    } else {
      bucket = 'Error'
    }
    counts[bucket] = (counts[bucket] || 0) + 1
  }
  return Object.entries(counts).map(([name, value]) => ({ name, value }))
})

const option = computed(() => {
  const textColor = theme.dark ? '#9ca3af' : '#6b7280'

  return {
    tooltip: {
      trigger: 'item',
      backgroundColor: theme.dark ? '#1f2937' : '#fff',
      borderColor: theme.dark ? '#374151' : '#e5e7eb',
      textStyle: { color: theme.dark ? '#e5e7eb' : '#111827', fontSize: 12 },
      formatter: '{b}: <b>{c}</b> ({d}%)',
    },
    legend: {
      bottom: 0,
      textStyle: { color: textColor, fontSize: 11 },
      itemWidth: 10,
      itemHeight: 10,
    },
    series: [
      {
        type: 'pie',
        radius: ['40%', '70%'],
        center: ['50%', '45%'],
        avoidLabelOverlap: true,
        itemStyle: {
          borderRadius: 4,
          borderColor: theme.dark ? '#1f2937' : '#fff',
          borderWidth: 2,
          color(params) {
            return statusColors[params.name] || '#6b7280'
          },
        },
        label: { show: false },
        emphasis: {
          label: { show: true, fontSize: 13, fontWeight: 'bold', color: textColor },
        },
        data: distribution.value,
      },
    ],
  }
})
</script>

<template>
  <v-chart :option="option" style="height: 220px" autoresize />
</template>
