<script setup lang="ts">
import { ref } from 'vue'

const searchQuery = ref('')
const searchResult = ref<any>(null)
const isSearching = ref(false)
const error = ref('')

async function search() {
  if (!searchQuery.value || searchQuery.value.length !== 10) {
    error.value = 'لطفا کد ملی ۱۰ رقمی وارد کنید'
    return
  }

  isSearching.value = true
  error.value = ''

  // Mock search
  setTimeout(() => {
    searchResult.value = {
      isEligible: true,
      member: {
        nationalCode: searchQuery.value,
        firstName: 'علی',
        lastName: 'محمدی',
        fatherName: 'رضا',
        birthDate: '1370/05/15',
        gender: 'مرد',
        dependencyType: 'بیمه‌شده اصلی'
      },
      supervisor: null,
      activePolicies: [
        { policyNumber: 'POL-123456', insurerName: 'بیمه ملی', startDate: '1403/01/01', endDate: '1403/12/29', status: 'فعال' }
      ]
    }
    isSearching.value = false
  }, 1000)
}
</script>

<template>
  <div class="members-view">
    <div class="page-header">
      <h1><i class="bi bi-people"></i> استعلام بیمه‌شدگان</h1>
    </div>

    <!-- Search -->
    <div class="card mb-4">
      <div class="card-body">
        <div class="row align-items-center">
          <div class="col-6">
            <div class="form-group mb-0">
              <label class="form-label">کد ملی بیمه‌شده</label>
              <input
                v-model="searchQuery"
                type="text"
                class="form-control"
                placeholder="کد ملی ۱۰ رقمی"
                maxlength="10"
                @keyup.enter="search"
              />
            </div>
          </div>
          <div class="col-2" style="margin-top: 28px;">
            <button class="btn btn-primary w-100" @click="search" :disabled="isSearching">
              <span v-if="isSearching" class="spinner"></span>
              <span v-else><i class="bi bi-search"></i> جستجو</span>
            </button>
          </div>
        </div>
        <div v-if="error" class="alert alert-danger mt-3 mb-0">{{ error }}</div>
      </div>
    </div>

    <!-- Result -->
    <div v-if="searchResult" class="row">
      <div class="col-6">
        <div class="card">
          <div class="card-header">
            <i class="bi bi-person"></i>
            اطلاعات بیمه‌شده
            <span class="badge badge-success" style="margin-right: auto;">دارای پوشش</span>
          </div>
          <div class="card-body">
            <div class="row">
              <div class="col-6 mb-3">
                <div class="info-label">نام</div>
                <div class="info-value">{{ searchResult.member.firstName }}</div>
              </div>
              <div class="col-6 mb-3">
                <div class="info-label">نام خانوادگی</div>
                <div class="info-value">{{ searchResult.member.lastName }}</div>
              </div>
              <div class="col-6 mb-3">
                <div class="info-label">نام پدر</div>
                <div class="info-value">{{ searchResult.member.fatherName }}</div>
              </div>
              <div class="col-6 mb-3">
                <div class="info-label">کد ملی</div>
                <div class="info-value">{{ searchResult.member.nationalCode }}</div>
              </div>
              <div class="col-6 mb-3">
                <div class="info-label">تاریخ تولد</div>
                <div class="info-value">{{ searchResult.member.birthDate }}</div>
              </div>
              <div class="col-6 mb-3">
                <div class="info-label">جنسیت</div>
                <div class="info-value">{{ searchResult.member.gender }}</div>
              </div>
              <div class="col-6">
                <div class="info-label">نسبت</div>
                <div class="info-value">{{ searchResult.member.dependencyType }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div class="col-6">
        <div class="card">
          <div class="card-header">
            <i class="bi bi-shield-check"></i>
            بیمه‌نامه‌های فعال
          </div>
          <div class="card-body">
            <div v-for="policy in searchResult.activePolicies" :key="policy.policyNumber" class="policy-item">
              <div class="policy-header">
                <span class="policy-number">{{ policy.policyNumber }}</span>
                <span class="badge badge-success">{{ policy.status }}</span>
              </div>
              <div class="policy-details">
                <span><i class="bi bi-building"></i> {{ policy.insurerName }}</span>
                <span><i class="bi bi-calendar"></i> {{ policy.startDate }} - {{ policy.endDate }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.info-label {
  font-size: 0.85rem;
  color: var(--text-muted);
  margin-bottom: 4px;
}

.info-value {
  font-weight: 600;
  color: var(--text-dark);
}

.policy-item {
  padding: 16px;
  background: var(--bg-light);
  border-radius: 12px;
  margin-bottom: 12px;
}

.policy-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.policy-number {
  font-weight: 700;
  color: var(--primary);
}

.policy-details {
  display: flex;
  gap: 20px;
  color: var(--text-muted);
  font-size: 0.9rem;
}

.policy-details i {
  margin-left: 6px;
}

.spinner {
  display: inline-block;
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
