<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { characterApi, type SessionDto } from '../services/characterApi'

const props = defineProps<{
  currentUser: { username: string; role: string; token: string; uuid?: string } | null
}>()

const emit = defineEmits<{
  (e: 'back'): void
  (e: 'logged-out'): void
}>()

const activeTab = ref<'sessions' | 'profile'>('sessions')
const sessions = ref<SessionDto[]>([])
const loading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')

// Dialog controls
const showRevokeDialog = ref(false)
const sessionToRevoke = ref<SessionDto | null>(null)
const showLogoutAllDialog = ref(false)
const actionInProgress = ref(false)

async function fetchSessions() {
  loading.value = true
  errorMessage.value = ''
  try {
    sessions.value = await characterApi.listSessions()
  } catch (err: any) {
    errorMessage.value = err?.message || 'Failed to load active sessions.'
  } finally {
    loading.value = false
  }
}

function formatDate(dateStr: string): string {
  if (!dateStr) return '—'
  try {
    const d = new Date(dateStr)
    return d.toLocaleString(undefined, {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    })
  } catch {
    return dateStr
  }
}

function getDeviceIcon(userAgent: string): string {
  const ua = (userAgent || '').toLowerCase()
  if (ua.includes('mobile') || ua.includes('android') || ua.includes('iphone')) {
    return 'mdi-cellphone'
  }
  if (ua.includes('ipad') || ua.includes('tablet')) {
    return 'mdi-tablet'
  }
  return 'mdi-laptop'
}

function getBrowserName(userAgent: string): string {
  if (!userAgent) return 'Unknown Device'
  const ua = userAgent
  if (ua.includes('Firefox/')) return 'Firefox Browser'
  if (ua.includes('Edg/')) return 'Microsoft Edge'
  if (ua.includes('Chrome/')) return 'Chrome Browser'
  if (ua.includes('Safari/')) return 'Safari Browser'
  return userAgent.length > 35 ? userAgent.substring(0, 32) + '...' : userAgent
}

function confirmRevokeSession(session: SessionDto) {
  sessionToRevoke.value = session
  showRevokeDialog.value = true
}

async function executeRevokeSession() {
  if (!sessionToRevoke.value) return
  actionInProgress.value = true
  const target = sessionToRevoke.value
  showRevokeDialog.value = false

  try {
    const ok = await characterApi.revokeSession(target.sessionUuid)
    if (ok) {
      if (target.isCurrent) {
        emit('logged-out')
      } else {
        successMessage.value = 'Session revoked successfully.'
        await fetchSessions()
      }
    } else {
      errorMessage.value = 'Failed to revoke session. Please try again.'
    }
  } catch (err: any) {
    errorMessage.value = err?.message || 'Error revoking session.'
  } finally {
    actionInProgress.value = false
    sessionToRevoke.value = null
  }
}

async function executeLogoutAll() {
  actionInProgress.value = true
  showLogoutAllDialog.value = false

  try {
    const ok = await characterApi.logoutAll()
    if (ok) {
      emit('logged-out')
    } else {
      errorMessage.value = 'Failed to sign out from all devices.'
    }
  } catch (err: any) {
    errorMessage.value = err?.message || 'Error signing out from devices.'
  } finally {
    actionInProgress.value = false
  }
}

onMounted(() => {
  fetchSessions()
})
</script>

