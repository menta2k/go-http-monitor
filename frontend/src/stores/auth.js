import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { login as apiLogin } from '../api/client'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || '')

  const isAuthenticated = computed(() => !!token.value)

  async function login(username, password) {
    const data = await apiLogin(username, password)
    token.value = data.token
    localStorage.setItem('token', data.token)
  }

  function logout() {
    token.value = ''
    localStorage.removeItem('token')
  }

  return { token, isAuthenticated, login, logout }
})
