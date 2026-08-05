<template>
  <div class="schedule-page">
    <div class="page-header">
      <h2>课表查询</h2>
      <div class="header-actions">
        <a-select v-model:value="selectedSemester" placeholder="选择学期" style="width: 200px" @change="fetchSchedule">
          <a-select-option v-for="s in SEMESTERS" :key="s" :value="s">{{ s }}</a-select-option>
        </a-select>
        <a-select v-model:value="selectedWeek" placeholder="选择周次" style="width: 120px">
          <a-select-option v-for="w in WEEKS" :key="w" :value="w">第{{ w }}周</a-select-option>
        </a-select>
        <a-button type="primary" @click="fetchSchedule" :loading="loading" :disabled="!hasCredentials">刷新</a-button>
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

    <!-- 当前周课程数统计 -->
    <div class="week-info" v-if="filteredCourses.length > 0">
      本周共 <strong>{{ filteredCourses.length }}</strong> 节课
    </div>

    <!-- 课表网格 -->
    <div class="schedule-grid">
      <div class="grid-header">
        <div class="time-col">节次</div>
        <div v-for="day in WEEK_DAYS" :key="day" class="day-col">{{ day }}</div>
      </div>

      <div v-for="section in SECTIONS" :key="section" class="grid-row">
        <div class="time-col">
          <div class="section-label">第{{ section }}节</div>
        </div>
        <div v-for="day in 7" :key="day" class="day-col">
          <div
            v-for="course in getCourses(day, section)"
            :key="course.kcmc + course.ksjc"
            class="course-card"
            :style="{ background: getCourseColor(course.kcmc) }"
          >
            <div class="course-name">{{ course.kcmc }}</div>
            <div class="course-info">{{ course.jsmc }}</div>
            <div class="course-info">{{ course.jsxm }}</div>
            <div class="course-week">第{{ course.zcmc }}</div>
          </div>
        </div>
      </div>
    </div>

    <!-- 课程列表视图 -->
    <a-divider>课程列表</a-divider>
    <a-table
      :columns="columns"
      :data-source="filteredCourses"
      :loading="loading"
      :row-key="(record: Course) => record.kcmc + record.xq + record.ksjc"
      :pagination="false"
      bordered
      size="small"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getSchedule } from '../api/jwxt'
import { useJwxtStore } from '../stores/jwxt'
import { SEMESTERS, DEFAULT_SEMESTER, WEEK_DAYS, SECTIONS, WEEKS, COURSE_COLORS } from '../constants'
import type { Course } from '../types'

const router = useRouter()
const jwxtStore = useJwxtStore()

const courses = ref<Course[]>([])
const loading = ref(false)
const selectedSemester = ref(DEFAULT_SEMESTER)
const selectedWeek = ref(1)
const hasCredentials = ref(jwxtStore.hasSavedCredentials())
const loginError = ref('')

const columns = [
  { title: '课程名称', dataIndex: 'kcmc', ellipsis: true },
  { title: '星期', dataIndex: 'xq', width: 80, align: 'center' as const },
  { title: '节次', dataIndex: 'ksjc', width: 100, align: 'center' as const },
  { title: '教室', dataIndex: 'jsmc', width: 160 },
  { title: '教师', dataIndex: 'jsxm', width: 100 },
  { title: '周次', dataIndex: 'zcmc', width: 160 },
]

// 解析 zcd 字段，如 "1-16周"、"1-8,10-16周"、"1-15周(单)"、"13周"
function parseWeeks(zcd: string): Set<number> {
  const weeks = new Set<number>()
  const cleaned = zcd.replace(/周/g, '')
  const isOdd = zcd.includes('(单)')
  const isEven = zcd.includes('(双)')

  // 按逗号分割多个范围
  for (const part of cleaned.split(',')) {
    const trimmed = part.trim().replace(/\(单\)|\(双\)/, '')
    if (trimmed.includes('-')) {
      const [start, end] = trimmed.split('-').map(Number)
      for (let w = start; w <= end; w++) {
        if (isOdd && w % 2 === 0) continue
        if (isEven && w % 2 !== 0) continue
        weeks.add(w)
      }
    } else {
      const w = parseInt(trimmed)
      if (!isNaN(w)) weeks.add(w)
    }
  }
  return weeks
}

