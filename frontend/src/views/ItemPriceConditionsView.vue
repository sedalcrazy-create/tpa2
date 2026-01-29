<script setup lang="ts">
import { ref, onMounted } from 'vue'
import api from '@/services/api'

interface PriceCondition {
  id: number
  title: string
  coverage_percentage?: number
  max_coverage_amount?: number
  franchise_percentage?: number
  priority: number
  is_active: boolean
}

const conditions = ref<PriceCondition[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const currentCondition = ref<Partial<PriceCondition>>({})

const fetchConditions = async () => {
  loading.value = true
  try {
    const response = await api.get('/item-price-conditions')
    conditions.value = response.data
  } catch (error) {
    console.error('Error fetching conditions:', error)
  } finally {
    loading.value = false
  }
}

const openDialog = (condition?: PriceCondition) => {
  if (condition) {
    currentCondition.value = { ...condition }
  } else {
    currentCondition.value = { is_active: true, priority: 0 }
  }
  dialogVisible.value = true
}

const saveCondition = async () => {
  try {
    if (currentCondition.value.id) {
      await api.put(`/item-price-conditions/${currentCondition.value.id}`, currentCondition.value)
    } else {
      await api.post('/item-price-conditions', currentCondition.value)
    }
    dialogVisible.value = false
    fetchConditions()
  } catch (error) {
    console.error('Error saving condition:', error)
  }
}

const deleteCondition = async (id: number) => {
  if (confirm('آیا از حذف این شرط قیمت‌گذاری اطمینان دارید?')) {
    try {
      await api.delete(`/item-price-conditions/${id}`)
      fetchConditions()
    } catch (error) {
      console.error('Error deleting condition:', error)
    }
  }
}

onMounted(() => {
  fetchConditions()
})
</script>

<template>
  <div class="page-container">
    <div class="page-header">
      <h1>شرایط قیمت‌گذاری</h1>
      <button class="btn btn-primary" @click="openDialog()">
        <i class="bi bi-plus-circle"></i>
        افزودن شرط جدید
      </button>
    </div>

    <div class="table-container" v-if="!loading">
      <table class="data-table">
        <thead>
          <tr>
            <th>عنوان</th>
            <th>درصد پوشش</th>
            <th>حداکثر پوشش</th>
            <th>فرانشیز</th>
            <th>اولویت</th>
            <th>وضعیت</th>
            <th>عملیات</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="condition in conditions" :key="condition.id">
            <td>{{ condition.title }}</td>
            <td>{{ condition.coverage_percentage }}%</td>
            <td>{{ condition.max_coverage_amount?.toLocaleString() }}</td>
            <td>{{ condition.franchise_percentage }}%</td>
            <td>{{ condition.priority }}</td>
            <td>
              <span :class="['badge', condition.is_active ? 'badge-success' : 'badge-danger']">
                {{ condition.is_active ? 'فعال' : 'غیرفعال' }}
              </span>
            </td>
            <td>
              <button class="btn btn-sm btn-info" @click="openDialog(condition)">
                <i class="bi bi-pencil"></i>
              </button>
              <button class="btn btn-sm btn-danger" @click="deleteCondition(condition.id)">
                <i class="bi bi-trash"></i>
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-else class="text-center">در حال بارگذاری...</div>

    <!-- Dialog -->
    <div v-if="dialogVisible" class="modal-overlay" @click="dialogVisible = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>{{ currentCondition.id ? 'ویرایش' : 'افزودن' }} شرط قیمت‌گذاری</h3>
          <button @click="dialogVisible = false">&times;</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>عنوان</label>
            <input v-model="currentCondition.title" type="text" class="form-control" />
          </div>
          <div class="form-group">
            <label>درصد پوشش</label>
            <input v-model.number="currentCondition.coverage_percentage" type="number" class="form-control" />
          </div>
          <div class="form-group">
            <label>حداکثر پوشش</label>
            <input v-model.number="currentCondition.max_coverage_amount" type="number" class="form-control" />
          </div>
          <div class="form-group">
            <label>درصد فرانشیز</label>
            <input v-model.number="currentCondition.franchise_percentage" type="number" class="form-control" />
          </div>
          <div class="form-group">
            <label>اولویت</label>
            <input v-model.number="currentCondition.priority" type="number" class="form-control" />
          </div>
          <div class="form-group">
            <label>
              <input v-model="currentCondition.is_active" type="checkbox" />
              فعال
            </label>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="dialogVisible = false">انصراف</button>
          <button class="btn btn-primary" @click="saveCondition">ذخیره</button>
        </div>
      </div>
    </div>
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

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  border-radius: 8px;
  width: 90%;
  max-width: 600px;
  max-height: 90vh;
  overflow-y: auto;
}

.modal-header {
  padding: 16px;
  border-bottom: 1px solid #eee;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-body {
  padding: 16px;
}

.modal-footer {
  padding: 16px;
  border-top: 1px solid #eee;
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  margin-bottom: 4px;
  font-weight: 500;
}

.form-control {
  width: 100%;
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.badge {
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.badge-success {
  background: #d4edda;
  color: #155724;
}

.badge-danger {
  background: #f8d7da;
  color: #721c24;
}
</style>
