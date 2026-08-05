import { defineStore } from 'pinia'
import { ref } from 'vue'
import * as portalApi from '../api/portal'

export const usePortalStore = defineStore('portal', () => {
  const portalLoggedIn = ref(false)
  const portalUsername = ref('')
  const loading = ref(false)

  async function ensureLogin(): Promise<boolean> {
    if (portalLoggedIn.value) return true

    const savedUser = localStorage.getItem('portalUsername')
    const savedPwd = localStorage.getItem('portalPassword')
    if (!savedUser || !savedPwd) return false

    loading.value = true
    try {
      const result = await portalApi.portalLogin(savedUser, savedPwd)
      if (result.success) {
        portalLoggedIn.value = true
        portalUsername.value = savedUser
        return true
      }
      return false
    } finally {
      loading.value = false
    }
  }

  function saveCredentials(username: string, password: string) {
    localStorage.setItem('portalUsername', username)
    localStorage.setItem('portalPassword', password)
    portalUsername.value = username
  }

  function clearCredentials() {
    localStorage.removeItem('portalUsername')
    localStorage.removeItem('portalPassword')
    portalLoggedIn.value = false
    portalUsername.value = ''
  }

  function hasSavedCredentials(): boolean {
    return !!(localStorage.getItem('portalUsername') && localStorage.getItem('portalPassword'))
  }

  return {
    portalLoggedIn,
    portalUsername,
    loading,
    ensureLogin,
    saveCredentials,
    clearCredentials,
    hasSavedCredentials,
  }
})
