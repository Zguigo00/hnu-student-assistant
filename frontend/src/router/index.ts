import { createRouter, createWebHashHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'Layout',
    component: () => import('../layouts/MainLayout.vue'),
    redirect: '/grades',
    children: [
      {
        path: 'grades',
        name: 'Grades',
        component: () => import('../pages/Grades.vue'),
      },
      {
        path: 'schedule',
        name: 'Schedule',
        component: () => import('../pages/Schedule.vue'),
      },
      {
        path: 'news',
        name: 'News',
        component: () => import('../pages/News.vue'),
      },
      {
        path: 'second-class',
        name: 'SecondClass',
        component: () => import('../pages/SecondClass.vue'),
      },
      {
        path: 'settings',
        name: 'Settings',
        component: () => import('../pages/Settings.vue'),
      },
      {
        path: 'help',
        name: 'Help',
        component: () => import('../pages/Help.vue'),
      },
    ],
  },
]

const router = createRouter({
  history: createWebHashHistory(),
  routes,
})

export default router
