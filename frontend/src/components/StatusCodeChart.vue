<script setup>
import { computed } from 'vue'
import { use } from 'echarts/core'
import { LineChart } from 'echarts/charts'
import { GridComponent, TooltipComponent, LegendComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import VChart from 'vue-echarts'
import { useThemeStore } from '../stores/theme'

use([LineChart, GridComponent, TooltipComponent, LegendComponent, CanvasRenderer])

const props = defineProps({
  timeline: { type: Array, default: () => [] },
})

const theme = useThemeStore()

const codeColors = {
  '2xx': '#10b981',
  '3xx': '#6366f1',
  '4xx': '#f59e0b',
  '5xx': '#ef4444',
  'Error': '#991b1b',
}

const allCodes = computed(() => {
  const codes = new Set()
  for (const p of props.timeline) {
    if (p.codes) {
      for (const c of Object.keys(p.codes)) {
        codes.add(c)
      }
    }
  }
  // Stable order
  const order = ['2xx', '3xx', '4xx', '5xx', 'Error']
  return order.filter((c) => codes.has(c))
})

const option = computed(() => {
  const textColor = theme.dark ? '#9ca3af' : '#6b7280'
  const lineColor = theme.dark ? '#374151' : '#e5e7eb'
  const data = props.timeline

  const series = allCodes.value.map((code) => ({
    name: code,
    type: 'line',
    smooth: true,
    symbol: 'circle',
    symbolSize: 4,
    lineStyle: { width: 2, color: codeColors[code] || '#6b7280' },
    itemStyle: { color: codeColors[code] || '#6b7280' },
    areaStyle: {
      color: {
        type: 'linear',
        x: 0, y: 0, x2: 0, y2: 1,
        colorStops: [
          { offset: 0, color: (codeColors[code] || '#6b7280') + (theme.dark ? '40' : '20') },
          { offset: 1, color: (codeColors[code] || '#6b7280') + '00' },
        ],
      },
    },
    data: data.map((p) => [p.timestamp, (p.codes && p.codes[code]) || 0]),
  }))

  return {
    tooltip: {
      trigger: 'axis',
      backgroundColor: theme.dark ? '#1f2937' : '#fff',
      borderColor: theme.dark ? '#374151' : '#e5e7eb',
      textStyle: { color: theme.dark ? '#e5e7eb' : '#111827', fontSize: 12 },
      formatter(params) {
        const t = new Date(params[0].axisValue).toLocaleString()
        let lines = `${t}`
        for (const p of params) {
          if (p.data[1] > 0) {
            lines += `<br/><span style="display:inline-block;width:8px;height:8px;border-radius:50%;background:${p.color};margin-right:4px"></span>${p.seriesName}: <b>${p.data[1]}</b>`
          }
        }
        return lines
      },
    },
    legend: {
      bottom: 0,
      textStyle: { color: textColor, fontSize: 11 },
      itemWidth: 14,
      itemHeight: 3,
    },
    grid: { left: 50, right: 16, top: 12, bottom: 40 },
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
      name: 'count',
      nameTextStyle: { color: textColor, fontSize: 11 },
      axisLabel: { color: textColor, fontSize: 11 },
      axisLine: { show: false },
      splitLine: { lineStyle: { color: lineColor, type: 'dashed' } },
      minInterval: 1,
    },
    series,
  }
})
</script>

<template>
  <v-chart :option="option" style="height: 220px" autoresize />
</template>
