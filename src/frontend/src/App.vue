<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useTheme } from 'vuetify'
import AuthDialog from './components/AuthDialog.vue'
import AdminPage from './views/AdminPage.vue'
import CharacterSheetPage from './views/CharacterSheetPage.vue'

const theme = useTheme()
const isDark = computed(() => theme.global.name.value === 'wfrpDark')

function toggleTheme() {
  theme.global.name.value = isDark.value ? 'wfrpLight' : 'wfrpDark'
}

// User & Auth State
const currentUser = ref<{ username: string; role: string; token: string; uuid?: string } | null>(null)
const showAuthDialog = ref(false)
const showKebabOverlay = ref(false)

// Page Navigation View State ('sheet' | 'admin')
const currentPage = ref<'sheet' | 'admin'>('sheet')
const sheetRef = ref<InstanceType<typeof CharacterSheetPage> | null>(null)

onMounted(async () => {
  if (window.location.pathname === '/admin') {
    currentPage.value = 'admin'
  }
  window.addEventListener('popstate', () => {
    currentPage.value = window.location.pathname === '/admin' ? 'admin' : 'sheet'
  })

  // Restore authenticated session from localStorage if present
  const storedToken = localStorage.getItem('auth_token')
  if (storedToken) {
    try {
      const res = await fetch('/api/user/my-data', {
        headers: { Authorization: `Bearer ${storedToken}` },
      })
      if (res.ok) {
        const userData = await res.json()
        currentUser.value = {
          username: userData.username,
          role: userData.role || 'user',
          token: storedToken,
          uuid: userData.uuid,
        }
      } else {
        localStorage.removeItem('auth_token')
      }
    } catch {
      // offline / demo state fallback
    }
  }
})

function navigateTo(page: 'sheet' | 'admin') {
  currentPage.value = page
  const targetPath = page === 'admin' ? '/admin' : '/'
  if (window.location.pathname !== targetPath) {
    window.history.pushState({}, '', targetPath)
  }
  showKebabOverlay.value = false
}

function onLoginSuccess(user: { username: string; role: string; token: string; uuid?: string }) {
  currentUser.value = user
}

async function logout() {
  if (currentUser.value?.token && currentUser.value.token !== 'demo-token') {
    try {
      await fetch('/api/auth/logout', {
        method: 'POST',
        headers: { Authorization: `Bearer ${currentUser.value.token}` },
      })
    } catch {
      // ignore network errors on logout
    }
  }
  localStorage.removeItem('auth_token')
  currentUser.value = null
}

const isAdmin = computed(() => currentUser.value?.role === 'admin')
</script>

