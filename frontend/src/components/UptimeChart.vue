<script setup>
import { computed } from 'vue'
import { use } from 'echarts/core'
import { BarChart } from 'echarts/charts'
import { GridComponent, TooltipComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import VChart from 'vue-echarts'
import { useThemeStore } from '../stores/theme'

use([BarChart, GridComponent, TooltipComponent, CanvasRenderer])

const props = defineProps({
  timeline: { type: Array, default: () => [] },
})

const theme = useThemeStore()

const option = computed(() => {
  const data = props.timeline
  const textColor = theme.dark ? '#9ca3af' : '#6b7280'
  const lineColor = theme.dark ? '#374151' : '#e5e7eb'

  return {
    tooltip: {
      trigger: 'axis',
      backgroundColor: theme.dark ? '#1f2937' : '#fff',
      borderColor: theme.dark ? '#374151' : '#e5e7eb',
      textStyle: { color: theme.dark ? '#e5e7eb' : '#111827', fontSize: 12 },
      formatter(params) {
        const ok = params[0]
        const fail = params[1]
        const t = new Date(ok.axisValue).toLocaleString()
        return `${t}<br/>Healthy: <b style="color:#10b981">${ok.data[1]}</b><br/>Failed: <b style="color:#ef4444">${fail.data[1]}</b>`
      },
    },
    grid: { left: 50, right: 16, top: 12, bottom: 50 },
    xAxis: {
      type: 'time',
      axisLabel: {
        color: textColor,
        fontSize: 10,
        hideOverlap: true,
        rotate: 30,
        formatter: '{HH}:{mm}',
      },
      axisLine: { lineStyle: { color: lineColor } },
      splitLine: { show: false },
      splitNumber: 5,
    },
    yAxis: {
      type: 'value',
      name: 'checks',
      nameTextStyle: { color: textColor, fontSize: 11 },
      axisLabel: { color: textColor, fontSize: 11 },
      axisLine: { show: false },
      splitLine: { lineStyle: { color: lineColor, type: 'dashed' } },
      minInterval: 1,
    },
    series: [
      {
        name: 'Healthy',
        type: 'bar',
        stack: 'checks',
        barMaxWidth: 12,
        itemStyle: { color: '#10b981', borderRadius: [0, 0, 0, 0] },
        data: data.map((p) => [p.timestamp, p.healthy]),
      },
      {
        name: 'Failed',
        type: 'bar',
        stack: 'checks',
        barMaxWidth: 12,
        itemStyle: { color: '#ef4444', borderRadius: [2, 2, 0, 0] },
        data: data.map((p) => [p.timestamp, p.total - p.healthy]),
      },
    ],
  }
})
</script>

<template>
  <v-chart :option="option" style="height: 220px" autoresize />
</template>
