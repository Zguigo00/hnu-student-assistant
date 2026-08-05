<template>
  <div class="grades-page">
    <div class="page-header">
      <h2>成绩查询</h2>
      <div class="header-actions">
        <a-select v-model:value="selectedSemester" placeholder="选择学期" style="width: 200px" @change="fetchGrades">
          <a-select-option value="">全部学期</a-select-option>
          <a-select-option v-for="s in SEMESTERS" :key="s" :value="s">{{ s }}</a-select-option>
        </a-select>
        <a-button type="primary" @click="fetchGrades" :loading="loading" :disabled="!hasCredentials">刷新</a-button>
      </div>
    </div>

    <a-alert v-if="!hasCredentials" type="warning" show-icon style="margin-bottom: 16px">
      <template #message>
        尚未配置教务系统账号，请前往
        <a @click="router.push('/settings')">设置</a>
        页面填写账号密码。
      </template>
    </a-alert>

    <a-alert v-if="loginError" type="error" show-icon closable style="margin-bottom: 16px" @close="loginError = ''">
      <template #message>{{ loginError }}</template>
    </a-alert>

    <!-- 统计卡片 -->
    <div class="stats-row" v-if="grades.length > 0">
      <a-card class="stat-card">
        <a-statistic title="总课程数" :value="grades.length" />
      </a-card>
      <a-card class="stat-card">
        <a-statistic title="平均绩点 (GPA)" :value="gpa" :precision="2" />
      </a-card>
      <a-card class="stat-card">
        <a-statistic title="总学分" :value="totalCredits" :precision="1" />
      </a-card>
    </div>

    <!-- 成绩表格 -->
    <a-table
      :columns="columns"
      :data-source="grades"
      :loading="loading"
      row-key="kcmc"
      :pagination="{ pageSize: 15 }"
      bordered
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.dataIndex === 'cj'">
          <span :class="{ 'fail-score': parseFloat(record.cj) < 60 }">
            {{ record.cj }}
          </span>
        </template>
        <template v-if="column.dataIndex === 'jd'">
          <span :class="{ 'fail-score': parseFloat(record.jd) < 1.0 }">
            {{ record.jd }}
          </span>
        </template>
      </template>
    </a-table>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { getGrades } from '../api/jwxt'
import { useJwxtStore } from '../stores/jwxt'
import { SEMESTERS } from '../constants'
import type { Grade } from '../types'

const router = useRouter()
const jwxtStore = useJwxtStore()

const grades = ref<Grade[]>([])
const loading = ref(false)
const gpa = ref('0.00')
const selectedSemester = ref('')
const hasCredentials = ref(jwxtStore.hasSavedCredentials())
const loginError = ref('')

const totalCredits = computed(() => {
  return grades.value.reduce((sum, g) => sum + (parseFloat(g.xf) || 0), 0)
})

const columns = [
  { title: '学年学期', dataIndex: 'xnxq', width: 140 },
  { title: '课程名称', dataIndex: 'kcmc', ellipsis: true },
  { title: '课程性质', dataIndex: 'kcxzmc', width: 100 },
  { title: '学分', dataIndex: 'xf', width: 80, align: 'center' as const },
  { title: '成绩', dataIndex: 'cj', width: 80, align: 'center' as const, sorter: (a: Grade, b: Grade) => parseFloat(a.cj) - parseFloat(b.cj) },
  { title: '绩点', dataIndex: 'jd', width: 80, align: 'center' as const },
  { title: '开课学院', dataIndex: 'kkxy', width: 160 },
]

async function fetchGrades() {
  if (!hasCredentials.value) {
    hasCredentials.value = jwxtStore.hasSavedCredentials()
    if (!hasCredentials.value) return
  }

  loading.value = true
  loginError.value = ''
  try {
    const loggedIn = await jwxtStore.ensureLogin()
    if (!loggedIn) {
      loginError.value = '教务系统登录失败，请检查设置中的账号密码'
      return
    }

    const result = await getGrades(selectedSemester.value)
    if (result.success) {
      grades.value = result.data || []
      gpa.value = result.gpa || '0.00'
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
  hasCredentials.value = jwxtStore.hasSavedCredentials()
  if (hasCredentials.value) {
    fetchGrades()
  }
})
</script>

<style scoped>
.grades-page {
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

.stats-row {
  display: flex;
  gap: 16px;
  margin-bottom: 24px;
}

.stat-card {
  flex: 1;
}

.fail-score {
  color: #ff4d4f;
  font-weight: bold;
}
</style>
