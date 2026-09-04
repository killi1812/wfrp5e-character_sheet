<script setup lang="ts">
import { ref, computed } from 'vue'

const props = defineProps<{
  currentUser: { username: string; role: string; token: string } | null
}>()

const emit = defineEmits<{
  (e: 'back'): void
}>()

const activeTab = ref<'users' | 'sheets' | 'system'>('users')
const searchUsers = ref('')
const searchSheets = ref('')

// User Management Data
const users = ref([
  { id: '1', username: 'admin', email: 'admin@wfrp.local', role: 'admin', status: 'Active', created: '2026-09-01' },
  { id: '2', username: 'gottfried', email: 'gottfried@altdorf.de', role: 'user', status: 'Active', created: '2026-09-02' },
  { id: '3', username: 'elena_shadow', email: 'elena@nuln.org', role: 'user', status: 'Active', created: '2026-09-03' },
  { id: '4', username: 'bardin_hammer', email: 'bardin@karaz-a-karak.com', role: 'user', status: 'Suspended', created: '2026-09-04' },
])

// Saved Character Sheets Data
const characterSheets = ref([
  { uuid: 'a1b2c3d4-e5f6-47a8-b9c0-d1e2f3a4b5c6', name: 'Gottfried von Altdorf', species: 'Human', career: 'Wizard', user: 'gottfried', updated: '2026-09-04 16:45' },
  { uuid: 'b2c3d4e5-f6a7-48b9-c0d1-e2f3a4b5c6d7', name: 'Elena Shadowstep', species: 'Elf', career: 'Spymaster', user: 'elena_shadow', updated: '2026-09-03 14:20' },
  { uuid: 'c3d4e5f6-a7b8-49c0-d1e2-f3a4b5c6d7e8', name: 'Bardin Ironfist', species: 'Dwarf', career: 'Ironbreaker', user: 'bardin_hammer', updated: '2026-09-04 11:10' },
])

const filteredUsers = computed(() => {
  if (!searchUsers.value) return users.value
  const q = searchUsers.value.toLowerCase()
  return users.value.filter(u => u.username.toLowerCase().includes(q) || u.email.toLowerCase().includes(q) || u.role.toLowerCase().includes(q))
})

const filteredSheets = computed(() => {
  if (!searchSheets.value) return characterSheets.value
  const q = searchSheets.value.toLowerCase()
  return characterSheets.value.filter(s => s.name.toLowerCase().includes(q) || s.user.toLowerCase().includes(q) || s.career.toLowerCase().includes(q))
})

// System Stats
const systemStats = ref({
  dbStatus: 'Connected (MongoDB)',
  totalUsers: users.value.length,
  totalSheets: characterSheets.value.length,
  serverUptime: '99.99%',
  apiLatency: '11ms',
  memoryUsage: '42.8 MB / 512 MB',
})

function deleteUser(id: string) {
  users.value = users.value.filter(u => u.id !== id)
}

function deleteSheet(uuidStr: string) {
  characterSheets.value = characterSheets.value.filter(s => s.uuid !== uuidStr)
}
</script>

