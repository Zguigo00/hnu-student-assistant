<template>
  <div class="settings-page">
    <div class="page-header">
      <h2>设置</h2>
    </div>

    <a-card title="教务系统账号" style="max-width: 500px; margin-bottom: 24px">
      <a-alert type="info" show-icon style="margin-bottom: 16px">
        <template #message>教务系统账号用于查询成绩和课表。</template>
      </a-alert>

      <a-form layout="vertical">
        <a-form-item label="账号">
          <a-input v-model:value="jwxtUsername" placeholder="请输入教务系统账号" />
        </a-form-item>
        <a-form-item label="密码">
          <a-input-password v-model:value="jwxtPassword" placeholder="请输入教务系统密码" />
        </a-form-item>
      </a-form>

      <div class="btn-group">
        <a-button type="primary" @click="handleJwxtSave">保存</a-button>
        <a-button @click="handleJwxtTest" :loading="jwxtTesting">测试连接</a-button>
        <a-button danger @click="handleJwxtClear" :disabled="!jwxtHasSaved">清除</a-button>
      </div>
    </a-card>

    <a-card title="校园门户账号" style="max-width: 500px">
      <a-alert type="info" show-icon style="margin-bottom: 16px">
        <template #message>门户账号与教务系统账号不同，用于访问校园新闻等功能。</template>
      </a-alert>

      <a-form layout="vertical">
        <a-form-item label="账号">
          <a-input v-model:value="portalUsername" placeholder="请输入门户账号" />
        </a-form-item>
        <a-form-item label="密码">
          <a-input-password v-model:value="portalPassword" placeholder="请输入门户密码" />
        </a-form-item>
      </a-form>

      <div class="btn-group">
        <a-button type="primary" @click="handlePortalSave">保存</a-button>
        <a-button @click="handlePortalTest" :loading="portalTesting">测试连接</a-button>
        <a-button danger @click="handlePortalClear" :disabled="!portalHasSaved">清除</a-button>
      </div>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { usePortalStore } from '../stores/portal'
import { useJwxtStore } from '../stores/jwxt'

const portalStore = usePortalStore()
const jwxtStore = useJwxtStore()

// 教务系统
const jwxtUsername = ref('')
const jwxtPassword = ref('')
const jwxtTesting = ref(false)
const jwxtHasSaved = ref(false)

// 校园门户
const portalUsername = ref('')
const portalPassword = ref('')
const portalTesting = ref(false)
const portalHasSaved = ref(false)

// === 教务系统 ===
function handleJwxtSave() {
  if (!jwxtUsername.value || !jwxtPassword.value) {
    message.warning('请填写完整的账号密码')
    return
  }
  jwxtStore.saveCredentials(jwxtUsername.value, jwxtPassword.value)
  jwxtHasSaved.value = true
  message.success('保存成功')
}

async function handleJwxtTest() {
  if (!jwxtUsername.value || !jwxtPassword.value) {
    message.warning('请先填写账号密码')
    return
  }
  jwxtTesting.value = true
  try {
    jwxtStore.saveCredentials(jwxtUsername.value, jwxtPassword.value)
    const success = await jwxtStore.ensureLogin()
    if (success) {
      message.success('连接成功')
    } else {
      message.error('连接失败，请检查账号密码')
    }
  } catch {
    message.error('连接失败')
  } finally {
    jwxtTesting.value = false
  }
}

function handleJwxtClear() {
  jwxtStore.clearCredentials()
  jwxtUsername.value = ''
  jwxtPassword.value = ''
  jwxtHasSaved.value = false
  message.success('已清除')
}

// === 校园门户 ===
function handlePortalSave() {
  if (!portalUsername.value || !portalPassword.value) {
    message.warning('请填写完整的账号密码')
    return
  }
  portalStore.saveCredentials(portalUsername.value, portalPassword.value)
  portalHasSaved.value = true
  message.success('保存成功')
}

async function handlePortalTest() {
  if (!portalUsername.value || !portalPassword.value) {
    message.warning('请先填写账号密码')
    return
  }
  portalTesting.value = true
  try {
    portalStore.saveCredentials(portalUsername.value, portalPassword.value)
    const success = await portalStore.ensureLogin()
    if (success) {
      message.success('连接成功')
    } else {
      message.error('连接失败，请检查账号密码')
    }
  } catch {
    message.error('连接失败')
  } finally {
    portalTesting.value = false
  }
}

function handlePortalClear() {
  portalStore.clearCredentials()
  portalUsername.value = ''
  portalPassword.value = ''
  portalHasSaved.value = false
  message.success('已清除')
}

onMounted(() => {
  // 教务系统
  const savedJwxtUser = localStorage.getItem('jwxtUsername')
  const savedJwxtPwd = localStorage.getItem('jwxtPassword')
  if (savedJwxtUser) {
    jwxtUsername.value = savedJwxtUser
    jwxtHasSaved.value = true
  }
  if (savedJwxtPwd) jwxtPassword.value = savedJwxtPwd

  // 校园门户
  const savedPortalUser = localStorage.getItem('portalUsername')
  const savedPortalPwd = localStorage.getItem('portalPassword')
  if (savedPortalUser) {
    portalUsername.value = savedPortalUser
    portalHasSaved.value = true
  }
  if (savedPortalPwd) portalPassword.value = savedPortalPwd
})
</script>

<style scoped>
.settings-page {
  padding: 0;
}

.page-header {
  margin-bottom: 24px;
}

.page-header h2 {
  margin: 0;
  font-size: 20px;
}

.btn-group {
  display: flex;
  gap: 12px;
}
</style>
