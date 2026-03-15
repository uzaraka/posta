<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { adminApi } from '../../api/admin'
import { analyticsApi } from '../../api/analytics'
import { useAuthStore } from '../../stores/auth'
import type { AdminMetrics, WorkerStatus, AnalyticsResponse } from '../../api/types'

const auth = useAuthStore()
const loading = ref(true)
const metrics = ref<AdminMetrics | null>(null)
const workerStatus = ref<WorkerStatus | null>(null)
const analytics = ref<AnalyticsResponse | null>(null)
let workerSSE: EventSource | null = null

onMounted(async () => {
  try {
    const [metricsRes, analyticsRes] = await Promise.all([
      adminApi.getMetrics(),
      analyticsApi.admin(),
    ])
    metrics.value = metricsRes.data.data
    analytics.value = analyticsRes.data.data
  } catch (e) {
    console.error('Failed to load metrics', e)
  } finally {
    loading.value = false
  }
  startWorkerStream()
})

onBeforeUnmount(() => {
  stopWorkerStream()
})

function startWorkerStream() {
  const baseUrl = import.meta.env.VITE_API_URL || '/api/v1'
  const token = auth.token
  if (!token) return
  const url = `${baseUrl}/admin/workers/stream?token=${encodeURIComponent(token)}`
  workerSSE = new EventSource(url)

  workerSSE.addEventListener('worker.status', (e) => {
    try {
      const status: WorkerStatus = JSON.parse((e as MessageEvent).data)
      workerStatus.value = status
      if (metrics.value) {
        metrics.value.active_workers = status.active_workers
      }
    } catch {
      // ignore parse errors
    }
  })

  workerSSE.onerror = () => {
    // Connection lost; will auto-reconnect via EventSource
  }
}

function stopWorkerStream() {
  if (workerSSE) {
    workerSSE.close()
    workerSSE = null
  }
}

const activeWorkers = computed(() => workerStatus.value?.active_workers ?? metrics.value?.active_workers ?? 0)

const deliveryRate = computed(() => {
  if (!metrics.value || metrics.value.total_emails === 0) return 0
  return ((metrics.value.sent_emails / metrics.value.total_emails) * 100)
})

const revokedKeys = computed(() => {
  if (!metrics.value) return 0
  return metrics.value.total_api_keys - metrics.value.active_api_keys
})

function failureColor(rate: number): string {
  if (rate <= 1) return 'var(--success-600, #16a34a)'
  if (rate <= 5) return 'var(--warning-600, #ca8a04)'
  return 'var(--danger-600, #dc2626)'
}

// Analytics helpers
const maxDailyCount = computed(() => {
  if (!analytics.value?.daily_counts?.length) return 1
  return Math.max(...analytics.value.daily_counts.map(d => d.count), 1)
})

function barHeight(count: number): string {
  return `${Math.max((count / maxDailyCount.value) * 100, 2)}%`
}

function formatShortDate(dateStr: string): string {
  const d = new Date(dateStr)
  return d.toLocaleDateString(undefined, { month: 'short', day: 'numeric' })
}

function statusColor(status: string): string {
  switch (status) {
    case 'sent': return 'var(--success-500, #22c55e)'
    case 'failed': return 'var(--danger-500, #ef4444)'
    case 'pending': return 'var(--warning-500, #f59e0b)'
    case 'queued': return 'var(--info-500, #3b82f6)'
    case 'suppressed': return 'var(--text-muted, #9ca3af)'
    default: return 'var(--text-muted)'
  }
}

const totalBreakdown = computed(() => analytics.value?.status_breakdown?.reduce((s, b) => s + b.count, 0) || 0)
</script>

