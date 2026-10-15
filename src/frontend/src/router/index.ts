import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import CharacterSheetPage from '../views/CharacterSheetPage.vue'
import ManagementLayout from '../layouts/ManagementLayout.vue'
import AccountPage from '../views/AccountPage.vue'
import AdminPage from '../views/AdminPage.vue'
import { useAuth } from '../composables/useAuth'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'sheet',
    component: CharacterSheetPage,
  },
  {
    path: '/',
    component: ManagementLayout,
    children: [
      {
        path: 'account',
        name: 'account',
        component: AccountPage,
        meta: { requiresAuth: true },
      },
      {
        path: 'admin',
        name: 'admin',
        component: AdminPage,
        meta: { requiresAuth: true, requiresAdmin: true },
      },
    ],
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/',
  },
]

export const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach(async (to, _from, next) => {
  const { ensureAuthReady, isAuthenticated, isAdmin, showAuthModal } = useAuth()
  await ensureAuthReady()

  if (to.meta.requiresAuth && !isAuthenticated.value) {
    showAuthModal.value = true
    return next({ path: '/' })
  }

  if (to.meta.requiresAdmin && !isAdmin.value) {
    return next({ path: '/' })
  }

  next()
})
