<template>
  <div class="contracts-view">
    <div class="page-header">
      <h1 class="page-title">قراردادها</h1>
      <p class="page-subtitle">مدیریت قراردادهای بیمه‌ای با کارفرمایان</p>
    </div>

    <div class="content-card">
      <div class="card-header">
        <h3>لیست قراردادها</h3>
        <button class="btn btn-primary" @click="showCreateModal = true">
          <i class="pi pi-plus"></i>
          افزودن قرارداد جدید
        </button>
      </div>

      <div class="filters-section">
        <div class="filter-group">
          <label>جستجو:</label>
          <input
            type="text"
            v-model="searchQuery"
            placeholder="جستجو در قراردادها..."
            class="form-control"
          />
        </div>
        <div class="filter-group">
          <label>وضعیت:</label>
          <select v-model="statusFilter" class="form-control">
            <option value="">همه</option>
            <option value="active">فعال</option>
            <option value="expired">منقضی شده</option>
            <option value="pending">در انتظار</option>
          </select>
        </div>
      </div>

      <div class="table-container">
        <table class="data-table">
          <thead>
            <tr>
              <th>شماره قرارداد</th>
              <th>کارفرما</th>
              <th>تاریخ شروع</th>
              <th>تاریخ پایان</th>
              <th>تعداد بیمه‌شدگان</th>
              <th>مبلغ کل حق بیمه</th>
              <th>وضعیت</th>
              <th>عملیات</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="loading">
              <td colspan="8" class="text-center">
                <div class="spinner"></div>
                در حال بارگذاری...
              </td>
            </tr>
            <tr v-else-if="contracts.length === 0">
              <td colspan="8" class="text-center text-muted">
                قراردادی یافت نشد
              </td>
            </tr>
            <tr v-else v-for="contract in filteredContracts" :key="contract.id">
              <td>{{ contract.contract_number }}</td>
              <td>{{ contract.employer_name }}</td>
              <td>{{ formatDate(contract.start_date) }}</td>
              <td>{{ formatDate(contract.end_date) }}</td>
              <td>{{ contract.total_insured }}</td>
              <td>{{ formatCurrency(contract.total_premium_amount) }}</td>
              <td>
                <span class="badge" :class="getStatusClass(contract.status)">
                  {{ getStatusLabel(contract.status) }}
                </span>
              </td>
              <td>
                <div class="action-buttons">
                  <button class="btn btn-sm btn-info" @click="viewContract(contract)" title="مشاهده">
                    <i class="pi pi-eye"></i>
                  </button>
                  <button class="btn btn-sm btn-warning" @click="editContract(contract)" title="ویرایش">
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

interface Contract {
  id: number
  contract_number: string
  employer_name: string
  start_date: string
  end_date?: string
  total_insured: number
  total_premium_amount?: number
  status: string
}

const contracts = ref<Contract[]>([])
const loading = ref(false)
const searchQuery = ref('')
const statusFilter = ref('')
const showCreateModal = ref(false)

const filteredContracts = computed(() => {
  let result = contracts.value

  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(contract =>
      contract.contract_number?.toLowerCase().includes(query) ||
      contract.employer_name?.toLowerCase().includes(query)
    )
  }

  if (statusFilter.value) {
    result = result.filter(contract => contract.status === statusFilter.value)
  }

  return result
})

const formatCurrency = (amount: number | undefined) => {
  if (!amount) return '-'
  return new Intl.NumberFormat('fa-IR').format(amount) + ' ریال'
}

const formatDate = (date: string | undefined) => {
  if (!date) return '-'
  return new Date(date).toLocaleDateString('fa-IR')
}

const getStatusLabel = (status: string) => {
  const labels: Record<string, string> = {
    active: 'فعال',
    expired: 'منقضی شده',
    pending: 'در انتظار',
    cancelled: 'لغو شده'
  }
  return labels[status] || status
}

const getStatusClass = (status: string) => {
  const classes: Record<string, string> = {
    active: 'badge-success',
    expired: 'badge-danger',
    pending: 'badge-warning',
    cancelled: 'badge-secondary'
  }
  return classes[status] || 'badge-secondary'
}

const fetchContracts = async () => {
  loading.value = true
  try {
    // TODO: Replace with actual API call
    // const response = await fetch('/api/v1/contracts')
    // contracts.value = await response.json()

    // Mock data for now
    contracts.value = []
  } catch (error) {
    console.error('Error fetching contracts:', error)
  } finally {
    loading.value = false
  }
}

const viewContract = (contract: Contract) => {
  console.log('View contract:', contract)
  // TODO: Implement view modal
}

const editContract = (contract: Contract) => {
  console.log('Edit contract:', contract)
  // TODO: Implement edit modal
}

onMounted(() => {
  fetchContracts()
})
</script>

<style scoped>
.contracts-view {
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
  display: flex;
  gap: 24px;
  flex-wrap: wrap;
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
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-family: 'Vazirmatn', sans-serif;
  min-width: 200px;
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

.badge-danger {
  background: #f8d7da;
  color: #721c24;
}

.badge-warning {
  background: #fff3cd;
  color: #856404;
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
