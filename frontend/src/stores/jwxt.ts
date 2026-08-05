import { defineStore } from 'pinia'
import { ref } from 'vue'
import { jwxtLogin } from '../api/jwxt'

export const useJwxtStore = defineStore('jwxt', () => {
  const jwxtLoggedIn = ref(false)
  const jwxtUsername = ref('')
  const loading = ref(false)

  async function ensureLogin(): Promise<boolean> {
    if (jwxtLoggedIn.value) return true

    const savedUser = localStorage.getItem('jwxtUsername')
    const savedPwd = localStorage.getItem('jwxtPassword')
    if (!savedUser || !savedPwd) return false

    loading.value = true
    try {
      const result = await jwxtLogin(savedUser, savedPwd)
      if (result.success) {
        jwxtLoggedIn.value = true
        jwxtUsername.value = savedUser
        return true
      }
      return false
    } finally {
      loading.value = false
    }
  }

  function saveCredentials(username: string, password: string) {
    localStorage.setItem('jwxtUsername', username)
    localStorage.setItem('jwxtPassword', password)
    jwxtUsername.value = username
  }

  function clearCredentials() {
    localStorage.removeItem('jwxtUsername')
    localStorage.removeItem('jwxtPassword')
    jwxtLoggedIn.value = false
    jwxtUsername.value = ''
  }

  function hasSavedCredentials(): boolean {
    return !!(localStorage.getItem('jwxtUsername') && localStorage.getItem('jwxtPassword'))
  }

  return {
    jwxtLoggedIn,
    jwxtUsername,
    loading,
    ensureLogin,
    saveCredentials,
    clearCredentials,
    hasSavedCredentials,
  }
})
