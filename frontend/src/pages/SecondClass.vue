<template>
  <div class="second-class-page">
    <div class="page-header">
      <h2>第二课堂成绩</h2>
      <div class="header-actions">
        <a-range-picker
          v-model:value="dateRange"
          :placeholder="['开始时间', '结束时间']"
          format="YYYY-MM-DD"
          style="width: 280px"
        />
        <a-button type="primary" @click="fetchScores" :loading="loading" :disabled="!hasCredentials">查询</a-button>
      </div>
    </div>

    <a-alert v-if="!hasCredentials" type="warning" show-icon style="margin-bottom: 16px">
      <template #message>
        尚未配置校园门户账号，请前往
        <a @click="router.push('/settings')">设置</a>
        页面填写账号密码。
      </template>
    </a-alert>

    <a-alert v-if="loginError" type="error" show-icon closable style="margin-bottom: 16px" @close="loginError = ''">
      <template #message>{{ loginError }}</template>
    </a-alert>

    <div v-if="scores.length > 0" class="score-cards">
      <a-card v-for="item in scores" :key="item.userCode" :title="item.userName || '我的成绩'" style="margin-bottom: 16px">
        <a-descriptions :column="1" bordered size="small">
          <a-descriptions-item label="思想政治与品德">{{ item.sxylyagrx?.toFixed(2) ?? '0' }}</a-descriptions-item>
          <a-descriptions-item label="专业技能与创新创业">{{ item.xskjycxcy?.toFixed(2) ?? '0' }}</a-descriptions-item>
          <a-descriptions-item label="体育健身运动">{{ item.tydlyydjn?.toFixed(2) ?? '0' }}</a-descriptions-item>
          <a-descriptions-item label="文化艺术修养">{{ item.rwskyyssy?.toFixed(2) ?? '0' }}</a-descriptions-item>
          <a-descriptions-item label="志愿服务与劳动实践">{{ item.shsjyzyfw?.toFixed(2) ?? '0' }}</a-descriptions-item>
        </a-descriptions>
      </a-card>
    </div>

    <a-empty v-if="!loading && scores.length === 0 && hasCredentials && queried" description="暂无数据，请尝试调整时间范围" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import type { Dayjs } from 'dayjs'
import dayjs from 'dayjs'
import { getSecondClassScores } from '../api/sc'
import { useScStore } from '../stores/sc'
import { usePortalStore } from '../stores/portal'
import type { SecondClassScore } from '../types'

const router = useRouter()
const scStore = useScStore()
const portalStore = usePortalStore()

const scores = ref<SecondClassScore[]>([])
const loading = ref(false)
const hasCredentials = ref(portalStore.hasSavedCredentials())
const loginError = ref('')
const queried = ref(false)

// 默认当前学年：2025-09-01 ~ 2026-08-31
const dateRange = ref<[Dayjs, Dayjs]>([dayjs('2025-09-01'), dayjs('2026-08-31')])

function formatTime(d: Dayjs): string {
  return d.format('YYYY-MM-DD') + ' 00:01:00'
}

function formatEndTime(d: Dayjs): string {
  return d.format('YYYY-MM-DD') + ' 23:59:59'
}

async function fetchScores() {
  if (!hasCredentials.value) {
    hasCredentials.value = portalStore.hasSavedCredentials()
    if (!hasCredentials.value) return
  }

  loading.value = true
  loginError.value = ''
  queried.value = true
  try {
    const loggedIn = await scStore.ensureLogin()
    if (!loggedIn) {
      loginError.value = '第二课堂登录失败，请检查设置中的门户账号密码'
      return
    }

    const startTime = formatTime(dateRange.value[0])
    const endTime = formatEndTime(dateRange.value[1])
    const result = await getSecondClassScores(startTime, endTime)
    if (result.success) {
      scores.value = result.data || []
    } else {
      loginError.value = result.message || '获取成绩失败'
    }
  } catch (e: any) {
    loginError.value = e.message || '请求失败'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  hasCredentials.value = portalStore.hasSavedCredentials()
  if (hasCredentials.value) {
    fetchScores()
  }
})
</script>

<style scoped>
.second-class-page {
  padding: 0;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.page-header h2 {
  margin: 0;
  font-size: 20px;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.score-cards {
  max-width: 600px;
}
</style>