<template>
  <v-app class="wfrp-full-app bg-background text-on-background">
    <!-- Floating Kebab Button -->
    <div class="kebab-fixed-pos">
      <v-tooltip text="System Options & Settings" location="left">
        <template #activator="{ props: tooltipProps }">
          <v-btn
            v-bind="tooltipProps"
            icon="mdi-dots-vertical"
            color="primary"
            elevation="8"
            size="large"
            class="kebab-fab"
            @click="showKebabOverlay = true"
          />
        </template>
      </v-tooltip>
    </div>

    <!-- Kebab POPOUT OVERLAY MODAL (Full Dialog, No dropdown menu) -->
    <v-dialog v-model="showKebabOverlay" max-width="500" transition="dialog-bottom-transition">
      <v-card color="surface" class="pa-6 rounded-xl elevation-24 border">
        <div class="d-flex justify-space-between align-center mb-4">
          <div class="d-flex align-center">
            <v-icon icon="mdi-shield-cog" color="primary" class="mr-2" size="large" />
            <h2 class="text-h5 font-weight-black mb-0 text-on-surface">Character Sheet Menu</h2>
          </div>
          <v-btn icon="mdi-close" variant="text" size="small" @click="showKebabOverlay = false" />
        </div>

        <!-- User Status Banner -->
        <div v-if="currentUser" class="pa-4 bg-surface-variant rounded-lg mb-4 d-flex align-center justify-space-between border">
          <div>
            <div class="text-subtitle-1 font-weight-bold text-on-surface-variant">{{ currentUser.username }}</div>
            <v-chip size="x-small" :color="isAdmin ? 'secondary' : 'info'" class="mt-1" variant="flat">
              {{ currentUser.role.toUpperCase() }}
            </v-chip>
          </div>
          <v-avatar color="primary" size="40">
            <v-icon icon="mdi-account" color="on-primary" />
          </v-avatar>
        </div>

        <v-list class="bg-transparent pa-0">
          <v-list-item
            :prepend-icon="isDark ? 'mdi-weather-sunny' : 'mdi-weather-night'"
            :title="isDark ? 'Switch to Light Theme' : 'Switch to Dark Theme'"
            class="mb-2 rounded-lg bg-surface-variant border"
            @click="toggleTheme"
          />

          <!-- Admin Panel Button (If logged in as Admin) -->
          <v-list-item
            v-if="isAdmin"
            prepend-icon="mdi-shield-crown"
            title="Admin Control Panel"
            subtitle="User & Database Management"
            color="secondary"
            class="mb-2 rounded-lg bg-secondary text-on-secondary font-weight-bold"
            @click="navigateTo('admin')"
          />

          <!-- Sheet Actions (when viewing character sheet) -->
          <template v-if="currentPage === 'sheet'">
            <v-divider class="my-3" />
            <div class="text-caption font-weight-bold text-primary mb-2 text-uppercase">Sheet Actions</div>

            <v-list-item
              prepend-icon="mdi-content-save"
              title="Save Character Sheet"
              subtitle="Save current changes"
              class="mb-2 rounded-lg bg-surface-variant border"
              @click="sheetRef?.saveSheet(); showKebabOverlay = false"
            />

            <v-list-item
              prepend-icon="mdi-file-plus-outline"
              title="New Blank Sheet"
              subtitle="Start with a blank sheet"
              class="mb-2 rounded-lg bg-surface-variant border"
              @click="sheetRef?.newBlankSheet(); showKebabOverlay = false"
            />

            <v-list-item
              prepend-icon="mdi-database-import"
              title="Load Demo Character"
              subtitle="Gottfried von Altdorf (Wizard)"
              class="mb-2 rounded-lg bg-surface-variant border"
              @click="sheetRef?.loadMockData(); showKebabOverlay = false"
            />
          </template>

          <!-- Auth Button -->
          <v-list-item
            v-if="!currentUser"
            prepend-icon="mdi-login"
            title="Sign In / Register"
            class="mt-4 rounded-lg bg-primary text-on-primary font-weight-bold"
            @click="showKebabOverlay = false; showAuthDialog = true"
          />
          <v-list-item
            v-else
            prepend-icon="mdi-logout"
            title="Log Out"
            class="mt-4 rounded-lg bg-error text-on-error font-weight-bold"
            @click="logout(); showKebabOverlay = false"
          />
        </v-list>
      </v-card>
    </v-dialog>

    <!-- PAGE CONDITIONAL RENDERING -->
    <template v-if="currentPage === 'sheet'">
      <CharacterSheetPage ref="sheetRef" />
    </template>

    <!-- STANDALONE ADMIN PAGE VIEW -->
    <template v-else-if="currentPage === 'admin'">
      <AdminPage :current-user="currentUser" @back="navigateTo('sheet')" />
    </template>

    <!-- Dialogs -->
    <AuthDialog v-model="showAuthDialog" @login-success="onLoginSuccess" />
  </v-app>
</template>

<style>
/* Clean Theme-Aware CSS Styles without any !important hacks */

/* Hide HTML & Chrome/Edge/Safari/Firefox number spinner arrows */
input[type='number']::-webkit-outer-spin-button,
input[type='number']::-webkit-inner-spin-button {
  -webkit-appearance: none;
  margin: 0;
}

input[type='number'] {
  -moz-appearance: textfield;
}

.kebab-fixed-pos {
  position: fixed;
  top: 18px;
  right: 18px;
  z-index: 1000;
}

.kebab-fab {
  box-shadow: 0 4px 14px rgba(0, 0, 0, 0.25);
}

.wfrp-full-app {
  min-height: 100vh;
  width: 100%;
}
</style>
