<script setup lang="ts">
import { ref } from 'vue'

const centers = ref([
  { id: 1, name: 'بیمارستان شفا', siamId: 'SIAM-001', type: 'بیمارستان', province: 'تهران', city: 'تهران', contractStatus: 'فعال', phone: '021-12345678' },
  { id: 2, name: 'داروخانه دکتر رضایی', siamId: 'SIAM-002', type: 'داروخانه', province: 'تهران', city: 'تهران', contractStatus: 'فعال', phone: '021-87654321' },
  { id: 3, name: 'آزمایشگاه پارس', siamId: 'SIAM-003', type: 'آزمایشگاه', province: 'اصفهان', city: 'اصفهان', contractStatus: 'فعال', phone: '031-12345678' },
  { id: 4, name: 'مرکز تصویربرداری نور', siamId: 'SIAM-004', type: 'تصویربرداری', province: 'مشهد', city: 'مشهد', contractStatus: 'در انتظار', phone: '051-12345678' }
])

function getStatusClass(status: string): string {
  return status === 'فعال' ? 'badge-success' : 'badge-warning'
}
</script>

<template>
  <div class="centers-view">
    <div class="page-header">
      <h1><i class="bi bi-hospital"></i> مراکز درمانی</h1>
      <button class="btn btn-primary">
        <i class="bi bi-plus-lg"></i>
        افزودن مرکز
      </button>
    </div>

    <div class="card">
      <div class="card-body">
        <table>
          <thead>
            <tr>
              <th>نام مرکز</th>
              <th>کد سیام</th>
              <th>نوع</th>
              <th>استان</th>
              <th>شهر</th>
              <th>تلفن</th>
              <th>وضعیت قرارداد</th>
              <th>عملیات</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="center in centers" :key="center.id">
              <td>{{ center.name }}</td>
              <td>{{ center.siamId }}</td>
              <td>{{ center.type }}</td>
              <td>{{ center.province }}</td>
              <td>{{ center.city }}</td>
              <td dir="ltr">{{ center.phone }}</td>
              <td><span class="badge" :class="getStatusClass(center.contractStatus)">{{ center.contractStatus }}</span></td>
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
