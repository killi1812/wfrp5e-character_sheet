<script setup lang="ts">
import { ref } from 'vue'

defineProps<{
  modelValue: boolean
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
}>()

const activeTab = ref<'users' | 'system'>('users')
const search = ref('')

const users = ref([
  { id: '1', username: 'admin', email: 'admin@wfrp.local', role: 'admin', status: 'Active', created: '2026-09-01' },
  { id: '2', username: 'gottfried', email: 'gottfried@altdorf.de', role: 'user', status: 'Active', created: '2026-09-02' },
  { id: '3', username: 'elena_shadow', email: 'elena@nuln.org', role: 'user', status: 'Active', created: '2026-09-03' },
  { id: '4', username: 'bardin_hammer', email: 'bardin@karaz-a-karak.com', role: 'user', status: 'Suspended', created: '2026-09-04' },
])

const systemStats = ref({
  dbStatus: 'Connected (MongoDB)',
  totalCharacters: 42,
  activeUsers: 4,
  serverUptime: '99.98%',
  apiLatency: '14ms',
})

function close() {
  emit('update:modelValue', false)
}
</script>

<template>
  <v-dialog :model-value="modelValue" max-width="850" @update:model-value="emit('update:modelValue', $event)">
    <v-card color="surface" class="pa-4 rounded-lg">
      <v-card-title class="d-flex justify-space-between align-center text-h5 font-weight-bold">
        <div class="d-flex align-center">
          <v-icon icon="mdi-shield-crown" color="primary" class="mr-2" size="large" />
          <span>Admin Control Panel</span>
        </div>
        <v-btn icon="mdi-close" variant="text" size="small" @click="close" />
      </v-card-title>

      <v-tabs v-model="activeTab" color="primary" class="mb-4">
        <v-tab value="users">
          <v-icon icon="mdi-account-group" start />
          User Management
        </v-tab>
        <v-tab value="system">
          <v-icon icon="mdi-server" start />
          System & Database
        </v-tab>
      </v-tabs>

      <v-card-text class="pa-0">
        <v-window v-model="activeTab">
          <!-- Users Management -->
          <v-window-item value="users">
            <div class="d-flex justify-space-between align-center mb-4">
              <v-text-field
                v-model="search"
                prepend-inner-icon="mdi-magnify"
                label="Search users..."
                variant="outlined"
                density="compact"
                hide-details
                style="max-width: 300px"
              />
              <v-chip color="primary" variant="tonal">Total Users: {{ users.length }}</v-chip>
            </div>

            <v-table hover>
              <thead>
                <tr>
                  <th>Username</th>
                  <th>Email</th>
                  <th>Role</th>
                  <th>Status</th>
                  <th>Created</th>
                  <th class="text-right">Actions</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="user in users" :key="user.id">
                  <td class="font-weight-bold">{{ user.username }}</td>
                  <td>{{ user.email }}</td>
                  <td>
                    <v-chip
                      size="small"
                      :color="user.role === 'admin' ? 'secondary' : 'info'"
                      variant="flat"
                    >
                      {{ user.role }}
                    </v-chip>
                  </td>
                  <td>
                    <v-badge
                      dot
                      :color="user.status === 'Active' ? 'success' : 'error'"
                      inline
                      class="mr-2"
                    />
                    {{ user.status }}
                  </td>
                  <td class="text-caption text-medium-emphasis">{{ user.created }}</td>
                  <td class="text-right">
                    <v-btn icon="mdi-pencil" size="x-small" variant="text" color="primary" class="mr-1" />
                    <v-btn icon="mdi-delete" size="x-small" variant="text" color="error" />
                  </td>
                </tr>
              </tbody>
            </v-table>
          </v-window-item>

          <!-- System Metrics -->
          <v-window-item value="system">
            <v-row>
              <v-col cols="12" sm="6" md="4">
                <v-card variant="tonal" color="success" class="pa-4 text-center">
                  <v-icon icon="mdi-database-check" size="36" class="mb-2" />
                  <div class="text-caption">Database Status</div>
                  <div class="text-h6 font-weight-bold">{{ systemStats.dbStatus }}</div>
                </v-card>
              </v-col>
              <v-col cols="12" sm="6" md="4">
                <v-card variant="tonal" color="primary" class="pa-4 text-center">
                  <v-icon icon="mdi-book-multiple" size="36" class="mb-2" />
                  <div class="text-caption">Saved Character Sheets</div>
                  <div class="text-h6 font-weight-bold">{{ systemStats.totalCharacters }}</div>
                </v-card>
              </v-col>
              <v-col cols="12" sm="6" md="4">
                <v-card variant="tonal" color="info" class="pa-4 text-center">
                  <v-icon icon="mdi-clock-fast" size="36" class="mb-2" />
                  <div class="text-caption">API Latency</div>
                  <div class="text-h6 font-weight-bold">{{ systemStats.apiLatency }}</div>
                </v-card>
              </v-col>
            </v-row>
          </v-window-item>
        </v-window>
      </v-card-text>
    </v-card>
  </v-dialog>
</template>
