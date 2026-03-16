<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useMonitorStore } from '../stores/monitors'
import { getMonitor } from '../api/client'

const router = useRouter()
const route = useRoute()
const store = useMonitorStore()

const isEdit = computed(() => !!route.params.id)

const form = ref({
  url: '',
  expected_status: 200,
  body_contains: '',
  interval_seconds: 60,
})
const error = ref('')
const loading = ref(false)

onMounted(async () => {
  if (isEdit.value) {
    try {
      const m = await getMonitor(route.params.id)
      form.value = {
        url: m.url,
        expected_status: m.expected_status,
        body_contains: m.body_contains,
        interval_seconds: m.interval_seconds,
      }
    } catch (e) {
      error.value = e.message
    }
  }
})

async function handleSubmit() {
  error.value = ''
  loading.value = true
  try {
    if (isEdit.value) {
      await store.update(Number(route.params.id), form.value)
    } else {
      await store.create(form.value)
    }
    router.push('/')
  } catch (e) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="max-w-xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
    <div class="mb-6">
      <router-link to="/" class="text-sm text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200 transition-colors">
        &larr; Back to monitors
      </router-link>
      <h1 class="text-lg font-semibold text-gray-900 dark:text-gray-100 mt-2">
        {{ isEdit ? 'Edit Monitor' : 'New Monitor' }}
      </h1>
    </div>

    <form @submit.prevent="handleSubmit" class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-6 space-y-5">
      <div v-if="error" class="bg-red-50 dark:bg-red-900/30 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-400 text-sm rounded-lg px-4 py-3">
        {{ error }}
      </div>

      <div>
        <label for="url" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">URL</label>
        <input
          id="url"
          v-model="form.url"
          type="url"
          required
          placeholder="https://example.com"
          class="w-full rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent placeholder-gray-400 dark:placeholder-gray-500"
        />
      </div>

      <div class="grid grid-cols-2 gap-4">
        <div>
          <label for="expected_status" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Expected Status</label>
          <input
            id="expected_status"
            v-model.number="form.expected_status"
            type="number"
            min="100"
            max="599"
            required
            class="w-full rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
          />
        </div>
        <div>
          <label for="interval" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Check Interval (seconds)</label>
          <input
            id="interval"
            v-model.number="form.interval_seconds"
            type="number"
            min="5"
            max="86400"
            required
            class="w-full rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
          />
        </div>
      </div>

      <div>
        <label for="body_contains" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
          Body Contains
          <span class="font-normal text-gray-400 dark:text-gray-500">(optional)</span>
        </label>
        <input
          id="body_contains"
          v-model="form.body_contains"
          type="text"
          placeholder="Expected string in response body"
          class="w-full rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent placeholder-gray-400 dark:placeholder-gray-500"
        />
        <p class="mt-1 text-xs text-gray-400 dark:text-gray-500">Leave empty to skip body content check</p>
      </div>

      <div class="flex items-center gap-3 pt-2">
        <button
          type="submit"
          :disabled="loading"
          class="bg-indigo-600 text-white text-sm font-medium rounded-lg px-4 py-2.5 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 dark:focus:ring-offset-gray-800 disabled:opacity-50 transition-colors"
        >
          {{ loading ? 'Saving...' : (isEdit ? 'Update Monitor' : 'Create Monitor') }}
        </button>
        <router-link to="/" class="text-sm text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200 transition-colors">
          Cancel
        </router-link>
      </div>
    </form>
  </div>
</template>
