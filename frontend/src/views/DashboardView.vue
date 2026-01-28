<script setup lang="ts">
import { ref, onMounted } from 'vue'

interface Stats {
  claims: { total: number; pending: number; approved: number; rejected: number }
  packages: { total: number; pending: number; paid: number }
  amounts: { requested: number; approved: number; deduction: number }
  centers: { total: number; active: number }
}

const stats = ref<Stats>({
  claims: { total: 0, pending: 0, approved: 0, rejected: 0 },
  packages: { total: 0, pending: 0, paid: 0 },
  amounts: { requested: 0, approved: 0, deduction: 0 },
  centers: { total: 0, active: 0 }
})

const isLoading = ref(true)

function formatNumber(num: number): string {
  return new Intl.NumberFormat('fa-IR').format(num)
}

function formatCurrency(num: number): string {
  return new Intl.NumberFormat('fa-IR').format(num / 10) + ' تومان'
}

onMounted(async () => {
  try {
    // For now, use mock data
    stats.value = {
      claims: { total: 15420, pending: 320, approved: 14800, rejected: 300 },
      packages: { total: 1250, pending: 45, paid: 1180 },
      amounts: { requested: 125000000000, approved: 118500000000, deduction: 6500000000 },
      centers: { total: 850, active: 720 }
    }
  } catch (error) {
    console.error('Failed to load stats:', error)
  } finally {
    isLoading.value = false
  }
})
</script>

<template>
  <div class="dashboard">
    <!-- Page Header -->
    <div class="page-header">
      <h1><i class="bi bi-speedometer2"></i> داشبورد</h1>
    </div>

    <!-- Loading -->
    <div v-if="isLoading" class="loading-spinner">
      <div class="spinner"></div>
    </div>

    <template v-else>
      <!-- Stats Cards -->
      <div class="row mb-4">
        <div class="col-3">
          <div class="stat-card primary">
            <i class="bi bi-file-earmark-medical stat-icon"></i>
            <div class="stat-value">{{ formatNumber(stats.claims.total) }}</div>
            <div class="stat-label">کل ادعاها</div>
          </div>
        </div>
        <div class="col-3">
          <div class="stat-card success">
            <i class="bi bi-box-seam stat-icon"></i>
            <div class="stat-value">{{ formatNumber(stats.packages.total) }}</div>
            <div class="stat-label">بسته‌های ارسالی</div>
          </div>
        </div>
        <div class="col-3">
          <div class="stat-card warning">
            <i class="bi bi-hospital stat-icon"></i>
            <div class="stat-value">{{ formatNumber(stats.centers.active) }}</div>
            <div class="stat-label">مرکز فعال</div>
          </div>
        </div>
        <div class="col-3">
          <div class="stat-card danger">
            <i class="bi bi-hourglass-split stat-icon"></i>
            <div class="stat-value">{{ formatNumber(stats.claims.pending) }}</div>
            <div class="stat-label">در انتظار ارزیابی</div>
          </div>
        </div>
      </div>

      <!-- Financial Summary -->
      <div class="row mb-4">
        <div class="col-12">
          <div class="card">
            <div class="card-header">
              <i class="bi bi-cash-stack"></i>
              خلاصه مالی
            </div>
            <div class="card-body">
              <div class="row">
                <div class="col-4">
                  <div class="financial-item">
                    <div class="financial-label">مبلغ درخواستی</div>
                    <div class="financial-value text-primary">{{ formatCurrency(stats.amounts.requested) }}</div>
                  </div>
                </div>
                <div class="col-4">
                  <div class="financial-item">
                    <div class="financial-label">مبلغ تایید شده</div>
                    <div class="financial-value text-success">{{ formatCurrency(stats.amounts.approved) }}</div>
                  </div>
                </div>
                <div class="col-4">
                  <div class="financial-item">
                    <div class="financial-label">کسورات</div>
                    <div class="financial-value text-danger">{{ formatCurrency(stats.amounts.deduction) }}</div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Quick Actions & Recent Activity -->
      <div class="row">
        <div class="col-6">
          <div class="card">
            <div class="card-header">
              <i class="bi bi-lightning-charge"></i>
              دسترسی سریع
            </div>
            <div class="card-body">
              <div class="quick-actions">
                <RouterLink :to="{ name: 'claims' }" class="quick-action-item">
                  <div class="action-icon primary">
                    <i class="bi bi-plus-lg"></i>
                  </div>
                  <div class="action-text">ثبت ادعای جدید</div>
                </RouterLink>
                <RouterLink :to="{ name: 'packages' }" class="quick-action-item">
                  <div class="action-icon success">
                    <i class="bi bi-box-seam"></i>
                  </div>
                  <div class="action-text">مدیریت بسته‌ها</div>
                </RouterLink>
                <RouterLink :to="{ name: 'members' }" class="quick-action-item">
                  <div class="action-icon warning">
                    <i class="bi bi-search"></i>
                  </div>
                  <div class="action-text">استعلام بیمه‌شده</div>
                </RouterLink>
                <RouterLink :to="{ name: 'reports' }" class="quick-action-item">
                  <div class="action-icon info">
                    <i class="bi bi-graph-up"></i>
                  </div>
                  <div class="action-text">گزارشات</div>
                </RouterLink>
              </div>
            </div>
          </div>
        </div>

        <div class="col-6">
          <div class="card">
            <div class="card-header">
              <i class="bi bi-clock-history"></i>
              فعالیت‌های اخیر
            </div>
            <div class="card-body">
              <div class="activity-list">
                <div class="activity-item">
                  <div class="activity-dot success"></div>
                  <div class="activity-content">
                    <div class="activity-text">بسته شماره ۱۲۵۶ تایید شد</div>
                    <div class="activity-time">۱۰ دقیقه پیش</div>
                  </div>
                </div>
                <div class="activity-item">
                  <div class="activity-dot primary"></div>
                  <div class="activity-content">
                    <div class="activity-text">۱۵ ادعای جدید ثبت شد</div>
                    <div class="activity-time">۳۰ دقیقه پیش</div>
                  </div>
                </div>
                <div class="activity-item">
                  <div class="activity-dot warning"></div>
                  <div class="activity-content">
                    <div class="activity-text">مرکز درمانی جدید اضافه شد</div>
                    <div class="activity-time">۱ ساعت پیش</div>
                  </div>
                </div>
                <div class="activity-item">
                  <div class="activity-dot info"></div>
                  <div class="activity-content">
                    <div class="activity-text">تسویه حساب مرکز شفا انجام شد</div>
                    <div class="activity-time">۲ ساعت پیش</div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<style scoped>
