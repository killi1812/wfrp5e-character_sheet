<script setup lang="ts">
import { ref } from 'vue'
import { useAuth } from '../composables/useAuth'

const props = defineProps<{
  modelValue: boolean
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'login-success', user: { username: string; role: string; token: string; uuid?: string }): void
}>()

const { login, register } = useAuth()

const tab = ref<'login' | 'register'>('login')
const loading = ref(false)
const errorMessage = ref('')

// Login fields
const loginUsername = ref('')
const loginPassword = ref('')

// Register fields
const regUsername = ref('')
const regEmail = ref('')
const regPassword = ref('')

function close() {
  emit('update:modelValue', false)
  errorMessage.value = ''
}

async function handleLogin() {
  if (!loginUsername.value || !loginPassword.value) {
    errorMessage.value = 'Please enter email/username and password.'
    return
  }
  loading.value = true
  errorMessage.value = ''

  try {
    const user = await login(loginUsername.value, loginPassword.value)
    emit('login-success', user)
    close()
  } catch (err: any) {
    errorMessage.value = err?.response?.data?.error || 'Invalid login credentials'
  } finally {
    loading.value = false
  }
}

async function handleRegister() {
  if (!regUsername.value || !regPassword.value) {
    errorMessage.value = 'Username and password are required.'
    return
  }
  if (regPassword.value.length < 6) {
    errorMessage.value = 'Password must be at least 6 characters long.'
    return
  }
  loading.value = true
  errorMessage.value = ''

  try {
    await register(regUsername.value, regPassword.value, regEmail.value || undefined)
    loginUsername.value = regUsername.value
    loginPassword.value = regPassword.value
    await handleLogin()
  } catch (err: any) {
    errorMessage.value = err?.response?.data?.error || err?.message || 'Registration failed.'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <v-dialog :model-value="modelValue" max-width="480" @update:model-value="emit('update:modelValue', $event)">
    <v-card color="surface" class="pa-4 rounded-lg border">
      <v-card-title class="d-flex justify-space-between align-center text-h5 font-weight-bold text-on-surface">
        <span>{{ tab === 'login' ? 'Welcome Back' : 'Create Account' }}</span>
        <v-btn icon="mdi-close" variant="text" size="small" @click="close" />
      </v-card-title>

      <v-tabs v-model="tab" color="primary" grow class="mb-4">
        <v-tab value="login">Login</v-tab>
        <v-tab value="register">Register</v-tab>
      </v-tabs>

      <v-alert v-if="errorMessage" type="error" variant="tonal" class="mb-4" density="compact" closable>
        {{ errorMessage }}
      </v-alert>

      <v-card-text class="pa-0">
        <v-window v-model="tab">
          <!-- Login Form -->
          <v-window-item value="login">
            <v-form @submit.prevent="handleLogin">
              <v-text-field
                v-model="loginUsername"
                label="Username"
                prepend-inner-icon="mdi-account"
                variant="outlined"
                density="comfortable"
                class="mb-3"
                required
              />
              <v-text-field
                v-model="loginPassword"
                label="Password"
                prepend-inner-icon="mdi-lock"
                type="password"
                variant="outlined"
                density="comfortable"
                class="mb-4"
                required
              />
              <v-btn
                type="submit"
                color="primary"
                block
                size="large"
                :loading="loading"
                elevation="2"
                class="font-weight-bold"
              >
                Sign In
              </v-btn>
            </v-form>
          </v-window-item>

          <!-- Register Form -->
          <v-window-item value="register">
            <v-form @submit.prevent="handleRegister">
              <v-text-field
                v-model="regUsername"
                label="Username"
                prepend-inner-icon="mdi-account-plus"
                variant="outlined"
                density="comfortable"
                class="mb-3"
                required
              />
              <v-text-field
                v-model="regEmail"
                label="Email Address (Optional)"
                prepend-inner-icon="mdi-email"
                type="email"
                variant="outlined"
                density="comfortable"
                class="mb-3"
              />
              <v-text-field
                v-model="regPassword"
                label="Password"
                prepend-inner-icon="mdi-lock"
                type="password"
                variant="outlined"
                density="comfortable"
                class="mb-4"
                required
              />
              <v-btn
                type="submit"
                color="primary"
                block
                size="large"
                :loading="loading"
                elevation="2"
                class="font-weight-bold"
              >
                Register Account
              </v-btn>
            </v-form>
          </v-window-item>
        </v-window>
      </v-card-text>
    </v-card>
  </v-dialog>
</template>
