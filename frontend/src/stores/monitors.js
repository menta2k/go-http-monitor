import { ref } from 'vue'
import { defineStore } from 'pinia'
import * as api from '../api/client'

export const useMonitorStore = defineStore('monitors', () => {
  const monitors = ref([])
  const loading = ref(false)
  const error = ref('')

  async function fetchMonitors() {
    loading.value = true
    error.value = ''
    try {
      monitors.value = await api.getMonitors()
    } catch (e) {
      error.value = e.message
    } finally {
      loading.value = false
    }
  }

  async function create(data) {
    const monitor = await api.createMonitor(data)
    monitors.value = [...monitors.value, monitor]
    return monitor
  }

  async function update(id, data) {
    const monitor = await api.updateMonitor(id, data)
    monitors.value = monitors.value.map((m) =>
      m.id === monitor.id ? monitor : m
    )
    return monitor
  }

  async function remove(id) {
    await api.deleteMonitor(id)
    monitors.value = monitors.value.filter((m) => m.id !== id)
  }

  return { monitors, loading, error, fetchMonitors, create, update, remove }
})