// 解析节次，如 "7-8节" → [7, 8]，"1-2节" → [1, 2]
function parseSections(jc: string): number[] {
  const match = jc.match(/(\d+)-(\d+)/)
  if (match) {
    const start = parseInt(match[1])
    const end = parseInt(match[2])
    return Array.from({ length: end - start + 1 }, (_, i) => start + i)
  }
  const single = parseInt(jc)
  return isNaN(single) ? [] : [single]
}

// 按选中周次过滤后的课程
const filteredCourses = computed(() => {
  const week = selectedWeek.value
  return courses.value.filter(c => {
    if (!c.zcmc) return true
    return parseWeeks(c.zcmc).has(week)
  })
})

function getCourseColor(name: string): string {
  let hash = 0
  for (let i = 0; i < name.length; i++) {
    hash = name.charCodeAt(i) + ((hash << 5) - hash)
  }
  return COURSE_COLORS[Math.abs(hash) % COURSE_COLORS.length]
}

function getCourses(day: number, section: string): Course[] {
  const [start, end] = section.split('-').map(Number)
  return filteredCourses.value.filter(c => {
    const cDay = parseInt(c.xq)
    const sections = parseSections(c.ksjc)
    return cDay === day && sections.some(s => s >= start && s <= end)
  })
}

async function fetchSchedule() {
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

    const result = await getSchedule(selectedSemester.value)
    if (result.success) {
      courses.value = result.data || []
      console.log('[课表] 总课程数:', courses.value.length)
      for (const c of courses.value) {
        const weeks = parseWeeks(c.zcmc)
        console.log(`[课表] ${c.kcmc} xq=${c.xq} ksjc=${c.ksjc} zcmc=${c.zcmc} → 周次:`, [...weeks])
      }
    } else {
      loginError.value = result.message || '获取课表失败'
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
    fetchSchedule()
  }
})
</script>

<style scoped>
.schedule-page {
  padding: 0;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.page-header h2 {
  margin: 0;
  font-size: 20px;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.week-info {
  margin-bottom: 16px;
  color: #666;
  font-size: 14px;
}

.week-info strong {
  color: #1890ff;
  font-size: 16px;
}

.schedule-grid {
  border: 1px solid #e8e8e8;
  border-radius: 8px;
  overflow: hidden;
}

.grid-header {
  display: flex;
  background: #fafafa;
  border-bottom: 1px solid #e8e8e8;
}

.grid-row {
  display: flex;
  border-bottom: 1px solid #f0f0f0;
}

.grid-row:last-child {
  border-bottom: none;
}

.time-col {
  width: 80px;
  min-width: 80px;
  padding: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #fafafa;
  border-right: 1px solid #e8e8e8;
  font-size: 12px;
  color: #666;
}

.day-col {
  flex: 1;
  min-height: 80px;
  padding: 4px;
  border-right: 1px solid #f0f0f0;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.day-col:last-child {
  border-right: none;
}

.grid-header .day-col {
  text-align: center;
  font-weight: 600;
  min-height: auto;
  padding: 8px;
}

.course-card {
  border-radius: 6px;
  padding: 6px 8px;
  font-size: 11px;
  line-height: 1.4;
  border: 1px solid rgba(0, 0, 0, 0.06);
}

.course-name {
  font-weight: 600;
  font-size: 12px;
  margin-bottom: 2px;
}

.course-info {
  color: #666;
  font-size: 11px;
}

.course-week {
  color: #999;
  font-size: 10px;
  margin-top: 2px;
}

.section-label {
  font-size: 12px;
}
</style>
