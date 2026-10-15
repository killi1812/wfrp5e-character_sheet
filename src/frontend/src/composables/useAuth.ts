import { ref, computed } from 'vue'
import { characterApi } from '../services/characterApi'

export interface AuthUser {
  username: string
  role: string
  token: string
  uuid?: string
}

const currentUser = ref<AuthUser | null>(null)
const isAuthInitializing = ref(false)
const isAuthReady = ref(false)
const showAuthModal = ref(false)

const isAuthenticated = computed(() => Boolean(currentUser.value?.token || localStorage.getItem('auth_token')))
const isAdmin = computed(() => currentUser.value?.role === 'admin')

async function ensureAuthReady(): Promise<void> {
  if (isAuthReady.value) return
  if (isAuthInitializing.value) {
    while (isAuthInitializing.value) {
      await new Promise((r) => setTimeout(r, 20))
    }
    return
  }
  await initAuth()
}

async function initAuth(): Promise<void> {
  isAuthInitializing.value = true
  const storedToken = localStorage.getItem('auth_token')
  if (storedToken) {
    try {
      const userData = await characterApi.getMyData()
      const activeToken = localStorage.getItem('auth_token') || storedToken
      currentUser.value = {
        username: userData.username,
        role: userData.role || 'user',
        token: activeToken,
        uuid: userData.uuid,
      }
    } catch (err: any) {
      if (err?.response?.status === 401) {
        localStorage.removeItem('auth_token')
        currentUser.value = null
      }
    }
  } else {
    currentUser.value = null
  }
  isAuthInitializing.value = false
  isAuthReady.value = true
}

async function login(username: string, password: string, sessionId?: string): Promise<AuthUser> {
  const reqBody: { username: string; password: string; sessionId?: string } = {
    username,
    password,
  }
  const storedSession = sessionId || characterApi.getSessionId()
  if (storedSession) {
    reqBody.sessionId = storedSession
  }

  const data = await characterApi.login(reqBody)
  if (data.accessToken) {
    localStorage.setItem('auth_token', data.accessToken)
  }
  if (data.sessionId) {
    characterApi.setSessionId(data.sessionId)
  }

  let role = username.toLowerCase() === 'admin' ? 'admin' : 'user'
  let uuid: string | undefined
  try {
    const userData = await characterApi.getMyData()
    if (userData.role) role = userData.role
    if (userData.uuid) uuid = userData.uuid
  } catch {
    // fallback to inferred role
  }

  const user: AuthUser = {
    username,
    role,
    token: data.accessToken,
    uuid,
  }
  currentUser.value = user
  return user
}

async function register(username: string, password: string, email?: string): Promise<void> {
  await characterApi.register({ username, password, email })
}

async function logout(): Promise<void> {
  await characterApi.logout()
  localStorage.removeItem('auth_token')
  characterApi.clearSessionId()
  currentUser.value = null
}

export function useAuth() {
  return {
    currentUser,
    isAuthenticated,
    isAdmin,
    isAuthReady,
    showAuthModal,
    ensureAuthReady,
    initAuth,
    login,
    register,
    logout,
  }
}
