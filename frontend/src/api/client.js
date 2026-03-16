const BASE_URL = ''

async function request(path, options = {}) {
  const token = localStorage.getItem('token')
  const headers = {
    'Content-Type': 'application/json',
    ...options.headers,
  }
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }

  const res = await fetch(`${BASE_URL}${path}`, {
    ...options,
    headers,
  })

  if (res.status === 401) {
    localStorage.removeItem('token')
    window.location.href = '/login'
    throw new Error('Unauthorized')
  }

  if (res.status === 204) {
    return null
  }

  const body = await res.json()

  if (!res.ok) {
    throw new Error(body.error || `Request failed with status ${res.status}`)
  }

  return body.data
}

export function login(username, password) {
  return request('/api/auth/login', {
    method: 'POST',
    body: JSON.stringify({ username, password }),
  })
}

export function getMonitors() {
  return request('/api/monitors')
}

export function getMonitor(id) {
  return request(`/api/monitors/${id}`)
}

export function createMonitor(data) {
  return request('/api/monitors', {
    method: 'POST',
    body: JSON.stringify(data),
  })
}

export function updateMonitor(id, data) {
  return request(`/api/monitors/${id}`, {
    method: 'PUT',
    body: JSON.stringify(data),
  })
}

export function deleteMonitor(id) {
  return request(`/api/monitors/${id}`, {
    method: 'DELETE',
  })
}

export function getMonitorStatus(id) {
  return request(`/api/monitors/${id}/status`)
}

export function getMonitorHistory(id, limit = 20, offset = 0) {
  return request(`/api/monitors/${id}/history?limit=${limit}&offset=${offset}`)
}

export function getNotifications(monitorId) {
  return request(`/api/monitors/${monitorId}/notifications`)
}

export function createNotification(monitorId, data) {
  return request(`/api/monitors/${monitorId}/notifications`, {
    method: 'POST',
    body: JSON.stringify(data),
  })
}

export function updateNotification(id, data) {
  return request(`/api/notifications/${id}`, {
    method: 'PUT',
    body: JSON.stringify(data),
  })
}

export function deleteNotification(id) {
  return request(`/api/notifications/${id}`, {
    method: 'DELETE',
  })
}

export function getMonitorStats(id, period = '24h') {
  return request(`/api/monitors/${id}/stats?period=${period}`)
}

export function getMonitorTimeline(id, period = '24h', buckets = 60) {
  return request(`/api/monitors/${id}/timeline?period=${period}&buckets=${buckets}`)
}

export function getMonitorStatusCodes(id, period = '24h') {
  return request(`/api/monitors/${id}/status-codes?period=${period}`)
}

export function getMonitorStatusCodeTimeline(id, period = '24h', buckets = 60) {
  return request(`/api/monitors/${id}/status-code-timeline?period=${period}&buckets=${buckets}`)
}
