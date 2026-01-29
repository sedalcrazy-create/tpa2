<script setup lang="ts">
import { ref, onMounted } from 'vue'

interface ImportHistory {
  id: number
  batchId: string
  importDate: string
  source: string
  totalRecords: number
  newRecords: number
  updatedRecords: number
  failedRecords: number
  status: string
  notes?: string
}

const importHistory = ref<ImportHistory[]>([])
const isLoading = ref(false)
const selectedFile = ref<File | null>(null)
const importType = ref<'csv' | 'excel'>('csv')
const currentPage = ref(1)
const itemsPerPage = 10

// Statistics
const stats = ref({
  totalEmployees: 1247,
  activeEmployees: 1189,
  familyMembers: 2834,
  lastSyncDate: '1403/10/25 - 14:30'
})

// Mock data for import history
const mockHistory: ImportHistory[] = [
  {
    id: 1,
    batchId: 'BATCH-2024-001',
    importDate: '1403/10/25 14:30',
    source: 'csv_file',
    totalRecords: 150,
    newRecords: 25,
    updatedRecords: 120,
    failedRecords: 5,
    status: 'completed',
    notes: 'به‌روزرسانی موفق اطلاعات کارمندان'
  },
  {
    id: 2,
    batchId: 'BATCH-2024-002',
    importDate: '1403/10/18 10:15',
    source: 'hr_server',
    totalRecords: 200,
    newRecords: 30,
    updatedRecords: 165,
    failedRecords: 5,
    status: 'completed',
    notes: 'همگام‌سازی خودکار از سرور منابع انسانی'
  },
  {
    id: 3,
    batchId: 'BATCH-2024-003',
    importDate: '1403/10/11 09:00',
    source: 'manual',
    totalRecords: 50,
    newRecords: 10,
    updatedRecords: 38,
    failedRecords: 2,
    status: 'completed'
  }
]

onMounted(() => {
  loadImportHistory()
})

function loadImportHistory() {
  isLoading.value = true
  setTimeout(() => {
    importHistory.value = mockHistory
    isLoading.value = false
  }, 500)
}

function handleFileSelect(event: Event) {
  const target = event.target as HTMLInputElement
  if (target.files && target.files.length > 0) {
    selectedFile.value = target.files[0]
  }
}

function clearFileSelection() {
  selectedFile.value = null
  const fileInput = document.getElementById('file-input') as HTMLInputElement
  if (fileInput) {
    fileInput.value = ''
  }
}

async function handleImport() {
  if (!selectedFile.value) {
    alert('لطفاً ابتدا فایل را انتخاب کنید')
    return
  }

  isLoading.value = true

  // TODO: Implement actual API call
  setTimeout(() => {
    alert('فایل با موفقیت آپلود و پردازش شد')
    clearFileSelection()
    loadImportHistory()
    isLoading.value = false
  }, 2000)
}

async function handleManualSync() {
  if (!confirm('آیا از همگام‌سازی با سرور منابع انسانی اطمینان دارید؟')) {
    return
  }

  isLoading.value = true

  // TODO: Implement actual API call to HR server
  setTimeout(() => {
    alert('همگام‌سازی با موفقیت انجام شد')
    loadImportHistory()
    isLoading.value = false
  }, 3000)
}

function getStatusBadgeClass(status: string): string {
  const classes: Record<string, string> = {
    completed: 'status-completed',
    processing: 'status-processing',
    failed: 'status-failed',
    pending: 'status-pending'
  }
  return classes[status] || 'status-pending'
}

function getStatusText(status: string): string {
  const texts: Record<string, string> = {
    completed: 'تکمیل شده',
    processing: 'در حال پردازش',
    failed: 'ناموفق',
    pending: 'در انتظار'
  }
  return texts[status] || status
}

function getSourceText(source: string): string {
  const sources: Record<string, string> = {
    hr_server: 'سرور منابع انسانی',
    csv_file: 'فایل CSV',
    excel_file: 'فایل Excel',
    manual: 'دستی'
  }
  return sources[source] || source
}

const paginatedHistory = ref<ImportHistory[]>([])
$: paginatedHistory.value = importHistory.value.slice(
  (currentPage.value - 1) * itemsPerPage,
  currentPage.value * itemsPerPage
)
</script>

