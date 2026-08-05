import { defineStore } from 'pinia'
import { ref } from 'vue'
import { scLogin } from '../api/sc'

export const useScStore = defineStore('sc', () => {
  const scLoggedIn = ref(false)
  const loading = ref(false)

  async function ensureLogin(): Promise<boolean> {
    if (scLoggedIn.value) return true

    const savedUser = localStorage.getItem('portalUsername')
    const savedPwd = localStorage.getItem('portalPassword')
    if (!savedUser || !savedPwd) return false

    loading.value = true
    try {
      const result = await scLogin(savedUser, savedPwd)
      if (result.success) {
        scLoggedIn.value = true
        return true
      }
      return false
    } finally {
      loading.value = false
    }
  }

  function hasSavedCredentials(): boolean {
    return !!(localStorage.getItem('portalUsername') && localStorage.getItem('portalPassword'))
  }

  return {
    scLoggedIn,
    loading,
    ensureLogin,
    hasSavedCredentials,
  }
})
