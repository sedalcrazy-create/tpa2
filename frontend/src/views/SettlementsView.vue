<script setup lang="ts">
import { ref } from 'vue'

const settlements = ref([
  { id: 1, centerName: 'بیمارستان شفا', periodStart: '1403/09/01', periodEnd: '1403/09/30', totalAmount: 2500000000, paidAmount: 2500000000, status: 'پرداخت شده' },
  { id: 2, centerName: 'داروخانه دکتر رضایی', periodStart: '1403/09/01', periodEnd: '1403/09/30', totalAmount: 850000000, paidAmount: 0, status: 'در انتظار تایید' }
])

function formatCurrency(amount: number): string {
  return new Intl.NumberFormat('fa-IR').format(amount / 10) + ' تومان'
}

function getStatusClass(status: string): string {
  return status === 'پرداخت شده' ? 'badge-success' : 'badge-warning'
}
</script>

<template>
  <div class="settlements-view">
    <div class="page-header">
      <h1><i class="bi bi-cash-stack"></i> تسویه حساب</h1>
      <button class="btn btn-primary">
        <i class="bi bi-plus-lg"></i>
        ایجاد تسویه
      </button>
    </div>

    <div class="card">
      <div class="card-body">
        <table>
          <thead>
            <tr>
              <th>مرکز درمانی</th>
              <th>دوره از</th>
              <th>دوره تا</th>
              <th>مبلغ کل</th>
              <th>مبلغ پرداختی</th>
              <th>وضعیت</th>
              <th>عملیات</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="s in settlements" :key="s.id">
              <td>{{ s.centerName }}</td>
              <td>{{ s.periodStart }}</td>
              <td>{{ s.periodEnd }}</td>
              <td>{{ formatCurrency(s.totalAmount) }}</td>
              <td>{{ s.paidAmount > 0 ? formatCurrency(s.paidAmount) : '-' }}</td>
              <td><span class="badge" :class="getStatusClass(s.status)">{{ s.status }}</span></td>
              <td>
                <div class="action-buttons">
                  <button class="btn btn-sm btn-secondary"><i class="bi bi-eye"></i></button>
                  <button class="btn btn-sm btn-success"><i class="bi bi-check"></i></button>
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
