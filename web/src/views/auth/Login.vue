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

const email = ref('')
const password = ref('')
const twoFactorCode = ref('')
const loading = ref(false)
const requires2FA = ref(false)
const registrationEnabled = ref(false)

onMounted(async () => {
  try {
    const res = await authApi.registrationStatus()
    registrationEnabled.value = res.data.data.registration_enabled
  } catch { /* ignore */ }
})

async function handleLogin() {
  if (!email.value || !password.value) {
    notification.error('Please fill in all fields.')
    return
  }
  if (requires2FA.value && !twoFactorCode.value) {
    notification.error('Please enter your 2FA code.')
    return
  }
  loading.value = true
  try {
    await auth.login(email.value, password.value, requires2FA.value ? twoFactorCode.value : undefined)
    router.push('/')
  } catch (err: any) {
    if (err?.requires2FA) {
      requires2FA.value = true
      notification.info('Please enter your two-factor authentication code.')
      return
    }
    // Check for 2FA required response (401 with requires_2fa flag)
    if (err?.response?.status === 401 && err?.response?.data?.data?.requires_2fa) {
      requires2FA.value = true
      notification.info('Please enter your two-factor authentication code.')
      return
    }
    if (err?.response?.status === 429) {
      notification.error('Too many login attempts. Please try again later.')
      return
    }
    const message = err?.response?.data?.error?.message || err?.response?.data?.error || err?.message || 'Login failed.'
    notification.error(message)
  } finally {
    loading.value = false
  }
}

function resetLogin() {
  requires2FA.value = false
  twoFactorCode.value = ''
}
</script>

<template>
  <div class="auth-page">
    <div class="auth-card">
      <div class="auth-header">
        <div class="auth-logo">
          <img src="/logo.png" alt="Posta" class="logo-img" />
          <span>Posta</span>
        </div>
        <p class="auth-subtitle">{{ requires2FA ? 'Two-Factor Authentication' : 'Sign in to your account' }}</p>
      </div>

      <form class="auth-form" @submit.prevent="handleLogin">
        <template v-if="!requires2FA">
          <div class="form-group">
            <label class="form-label" for="email">Email</label>
            <input id="email" v-model="email" type="email" class="form-input" placeholder="you@example.com" autocomplete="email" />
          </div>
          <div class="form-group">
            <label class="form-label" for="password">Password</label>
            <input id="password" v-model="password" type="password" class="form-input" placeholder="Enter your password" autocomplete="current-password" />
          </div>
        </template>
        <template v-else>
          <div class="form-group">
            <label class="form-label" for="2fa-code">Authentication Code</label>
            <input
              id="2fa-code"
              v-model="twoFactorCode"
              type="text"
              class="form-input totp-input"
              placeholder="000000"
              maxlength="6"
              inputmode="numeric"
              autocomplete="one-time-code"
              autofocus
            />
            <small class="form-hint">Enter the 6-digit code from your authenticator app</small>
          </div>
        </template>
        <button type="submit" class="btn btn-primary auth-btn" :disabled="loading">
          <span v-if="loading" class="spinner"></span>
          {{ loading ? 'Signing in...' : 'Sign in' }}
        </button>
        <button v-if="requires2FA" type="button" class="btn btn-secondary auth-btn" style="margin-top: 8px" @click="resetLogin">
          Back to Login
        </button>
      </form>

      <div class="auth-footer">
        <template v-if="registrationEnabled">
          <span>Don't have an account?</span>
          <router-link to="/register">Sign up</router-link>
        </template>
        <span v-else>Contact your administrator for an account.</span>
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
  width: 100px;
  height: 100px;
  object-fit: contain;
}
.auth-logo span {
  font-size: 32px;
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

.totp-input {
  font-size: 24px;
  text-align: center;
  letter-spacing: 8px;
  font-family: 'JetBrains Mono', 'Fira Code', monospace;
}

.form-hint {
  font-size: 12px;
  color: var(--text-muted);
  margin-top: 6px;
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