<template>
  <div class="admin-page-container bg-background min-vh-100 pa-4 pa-md-6">
    <!-- Top Header Bar -->
    <v-card color="surface" elevation="3" class="pa-4 rounded-xl mb-6 border">
      <div class="d-flex flex-wrap justify-space-between align-center gap-3">
        <div class="d-flex align-center">
          <v-btn
            icon="mdi-arrow-left"
            color="primary"
            variant="tonal"
            class="mr-4"
            @click="emit('back')"
          />
          <div>
            <h1 class="text-h4 font-weight-black text-primary mb-0 d-flex align-center">
              <v-icon icon="mdi-shield-crown" color="primary" class="mr-2" />
              Admin Management Console
            </h1>
            <div class="text-subtitle-2 text-medium-emphasis">
              System Control, Database & User Administration
            </div>
          </div>
        </div>

        <div class="d-flex align-center gap-2">
          <v-chip color="secondary" variant="flat" size="medium" prepend-icon="mdi-shield-account">
            Admin: {{ currentUser?.username || 'Administrator' }}
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

    <!-- Navigation Tabs -->
    <v-card color="surface" elevation="2" class="rounded-xl border">
      <v-tabs v-model="activeTab" color="primary" class="border-bottom px-4">
        <v-tab value="users">
          <v-icon icon="mdi-account-group" start />
          User Accounts ({{ users.length }})
        </v-tab>
        <v-tab value="sheets">
          <v-icon icon="mdi-book-multiple" start />
          Saved Character Sheets ({{ characterSheets.length }})
        </v-tab>
        <v-tab value="system">
          <v-icon icon="mdi-server" start />
          System Health & MongoDB
        </v-tab>
      </v-tabs>

      <v-card-text class="pa-6">
        <v-window v-model="activeTab">
          
          <!-- TAB 1: USER MANAGEMENT -->
          <v-window-item value="users">
            <div class="d-flex flex-wrap justify-space-between align-center gap-3 mb-4">
              <v-text-field
                v-model="searchUsers"
                prepend-inner-icon="mdi-magnify"
                label="Search users by name, email or role..."
                variant="outlined"
                density="compact"
                hide-details
                style="max-width: 360px;"
              />
              <v-btn color="primary" prepend-icon="mdi-account-plus">
                Create User Account
              </v-btn>
            </div>

            <v-table hover class="bg-transparent border rounded-lg">
              <thead>
                <tr>
                  <th>Username</th>
                  <th>Email</th>
                  <th>Role</th>
                  <th>Status</th>
                  <th>Joined Date</th>
                  <th class="text-right">Actions</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="user in filteredUsers" :key="user.id">
                  <td class="font-weight-bold">{{ user.username }}</td>
                  <td>{{ user.email }}</td>
                  <td>
                    <v-chip size="small" :color="user.role === 'admin' ? 'secondary' : 'info'" variant="flat">
                      {{ user.role }}
                    </v-chip>
                  </td>
                  <td>
                    <v-badge dot :color="user.status === 'Active' ? 'success' : 'error'" inline class="mr-2" />
                    {{ user.status }}
                  </td>
                  <td class="text-caption text-medium-emphasis">{{ user.created }}</td>
                  <td class="text-right">
                    <v-btn icon="mdi-pencil" size="small" variant="text" color="primary" class="mr-1" />
                    <v-btn icon="mdi-delete" size="small" variant="text" color="error" @click="deleteUser(user.id)" />
                  </td>
                </tr>
              </tbody>
            </v-table>
          </v-window-item>

          <!-- TAB 2: SAVED CHARACTER SHEETS -->
          <v-window-item value="sheets">
            <div class="d-flex justify-space-between align-center mb-4">
              <v-text-field
                v-model="searchSheets"
                prepend-inner-icon="mdi-magnify"
                label="Search character sheets..."
                variant="outlined"
                density="compact"
                hide-details
                style="max-width: 360px;"
              />
              <v-chip color="primary" variant="tonal">Total Sheets: {{ characterSheets.length }}</v-chip>
            </div>

            <v-table hover class="bg-transparent border rounded-lg">
              <thead>
                <tr>
                  <th>Character Name</th>
                  <th>Species</th>
                  <th>Career</th>
                  <th>Owner User</th>
                  <th>UUID</th>
                  <th>Last Modified</th>
                  <th class="text-right">Actions</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="sheet in filteredSheets" :key="sheet.uuid">
                  <td class="font-weight-bold text-primary">{{ sheet.name }}</td>
                  <td>{{ sheet.species }}</td>
                  <td>{{ sheet.career }}</td>
                  <td>
                    <v-chip size="x-small" color="info" variant="outlined">{{ sheet.user }}</v-chip>
                  </td>
                  <td class="text-caption font-monospace text-medium-emphasis">{{ sheet.uuid }}</td>
                  <td class="text-caption">{{ sheet.updated }}</td>
                  <td class="text-right">
                    <v-btn icon="mdi-eye" size="small" variant="text" color="primary" class="mr-1" />
                    <v-btn icon="mdi-delete" size="small" variant="text" color="error" @click="deleteSheet(sheet.uuid)" />
                  </td>
                </tr>
              </tbody>
            </v-table>
          </v-window-item>

          <!-- TAB 3: SYSTEM METRICS -->
          <v-window-item value="system">
            <v-row>
              <v-col cols="12" sm="6" md="4">
                <v-card variant="tonal" color="success" class="pa-4 text-center rounded-lg border">
                  <v-icon icon="mdi-database-check" size="40" class="mb-2" />
                  <div class="text-caption font-weight-bold text-uppercase">MongoDB Connection</div>
                  <div class="text-h6 font-weight-black">{{ systemStats.dbStatus }}</div>
                </v-card>
              </v-col>
              <v-col cols="12" sm="6" md="4">
                <v-card variant="tonal" color="primary" class="pa-4 text-center rounded-lg border">
                  <v-icon icon="mdi-account-multiple-check" size="40" class="mb-2" />
                  <div class="text-caption font-weight-bold text-uppercase">Registered Users</div>
                  <div class="text-h6 font-weight-black">{{ systemStats.totalUsers }}</div>
                </v-card>
              </v-col>
              <v-col cols="12" sm="6" md="4">
                <v-card variant="tonal" color="info" class="pa-4 text-center rounded-lg border">
                  <v-icon icon="mdi-clock-fast" size="40" class="mb-2" />
                  <div class="text-caption font-weight-bold text-uppercase">API Response Time</div>
                  <div class="text-h6 font-weight-black">{{ systemStats.apiLatency }}</div>
                </v-card>
              </v-col>
              <v-col cols="12" sm="6" md="4">
                <v-card variant="tonal" color="warning" class="pa-4 text-center rounded-lg border">
                  <v-icon icon="mdi-server-network" size="40" class="mb-2" />
                  <div class="text-caption font-weight-bold text-uppercase">Server Uptime</div>
                  <div class="text-h6 font-weight-black">{{ systemStats.serverUptime }}</div>
                </v-card>
              </v-col>
              <v-col cols="12" sm="6" md="4">
                <v-card variant="tonal" color="secondary" class="pa-4 text-center rounded-lg border">
                  <v-icon icon="mdi-memory" size="40" class="mb-2" />
                  <div class="text-caption font-weight-bold text-uppercase">Memory Usage</div>
                  <div class="text-h6 font-weight-black">{{ systemStats.memoryUsage }}</div>
                </v-card>
              </v-col>
            </v-row>
          </v-window-item>

        </v-window>
      </v-card-text>
    </v-card>
  </div>
</template>

<style scoped>
.admin-page-container {
  width: 100%;
  min-height: 100vh;
}

.gap-2 { gap: 8px; }
.gap-3 { gap: 12px; }

.font-monospace {
  font-family: monospace;
}
</style>
