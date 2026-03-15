<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { analyticsApi } from '../../api/analytics'
import type { DailyCount, StatusBreakdown } from '../../api/types'

const loading = ref(true)
const dailyCounts = ref<DailyCount[]>([])
const statusBreakdown = ref<StatusBreakdown[]>([])

const fromDate = ref('')
const toDate = ref('')
const statusFilter = ref('')

// Set default range: last 30 days
const now = new Date()
const thirtyDaysAgo = new Date(now)
thirtyDaysAgo.setDate(thirtyDaysAgo.getDate() - 30)
fromDate.value = thirtyDaysAgo.toISOString().slice(0, 10)
toDate.value = now.toISOString().slice(0, 10)

async function loadAnalytics() {
  loading.value = true
  try {
    const res = await analyticsApi.user(fromDate.value, toDate.value, statusFilter.value || undefined)
    dailyCounts.value = res.data.data.daily_counts || []
    statusBreakdown.value = res.data.data.status_breakdown || []
  } catch (e) {
    console.error('Failed to load analytics', e)
  } finally {
    loading.value = false
  }
}

onMounted(loadAnalytics)

const totalEmails = computed(() => dailyCounts.value.reduce((sum, d) => sum + d.count, 0))
const maxCount = computed(() => Math.max(...dailyCounts.value.map(d => d.count), 1))

function barHeight(count: number): string {
  return `${Math.max((count / maxCount.value) * 100, 2)}%`
}

function statusColor(status: string): string {
  switch (status) {
    case 'sent': return 'var(--success-500, #22c55e)'
    case 'failed': return 'var(--danger-500, #ef4444)'
    case 'pending': return 'var(--warning-500, #f59e0b)'
    case 'queued': return 'var(--info-500, #3b82f6)'
    case 'processing': return 'var(--warning-400, #fbbf24)'
    case 'suppressed': return 'var(--text-muted, #9ca3af)'
    case 'scheduled': return 'var(--primary-500, #6366f1)'
    default: return 'var(--text-muted)'
  }
}

function formatDate(dateStr: string): string {
  const d = new Date(dateStr)
  return d.toLocaleDateString(undefined, { month: 'short', day: 'numeric' })
}

const totalBreakdown = computed(() => statusBreakdown.value.reduce((sum, s) => sum + s.count, 0))

function breakdownPercent(count: number): string {
  if (totalBreakdown.value === 0) return '0'
  return ((count / totalBreakdown.value) * 100).toFixed(1)
}
</script>

