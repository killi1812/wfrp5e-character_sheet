<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useTheme } from 'vuetify'
import { useAuth } from '../composables/useAuth'

const theme = useTheme()
const isDark = computed(() => theme.global.name.value === 'wfrpDark')
function toggleTheme() {
  theme.global.name.value = isDark.value ? 'wfrpLight' : 'wfrpDark'
}

const router = useRouter()
const { currentUser, isAdmin, logout } = useAuth()

async function handleLogout() {
  await logout()
  router.push('/')
}
</script>

<template>
  <div class="management-layout">
    <!-- Standard Navigation App Bar for Management Pages (Account & Admin) -->
    <v-app-bar color="surface" elevation="2" class="border-b px-2 px-sm-4" density="comfortable">
      <div class="d-flex align-center cursor-pointer mr-4" @click="router.push('/')">
        <v-icon icon="mdi-shield-crown" color="primary" class="mr-2" size="28" />
        <span class="text-h6 font-weight-black text-primary text-uppercase tracking-wider">WFRP 5e</span>
      </div>

      <!-- Nav links -->
      <div class="d-flex align-center gap-1">
        <v-btn
          to="/"
          variant="text"
          prepend-icon="mdi-file-document-outline"
          class="font-weight-bold"
        >
          Character Sheet
        </v-btn>

        <v-btn
          to="/account"
          variant="text"
          prepend-icon="mdi-account-cog"
          :color="$route.path === '/account' ? 'primary' : undefined"
          class="font-weight-bold"
        >
          Account Settings
        </v-btn>

        <v-btn
          v-if="isAdmin"
          to="/admin"
          variant="text"
          prepend-icon="mdi-shield-crown"
          :color="$route.path === '/admin' ? 'secondary' : undefined"
          class="font-weight-bold"
        >
          Admin Console
        </v-btn>
      </div>

      <v-spacer />

      <!-- Right items: Theme toggle, User badge, Log out -->
      <div class="d-flex align-center gap-2">
        <v-tooltip :text="isDark ? 'Switch to Light Theme' : 'Switch to Dark Theme'" location="bottom">
          <template #activator="{ props: tProps }">
            <v-btn
              v-bind="tProps"
              icon
              variant="text"
              size="small"
              @click="toggleTheme"
            >
              <v-icon :icon="isDark ? 'mdi-weather-sunny' : 'mdi-weather-night'" />
            </v-btn>
          </template>
        </v-tooltip>

        <v-chip
          v-if="currentUser"
          color="primary"
          variant="flat"
          size="small"
          prepend-icon="mdi-account-circle"
          class="font-weight-bold text-uppercase d-none d-sm-inline-flex"
        >
          {{ currentUser.username }} ({{ currentUser.role }})
        </v-chip>

        <v-btn
          color="error"
          variant="tonal"
          size="small"
          prepend-icon="mdi-logout"
          @click="handleLogout"
        >
          Log Out
        </v-btn>
      </div>
    </v-app-bar>

    <v-main class="bg-background min-vh-100 pa-3 pa-md-5">
      <router-view />
    </v-main>
  </div>
</template>

<style scoped>
.management-layout {
  min-height: 100vh;
  width: 100%;
}

.tracking-wider {
  letter-spacing: 0.05em;
}

.gap-1 {
  gap: 4px;
}

.gap-2 {
  gap: 8px;
}

.cursor-pointer {
  cursor: pointer;
}
</style>