<template>
  <div class="account-page-container bg-background min-vh-100 pa-4 pa-md-6">
    <!-- Top Header Bar -->
    <v-card color="surface" elevation="3" class="pa-4 rounded-xl mb-6 border">
      <div class="d-flex flex-wrap justify-space-between align-center gap-3">
        <div class="d-flex align-center">
          <v-btn
            icon="mdi-arrow-left"
            color="primary"
            variant="tonal"
            class="mr-4"
            aria-label="Back to character sheet"
            @click="emit('back')"
          />
          <div>
            <h1 class="text-h4 font-weight-black text-primary mb-0 d-flex align-center">
              <v-icon icon="mdi-account-cog" color="primary" class="mr-2" />
              Account Settings
            </h1>
            <div class="text-subtitle-2 text-medium-emphasis">
              Manage security, active devices, and account configuration
            </div>
          </div>
        </div>

        <div class="d-flex align-center gap-2">
          <v-chip color="primary" variant="flat" size="medium" prepend-icon="mdi-account-circle">
            {{ currentUser?.username || 'User' }} ({{ currentUser?.role || 'user' }})
          </v-chip>
          <v-btn
            color="primary"
            variant="outlined"
            prepend-icon="mdi-file-document-outline"
            @click="emit('back')"
          >
            Back to Character Sheet
          </v-btn>
        </div>
      </div>
    </v-card>

    <!-- Alerts -->
    <v-alert
      v-if="errorMessage"
      type="error"
      variant="tonal"
      class="mb-4 rounded-lg"
      density="comfortable"
      closable
      @click:close="errorMessage = ''"
    >
      {{ errorMessage }}
    </v-alert>

    <v-alert
      v-if="successMessage"
      type="success"
      variant="tonal"
      class="mb-4 rounded-lg"
      density="comfortable"
      closable
      @click:close="successMessage = ''"
    >
      {{ successMessage }}
    </v-alert>

    <!-- Navigation Tabs -->
    <v-card color="surface" elevation="2" class="rounded-xl border">
      <v-tabs v-model="activeTab" color="primary" class="border-bottom px-4">
        <v-tab value="sessions">
          <v-icon icon="mdi-devices" start />
          Active Sessions & Security ({{ sessions.length }})
        </v-tab>
        <v-tab value="profile">
          <v-icon icon="mdi-card-account-details-outline" start />
          Account Profile
        </v-tab>
      </v-tabs>

      <v-card-text class="pa-6">
        <v-window v-model="activeTab">
          <!-- TAB 1: SESSIONS & SECURITY -->
          <v-window-item value="sessions">
            <div class="d-flex flex-wrap justify-space-between align-center gap-3 mb-4">
              <div>
                <h2 class="text-h6 font-weight-bold text-on-surface mb-1">
                  Active Logged-In Devices
                </h2>
                <p class="text-body-2 text-medium-emphasis mb-0">
                  These devices and browsers currently have active session access. Refresh tokens last 30 days and rotate automatically.
                </p>
              </div>

              <div class="d-flex align-center gap-2">
                <v-btn
                  color="secondary"
                  variant="outlined"
                  prepend-icon="mdi-refresh"
                  :loading="loading"
                  @click="fetchSessions"
                >
                  Refresh
                </v-btn>
                <v-btn
                  color="error"
                  variant="flat"
                  prepend-icon="mdi-logout-variant"
                  :disabled="sessions.length === 0"
                  @click="showLogoutAllDialog = true"
                >
                  Sign Out All Devices
                </v-btn>
              </div>
            </div>

            <!-- Loading Spinner -->
            <div v-if="loading" class="d-flex justify-center my-8">
              <v-progress-circular indeterminate color="primary" size="48" />
            </div>

            <!-- Empty Sessions State -->
            <v-card
              v-else-if="sessions.length === 0"
              color="surface-variant"
              class="pa-8 text-center rounded-lg border"
            >
              <v-icon icon="mdi-shield-check" size="48" color="medium-emphasis" class="mb-2" />
              <div class="text-subtitle-1 font-weight-medium">No active sessions found</div>
              <div class="text-body-2 text-medium-emphasis">Try refreshing the session list.</div>
            </v-card>

            <!-- Sessions Table -->
            <v-table v-else hover class="bg-transparent border rounded-lg">
              <thead>
                <tr>
                  <th class="font-weight-bold text-primary">Device & Browser</th>
                  <th class="font-weight-bold text-primary">IP Address</th>
                  <th class="font-weight-bold text-primary">First Signed In</th>
                  <th class="font-weight-bold text-primary">Last Active</th>
                  <th class="font-weight-bold text-primary">Status</th>
                  <th class="font-weight-bold text-primary text-right">Actions</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="s in sessions" :key="s.sessionUuid">
                  <td>
                    <div class="d-flex align-center">
                      <v-avatar color="primary" variant="tonal" size="36" class="mr-3">
                        <v-icon :icon="getDeviceIcon(s.userAgent)" size="20" />
                      </v-avatar>
                      <div>
                        <div class="font-weight-bold text-on-surface">
                          {{ getBrowserName(s.userAgent) }}
                        </div>
                        <div class="text-caption text-medium-emphasis text-truncate" style="max-width: 280px;" :title="s.userAgent">
                          {{ s.userAgent || 'Unknown user agent' }}
                        </div>
                      </div>
                    </div>
                  </td>
                  <td>
                    <code class="text-body-2 bg-surface-variant pa-1 rounded font-weight-medium">
                      {{ s.clientIp || '—' }}
                    </code>
                  </td>
                  <td class="text-body-2 text-medium-emphasis">
                    {{ formatDate(s.createdAt) }}
                  </td>
                  <td class="text-body-2 text-medium-emphasis">
                    {{ formatDate(s.updatedAt) }}
                  </td>
                  <td>
                    <v-chip
                      v-if="s.isCurrent"
                      color="primary"
                      variant="flat"
                      size="small"
                      prepend-icon="mdi-check-circle"
                      class="font-weight-bold"
                    >
                      Current Device
                    </v-chip>
                    <v-chip
                      v-else
                      color="secondary"
                      variant="tonal"
                      size="small"
                    >
                      Active
                    </v-chip>
                  </td>
                  <td class="text-right">
                    <v-btn
                      color="error"
                      variant="tonal"
                      size="small"
                      prepend-icon="mdi-cancel"
                      :disabled="actionInProgress"
                      @click="confirmRevokeSession(s)"
                    >
                      Revoke
                    </v-btn>
                  </td>
                </tr>
              </tbody>
            </v-table>
          </v-window-item>

          <!-- TAB 2: PROFILE & PREFERENCES (Extensible) -->
          <v-window-item value="profile">
            <v-card color="surface-variant" class="pa-6 rounded-lg border">
              <h2 class="text-h6 font-weight-bold text-on-surface mb-3">
                User Profile Summary
              </h2>
              <v-row>
                <v-col cols="12" sm="6">
                  <div class="text-caption text-medium-emphasis text-uppercase font-weight-bold">Username</div>
                  <div class="text-body-1 font-weight-medium text-on-surface mb-3">
                    {{ currentUser?.username || '—' }}
                  </div>
                </v-col>
                <v-col cols="12" sm="6">
                  <div class="text-caption text-medium-emphasis text-uppercase font-weight-bold">Account Role</div>
                  <div class="text-body-1 font-weight-medium text-on-surface mb-3">
                    <v-chip color="secondary" size="small" class="text-uppercase font-weight-bold">
                      {{ currentUser?.role || 'user' }}
                    </v-chip>
                  </div>
                </v-col>
                <v-col cols="12" sm="6">
                  <div class="text-caption text-medium-emphasis text-uppercase font-weight-bold">User UUID</div>
                  <div class="text-body-2 text-medium-emphasis mb-3 font-mono">
                    {{ currentUser?.uuid || 'Available via JWT token' }}
                  </div>
                </v-col>
                <v-col cols="12" sm="6">
                  <div class="text-caption text-medium-emphasis text-uppercase font-weight-bold">Session Security</div>
                  <div class="text-body-2 text-medium-emphasis mb-3">
                    Tokens protected by HTTP-only Lax cookies and SHA-256 rotating hashes.
                  </div>
                </v-col>
              </v-row>
            </v-card>
          </v-window-item>
        </v-window>
      </v-card-text>
    </v-card>

    <!-- Confirmation Dialog: Revoke Session -->
    <v-dialog v-model="showRevokeDialog" max-width="450">
      <v-card color="surface" class="pa-4 rounded-xl border">
        <v-card-title class="d-flex align-center font-weight-bold text-h6">
          <v-icon icon="mdi-alert-circle" color="warning" class="mr-2" />
          Revoke Device Session?
        </v-card-title>
        <v-card-text class="pt-2 text-body-1">
          Are you sure you want to revoke access for this device
          <strong>({{ getBrowserName(sessionToRevoke?.userAgent || '') }})</strong>?
          <div v-if="sessionToRevoke?.isCurrent" class="mt-2 text-error font-weight-bold">
            Note: This is your current device. You will be logged out immediately.
          </div>
        </v-card-text>
        <v-card-actions class="justify-end gap-2">
          <v-btn variant="text" :disabled="actionInProgress" @click="showRevokeDialog = false">
            Cancel
          </v-btn>
          <v-btn
            color="error"
            variant="flat"
            :loading="actionInProgress"
            @click="executeRevokeSession"
          >
            Revoke Access
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Confirmation Dialog: Sign Out All Devices -->
    <v-dialog v-model="showLogoutAllDialog" max-width="480">
      <v-card color="surface" class="pa-4 rounded-xl border">
        <v-card-title class="d-flex align-center font-weight-bold text-h6">
          <v-icon icon="mdi-logout-variant" color="error" class="mr-2" />
          Sign Out All Devices?
        </v-card-title>
        <v-card-text class="pt-2 text-body-1">
          This will immediately invalidate all active sessions across all devices and browsers. You will need to sign in again everywhere.
        </v-card-text>
        <v-card-actions class="justify-end gap-2">
          <v-btn variant="text" :disabled="actionInProgress" @click="showLogoutAllDialog = false">
            Cancel
          </v-btn>
          <v-btn
            color="error"
            variant="flat"
            :loading="actionInProgress"
            @click="executeLogoutAll"
          >
            Sign Out Everywhere
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<style scoped>
.account-page-container {
  max-width: 1400px;
  margin: 0 auto;
}

.border-bottom {
  border-bottom: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
}

.gap-2 {
  gap: 8px;
}

.gap-3 {
  gap: 12px;
}

.font-mono {
  font-family: monospace;
}
</style>