<template>
  <div>
    <div class="page-header">
      <h1>Platform Metrics</h1>
    </div>

    <div v-if="loading" class="loading-page">
      <div class="spinner"></div>
    </div>

    <template v-else-if="metrics">
      <!-- Overview -->
      <div class="metrics-section-label">Overview</div>
      <div class="stats-grid">
        <div class="stat-card">
          <div class="stat-header">
            <div class="stat-label">Total Users</div>
            <div class="stat-icon stat-icon-primary">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg>
            </div>
          </div>
          <div class="stat-value">{{ metrics.total_users }}</div>
        </div>
        <div class="stat-card">
          <div class="stat-header">
            <div class="stat-label">Total Emails</div>
            <div class="stat-icon stat-icon-primary">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"/><polyline points="22,6 12,13 2,6"/></svg>
            </div>
          </div>
          <div class="stat-value">{{ metrics.total_emails }}</div>
        </div>
        <div class="stat-card">
          <div class="stat-header">
            <div class="stat-label">Active Workers</div>
            <div class="stat-icon" :class="activeWorkers > 0 ? 'stat-icon-success' : 'stat-icon-danger'">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="2" width="20" height="8" rx="2"/><rect x="2" y="14" width="20" height="8" rx="2"/><circle cx="6" cy="6" r="1" fill="currentColor"/><circle cx="6" cy="18" r="1" fill="currentColor"/></svg>
            </div>
          </div>
          <div class="stat-value">{{ activeWorkers }}</div>
          <div class="stat-sub">{{ activeWorkers > 0 ? 'Processing emails' : 'No workers running' }}</div>
        </div>
        <div class="stat-card">
          <div class="stat-header">
            <div class="stat-label">Shared SMTP Servers</div>
            <div class="stat-icon stat-icon-primary">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="2" width="20" height="8" rx="2"/><rect x="2" y="14" width="20" height="8" rx="2"/><line x1="6" y1="6" x2="6.01" y2="6"/><line x1="6" y1="18" x2="6.01" y2="18"/></svg>
            </div>
          </div>
          <div class="stat-value">{{ metrics.shared_smtp_servers }}</div>
        </div>
        <div class="stat-card">
          <div class="stat-header">
            <div class="stat-label">Total Domains</div>
            <div class="stat-icon stat-icon-primary">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><line x1="2" y1="12" x2="22" y2="12"/><path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"/></svg>
            </div>
          </div>
          <div class="stat-value">{{ metrics.total_domains }}</div>
        </div>
      </div>

      <!-- Worker Details -->
      <template v-if="workerStatus && workerStatus.workers.length > 0">
        <div class="metrics-section-label">Workers</div>
        <div class="card" style="margin-bottom: 28px;">
          <div class="table-wrapper">
            <table>
              <thead>
                <tr>
                  <th>Host</th>
                  <th>PID</th>
                  <th>Type</th>
                  <th>Queues</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="w in workerStatus.workers" :key="`${w.host}:${w.pid}`">
                  <td>{{ w.host }}</td>
                  <td>{{ w.pid }}</td>
                  <td>
                    <span class="badge" :class="w.type === 'embedded' ? 'badge-info' : 'badge-success'">
                      {{ w.type }}
                    </span>
                  </td>
                  <td>
                    <span v-for="(concurrency, queue) in w.queues" :key="queue" class="badge badge-neutral" style="margin-right: 6px;">
                      {{ queue }}: {{ concurrency }}
                    </span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </template>

      <!-- Email Delivery -->
      <div class="metrics-section-label">Email Delivery</div>
      <div class="stats-grid">
        <div class="stat-card">
          <div class="stat-header">
            <div class="stat-label">Sent</div>
            <div class="stat-icon stat-icon-success">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"/></svg>
            </div>
          </div>
          <div class="stat-value">{{ metrics.sent_emails }}</div>
          <div class="stat-sub">{{ deliveryRate.toFixed(1) }}% delivery rate</div>
        </div>
        <div v-if="metrics.queued_emails > 0 || metrics.processing_emails > 0" class="stat-card">
          <div class="stat-header">
            <div class="stat-label">In Queue</div>
            <div class="stat-icon stat-icon-info">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
            </div>
          </div>
          <div class="stat-value">{{ metrics.queued_emails + metrics.processing_emails }}</div>
          <div class="stat-sub">{{ metrics.queued_emails }} queued, {{ metrics.processing_emails }} processing</div>
        </div>
        <div class="stat-card">
          <div class="stat-header">
            <div class="stat-label">Failed</div>
            <div class="stat-icon stat-icon-danger">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><line x1="15" y1="9" x2="9" y2="15"/><line x1="9" y1="9" x2="15" y2="15"/></svg>
            </div>
          </div>
          <div class="stat-value">{{ metrics.failed_emails }}</div>
        </div>
        <div class="stat-card">
          <div class="stat-header">
            <div class="stat-label">Suppressed</div>
            <div class="stat-icon stat-icon-secondary">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><line x1="4.93" y1="4.93" x2="19.07" y2="19.07"/></svg>
            </div>
          </div>
          <div class="stat-value">{{ metrics.suppressed_emails }}</div>
        </div>
        <div class="stat-card">
          <div class="stat-header">
            <div class="stat-label">Failure Rate</div>
            <div class="stat-icon stat-icon-warning">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>
            </div>
          </div>
          <div class="stat-value" :style="{ color: failureColor(metrics.failure_rate) }">{{ metrics.failure_rate.toFixed(1) }}%</div>
        </div>
      </div>

      <!-- Email delivery progress bar -->
      <div class="card" style="margin-bottom: 28px;">
        <div class="card-body">
          <div class="delivery-bar-label">
            <span>Email Status Distribution</span>
            <span class="delivery-bar-total">{{ metrics.total_emails }} total</span>
          </div>
          <div class="delivery-bar" v-if="metrics.total_emails > 0">
            <div class="delivery-segment delivery-sent" :style="{ width: (metrics.sent_emails / metrics.total_emails * 100) + '%' }" :title="`Sent: ${metrics.sent_emails}`"></div>
            <div class="delivery-segment delivery-queued" :style="{ width: ((metrics.queued_emails + metrics.processing_emails) / metrics.total_emails * 100) + '%' }" :title="`In queue: ${metrics.queued_emails + metrics.processing_emails}`"></div>
            <div class="delivery-segment delivery-failed" :style="{ width: (metrics.failed_emails / metrics.total_emails * 100) + '%' }" :title="`Failed: ${metrics.failed_emails}`"></div>
            <div class="delivery-segment delivery-suppressed" :style="{ width: (metrics.suppressed_emails / metrics.total_emails * 100) + '%' }" :title="`Suppressed: ${metrics.suppressed_emails}`"></div>
          </div>
          <div v-else class="delivery-bar">
            <div class="delivery-segment delivery-empty" style="width: 100%"></div>
          </div>
          <div class="delivery-legend">
            <span class="legend-item"><span class="legend-dot delivery-sent"></span> Sent</span>
            <span class="legend-item"><span class="legend-dot delivery-queued"></span> Queued</span>
            <span class="legend-item"><span class="legend-dot delivery-failed"></span> Failed</span>
            <span class="legend-item"><span class="legend-dot delivery-suppressed"></span> Suppressed</span>
          </div>
        </div>
      </div>

      <!-- Analytics (last 30 days) -->
      <template v-if="analytics">
        <div class="metrics-section-label">Email Volume (Last 30 Days)</div>
        <div class="card" style="margin-bottom: 28px;">
          <div class="card-body">
            <div v-if="analytics.daily_counts && analytics.daily_counts.length > 0" class="admin-chart">
              <div
                v-for="day in analytics.daily_counts"
                :key="day.date"
                class="admin-chart-bar-group"
                :title="`${formatShortDate(day.date)}: ${day.count}`"
              >
                <div class="admin-chart-bar" :style="{ height: barHeight(day.count) }">
                  <span v-if="day.count > 0" class="admin-chart-bar-label">{{ day.count }}</span>
                </div>
                <div class="admin-chart-bar-date">{{ formatShortDate(day.date) }}</div>
              </div>
            </div>
            <div v-else style="text-align: center; color: var(--text-muted); padding: 24px;">No data</div>
          </div>
        </div>

        <div v-if="analytics.status_breakdown && analytics.status_breakdown.length > 0" class="card" style="margin-bottom: 28px;">
          <div class="card-header"><h2>Status Breakdown</h2></div>
          <div class="card-body">
            <div class="admin-breakdown">
              <div v-for="s in analytics.status_breakdown" :key="s.status" class="admin-breakdown-row">
                <div class="admin-breakdown-label">
                  <span class="admin-breakdown-dot" :style="{ background: statusColor(s.status) }"></span>
                  <span style="text-transform: capitalize">{{ s.status }}</span>
                </div>
                <div class="admin-breakdown-bar-track">
                  <div
                    class="admin-breakdown-bar-fill"
                    :style="{ width: (totalBreakdown > 0 ? (s.count / totalBreakdown * 100) : 0) + '%', background: statusColor(s.status) }"
                  ></div>
                </div>
                <div class="admin-breakdown-value">{{ s.count }}</div>
              </div>
            </div>
          </div>
        </div>
      </template>

      <!-- Webhook Deliveries -->
      <template v-if="metrics.webhook_deliveries && metrics.webhook_deliveries.total_deliveries > 0">
        <div class="metrics-section-label">Webhook Deliveries</div>
        <div class="stats-grid">
          <div class="stat-card">
            <div class="stat-header">
              <div class="stat-label">Total Deliveries</div>
              <div class="stat-icon stat-icon-primary">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M18 16.98h-5.99c-1.1 0-1.95.68-2.95 1.76C8.07 19.82 6.22 20 5 20c-1.22 0-2.2-.38-3-1"/><path d="M18 16.98h-5.99c-1.66 0-2.61-1.22-3.15-2.59C8.23 12.64 8 10.66 8 9c0-3.87 3.13-7 7-7"/><circle cx="12" cy="12" r="2"/></svg>
              </div>
            </div>
            <div class="stat-value">{{ metrics.webhook_deliveries.total_deliveries }}</div>
          </div>
          <div class="stat-card">
            <div class="stat-header">
              <div class="stat-label">Successful</div>
              <div class="stat-icon stat-icon-success">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"/></svg>
              </div>
            </div>
            <div class="stat-value">{{ metrics.webhook_deliveries.success_deliveries }}</div>
          </div>
          <div class="stat-card">
            <div class="stat-header">
              <div class="stat-label">Failed</div>
              <div class="stat-icon stat-icon-danger">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><line x1="15" y1="9" x2="9" y2="15"/><line x1="9" y1="9" x2="15" y2="15"/></svg>
              </div>
            </div>
            <div class="stat-value">{{ metrics.webhook_deliveries.failed_deliveries }}</div>
          </div>
          <div class="stat-card">
            <div class="stat-header">
              <div class="stat-label">Success Rate</div>
              <div class="stat-icon stat-icon-info">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 12h-4l-3 9L9 3l-3 9H2"/></svg>
              </div>
            </div>
            <div class="stat-value" :style="{ color: failureColor(100 - metrics.webhook_deliveries.success_rate) }">
              {{ metrics.webhook_deliveries.success_rate.toFixed(1) }}%
            </div>
          </div>
        </div>
        <div class="card" style="margin-bottom: 28px;">
          <div class="card-body">
            <div class="delivery-bar-label">
              <span>Webhook Delivery Distribution</span>
              <span class="delivery-bar-total">{{ metrics.webhook_deliveries.total_deliveries }} total</span>
            </div>
            <div class="delivery-bar">
              <div class="delivery-segment delivery-sent" :style="{ width: metrics.webhook_deliveries.success_rate + '%' }" :title="`Success: ${metrics.webhook_deliveries.success_deliveries}`"></div>
              <div class="delivery-segment delivery-failed" :style="{ width: (100 - metrics.webhook_deliveries.success_rate) + '%' }" :title="`Failed: ${metrics.webhook_deliveries.failed_deliveries}`"></div>
            </div>
            <div class="delivery-legend">
              <span class="legend-item"><span class="legend-dot delivery-sent"></span> Success</span>
              <span class="legend-item"><span class="legend-dot delivery-failed"></span> Failed</span>
            </div>
          </div>
        </div>
      </template>

      <!-- API Keys & Reputation -->
      <div class="metrics-section-label">API Keys & Reputation</div>
      <div class="stats-grid">
        <div class="stat-card">
          <div class="stat-header">
            <div class="stat-label">Active API Keys</div>
            <div class="stat-icon stat-icon-success">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 2l-2 2m-7.61 7.61a5.5 5.5 0 1 1-7.78 7.78 5.5 5.5 0 0 1 7.78-7.78zm0 0L15.5 7.5m0 0l3 3L22 7l-3-3m-3.5 3.5L19 4"/></svg>
            </div>
          </div>
          <div class="stat-value">{{ metrics.active_api_keys }}</div>
          <div class="stat-sub">{{ revokedKeys }} revoked</div>
        </div>
        <div class="stat-card">
          <div class="stat-header">
            <div class="stat-label">Total Bounces</div>
            <div class="stat-icon stat-icon-warning">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>
            </div>
          </div>
          <div class="stat-value">{{ metrics.total_bounces }}</div>
        </div>
        <div class="stat-card">
          <div class="stat-header">
            <div class="stat-label">Suppressions</div>
            <div class="stat-icon stat-icon-secondary">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><line x1="4.93" y1="4.93" x2="19.07" y2="19.07"/></svg>
            </div>
          </div>
          <div class="stat-value">{{ metrics.total_suppressions }}</div>
        </div>
      </div>
    </template>
  </div>
