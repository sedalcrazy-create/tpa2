<script setup lang="ts">
import { ref } from 'vue'

const packages = ref([
  { id: 1, title: 'بسته شماره ۱۲۵۶', centerName: 'بیمارستان شفا', claimCount: 45, totalAmount: 1250000000, approvedAmount: 1180000000, status: 'پرداخت شده', letterNumber: 'LT-1001', letterDate: '1403/10/10' },
  { id: 2, title: 'بسته شماره ۱۲۵۷', centerName: 'داروخانه دکتر رضایی', claimCount: 120, totalAmount: 450000000, approvedAmount: 420000000, status: 'منتظر پرداخت', letterNumber: 'LT-1002', letterDate: '1403/10/12' },
  { id: 3, title: 'بسته شماره ۱۲۵۸', centerName: 'آزمایشگاه پارس', claimCount: 85, totalAmount: 320000000, approvedAmount: 0, status: 'منتظر ارزیابی', letterNumber: 'LT-1003', letterDate: '1403/10/14' }
])

function formatCurrency(amount: number): string {
  return new Intl.NumberFormat('fa-IR').format(amount / 10) + ' تومان'
}

function getStatusClass(status: string): string {
  const classes: Record<string, string> = {
    'پرداخت شده': 'badge-success',
    'منتظر پرداخت': 'badge-warning',
    'منتظر ارزیابی': 'badge-info',
    'برگشت خورده': 'badge-danger'
  }
  return classes[status] || 'badge-secondary'
}
</script>

<template>
  <div class="packages-view">
    <div class="page-header">
      <h1><i class="bi bi-box-seam"></i> بسته‌های اسناد</h1>
      <button class="btn btn-primary">
        <i class="bi bi-plus-lg"></i>
        ایجاد بسته جدید
      </button>
    </div>

    <div class="card">
      <div class="card-body">
        <table>
          <thead>
            <tr>
              <th>عنوان</th>
              <th>مرکز درمانی</th>
              <th>شماره نامه</th>
              <th>تاریخ نامه</th>
              <th>تعداد ادعا</th>
              <th>مبلغ کل</th>
              <th>مبلغ تایید</th>
              <th>وضعیت</th>
              <th>عملیات</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="pkg in packages" :key="pkg.id">
              <td>{{ pkg.title }}</td>
              <td>{{ pkg.centerName }}</td>
              <td>{{ pkg.letterNumber }}</td>
              <td>{{ pkg.letterDate }}</td>
              <td>{{ pkg.claimCount }}</td>
              <td>{{ formatCurrency(pkg.totalAmount) }}</td>
              <td>{{ pkg.approvedAmount > 0 ? formatCurrency(pkg.approvedAmount) : '-' }}</td>
              <td><span class="badge" :class="getStatusClass(pkg.status)">{{ pkg.status }}</span></td>
              <td>
                <div class="action-buttons">
                  <button class="btn btn-sm btn-secondary"><i class="bi bi-eye"></i></button>
                  <button class="btn btn-sm btn-info"><i class="bi bi-pencil"></i></button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<style scoped>
.action-buttons {
  display: flex;
  gap: 8px;
}
</style>
