<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { adminApi } from '../../api/admin'
import type { CronJob } from '../../api/types'

const loading = ref(true)
const jobs = ref<CronJob[]>([])
let refreshTimer: ReturnType<typeof setInterval> | null = null

onMounted(() => {
  loadJobs()
  refreshTimer = setInterval(loadJobs, 30000)
})

onBeforeUnmount(() => {
  if (refreshTimer) {
    clearInterval(refreshTimer)
  }
})

async function loadJobs() {
  try {
    const res = await adminApi.listJobs()
    jobs.value = res.data.data
  } catch (e) {
    console.error('Failed to load jobs', e)
  } finally {
    loading.value = false
  }
}

function formatDate(date: string | null) {
  if (!date) return '-'
  return new Date(date).toLocaleString()
}

function timeUntil(date: string | null) {
  if (!date) return '-'
  const diff = new Date(date).getTime() - Date.now()
  if (diff <= 0) return 'now'
  const minutes = Math.floor(diff / 60000)
  const hours = Math.floor(minutes / 60)
  if (hours > 0) return `${hours}h ${minutes % 60}m`
  return `${minutes}m`
}

function scheduleName(schedule: string) {
  const known: Record<string, string> = {
    '0 3 * * *': 'Daily at 03:00 UTC',
    '0 7 * * *': 'Daily at 07:00 UTC',
  }
  return known[schedule] || schedule
}
</script>

<template>
  <div>
    <div class="page-header">
      <h1>Scheduled Jobs</h1>
      <button class="btn btn-sm btn-secondary" @click="loadJobs" :disabled="loading">Refresh</button>
    </div>

    <div v-if="loading" class="loading-page">
      <div class="spinner"></div>
    </div>

    <template v-else>
      <div v-if="jobs.length === 0" class="empty-state">
        <h3>No scheduled jobs</h3>
        <p>No cron jobs are currently registered. Jobs are only available when the server is running in production mode.</p>
      </div>

      <div v-else class="jobs-grid">
        <div v-for="job in jobs" :key="job.name" class="card job-card">
          <div class="card-body">
            <div class="job-header">
              <h3 class="job-name">{{ job.name }}</h3>
              <span v-if="job.running" class="badge badge-running">Running</span>
              <span v-else-if="job.last_error" class="badge badge-error">Failed</span>
              <span v-else-if="job.last_run_at" class="badge badge-success">OK</span>
              <span v-else class="badge badge-neutral">Pending</span>
            </div>

            <div class="job-schedule">
              <code>{{ job.schedule }}</code>
              <span class="schedule-label">{{ scheduleName(job.schedule) }}</span>
            </div>

            <div class="job-details">
              <div class="job-detail">
                <span class="detail-label">Last run</span>
                <span class="detail-value">{{ formatDate(job.last_run_at) }}</span>
              </div>
              <div class="job-detail">
                <span class="detail-label">Next run</span>
                <span class="detail-value">
                  {{ formatDate(job.next_run_at) }}
                  <span v-if="job.next_run_at" class="time-until">(in {{ timeUntil(job.next_run_at) }})</span>
                </span>
              </div>
            </div>

            <div v-if="job.last_error" class="job-error">
              <span class="error-label">Last error</span>
              <code class="error-message">{{ job.last_error }}</code>
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<style scoped>
.jobs-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
  gap: 1rem;
}

.job-card {
  border: 1px solid var(--border-primary);
}

.job-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 0.75rem;
}

.job-name {
  font-size: 16px;
  font-weight: 600;
  margin: 0;
}

.badge-running {
  background: var(--info-50, #eff6ff);
  color: var(--info-700, #1d4ed8);
}

.badge-error {
  background: var(--danger-50, #fef2f2);
  color: var(--danger-600, #dc2626);
}

.badge-success {
  background: var(--success-50, #f0fdf4);
  color: var(--success-700, #15803d);
}

.badge-neutral {
  background: var(--bg-secondary);
  color: var(--text-muted);
}

.job-schedule {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-bottom: 1rem;
  padding-bottom: 0.75rem;
  border-bottom: 1px solid var(--border-primary);
}

.job-schedule code {
  font-size: 12px;
  padding: 2px 8px;
  background: var(--bg-secondary);
  border-radius: var(--radius-sm, 4px);
}

.schedule-label {
  font-size: 13px;
  color: var(--text-muted);
}

.job-details {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.job-detail {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 13px;
}

.detail-label {
  color: var(--text-muted);
  font-weight: 500;
}

.detail-value {
  color: var(--text-primary);
}

.time-until {
  color: var(--text-muted);
  font-size: 12px;
  margin-left: 4px;
}

.job-error {
  margin-top: 0.75rem;
  padding-top: 0.75rem;
  border-top: 1px solid var(--border-primary);
}

.error-label {
  display: block;
  font-size: 12px;
  font-weight: 600;
  color: var(--danger-600, #dc2626);
  margin-bottom: 4px;
}

.error-message {
  display: block;
  font-size: 12px;
  padding: 6px 10px;
  background: var(--danger-50, #fef2f2);
  color: var(--danger-600, #dc2626);
  border-radius: var(--radius-sm, 4px);
  word-break: break-word;
}
</style>