.financial-item {
  text-align: center;
  padding: 20px;
  background: var(--bg-light);
  border-radius: 12px;
}

.financial-label {
  font-size: 0.9rem;
  color: var(--text-muted);
  margin-bottom: 8px;
}

.financial-value {
  font-size: 1.3rem;
  font-weight: 700;
}

.text-primary { color: var(--primary); }
.text-success { color: var(--secondary); }
.text-danger { color: var(--danger); }

.quick-actions {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.quick-action-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px;
  background: var(--bg-light);
  border-radius: 12px;
  text-decoration: none;
  color: var(--text-dark);
  transition: var(--transition);
}

.quick-action-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.action-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.3rem;
  color: #fff;
}

.action-icon.primary { background: linear-gradient(135deg, var(--primary) 0%, var(--primary-dark) 100%); }
.action-icon.success { background: linear-gradient(135deg, var(--secondary) 0%, var(--secondary-dark) 100%); }
.action-icon.warning { background: linear-gradient(135deg, var(--accent) 0%, var(--accent-dark) 100%); }
.action-icon.info { background: linear-gradient(135deg, var(--info) 0%, #2563eb 100%); }

.action-text {
  font-weight: 600;
  font-size: 0.95rem;
}

.activity-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.activity-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
}

.activity-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  margin-top: 6px;
  flex-shrink: 0;
}

.activity-dot.success { background: var(--success); }
.activity-dot.primary { background: var(--primary); }
.activity-dot.warning { background: var(--warning); }
.activity-dot.info { background: var(--info); }

.activity-content {
  flex: 1;
}

.activity-text {
  font-size: 0.9rem;
  color: var(--text-dark);
  margin-bottom: 4px;
}

.activity-time {
  font-size: 0.8rem;
  color: var(--text-muted);
}
</style>
