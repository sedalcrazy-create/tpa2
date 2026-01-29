<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import api from '@/services/api'

interface Employee {
  id: number
  personnel_code: string
  national_code: string
  first_name: string
  last_name: string
  parent_id?: number
  relation_type?: string
  is_active: boolean
  status: string
}

interface Stats {
  total: number
  active: number
  family_members: number
  retired: number
}

const employees = ref<Employee[]>([])
const stats = ref<Stats>({
  total: 0,
  active: 0,
  family_members: 0,
  retired: 0
})
const searchTerm = ref('')
const filterStatus = ref('all')
const currentPage = ref(1)
const itemsPerPage = 10
const isLoading = ref(false)

onMounted(() => {
  loadEmployees()
  loadStats()
})

async function loadEmployees() {
  isLoading.value = true
  try {
    const response = await api.get('/employees', {
      params: {
        search: searchTerm.value,
        status: filterStatus.value,
        page: currentPage.value,
        limit: itemsPerPage
      }
    })

    if (response.data.success) {
      employees.value = response.data.data.employees || []
    }
  } catch (error) {
    console.error('Error loading employees:', error)
  } finally {
    isLoading.value = false
  }
}

async function loadStats() {
  try {
    const response = await api.get('/employees/stats')
    if (response.data.success) {
      stats.value = response.data.data
    }
  } catch (error) {
    console.error('Error loading stats:', error)
  }
}

function getStatusBadgeClass(status: string): string {
  const classes: Record<string, string> = {
    active: 'status-active',
    inactive: 'status-inactive',
    retired: 'status-retired'
  }
  return classes[status] || 'status-inactive'
}

function getStatusText(status: string): string {
  const texts: Record<string, string> = {
    active: 'فعال',
    inactive: 'غیرفعال',
    retired: 'بازنشسته'
  }
  return texts[status] || status
}

const filteredEmployees = computed(() => {
  let result = employees.value

  // Filter by search term
  if (searchTerm.value) {
    result = result.filter(emp =>
      emp.first_name.includes(searchTerm.value) ||
      emp.last_name.includes(searchTerm.value) ||
      emp.personnel_code.includes(searchTerm.value) ||
      emp.national_code.includes(searchTerm.value)
    )
  }

  // Filter by status
  if (filterStatus.value !== 'all') {
    result = result.filter(emp => emp.status === filterStatus.value)
  }

  return result
})

const paginatedEmployees = computed(() => {
  return filteredEmployees.value.slice(
    (currentPage.value - 1) * itemsPerPage,
    currentPage.value * itemsPerPage
  )
})

// Watch for search/filter changes and reload from API
watch([searchTerm, filterStatus], () => {
  currentPage.value = 1
  loadEmployees()
})
</script>

