<script setup lang="ts">
import { ref } from 'vue'

const users = ref([
  { id: 1, username: 'admin', email: 'admin@tpa.ir', firstName: 'مدیر', lastName: 'سیستم', role: 'مدیر سیستم', isActive: true },
  { id: 2, username: 'examiner1', email: 'examiner1@tpa.ir', firstName: 'علی', lastName: 'محمدی', role: 'ارزیاب', isActive: true },
  { id: 3, username: 'financial1', email: 'financial1@tpa.ir', firstName: 'مریم', lastName: 'احمدی', role: 'کارشناس مالی', isActive: true }
])
</script>

<template>
  <div class="users-view">
    <div class="page-header">
      <h1><i class="bi bi-person-gear"></i> مدیریت کاربران</h1>
      <button class="btn btn-primary">
        <i class="bi bi-plus-lg"></i>
        افزودن کاربر
      </button>
    </div>

    <div class="card">
      <div class="card-body">
        <table>
          <thead>
            <tr>
              <th>نام کاربری</th>
              <th>ایمیل</th>
              <th>نام</th>
              <th>نقش</th>
              <th>وضعیت</th>
              <th>عملیات</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="user in users" :key="user.id">
              <td>{{ user.username }}</td>
              <td dir="ltr">{{ user.email }}</td>
              <td>{{ user.firstName }} {{ user.lastName }}</td>
              <td>{{ user.role }}</td>
              <td><span class="badge" :class="user.isActive ? 'badge-success' : 'badge-danger'">{{ user.isActive ? 'فعال' : 'غیرفعال' }}</span></td>
              <td>
                <div class="action-buttons">
                  <button class="btn btn-sm btn-info"><i class="bi bi-pencil"></i></button>
                  <button class="btn btn-sm btn-danger"><i class="bi bi-trash"></i></button>
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
