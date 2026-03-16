import { ref, watch } from 'vue'
import { defineStore } from 'pinia'

export const useThemeStore = defineStore('theme', () => {
  const stored = localStorage.getItem('theme')
  const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
  const dark = ref(stored ? stored === 'dark' : prefersDark)

  function apply() {
    document.documentElement.classList.toggle('dark', dark.value)
  }

  function toggle() {
    dark.value = !dark.value
  }

  watch(dark, (val) => {
    localStorage.setItem('theme', val ? 'dark' : 'light')
    apply()
  })

  // Apply on init
  apply()

  return { dark, toggle }
})
