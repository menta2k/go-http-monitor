<script setup>
import { computed } from 'vue'
import { use } from 'echarts/core'
import { LineChart } from 'echarts/charts'
import { GridComponent, TooltipComponent, DataZoomComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import VChart from 'vue-echarts'
import { useThemeStore } from '../stores/theme'

use([LineChart, GridComponent, TooltipComponent, DataZoomComponent, CanvasRenderer])

const props = defineProps({
  timeline: { type: Array, default: () => [] },
})

const theme = useThemeStore()

const option = computed(() => {
  const data = props.timeline.filter((p) => p.total > 0)
  const textColor = theme.dark ? '#9ca3af' : '#6b7280'
  const lineColor = theme.dark ? '#374151' : '#e5e7eb'

  return {
    tooltip: {
      trigger: 'axis',
      backgroundColor: theme.dark ? '#1f2937' : '#fff',
      borderColor: theme.dark ? '#374151' : '#e5e7eb',
      textStyle: { color: theme.dark ? '#e5e7eb' : '#111827', fontSize: 12 },
      formatter(params) {
        const p = params[0]
        const t = new Date(p.axisValue).toLocaleString()
        return `${t}<br/>Avg: <b>${p.data[1].toFixed(1)}ms</b>`
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
      name: 'ms',
      nameTextStyle: { color: textColor, fontSize: 11 },
      axisLabel: { color: textColor, fontSize: 11 },
      axisLine: { show: false },
      splitLine: { lineStyle: { color: lineColor, type: 'dashed' } },
    },
    series: [
      {
        type: 'line',
        smooth: true,
        symbol: 'none',
        lineStyle: { width: 2, color: '#6366f1' },
        areaStyle: {
          color: {
            type: 'linear',
            x: 0, y: 0, x2: 0, y2: 1,
            colorStops: [
              { offset: 0, color: theme.dark ? 'rgba(99,102,241,0.3)' : 'rgba(99,102,241,0.15)' },
              { offset: 1, color: 'rgba(99,102,241,0)' },
            ],
          },
        },
        data: data.map((p) => [p.timestamp, p.avg_response_ms]),
      },
    ],
  }
})
</script>

<template>
  <v-chart :option="option" style="height: 220px" autoresize />
</template>