</template>

<style scoped>
.metrics-section-label {
  font-size: 12px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: var(--text-muted);
  margin-bottom: 12px;
}

.delivery-bar-label {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
  font-size: 13px;
  font-weight: 600;
  color: var(--text-primary);
}

.delivery-bar-total {
  font-weight: 500;
  color: var(--text-muted);
}

.delivery-bar {
  display: flex;
  height: 12px;
  border-radius: 6px;
  overflow: hidden;
  background: var(--bg-tertiary);
  margin-bottom: 12px;
}

.delivery-segment {
  transition: width 0.6s ease;
  min-width: 0;
}

.delivery-sent { background: var(--success-500, #22c55e); }
.delivery-queued { background: var(--primary-400, #60a5fa); }
.delivery-failed { background: var(--danger-500, #ef4444); }
.delivery-suppressed { background: var(--text-muted); }
.delivery-empty { background: var(--bg-tertiary); }

.delivery-legend {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: var(--text-secondary);
}

.legend-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  flex-shrink: 0;
}

.admin-chart {
  display: flex;
  align-items: flex-end;
  gap: 4px;
  height: 180px;
  min-width: fit-content;
  overflow-x: auto;
}

.admin-chart-bar-group {
  display: flex;
  flex-direction: column;
  align-items: center;
  flex: 1;
  min-width: 24px;
  max-width: 40px;
  height: 100%;
  justify-content: flex-end;
}

.admin-chart-bar {
  width: 100%;
  background: var(--primary-500, #6366f1);
  border-radius: 3px 3px 0 0;
  min-height: 2px;
  position: relative;
  transition: height 0.3s ease;
}

.admin-chart-bar-label {
  position: absolute;
  top: -18px;
  left: 50%;
  transform: translateX(-50%);
  font-size: 10px;
  color: var(--text-muted);
  white-space: nowrap;
}

.admin-chart-bar-date {
  font-size: 9px;
  color: var(--text-muted);
  margin-top: 4px;
  white-space: nowrap;
  transform: rotate(-45deg);
  transform-origin: top center;
}

.admin-breakdown {
  display: grid;
  gap: 10px;
}

.admin-breakdown-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.admin-breakdown-label {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 100px;
  font-size: 13px;
  font-weight: 500;
}

.admin-breakdown-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  flex-shrink: 0;
}

.admin-breakdown-bar-track {
  flex: 1;
  height: 8px;
  background: var(--bg-secondary);
  border-radius: 4px;
  overflow: hidden;
}

.admin-breakdown-bar-fill {
  height: 100%;
  border-radius: 4px;
  transition: width 0.3s ease;
}

.admin-breakdown-value {
  min-width: 60px;
  text-align: right;
  font-size: 13px;
  font-weight: 600;
}
</style>
