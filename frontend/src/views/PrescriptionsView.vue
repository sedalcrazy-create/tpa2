<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/services/api'

interface Prescription {
  id: number
  prescription_number: string
  prescription_date: string
  physician_name: string
  status: string
  employee: {
    first_name: string
    last_name: string
  }
}

const prescriptions = ref<Prescription[]>([])
const loading = ref(false)
const router = useRouter()

const statusLabels: Record<string, string> = {
  DRAFT: 'پیش‌نویس',
  SUBMITTED: 'ارسال شده',
  VALIDATED: 'تایید شده',
  REJECTED: 'رد شده',
  CONVERTED: 'تبدیل به ادعا',
}

const statusColors: Record<string, string> = {
  DRAFT: 'secondary',
  SUBMITTED: 'info',
  VALIDATED: 'success',
  REJECTED: 'danger',
  CONVERTED: 'primary',
}

const fetchPrescriptions = async () => {
  loading.value = true
  try {
    const response = await api.get('/prescriptions')
    prescriptions.value = response.data
  } catch (error) {
    console.error('Error fetching prescriptions:', error)
  } finally {
    loading.value = false
  }
}

const convertToClaim = async (id: number) => {
  if (confirm('آیا از تبدیل این نسخه به ادعا اطمینان دارید?')) {
    try {
      await api.post(`/prescriptions/${id}/convert-to-claim`)
      fetchPrescriptions()
    } catch (error) {
      console.error('Error converting prescription:', error)
    }
  }
}

const viewDetails = (id: number) => {
  router.push({ name: 'prescription-detail', params: { id } })
}

onMounted(() => {
  fetchPrescriptions()
})
</script>

<template>
  <div class="page-container">
    <div class="page-header">
      <h1>مدیریت نسخه‌های پزشکی</h1>
      <button class="btn btn-primary" @click="router.push({ name: 'prescription-create' })">
        <i class="bi bi-plus-circle"></i>
        ثبت نسخه جدید
      </button>
    </div>

    <div class="filters">
      <input type="text" class="form-control" placeholder="جستجو..." />
      <select class="form-control">
        <option value="">همه وضعیت‌ها</option>
        <option value="DRAFT">پیش‌نویس</option>
        <option value="VALIDATED">تایید شده</option>
      </select>
    </div>

    <div class="table-container" v-if="!loading">
      <table class="data-table">
        <thead>
          <tr>
            <th>شماره نسخه</th>
            <th>تاریخ</th>
            <th>بیمه‌شده</th>
            <th>پزشک</th>
            <th>وضعیت</th>
            <th>عملیات</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="rx in prescriptions" :key="rx.id">
            <td>{{ rx.prescription_number }}</td>
            <td>{{ new Date(rx.prescription_date).toLocaleDateString('fa-IR') }}</td>
            <td>{{ rx.employee.first_name }} {{ rx.employee.last_name }}</td>
            <td>{{ rx.physician_name }}</td>
            <td>
              <span :class="['badge', `badge-${statusColors[rx.status]}`]">
                {{ statusLabels[rx.status] }}
              </span>
            </td>
            <td>
              <button class="btn btn-sm btn-info" @click="viewDetails(rx.id)" title="مشاهده جزئیات">
                <i class="bi bi-eye"></i>
              </button>
              <button
                v-if="rx.status === 'VALIDATED'"
                class="btn btn-sm btn-success"
                @click="convertToClaim(rx.id)"
                title="تبدیل به ادعا"
              >
                <i class="bi bi-arrow-right-circle"></i>
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-else class="text-center">در حال بارگذاری...</div>
  </div>
</template>

<style scoped>
.page-container {
  padding: 24px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.filters {
  display: flex;
  gap: 16px;
  margin-bottom: 16px;
}

.filters .form-control {
  flex: 1;
}

.table-container {
  background: white;
  border-radius: 8px;
  padding: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.data-table {
  width: 100%;
  border-collapse: collapse;
}

.data-table th,
.data-table td {
  padding: 12px;
  text-align: right;
  border-bottom: 1px solid #eee;
}

.data-table th {
  background: #f5f5f5;
  font-weight: 600;
}

.badge {
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.badge-primary { background: #007bff; color: white; }
.badge-secondary { background: #6c757d; color: white; }
.badge-success { background: #28a745; color: white; }
.badge-danger { background: #dc3545; color: white; }
.badge-info { background: #17a2b8; color: white; }
</style>
