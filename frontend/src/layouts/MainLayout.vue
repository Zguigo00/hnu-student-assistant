<template>
  <a-layout class="app-layout">
    <a-layout-sider v-model:collapsed="collapsed" collapsible :width="200">
      <div class="logo">
        <h3 v-if="!collapsed">学生助手</h3>
        <h3 v-else>助手</h3>
      </div>
      <a-menu v-model:selectedKeys="selectedKeys" theme="dark" mode="inline">
        <a-menu-item key="grades" @click="router.push('/grades')">
          <template #icon><BarChartOutlined /></template>
          成绩查询
        </a-menu-item>
        <a-menu-item key="schedule" @click="router.push('/schedule')">
          <template #icon><CalendarOutlined /></template>
          课表查询
        </a-menu-item>
        <a-menu-item key="news" @click="router.push('/news')">
          <template #icon><ReadOutlined /></template>
          校园新闻
        </a-menu-item>
        <a-menu-item key="second-class" @click="router.push('/second-class')">
          <template #icon><TrophyOutlined /></template>
          第二课堂
        </a-menu-item>
        <a-menu-item key="settings" @click="router.push('/settings')">
          <template #icon><SettingOutlined /></template>
          设置
        </a-menu-item>
        <a-menu-item key="help" @click="router.push('/help')">
          <template #icon><QuestionCircleOutlined /></template>
          使用说明
        </a-menu-item>
      </a-menu>
    </a-layout-sider>

    <a-layout>
      <a-layout-header class="app-header">
        <div class="header-left">
          <span class="page-title">{{ pageTitle }}</span>
        </div>
      </a-layout-header>

      <a-layout-content class="app-content">
        <router-view />
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import {
  BarChartOutlined,
  CalendarOutlined,
  ReadOutlined,
  TrophyOutlined,
  SettingOutlined,
  QuestionCircleOutlined,
} from '@ant-design/icons-vue'

const router = useRouter()
const route = useRoute()
const collapsed = ref(false)

const selectedKeys = computed(() => {
  const path = route.path
  if (path.includes('schedule')) return ['schedule']
  if (path.includes('news')) return ['news']
  if (path.includes('second-class')) return ['second-class']
  if (path.includes('settings')) return ['settings']
  if (path.includes('help')) return ['help']
  return ['grades']
})

const pageTitle = computed(() => {
  const path = route.path
  if (path.includes('schedule')) return '课表查询'
  if (path.includes('news')) return '校园新闻'
  if (path.includes('second-class')) return '第二课堂成绩'
  if (path.includes('settings')) return '设置'
  if (path.includes('help')) return '使用说明'
  return '成绩查询'
})
</script>

<style scoped>
.app-layout {
  min-height: 100vh;
}

.logo {
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.logo h3 {
  margin: 0;
  font-size: 16px;
  color: white;
}

.app-header {
  background: white;
  padding: 0 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);
}

.page-title {
  font-size: 18px;
  font-weight: 600;
}

.app-content {
  margin: 24px;
  padding: 24px;
  background: white;
  border-radius: 8px;
  min-height: 280px;
}
</style>
