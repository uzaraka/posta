<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { authApi } from '../../api/auth'
import { useAuthStore } from '../../stores/auth'
import { useNotificationStore } from '../../stores/notification'
import { useThemeStore } from '../../stores/theme'

const router = useRouter()
const auth = useAuthStore()
const notification = useNotificationStore()
const theme = useThemeStore()

const name = ref('')
const email = ref('')
const password = ref('')
const confirmPassword = ref('')
const loading = ref(false)
const registrationEnabled = ref<boolean | null>(null)

onMounted(async () => {
  try {
    const res = await authApi.registrationStatus()
    registrationEnabled.value = res.data.data.registration_enabled
    if (!res.data.data.registration_enabled) {
      router.replace('/login')
    }
  } catch {
    router.replace('/login')
  }
})

async function handleRegister() {
  if (!name.value.trim() || !email.value || !password.value || !confirmPassword.value) {
    notification.error('Please fill in all fields.')
    return
  }
  if (password.value.length < 8) {
    notification.error('Password must be at least 8 characters.')
    return
  }
  if (password.value !== confirmPassword.value) {
    notification.error('Passwords do not match.')
    return
  }
  loading.value = true
  try {
    const res = await authApi.register(name.value.trim(), email.value.trim(), password.value)
    // Auto-login with the returned token
    const data = res.data.data
    localStorage.setItem('posta_token', data.token)
    localStorage.setItem('posta_user', JSON.stringify(data.user))
    // Reload to pick up auth state
    window.location.href = '/'
  } catch (err: any) {
    const message = err?.response?.data?.error?.message || err?.response?.data?.error || err?.message || 'Registration failed.'
    notification.error(message)
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="auth-page" v-if="registrationEnabled">
    <div class="auth-card">
      <div class="auth-header">
        <div class="auth-logo">
          <img src="/logo.png" alt="Posta" class="logo-img" />
          <span>Posta</span>
        </div>
        <p class="auth-subtitle">Create your account</p>
      </div>

      <form class="auth-form" @submit.prevent="handleRegister">
        <div class="form-group">
          <label class="form-label" for="name">Name</label>
          <input id="name" v-model="name" type="text" class="form-input" placeholder="Your name" autocomplete="name" />
        </div>
        <div class="form-group">
          <label class="form-label" for="email">Email</label>
          <input id="email" v-model="email" type="email" class="form-input" placeholder="you@example.com" autocomplete="email" />
        </div>
        <div class="form-group">
          <label class="form-label" for="password">Password</label>
          <input id="password" v-model="password" type="password" class="form-input" placeholder="Minimum 8 characters" autocomplete="new-password" />
        </div>
        <div class="form-group">
          <label class="form-label" for="confirm-password">Confirm Password</label>
          <input id="confirm-password" v-model="confirmPassword" type="password" class="form-input" placeholder="Re-enter your password" autocomplete="new-password" />
        </div>
        <button type="submit" class="btn btn-primary auth-btn" :disabled="loading">
          <span v-if="loading" class="spinner"></span>
          {{ loading ? 'Creating account...' : 'Create Account' }}
        </button>
      </form>

      <div class="auth-footer">
        <span>Already have an account?</span>
        <router-link to="/login">Sign in</router-link>
      </div>
    </div>

    <button class="theme-btn" @click="theme.toggle()" :title="theme.isDark ? 'Light mode' : 'Dark mode'">
      <svg v-if="theme.isDark" width="18" height="18" viewBox="0 0 16 16" fill="none"><circle cx="8" cy="8" r="3" stroke="currentColor" stroke-width="1.5"/><path d="M8 1v2M8 13v2M1 8h2M13 8h2M3.05 3.05l1.41 1.41M11.54 11.54l1.41 1.41M3.05 12.95l1.41-1.41M11.54 4.46l1.41-1.41" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/></svg>
      <svg v-else width="18" height="18" viewBox="0 0 16 16" fill="none"><path d="M14 9.5A6.5 6.5 0 016.5 2 6.5 6.5 0 1014 9.5z" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/></svg>
    </button>
  </div>
</template>

<style scoped>
.auth-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-secondary);
  padding: 20px;
  position: relative;
}

.auth-card {
  width: 100%;
  max-width: 400px;
  background: var(--bg-primary);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-lg);
}

.auth-header { text-align: center; padding: 36px 32px 0; }

.auth-logo {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  margin-bottom: 10px;
}
.auth-logo .logo-img {
  width: 80px;
  height: 80px;
  object-fit: contain;
}
.auth-logo span {
  font-size: 26px;
  font-weight: 800;
  color: var(--text-primary);
  letter-spacing: -0.5px;
}

.auth-subtitle { font-size: 14px; color: var(--text-muted); }

.auth-form { padding: 28px 32px 20px; }

.auth-btn {
  width: 100%;
  padding: 11px 18px;
  font-size: 15px;
  margin-top: 4px;
}

.auth-footer {
  text-align: center;
  padding: 0 32px 28px;
  font-size: 14px;
  color: var(--text-muted);
  display: flex;
  gap: 6px;
  justify-content: center;
}
.auth-footer a { color: var(--primary-500); font-weight: 500; }

.theme-btn {
  position: fixed;
  top: 20px;
  right: 20px;
  background: var(--bg-primary);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius);
  padding: 10px;
  cursor: pointer;
  color: var(--text-tertiary);
  display: flex;
  align-items: center;
  transition: all var(--transition);
  box-shadow: var(--shadow-sm);
}
.theme-btn:hover { color: var(--text-primary); border-color: var(--border-input); }
</style>
