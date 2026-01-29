<template>
  <div class="insurance-rules-view">
    <div class="page-header">
      <h1 class="page-title">قوانین بیمه</h1>
      <p class="page-subtitle">مدیریت قوانین و محدودیت‌های پوشش بیمه‌ای</p>
    </div>

    <div class="content-card">
      <div class="card-header">
        <h3>لیست قوانین بیمه</h3>
        <button class="btn btn-primary" @click="showCreateModal = true">
          <i class="pi pi-plus"></i>
          افزودن قانون جدید
        </button>
      </div>

      <div class="filters-section">
        <div class="filter-group">
          <label>جستجو:</label>
          <input
            type="text"
            v-model="searchQuery"
            placeholder="جستجو در قوانین..."
            class="form-control"
          />
        </div>
      </div>

      <div class="table-container">
        <table class="data-table">
          <thead>
            <tr>
              <th>کد</th>
              <th>نوع قانون</th>
              <th>سقف سالانه</th>
              <th>سقف هر ادعا</th>
              <th>دوره انتظار</th>
              <th>وضعیت</th>
              <th>عملیات</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="loading">
              <td colspan="7" class="text-center">
                <div class="spinner"></div>
                در حال بارگذاری...
              </td>
            </tr>
            <tr v-else-if="rules.length === 0">
              <td colspan="7" class="text-center text-muted">
                قانونی یافت نشد
              </td>
            </tr>
            <tr v-else v-for="rule in filteredRules" :key="rule.id">
              <td>{{ rule.code }}</td>
              <td>{{ rule.rule_type }}</td>
              <td>{{ formatCurrency(rule.annual_limit) }}</td>
              <td>{{ formatCurrency(rule.per_claim_limit) }}</td>
              <td>{{ rule.general_waiting_days || '-' }} روز</td>
              <td>
                <span class="badge" :class="rule.is_active ? 'badge-success' : 'badge-secondary'">
                  {{ rule.is_active ? 'فعال' : 'غیرفعال' }}
                </span>
              </td>
              <td>
                <div class="action-buttons">
                  <button class="btn btn-sm btn-info" @click="viewRule(rule)" title="مشاهده">
                    <i class="pi pi-eye"></i>
                  </button>
                  <button class="btn btn-sm btn-warning" @click="editRule(rule)" title="ویرایش">
                    <i class="pi pi-pencil"></i>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'

interface InsuranceRule {
  id: number
  code: string
  rule_type: string
  annual_limit: number
  per_claim_limit: number
  lifetime_limit?: number
  general_waiting_days?: number
  dental_waiting_days?: number
  co_payment_percentage?: number
  requires_pre_auth: boolean
  is_active: boolean
}

const rules = ref<InsuranceRule[]>([])
const loading = ref(false)
const searchQuery = ref('')
const showCreateModal = ref(false)

const filteredRules = computed(() => {
  if (!searchQuery.value) return rules.value

  const query = searchQuery.value.toLowerCase()
  return rules.value.filter(rule =>
    rule.code?.toLowerCase().includes(query) ||
    rule.rule_type?.toLowerCase().includes(query)
  )
})

const formatCurrency = (amount: number | undefined) => {
  if (!amount) return '-'
  return new Intl.NumberFormat('fa-IR').format(amount) + ' ریال'
}

const fetchRules = async () => {
  loading.value = true
  try {
    // TODO: Replace with actual API call
    // const response = await fetch('/api/v1/insurance-rules')
    // rules.value = await response.json()

    // Mock data for now
    rules.value = []
  } catch (error) {
    console.error('Error fetching rules:', error)
  } finally {
    loading.value = false
  }
}

const viewRule = (rule: InsuranceRule) => {
  console.log('View rule:', rule)
  // TODO: Implement view modal
}

const editRule = (rule: InsuranceRule) => {
  console.log('Edit rule:', rule)
  // TODO: Implement edit modal
}

onMounted(() => {
  fetchRules()
})
</script>

<style scoped>
.insurance-rules-view {
  padding: 24px;
}

.page-header {
  margin-bottom: 32px;
}

.page-title {
  font-size: 28px;
  font-weight: 700;
  color: #1a1a1a;
  margin: 0 0 8px 0;
}

.page-subtitle {
  font-size: 16px;
  color: #666;
  margin: 0;
}

.content-card {
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  border-bottom: 1px solid #e0e0e0;
}

.card-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
}

.filters-section {
  padding: 20px 24px;
  background: #f8f9fa;
  border-bottom: 1px solid #e0e0e0;
}

.filter-group {
  display: flex;
  align-items: center;
  gap: 12px;
}

.filter-group label {
  font-weight: 500;
  min-width: 60px;
}

.form-control {
  flex: 1;
  max-width: 400px;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-family: 'Vazirmatn', sans-serif;
}

.table-container {
  overflow-x: auto;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
}

.data-table th,
.data-table td {
  padding: 12px 16px;
  text-align: right;
}

.data-table thead th {
  background: #f8f9fa;
  font-weight: 600;
  border-bottom: 2px solid #e0e0e0;
}

.data-table tbody tr {
  border-bottom: 1px solid #f0f0f0;
}

.data-table tbody tr:hover {
  background: #f8f9fa;
}

.badge {
  display: inline-block;
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.badge-success {
  background: #d4edda;
  color: #155724;
}

.badge-secondary {
  background: #e2e3e5;
  color: #383d41;
}

.action-buttons {
  display: flex;
  gap: 8px;
}

.btn {
  padding: 8px 16px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-family: 'Vazirmatn', sans-serif;
  font-weight: 500;
  transition: all 0.2s;
}

.btn-primary {
  background: #ff6b6b;
  color: white;
}

.btn-primary:hover {
  background: #ff5252;
}

.btn-sm {
  padding: 6px 12px;
  font-size: 14px;
}

.btn-info {
  background: #00d2d3;
  color: white;
}

.btn-info:hover {
  background: #00b8b9;
}

.btn-warning {
  background: #feca57;
  color: #1a1a1a;
}

.btn-warning:hover {
  background: #feb236;
}

.spinner {
  display: inline-block;
  width: 20px;
  height: 20px;
  border: 3px solid #f3f3f3;
  border-top: 3px solid #ff6b6b;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-left: 8px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.text-center {
  text-align: center;
}

.text-muted {
  color: #999;
}
</style>