<template>
  <div>
    <div class="page-header">
      <h1>Analytics</h1>
    </div>

    <!-- Filters -->
    <div class="card" style="margin-bottom: 24px">
      <div class="card-body">
        <div class="filters">
          <div class="form-group">
            <label class="form-label">From</label>
            <input v-model="fromDate" type="date" class="form-input" />
          </div>
          <div class="form-group">
            <label class="form-label">To</label>
            <input v-model="toDate" type="date" class="form-input" />
          </div>
          <div class="form-group">
            <label class="form-label">Status</label>
            <select v-model="statusFilter" class="form-input">
              <option value="">All</option>
              <option value="sent">Sent</option>
              <option value="failed">Failed</option>
              <option value="pending">Pending</option>
              <option value="queued">Queued</option>
              <option value="suppressed">Suppressed</option>
              <option value="scheduled">Scheduled</option>
            </select>
          </div>
          <div class="form-group" style="align-self: flex-end">
            <button class="btn btn-primary" @click="loadAnalytics" :disabled="loading">
              {{ loading ? 'Loading...' : 'Apply' }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <div v-if="loading" class="loading-page">
      <div class="spinner"></div>
    </div>

    <template v-else>
      <!-- Summary -->
      <div class="stats-grid" style="margin-bottom: 24px">
        <div class="stat-card">
          <div class="stat-header">
            <div class="stat-label">Total in Period</div>
          </div>
          <div class="stat-value">{{ totalEmails }}</div>
        </div>
        <div v-for="s in statusBreakdown" :key="s.status" class="stat-card">
          <div class="stat-header">
            <div class="stat-label" style="text-transform: capitalize">{{ s.status }}</div>
          </div>
          <div class="stat-value">{{ s.count }}</div>
          <div class="stat-sub">{{ breakdownPercent(s.count) }}%</div>
        </div>
      </div>

      <!-- Daily Chart -->
      <div class="card" style="margin-bottom: 24px">
        <div class="card-header"><h2>Daily Email Volume</h2></div>
        <div class="card-body">
          <div v-if="dailyCounts.length === 0" class="empty-state">
            <h3>No data</h3>
            <p>No emails found in the selected date range.</p>
          </div>
          <div v-else class="chart-container">
            <div class="chart">
              <div
                v-for="day in dailyCounts"
                :key="day.date"
                class="chart-bar-group"
                :title="`${formatDate(day.date)}: ${day.count} emails`"
              >
                <div class="chart-bar" :style="{ height: barHeight(day.count) }">
                  <span v-if="day.count > 0" class="chart-bar-label">{{ day.count }}</span>
                </div>
                <div class="chart-bar-date">{{ formatDate(day.date) }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Status Breakdown -->
      <div class="card">
        <div class="card-header"><h2>Status Breakdown</h2></div>
        <div class="card-body">
          <div v-if="statusBreakdown.length === 0" class="empty-state">
            <h3>No data</h3>
            <p>No emails found in the selected date range.</p>
          </div>
          <div v-else class="breakdown">
            <div v-for="s in statusBreakdown" :key="s.status" class="breakdown-row">
              <div class="breakdown-label">
                <span class="breakdown-dot" :style="{ background: statusColor(s.status) }"></span>
                <span style="text-transform: capitalize">{{ s.status }}</span>
              </div>
              <div class="breakdown-bar-track">
                <div
                  class="breakdown-bar-fill"
                  :style="{ width: breakdownPercent(s.count) + '%', background: statusColor(s.status) }"
                ></div>
              </div>
              <div class="breakdown-value">{{ s.count }} <span class="breakdown-pct">({{ breakdownPercent(s.count) }}%)</span></div>
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<style scoped>
.filters {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
  align-items: flex-end;
}

.filters .form-group {
  min-width: 140px;
}

.stat-sub {
  font-size: 12px;
  color: var(--text-muted);
  margin-top: 2px;
}

.chart-container {
  overflow-x: auto;
  padding-bottom: 8px;
}

.chart {
  display: flex;
  align-items: flex-end;
  gap: 4px;
  height: 220px;
  min-width: fit-content;
}

.chart-bar-group {
  display: flex;
  flex-direction: column;
  align-items: center;
  flex: 1;
  min-width: 28px;
  max-width: 48px;
  height: 100%;
  justify-content: flex-end;
}

.chart-bar {
  width: 100%;
  background: var(--primary-500, #6366f1);
  border-radius: 4px 4px 0 0;
  min-height: 2px;
  position: relative;
  transition: height 0.3s ease;
}

.chart-bar-label {
  position: absolute;
  top: -20px;
  left: 50%;
  transform: translateX(-50%);
  font-size: 10px;
  color: var(--text-muted);
  white-space: nowrap;
}

.chart-bar-date {
  font-size: 10px;
  color: var(--text-muted);
  margin-top: 6px;
  white-space: nowrap;
  transform: rotate(-45deg);
  transform-origin: top center;
}

.breakdown {
  display: grid;
  gap: 12px;
}

.breakdown-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.breakdown-label {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 110px;
  font-size: 13px;
  font-weight: 500;
  color: var(--text-primary);
}

.breakdown-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  flex-shrink: 0;
}

.breakdown-bar-track {
  flex: 1;
  height: 8px;
  background: var(--bg-secondary);
  border-radius: 4px;
  overflow: hidden;
}

.breakdown-bar-fill {
  height: 100%;
  border-radius: 4px;
  transition: width 0.3s ease;
}

.breakdown-value {
  min-width: 100px;
  text-align: right;
  font-size: 13px;
  font-weight: 600;
  color: var(--text-primary);
}

.breakdown-pct {
  font-weight: 400;
  color: var(--text-muted);
}
</style>
