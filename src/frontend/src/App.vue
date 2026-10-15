<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useTheme } from 'vuetify'
import AuthDialog from './components/AuthDialog.vue'
import GuestUploadDialog from './components/GuestUploadDialog.vue'
import CharacterSheetPage from './views/CharacterSheetPage.vue'
import { characterApi } from './services/characterApi'
import { useAuth } from './composables/useAuth'

const theme = useTheme()
const isDark = computed(() => theme.global.name.value === 'wfrpDark')

function toggleTheme() {
  theme.global.name.value = isDark.value ? 'wfrpLight' : 'wfrpDark'
}

const router = useRouter()
const route = useRoute()
const { currentUser, isAdmin, showAuthModal, initAuth, logout } = useAuth()

// Dialog & UI State
const showKebabOverlay = ref(false)
const showGuestUploadDialog = ref(false)
const guestCharacterName = ref('')

const sheetRef = ref<InstanceType<typeof CharacterSheetPage> | null>(null)
const fileInputRef = ref<HTMLInputElement | null>(null)

async function handleFileUpload(e: Event) {
  const target = e.target as HTMLInputElement
  const file = target.files?.[0]
  if (!file) return
  if (route.path !== '/') {
    await router.push('/')
  }
  await sheetRef.value?.importJson(file)
  target.value = ''
}

function handleGlobalKeydown(e: KeyboardEvent) {
  if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === 's') {
    if (route.path === '/') {
      e.preventDefault()
      sheetRef.value?.saveSheet()
    }
  }
}

onMounted(async () => {
  window.addEventListener('keydown', handleGlobalKeydown)
  await initAuth()
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleGlobalKeydown)
})

async function onLoginSuccess() {
  // Check if an offline character exists in localStorage to prompt migration
  const guestData = characterApi.loadGuestSheet()
  if (guestData && (guestData.character.name || guestData.character.career)) {
    guestCharacterName.value = guestData.character.name || 'Unnamed Character'
    showGuestUploadDialog.value = true
  } else {
    sheetRef.value?.loadUserCharacter()
  }
}

async function uploadGuestSheet() {
  const guestData = characterApi.loadGuestSheet()
  if (guestData) {
    try {
      const res = await characterApi.saveCharacter(guestData)
      characterApi.clearGuestSheet()
      if (route.path === '/' && sheetRef.value) {
        sheetRef.value.loadCharacterSheet(res.uuid)
      }
    } catch (err) {
      console.warn('Failed to upload guest sheet to cloud:', err)
    }
  }
}

async function handleLogout() {
  await logout()
  if (route.path === '/account' || route.path === '/admin') {
    router.push('/')
  }
}
</script>

<template>
  <v-app class="wfrp-full-app bg-background text-on-background">
    <!-- Kebab POPOUT OVERLAY MODAL (Full Dialog for Character Sheet actions) -->
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
        <div v-if="currentUser"
          class="pa-4 bg-surface-variant rounded-lg mb-4 d-flex align-center justify-space-between border">
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
          <v-list-item :prepend-icon="isDark ? 'mdi-weather-sunny' : 'mdi-weather-night'"
            :title="isDark ? 'Switch to Light Theme' : 'Switch to Dark Theme'"
            class="mb-2 rounded-lg bg-surface-variant border" @click="toggleTheme" />

          <!-- Admin Panel Button (If logged in as Admin) -->
          <v-list-item v-if="isAdmin" prepend-icon="mdi-shield-crown" title="Admin Control Panel"
            subtitle="User & Database Management" color="secondary"
            class="mb-2 rounded-lg bg-secondary text-on-secondary font-weight-bold"
            @click="router.push('/admin'); showKebabOverlay = false" />

          <!-- Account Settings Button (If logged in) -->
          <v-list-item v-if="currentUser" prepend-icon="mdi-account-cog" title="Account Settings"
            subtitle="Active devices & session security"
            class="mb-2 rounded-lg bg-surface-variant border"
            @click="router.push('/account'); showKebabOverlay = false" />

          <!-- Sheet Actions (when viewing character sheet) -->
          <template v-if="route.path === '/'">
            <v-divider class="my-3" />
            <div class="text-caption font-weight-bold text-primary mb-2 text-uppercase">Sheet Actions</div>

            <v-list-item prepend-icon="mdi-content-save" title="Save Character Sheet"
              subtitle="Save current changes (Ctrl+S)" class="mb-2 rounded-lg bg-surface-variant border"
              @click="sheetRef?.saveSheet(); showKebabOverlay = false" />

            <v-list-item prepend-icon="mdi-download" title="Download Data (JSON)"
              subtitle="Export as [date]-[character-name].json" class="mb-2 rounded-lg bg-surface-variant border"
              @click="sheetRef?.exportJson(); showKebabOverlay = false" />

            <v-list-item prepend-icon="mdi-upload" title="Upload Data (JSON)"
              subtitle="Import character sheet from JSON file" class="mb-2 rounded-lg bg-surface-variant border"
              @click="fileInputRef?.click(); showKebabOverlay = false" />

            <v-list-item prepend-icon="mdi-file-plus-outline" title="New Blank Sheet"
              subtitle="Start with a blank sheet" class="mb-2 rounded-lg bg-surface-variant border"
              @click="sheetRef?.newBlankSheet(); showKebabOverlay = false" />

            <v-list-item prepend-icon="mdi-database-import" title="Load Demo Character"
              subtitle="Gottfried von Altdorf (Wizard)" class="mb-2 rounded-lg bg-surface-variant border"
              @click="sheetRef?.loadMockData(); showKebabOverlay = false" />
          </template>

          <!-- Auth Button -->
          <v-list-item v-if="!currentUser" prepend-icon="mdi-login" title="Sign In / Register"
            class="mt-4 rounded-lg bg-primary text-on-primary font-weight-bold"
            @click="showKebabOverlay = false; showAuthModal = true" />
          <v-list-item v-else prepend-icon="mdi-logout" title="Log Out"
            class="mt-4 rounded-lg bg-error text-on-error font-weight-bold"
            @click="handleLogout(); showKebabOverlay = false" />
        </v-list>
      </v-card>
    </v-dialog>

    <!-- Hidden File Input for JSON Upload -->
    <input ref="fileInputRef" type="file" accept=".json,application/json" class="d-none" @change="handleFileUpload" />

    <!-- ROUTER VIEW -->
    <router-view v-slot="{ Component }">
      <component
        :is="Component"
        ref="sheetRef"
        :current-user="currentUser"
        @open-menu="showKebabOverlay = true"
        @logged-out="handleLogout"
      />
    </router-view>

    <!-- Dialogs -->
    <AuthDialog v-model="showAuthModal" @login-success="onLoginSuccess" />
    <GuestUploadDialog v-model="showGuestUploadDialog" :character-name="guestCharacterName"
      @confirm-upload="uploadGuestSheet" />
  </v-app>
</template>

<style>
/* Hide HTML & Chrome/Edge/Safari/Firefox number spinner arrows */
input[type='number']::-webkit-outer-spin-button,
input[type='number']::-webkit-inner-spin-button {
  -webkit-appearance: none;
  margin: 0;
}

input[type='number'] {
  -moz-appearance: textfield;
}

.wfrp-full-app {
  min-height: 100vh;
  width: 100%;
}
</style>
