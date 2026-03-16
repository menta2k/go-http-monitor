<script setup>
import { ref, onMounted } from 'vue'
import { getNotifications, createNotification, updateNotification, deleteNotification } from '../api/client'

const props = defineProps({
  monitorId: { type: Number, required: true },
})

const notifications = ref([])
const loading = ref(true)
const error = ref('')
const showForm = ref(false)

const form = ref({
  type: 'email',
  target: '',
  enabled: true,
})
const editingId = ref(null)
const formError = ref('')
const saving = ref(false)

async function fetchNotifications() {
  try {
    notifications.value = await getNotifications(props.monitorId)
  } catch (e) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editingId.value = null
  form.value = { type: 'email', target: '', enabled: true }
  formError.value = ''
  showForm.value = true
}

function openEdit(n) {
  editingId.value = n.id
  form.value = { type: n.type, target: n.target, enabled: n.enabled }
  formError.value = ''
  showForm.value = true
}

function cancelForm() {
  showForm.value = false
  editingId.value = null
}

async function handleSubmit() {
  formError.value = ''
  saving.value = true
  try {
    if (editingId.value) {
      const updated = await updateNotification(editingId.value, form.value)
      notifications.value = notifications.value.map((n) =>
        n.id === updated.id ? updated : n
      )
    } else {
      const created = await createNotification(props.monitorId, form.value)
      notifications.value = [...notifications.value, created]
    }
    showForm.value = false
    editingId.value = null
  } catch (e) {
    formError.value = e.message
  } finally {
    saving.value = false
  }
}

async function handleToggle(n) {
  try {
    const updated = await updateNotification(n.id, {
      type: n.type,
      target: n.target,
      enabled: !n.enabled,
    })
    notifications.value = notifications.value.map((x) =>
      x.id === updated.id ? updated : x
    )
  } catch (e) {
    alert(e.message)
  }
}

async function handleDelete(id) {
  if (!confirm('Delete this notification?')) return
  try {
    await deleteNotification(id)
    notifications.value = notifications.value.filter((n) => n.id !== id)
  } catch (e) {
    alert(e.message)
  }
}

onMounted(fetchNotifications)
</script>

<template>
  <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700">
    <div class="px-6 py-4 border-b border-gray-200 dark:border-gray-700 flex items-center justify-between">
      <h2 class="text-sm font-semibold text-gray-900 dark:text-gray-100">Notifications</h2>
      <button
        @click="openCreate"
        class="inline-flex items-center gap-1 text-xs font-medium text-indigo-600 dark:text-indigo-400 hover:text-indigo-700 dark:hover:text-indigo-300 transition-colors"
      >
        <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        Add
      </button>
    </div>

    <!-- Form -->
    <div v-if="showForm" class="px-6 py-4 border-b border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-900/50">
      <form @submit.prevent="handleSubmit" class="space-y-3">
        <div v-if="formError" class="bg-red-50 dark:bg-red-900/30 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-400 text-xs rounded-lg px-3 py-2">
          {{ formError }}
        </div>
        <div class="grid grid-cols-2 gap-3">
          <div>
            <label class="block text-xs font-medium text-gray-700 dark:text-gray-300 mb-1">Type</label>
            <select
              v-model="form.type"
              class="w-full rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 px-3 py-1.5 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
            >
              <option value="email">Email</option>
              <option value="slack">Slack Webhook</option>
            </select>
          </div>
          <div>
            <label class="block text-xs font-medium text-gray-700 dark:text-gray-300 mb-1">
              {{ form.type === 'email' ? 'Email Address' : 'Webhook URL' }}
            </label>
            <input
              v-model="form.target"
              :type="form.type === 'email' ? 'email' : 'url'"
              required
              :placeholder="form.type === 'email' ? 'alerts@example.com' : 'https://hooks.slack.com/...'"
              class="w-full rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 px-3 py-1.5 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent placeholder-gray-400 dark:placeholder-gray-500"
            />
          </div>
        </div>
        <div class="flex items-center justify-between">
          <label class="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300">
            <input
              v-model="form.enabled"
              type="checkbox"
              class="rounded border-gray-300 dark:border-gray-600 text-indigo-600 focus:ring-indigo-500"
            />
            Enabled
          </label>
          <div class="flex items-center gap-2">
            <button
              type="button"
              @click="cancelForm"
              class="text-xs text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200 transition-colors"
            >
              Cancel
            </button>
            <button
              type="submit"
              :disabled="saving"
              class="bg-indigo-600 text-white text-xs font-medium rounded-lg px-3 py-1.5 hover:bg-indigo-700 disabled:opacity-50 transition-colors"
            >
              {{ saving ? 'Saving...' : (editingId ? 'Update' : 'Create') }}
            </button>
          </div>
        </div>
      </form>
    </div>

    <!-- List -->
    <div v-if="loading" class="px-6 py-6 text-center text-sm text-gray-500 dark:text-gray-400">
      Loading...
    </div>
    <div v-else-if="error" class="px-6 py-4 text-sm text-red-600 dark:text-red-400">
      {{ error }}
    </div>
    <div v-else-if="notifications.length === 0 && !showForm" class="px-6 py-6 text-center text-sm text-gray-500 dark:text-gray-400">
      No notifications configured
    </div>
    <div v-else class="divide-y divide-gray-100 dark:divide-gray-700">
      <div
        v-for="n in notifications"
        :key="n.id"
        class="px-6 py-3 flex items-center justify-between gap-4"
      >
        <div class="min-w-0 flex-1">
          <div class="flex items-center gap-2">
            <span
              :class="[
                'inline-flex items-center px-2 py-0.5 rounded text-xs font-medium',
                n.type === 'email'
                  ? 'bg-blue-100 text-blue-700 dark:bg-blue-900/40 dark:text-blue-400'
                  : 'bg-purple-100 text-purple-700 dark:bg-purple-900/40 dark:text-purple-400'
              ]"
            >
              {{ n.type === 'email' ? 'Email' : 'Slack' }}
            </span>
            <span
              :class="[
                'inline-flex items-center px-2 py-0.5 rounded text-xs font-medium',
                n.enabled
                  ? 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-400'
                  : 'bg-gray-100 text-gray-500 dark:bg-gray-700 dark:text-gray-400'
              ]"
            >
              {{ n.enabled ? 'Active' : 'Disabled' }}
            </span>
          </div>
          <p class="mt-1 text-sm text-gray-700 dark:text-gray-300 truncate">{{ n.target }}</p>
        </div>
        <div class="flex items-center gap-1 shrink-0">
          <button
            @click="handleToggle(n)"
            :title="n.enabled ? 'Disable' : 'Enable'"
            class="p-1.5 text-gray-400 dark:text-gray-500 hover:text-indigo-600 dark:hover:text-indigo-400 transition-colors rounded"
          >
            <svg v-if="n.enabled" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636" />
            </svg>
            <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </button>
          <button
            @click="openEdit(n)"
            class="p-1.5 text-gray-400 dark:text-gray-500 hover:text-indigo-600 dark:hover:text-indigo-400 transition-colors rounded"
            title="Edit"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
            </svg>
          </button>
          <button
            @click="handleDelete(n.id)"
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
  </div>
</template>
