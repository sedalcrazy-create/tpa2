<script setup lang="ts">
// Placeholder view
</script>

<template>
  <div class="reports-view">
    <div class="page-header">
      <h1><i class="bi bi-bar-chart-line"></i> گزارشات</h1>
    </div>

    <div class="row">
      <div class="col-4">
        <div class="report-card">
          <div class="report-icon primary">
            <i class="bi bi-file-earmark-medical"></i>
          </div>
          <h3>گزارش ادعاها</h3>
          <p>گزارش جامع از ادعاهای ثبت شده</p>
          <button class="btn btn-primary">مشاهده</button>
        </div>
      </div>
      <div class="col-4">
        <div class="report-card">
          <div class="report-icon success">
            <i class="bi bi-cash-stack"></i>
          </div>
          <h3>گزارش پرداخت‌ها</h3>
          <p>گزارش تسویه حساب با مراکز</p>
          <button class="btn btn-success">مشاهده</button>
        </div>
      </div>
      <div class="col-4">
        <div class="report-card">
          <div class="report-icon warning">
            <i class="bi bi-hospital"></i>
          </div>
          <h3>گزارش مراکز</h3>
          <p>عملکرد مراکز درمانی</p>
          <button class="btn btn-warning">مشاهده</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.report-card {
  background: var(--bg-white);
  border-radius: 16px;
  padding: 32px;
  text-align: center;
  border: 1px solid var(--border-color);
  transition: var(--transition);
}

.report-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
}

.report-icon {
  width: 80px;
  height: 80px;
  margin: 0 auto 20px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 2rem;
  color: #fff;
}

.report-icon.primary { background: linear-gradient(135deg, var(--primary) 0%, var(--primary-dark) 100%); }
.report-icon.success { background: linear-gradient(135deg, var(--secondary) 0%, var(--secondary-dark) 100%); }
.report-icon.warning { background: linear-gradient(135deg, var(--accent) 0%, var(--accent-dark) 100%); }

.report-card h3 {
  font-size: 1.2rem;
  font-weight: 700;
  margin-bottom: 8px;
}

.report-card p {
  color: var(--text-muted);
  margin-bottom: 20px;
}
</style>