<template>
  <div class="employee-sync-view">
    <div class="page-header">
      <h1 class="page-title">به‌روزرسانی کارمندان</h1>
      <p class="page-subtitle">مدیریت و همگام‌سازی اطلاعات کارمندان و افراد تحت تکفل</p>
    </div>

    <!-- Statistics Cards -->
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon" style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%)">
          <i class="bi bi-people-fill"></i>
        </div>
        <div class="stat-content">
          <div class="stat-label">کل کارمندان</div>
          <div class="stat-value">{{ stats.totalEmployees.toLocaleString('fa-IR') }}</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon" style="background: linear-gradient(135deg, #06b6d4 0%, #3b82f6 100%)">
          <i class="bi bi-person-check-fill"></i>
        </div>
        <div class="stat-content">
          <div class="stat-label">کارمندان فعال</div>
          <div class="stat-value">{{ stats.activeEmployees.toLocaleString('fa-IR') }}</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon" style="background: linear-gradient(135deg, #f59e0b 0%, #ef4444 100%)">
          <i class="bi bi-person-hearts"></i>
        </div>
        <div class="stat-content">
          <div class="stat-label">افراد تحت تکفل</div>
          <div class="stat-value">{{ stats.familyMembers.toLocaleString('fa-IR') }}</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon" style="background: linear-gradient(135deg, #10b981 0%, #059669 100%)">
          <i class="bi bi-clock-history"></i>
        </div>
        <div class="stat-content">
          <div class="stat-label">آخرین به‌روزرسانی</div>
          <div class="stat-value small">{{ stats.lastSyncDate }}</div>
        </div>
      </div>
    </div>

    <!-- Import Section -->
    <div class="content-card">
      <div class="card-header">
        <h2 class="card-title">
          <i class="bi bi-cloud-upload"></i>
          بارگذاری فایل کارمندان
        </h2>
      </div>

      <div class="import-section">
        <div class="import-options">
          <div class="option-group">
            <label class="option-label">نوع فایل:</label>
            <div class="radio-group">
              <label class="radio-option">
                <input type="radio" v-model="importType" value="csv" />
                <span>CSV</span>
              </label>
              <label class="radio-option">
                <input type="radio" v-model="importType" value="excel" />
                <span>Excel</span>
              </label>
            </div>
          </div>

          <div class="file-upload-section">
            <input
              type="file"
              id="file-input"
              :accept="importType === 'csv' ? '.csv' : '.xlsx,.xls'"
              @change="handleFileSelect"
              style="display: none"
            />
            <label for="file-input" class="btn btn-select-file">
              <i class="bi bi-file-earmark-arrow-up"></i>
              انتخاب فایل
            </label>
            <span v-if="selectedFile" class="file-name">
              {{ selectedFile.name }}
              <button @click="clearFileSelection" class="btn-clear-file">
                <i class="bi bi-x"></i>
              </button>
            </span>
          </div>
        </div>

        <div class="action-buttons">
          <button
            @click="handleImport"
            class="btn btn-primary"
            :disabled="!selectedFile || isLoading"
          >
            <i class="bi bi-upload"></i>
            آپلود و پردازش فایل
          </button>
          <button
            @click="handleManualSync"
            class="btn btn-secondary"
            :disabled="isLoading"
          >
            <i class="bi bi-arrow-repeat"></i>
            همگام‌سازی با سرور HR
          </button>
          <a href="#" class="btn btn-outline" download>
            <i class="bi bi-download"></i>
            دانلود فایل نمونه
          </a>
        </div>

        <div class="import-help">
          <i class="bi bi-info-circle"></i>
          <div class="help-content">
            <p><strong>نکات مهم:</strong></p>
            <ul>
              <li>فایل باید شامل ستون‌های: کد پرسنلی، کد ملی، نام، نام خانوادگی باشد</li>
              <li>برای افراد تحت تکفل، کد پرسنلی والد و نوع نسبت الزامی است</li>
              <li>فرمت تاریخ: YYYY-MM-DD (میلادی)</li>
              <li>حداکثر حجم فایل: 10 مگابایت</li>
            </ul>
          </div>
        </div>
      </div>
    </div>

    <!-- Import History -->
    <div class="content-card">
      <div class="card-header">
        <h2 class="card-title">
          <i class="bi bi-clock-history"></i>
          تاریخچه به‌روزرسانی‌ها
        </h2>
      </div>

      <div class="table-container">
        <table class="data-table">
          <thead>
            <tr>
              <th>شناسه دسته</th>
              <th>تاریخ</th>
              <th>منبع</th>
              <th>کل رکوردها</th>
              <th>جدید</th>
              <th>به‌روزرسانی</th>
              <th>ناموفق</th>
              <th>وضعیت</th>
              <th>یادداشت</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="isLoading">
              <td colspan="9" class="text-center">
                <div class="loading-spinner">در حال بارگذاری...</div>
              </td>
            </tr>
            <tr v-else-if="importHistory.length === 0">
              <td colspan="9" class="text-center empty-state">
                <i class="bi bi-inbox"></i>
                <p>هیچ تاریخچه‌ای یافت نشد</p>
              </td>
            </tr>
            <tr v-else v-for="item in paginatedHistory" :key="item.id">
              <td>
                <span class="batch-id">{{ item.batchId }}</span>
              </td>
              <td>{{ item.importDate }}</td>
              <td>
                <span class="source-badge">{{ getSourceText(item.source) }}</span>
              </td>
              <td>{{ item.totalRecords.toLocaleString('fa-IR') }}</td>
              <td>
                <span class="count-badge new">{{ item.newRecords.toLocaleString('fa-IR') }}</span>
              </td>
              <td>
                <span class="count-badge updated">{{ item.updatedRecords.toLocaleString('fa-IR') }}</span>
              </td>
              <td>
                <span class="count-badge failed">{{ item.failedRecords.toLocaleString('fa-IR') }}</span>
              </td>
              <td>
                <span :class="['status-badge', getStatusBadgeClass(item.status)]">
                  {{ getStatusText(item.status) }}
                </span>
              </td>
              <td>
                <span class="notes">{{ item.notes || '-' }}</span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div v-if="importHistory.length > itemsPerPage" class="pagination">
        <button
          @click="currentPage--"
          :disabled="currentPage === 1"
          class="btn-page"
        >
          <i class="bi bi-chevron-right"></i>
        </button>
        <span class="page-info">
          صفحه {{ currentPage.toLocaleString('fa-IR') }} از
          {{ Math.ceil(importHistory.length / itemsPerPage).toLocaleString('fa-IR') }}
        </span>
        <button
          @click="currentPage++"
          :disabled="currentPage >= Math.ceil(importHistory.length / itemsPerPage)"
          class="btn-page"
        >
          <i class="bi bi-chevron-left"></i>
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.employee-sync-view {
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

/* Statistics Grid */
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

.stat-value.small {
  font-size: 1.1rem;
}

/* Content Card */
.content-card {
  background: white;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  margin-bottom: 2rem;
  overflow: hidden;
}

.card-header {
  padding: 1.5rem;
  border-bottom: 1px solid #e2e8f0;
  background: #f8fafc;
}

.card-title {
  margin: 0;
  font-size: 1.25rem;
  font-weight: 600;
  color: #1e293b;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

/* Import Section */
.import-section {
  padding: 2rem;
}

.import-options {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.option-group {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.option-label {
  font-weight: 600;
  color: #475569;
  min-width: 100px;
}

.radio-group {
  display: flex;
  gap: 1.5rem;
}

.radio-option {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  cursor: pointer;
  font-size: 0.95rem;
}

.radio-option input[type="radio"] {
  width: 18px;
  height: 18px;
  cursor: pointer;
}

.file-upload-section {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.btn-select-file {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  background: #f1f5f9;
  border: 2px dashed #cbd5e1;
  border-radius: 8px;
  color: #475569;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.btn-select-file:hover {
  background: #e2e8f0;
  border-color: #94a3b8;
}

.file-name {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 1rem;
  background: #f1f5f9;
  border-radius: 6px;
  color: #334155;
  font-size: 0.9rem;
}

.btn-clear-file {
  background: none;
  border: none;
  color: #ef4444;
  cursor: pointer;
  padding: 0;
  display: flex;
  align-items: center;
  font-size: 1.2rem;
}

.action-buttons {
  display: flex;
  gap: 1rem;
  margin-bottom: 2rem;
  flex-wrap: wrap;
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
  text-decoration: none;
}

.btn-primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

.btn-secondary {
  background: linear-gradient(135deg, #06b6d4 0%, #3b82f6 100%);
  color: white;
}

.btn-secondary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(6, 182, 212, 0.4);
}

.btn-outline {
  background: white;
  color: #475569;
  border: 2px solid #e2e8f0;
}

.btn-outline:hover {
  background: #f8fafc;
  border-color: #cbd5e1;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.import-help {
  display: flex;
  gap: 1rem;
  padding: 1rem;
  background: #eff6ff;
  border-right: 4px solid #3b82f6;
  border-radius: 8px;
  font-size: 0.9rem;
}

.import-help i {
  color: #3b82f6;
  font-size: 1.25rem;
  flex-shrink: 0;
}

.help-content p {
  margin: 0 0 0.5rem 0;
  font-weight: 600;
  color: #1e40af;
}

.help-content ul {
  margin: 0;
  padding-right: 1.5rem;
  color: #1e3a8a;
}

.help-content li {
  margin-bottom: 0.25rem;
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

.batch-id {
  font-family: 'Courier New', monospace;
  font-size: 0.85rem;
  background: #f1f5f9;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-weight: 600;
}

.source-badge {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  background: #e0e7ff;
  color: #3730a3;
  border-radius: 12px;
  font-size: 0.8rem;
  font-weight: 600;
}

.count-badge {
  display: inline-block;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-size: 0.85rem;
  font-weight: 600;
}

.count-badge.new {
  background: #dcfce7;
  color: #166534;
}

.count-badge.updated {
  background: #dbeafe;
  color: #1e40af;
}

.count-badge.failed {
  background: #fee2e2;
  color: #991b1b;
}

.status-badge {
  display: inline-block;
  padding: 0.375rem 0.75rem;
  border-radius: 12px;
  font-size: 0.8rem;
  font-weight: 600;
}

.status-completed {
  background: #dcfce7;
  color: #166534;
}

.status-processing {
  background: #fef3c7;
  color: #92400e;
}

.status-failed {
  background: #fee2e2;
  color: #991b1b;
}

.status-pending {
  background: #e0e7ff;
  color: #3730a3;
}

.notes {
  color: #64748b;
  font-size: 0.85rem;
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