<template>
  <div class="employees-view">
    <div class="page-header">
      <h1 class="page-title">کارمندان و افراد تحت تکفل</h1>
      <p class="page-subtitle">مدیریت اطلاعات کارمندان و اعضای خانواده</p>
    </div>

    <!-- Filters -->
    <div class="content-card">
      <div class="filters-section">
        <div class="search-box">
          <i class="bi bi-search"></i>
          <input
            v-model="searchTerm"
            type="text"
            placeholder="جستجو بر اساس نام، کد پرسنلی، کد ملی..."
            class="search-input"
          />
        </div>

        <div class="filter-group">
          <label class="filter-label">وضعیت:</label>
          <select v-model="filterStatus" class="filter-select">
            <option value="all">همه</option>
            <option value="active">فعال</option>
            <option value="inactive">غیرفعال</option>
            <option value="retired">بازنشسته</option>
          </select>
        </div>

        <button class="btn btn-primary">
          <i class="bi bi-plus-circle"></i>
          افزودن کارمند جدید
        </button>
      </div>
    </div>

    <!-- Statistics -->
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon" style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%)">
          <i class="bi bi-people-fill"></i>
        </div>
        <div class="stat-content">
          <div class="stat-label">کل کارمندان</div>
          <div class="stat-value">{{ stats.total.toLocaleString('fa-IR') }}</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon" style="background: linear-gradient(135deg, #06b6d4 0%, #3b82f6 100%)">
          <i class="bi bi-person-check-fill"></i>
        </div>
        <div class="stat-content">
          <div class="stat-label">کارمندان فعال</div>
          <div class="stat-value">{{ stats.active.toLocaleString('fa-IR') }}</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon" style="background: linear-gradient(135deg, #f59e0b 0%, #ef4444 100%)">
          <i class="bi bi-person-hearts"></i>
        </div>
        <div class="stat-content">
          <div class="stat-label">افراد تحت تکفل</div>
          <div class="stat-value">{{ stats.family_members.toLocaleString('fa-IR') }}</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon" style="background: linear-gradient(135deg, #10b981 0%, #059669 100%)">
          <i class="bi bi-person-x-fill"></i>
        </div>
        <div class="stat-content">
          <div class="stat-label">بازنشستگان</div>
          <div class="stat-value">{{ stats.retired.toLocaleString('fa-IR') }}</div>
        </div>
      </div>
    </div>

    <!-- Table -->
    <div class="content-card">
      <div class="table-container">
        <table class="data-table">
          <thead>
            <tr>
              <th>کد پرسنلی</th>
              <th>کد ملی</th>
              <th>نام و نام خانوادگی</th>
              <th>نسبت</th>
              <th>وضعیت</th>
              <th>عملیات</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="isLoading">
              <td colspan="6" class="text-center">
                <div class="loading-spinner">در حال بارگذاری...</div>
              </td>
            </tr>
            <tr v-else-if="paginatedEmployees.length === 0">
              <td colspan="6" class="text-center empty-state">
                <i class="bi bi-inbox"></i>
                <p>کارمندی یافت نشد</p>
              </td>
            </tr>
            <tr v-else v-for="employee in paginatedEmployees" :key="employee.id">
              <td>
                <span class="code-badge">{{ employee.personnel_code }}</span>
              </td>
              <td>
                <span class="national-code">{{ employee.national_code }}</span>
              </td>
              <td>
                <div class="employee-info">
                  <span class="employee-name">{{ employee.first_name }} {{ employee.last_name }}</span>
                  <span v-if="employee.parent_id" class="family-badge">
                    <i class="bi bi-link-45deg"></i>
                    تبعی
                  </span>
                </div>
              </td>
              <td>
                <span v-if="employee.relation_type" class="relation-badge">
                  {{ employee.relation_type }}
                </span>
                <span v-else class="relation-badge main">
                  کارمند اصلی
                </span>
              </td>
              <td>
                <span :class="['status-badge', getStatusBadgeClass(employee.status)]">
                  {{ getStatusText(employee.status) }}
                </span>
              </td>
              <td>
                <div class="action-buttons">
                  <button class="btn-action" title="مشاهده">
                    <i class="bi bi-eye"></i>
                  </button>
                  <button class="btn-action" title="ویرایش">
                    <i class="bi bi-pencil"></i>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div v-if="filteredEmployees.length > itemsPerPage" class="pagination">
        <button
          @click="currentPage--"
          :disabled="currentPage === 1"
          class="btn-page"
        >
          <i class="bi bi-chevron-right"></i>
        </button>
        <span class="page-info">
          صفحه {{ currentPage.toLocaleString('fa-IR') }} از
          {{ Math.ceil(filteredEmployees.length / itemsPerPage).toLocaleString('fa-IR') }}
        </span>
        <button
          @click="currentPage++"
          :disabled="currentPage >= Math.ceil(filteredEmployees.length / itemsPerPage)"
          class="btn-page"
        >
          <i class="bi bi-chevron-left"></i>
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.employees-view {
  padding: 2rem;
}

.page-header {
  margin-bottom: 2rem;
}

.page-title {
  font-size: 1.75rem;
  font-weight: 700;
  color: #1e293b;
  margin: 0 0 0.5rem 0;
}

.page-subtitle {
  color: #64748b;
  font-size: 0.95rem;
  margin: 0;
}

.content-card {
  background: white;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  margin-bottom: 2rem;
  overflow: hidden;
}

.filters-section {
  padding: 1.5rem;
  display: flex;
  gap: 1rem;
  align-items: center;
  flex-wrap: wrap;
}

.search-box {
  flex: 1;
  min-width: 250px;
  position: relative;
}

.search-box i {
  position: absolute;
  right: 1rem;
  top: 50%;
  transform: translateY(-50%);
  color: #94a3b8;
}

.search-input {
  width: 100%;
  padding: 0.75rem 1rem 0.75rem 2.5rem;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  font-size: 0.95rem;
  transition: all 0.3s ease;
}

.search-input:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.filter-group {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.filter-label {
  font-weight: 600;
  color: #475569;
  font-size: 0.9rem;
}

.filter-select {
  padding: 0.75rem 1rem;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  font-size: 0.95rem;
  background: white;
  cursor: pointer;
}

.btn {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  border-radius: 8px;
  font-weight: 600;
  font-size: 0.95rem;
  cursor: pointer;
  transition: all 0.3s ease;
  border: none;
}

.btn-primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

/* Stats Grid */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.stat-card {
  background: white;
  border-radius: 12px;
  padding: 1.5rem;
  display: flex;
  align-items: center;
  gap: 1rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 1.5rem;
}

.stat-content {
  flex: 1;
}

.stat-label {
  font-size: 0.875rem;
  color: #64748b;
  margin-bottom: 0.25rem;
}

.stat-value {
  font-size: 1.75rem;
  font-weight: 700;
  color: #1e293b;
}

/* Table */
.table-container {
  overflow-x: auto;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
}

.data-table thead {
  background: #f8fafc;
  border-bottom: 2px solid #e2e8f0;
}

.data-table th {
  padding: 1rem;
  text-align: right;
  font-weight: 600;
  color: #475569;
  font-size: 0.875rem;
  white-space: nowrap;
}

.data-table td {
  padding: 1rem;
  border-bottom: 1px solid #f1f5f9;
  color: #334155;
  font-size: 0.9rem;
}

.data-table tbody tr:hover {
  background: #f8fafc;
}

.text-center {
  text-align: center !important;
}

.loading-spinner {
  padding: 2rem;
  color: #64748b;
}

.empty-state {
  padding: 3rem !important;
  color: #94a3b8;
}

.empty-state i {
  font-size: 3rem;
  margin-bottom: 1rem;
  display: block;
}

.code-badge {
  font-family: 'Courier New', monospace;
  font-size: 0.85rem;
  background: #f1f5f9;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-weight: 600;
}

.national-code {
  font-family: 'Courier New', monospace;
  font-size: 0.85rem;
  color: #64748b;
}

.employee-info {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.employee-name {
  font-weight: 600;
}

.family-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
  padding: 0.125rem 0.5rem;
  background: #e0e7ff;
  color: #3730a3;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 600;
}

.relation-badge {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  background: #dbeafe;
  color: #1e40af;
  border-radius: 12px;
  font-size: 0.8rem;
  font-weight: 600;
}

.relation-badge.main {
  background: #dcfce7;
  color: #166534;
}

.status-badge {
  display: inline-block;
  padding: 0.375rem 0.75rem;
  border-radius: 12px;
  font-size: 0.8rem;
  font-weight: 600;
}

.status-active {
  background: #dcfce7;
  color: #166534;
}

.status-inactive {
  background: #fee2e2;
  color: #991b1b;
}

.status-retired {
  background: #fef3c7;
  color: #92400e;
}

.action-buttons {
  display: flex;
  gap: 0.5rem;
}

.btn-action {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid #e2e8f0;
  background: white;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
  color: #64748b;
}

.btn-action:hover {
  background: #f8fafc;
  border-color: #cbd5e1;
  color: #334155;
}

/* Pagination */
.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 1rem;
  padding: 1.5rem;
  border-top: 1px solid #e2e8f0;
}

.btn-page {
  background: white;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-page:hover:not(:disabled) {
  background: #f8fafc;
  border-color: #cbd5e1;
}

.btn-page:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.page-info {
  font-size: 0.9rem;
  color: #64748b;
}
</style>
