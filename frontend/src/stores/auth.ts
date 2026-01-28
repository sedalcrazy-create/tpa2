import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/services/api'

interface User {
  id: number
  username: string
  email: string
  firstName: string
  lastName: string
  roleName: string
  roleTitle: string
  tenantId: number
  centerId?: number
  centerTitle?: string
}

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const accessToken = ref<string | null>(null)
  const refreshToken = ref<string | null>(null)

  const isAuthenticated = computed(() => !!accessToken.value)
  const fullName = computed(() => user.value ? `${user.value.firstName} ${user.value.lastName}` : '')
  const initials = computed(() => user.value ? user.value.firstName.charAt(0) : '')

  function hasRole(roles: string[]): boolean {
    return user.value ? roles.includes(user.value.roleName) : false
  }

  async function login(email: string, password: string) {
    try {
      const response = await api.post('/auth/login', { email, password })
      accessToken.value = response.data.token
      user.value = response.data.user
      api.defaults.headers.common['Authorization'] = `Bearer ${accessToken.value}`
      return true
    } catch (error) {
      console.error('Login failed:', error)
      return false
    }
  }

  async function logout() {
    accessToken.value = null
    refreshToken.value = null
    user.value = null
    delete api.defaults.headers.common['Authorization']
  }

  async function refreshAccessToken() {
    if (!refreshToken.value) return false

    try {
      const response = await api.post('/auth/refresh', { refresh_token: refreshToken.value })
      accessToken.value = response.data.access_token
      api.defaults.headers.common['Authorization'] = `Bearer ${accessToken.value}`
      return true
    } catch (error) {
      logout()
      return false
    }
  }

  function initialize() {
    if (accessToken.value) {
      api.defaults.headers.common['Authorization'] = `Bearer ${accessToken.value}`
    }
  }

  return {
    user,
    accessToken,
    refreshToken,
    isAuthenticated,
    fullName,
    initials,
    hasRole,
    login,
    logout,
    refreshAccessToken,
    initialize
  }
}, {
  persist: true
})
