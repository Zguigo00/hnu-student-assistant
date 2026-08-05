<template>
  <div class="news-page">
    <div class="page-header">
      <h2>校园新闻</h2>
      <div class="header-actions">
        <a-button type="primary" @click="fetchNews" :loading="loading" :disabled="!hasCredentials">刷新</a-button>
      </div>
    </div>

    <a-alert v-if="!hasCredentials" type="warning" show-icon style="margin-bottom: 16px">
      <template #message>
        尚未配置校园门户账号，请前往
        <a @click="router.push('/settings')">设置</a>
        页面填写门户账号密码。
      </template>
    </a-alert>

    <a-alert v-if="loginError" type="error" show-icon closable style="margin-bottom: 16px" @close="loginError = ''">
      <template #message>{{ loginError }}</template>
    </a-alert>

    <a-table
      :columns="columns"
      :data-source="newsItems"
      :loading="loading"
      row-key="id"
      :pagination="{ pageSize: 15, showSizeChanger: false }"
      bordered
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.dataIndex === 'title'">
          <a @click="viewDetail(record)">{{ record.title }}</a>
        </template>
      </template>
    </a-table>

    <a-modal
      v-model:open="modalVisible"
      :title="selectedNews?.title"
      width="800px"
    >
      <div class="news-meta" v-if="selectedNews">
        <span>来源：{{ selectedNews.createBy }}</span>
        <a-divider type="vertical" />
        <span>{{ selectedNews.createTime }}</span>
        <a-divider type="vertical" />
        <span>浏览：{{ selectedNews.look }}</span>
      </div>
      <a-divider />
      <div class="news-content" v-html="selectedNews?.content"></div>
      <template #footer>
        <a-button type="primary" @click="openInBrowser">在浏览器中打开</a-button>
      </template>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import { BrowserOpenURL } from '../../wailsjs/runtime/runtime'
import { getNews } from '../api/portal'
import { usePortalStore } from '../stores/portal'
import type { NewsItem } from '../types'

const router = useRouter()
const portalStore = usePortalStore()

const newsItems = ref<NewsItem[]>([])
const loading = ref(false)
const hasCredentials = ref(portalStore.hasSavedCredentials())
const loginError = ref('')
const modalVisible = ref(false)
const selectedNews = ref<NewsItem | null>(null)

const columns = [
  { title: '标题', dataIndex: 'title', ellipsis: true },
  { title: '来源', dataIndex: 'createBy', width: 160, align: 'center' as const },
  { title: '发布时间', dataIndex: 'createTime', width: 180, align: 'center' as const },
  { title: '浏览', dataIndex: 'look', width: 80, align: 'center' as const, sorter: (a: NewsItem, b: NewsItem) => a.look - b.look },
]

async function fetchNews() {
  if (!hasCredentials.value) {
    hasCredentials.value = portalStore.hasSavedCredentials()
    if (!hasCredentials.value) return
  }

  loading.value = true
  loginError.value = ''
  try {
    const loggedIn = await portalStore.ensureLogin()
    if (!loggedIn) {
      loginError.value = '门户登录失败，请检查设置中的账号密码'
      return
    }

    const result = await getNews(0, 50)
    if (result.success) {
      newsItems.value = result.data || []
    } else {
      loginError.value = result.message || '获取新闻失败'
    }
  } catch (e: any) {
    loginError.value = e.message || '请求失败'
  } finally {
    loading.value = false
  }
}

function viewDetail(item: NewsItem) {
  selectedNews.value = item
  modalVisible.value = true
}

function openInBrowser() {
  if (selectedNews.value) {
    BrowserOpenURL(`https://oshall.chnu.edu.cn/zhxy-new/campusInformation/detail?name=${selectedNews.value.id}`)
  }
}

onMounted(() => {
  hasCredentials.value = portalStore.hasSavedCredentials()
  if (hasCredentials.value) {
    fetchNews()
  }
})
</script>

<style scoped>
.news-page {
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

.news-meta {
  color: #666;
  font-size: 14px;
  margin-bottom: 8px;
}

.news-content {
  line-height: 1.8;
  font-size: 15px;
}

.news-content :deep(img) {
  max-width: 100%;
  height: auto;
  margin: 8px 0;
}

.news-content :deep(p) {
  margin-bottom: 12px;
}
</style>
