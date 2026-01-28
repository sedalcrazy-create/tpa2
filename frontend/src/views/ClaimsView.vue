<script setup lang="ts">
import { ref, onMounted } from 'vue'

interface Claim {
  id: number
  trackingCode: string
  memberName: string
  centerName: string
  claimType: string
  status: string
  statusLabel: string
  requestAmount: number
  approvedAmount: number
  serviceDate: string
}

const claims = ref<Claim[]>([])
const isLoading = ref(true)
const searchQuery = ref('')
const statusFilter = ref('')

const statusOptions = [
  { value: '', label: 'همه وضعیت‌ها' },
  { value: '3', label: 'منتظر ارزیابی' },
  { value: '4', label: 'منتظر تایید' },
  { value: '5', label: 'منتظر ارسال به مالی' },
  { value: '6', label: 'آرشیو' },
  { value: '1', label: 'عودت شده' }
]

function formatCurrency(amount: number): string {
  return new Intl.NumberFormat('fa-IR').format(amount / 10) + ' تومان'
}

function getStatusClass(status: string): string {
  const classes: Record<string, string> = {
    '3': 'badge-warning',
    '4': 'badge-info',
    '5': 'badge-success',
    '6': 'badge-secondary',
    '1': 'badge-danger'
  }
  return classes[status] || 'badge-secondary'
}

onMounted(async () => {
  // Mock data
  claims.value = [
    { id: 1, trackingCode: 'CLM-1001', memberName: 'علی محمدی', centerName: 'بیمارستان شفا', claimType: 'بستری', status: '3', statusLabel: 'منتظر ارزیابی', requestAmount: 125000000, approvedAmount: 0, serviceDate: '1403/10/15' },
    { id: 2, trackingCode: 'CLM-1002', memberName: 'مریم احمدی', centerName: 'داروخانه دکتر رضایی', claimType: 'داروخانه', status: '4', statusLabel: 'منتظر تایید', requestAmount: 8500000, approvedAmount: 8200000, serviceDate: '1403/10/14' },
    { id: 3, trackingCode: 'CLM-1003', memberName: 'رضا کریمی', centerName: 'آزمایشگاه پارس', claimType: 'آزمایشگاه', status: '5', statusLabel: 'منتظر ارسال به مالی', requestAmount: 4500000, approvedAmount: 4500000, serviceDate: '1403/10/13' },
    { id: 4, trackingCode: 'CLM-1004', memberName: 'زهرا حسینی', centerName: 'مطب دکتر امینی', claimType: 'ویزیت', status: '6', statusLabel: 'آرشیو', requestAmount: 1200000, approvedAmount: 1200000, serviceDate: '1403/10/12' },
    { id: 5, trackingCode: 'CLM-1005', memberName: 'محمد رضایی', centerName: 'مرکز تصویربرداری نور', claimType: 'تصویربرداری', status: '1', statusLabel: 'عودت شده', requestAmount: 15000000, approvedAmount: 0, serviceDate: '1403/10/11' }
  ]
  isLoading.value = false
})
</script>

<template>
  <div class="claims-view">
    <!-- Page Header -->
    <div class="page-header">
      <h1><i class="bi bi-file-earmark-medical"></i> مدیریت ادعاها</h1>
      <button class="btn btn-primary">
        <i class="bi bi-plus-lg"></i>
        ثبت ادعای جدید
      </button>
    </div>

    <!-- Filters -->
    <div class="card mb-4">
      <div class="card-body">
        <div class="row align-items-center">
          <div class="col-4">
            <div class="input-wrapper">
              <i class="bi bi-search"></i>
              <input
                v-model="searchQuery"
                type="text"
                class="form-control"
                placeholder="جستجو در ادعاها..."
              />
            </div>
          </div>
          <div class="col-3">
            <select v-model="statusFilter" class="form-control">
              <option v-for="opt in statusOptions" :key="opt.value" :value="opt.value">
                {{ opt.label }}
              </option>
            </select>
          </div>
          <div class="col-3">
            <input type="date" class="form-control" placeholder="تاریخ خدمت" />
          </div>
          <div class="col-2">
            <button class="btn btn-secondary w-100">
              <i class="bi bi-funnel"></i>
              اعمال فیلتر
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="isLoading" class="loading-spinner">
      <div class="spinner"></div>
    </div>

    <!-- Claims Table -->
    <div v-else class="card">
      <div class="card-body">
        <div class="table-container">
          <table>
            <thead>
              <tr>
                <th>کد پیگیری</th>
                <th>بیمه‌شده</th>
                <th>مرکز درمانی</th>
                <th>نوع ادعا</th>
                <th>تاریخ خدمت</th>
                <th>مبلغ درخواستی</th>
                <th>مبلغ تایید شده</th>
                <th>وضعیت</th>
                <th>عملیات</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="claim in claims" :key="claim.id">
                <td>
                  <RouterLink :to="{ name: 'claim-detail', params: { id: claim.id } }" class="tracking-link">
                    {{ claim.trackingCode }}
                  </RouterLink>
                </td>
                <td>{{ claim.memberName }}</td>
                <td>{{ claim.centerName }}</td>
                <td>{{ claim.claimType }}</td>
                <td>{{ claim.serviceDate }}</td>
                <td>{{ formatCurrency(claim.requestAmount) }}</td>
                <td>{{ claim.approvedAmount > 0 ? formatCurrency(claim.approvedAmount) : '-' }}</td>
                <td>
                  <span class="badge" :class="getStatusClass(claim.status)">
                    {{ claim.statusLabel }}
                  </span>
                </td>
                <td>
                  <div class="action-buttons">
                    <RouterLink :to="{ name: 'claim-detail', params: { id: claim.id } }" class="btn btn-sm btn-secondary">
                      <i class="bi bi-eye"></i>
                    </RouterLink>
                    <button class="btn btn-sm btn-info">
                      <i class="bi bi-pencil"></i>
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Pagination -->
        <div class="pagination-wrapper">
          <div class="pagination-info">
            نمایش ۱ تا ۵ از ۱۵۴۲۰ ادعا
          </div>
          <div class="pagination-buttons">
            <button class="btn btn-sm btn-secondary" disabled>
              <i class="bi bi-chevron-right"></i>
            </button>
            <button class="btn btn-sm btn-primary">۱</button>
            <button class="btn btn-sm btn-secondary">۲</button>
            <button class="btn btn-sm btn-secondary">۳</button>
            <button class="btn btn-sm btn-secondary">
              <i class="bi bi-chevron-left"></i>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.input-wrapper {
  position: relative;
}

.input-wrapper i {
  position: absolute;
  right: 14px;
  top: 50%;
  transform: translateY(-50%);
  color: var(--text-muted);
}

.input-wrapper .form-control {
  padding-right: 44px;
}

.tracking-link {
  color: var(--primary);
  text-decoration: none;
  font-weight: 600;
}

.tracking-link:hover {
  text-decoration: underline;
}

.action-buttons {
  display: flex;
  gap: 8px;
}

.pagination-wrapper {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 24px;
  padding-top: 20px;
  border-top: 1px solid var(--border-color);
}

.pagination-info {
  color: var(--text-muted);
  font-size: 0.9rem;
}

.pagination-buttons {
  display: flex;
  gap: 8px;
}
</style>
