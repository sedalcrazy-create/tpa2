<script setup lang="ts">
import { ref } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const claimId = route.params.id

// Mock claim data
const claim = ref({
  id: claimId,
  trackingCode: 'CLM-1001',
  status: 'منتظر ارزیابی',
  claimType: 'بستری',
  admissionType: 'بستری',
  serviceDate: '1403/10/15',
  admissionDate: '1403/10/13',
  dischargeDate: '1403/10/15',
  requestAmount: 125000000,
  approvedAmount: 0,
  deduction: 0,
  member: {
    name: 'علی محمدی',
    nationalCode: '0012345678',
    policyNumber: 'POL-123456',
    relation: 'بیمه‌شده اصلی'
  },
  center: {
    name: 'بیمارستان شفا',
    code: 'CTR-001',
    type: 'بیمارستان'
  },
  items: [
    { id: 1, title: 'هزینه اتاق', code: 'SRV-001', count: 3, requestPrice: 45000000, confirmedPrice: 0 },
    { id: 2, title: 'ویزیت پزشک', code: 'SRV-002', count: 5, requestPrice: 25000000, confirmedPrice: 0 },
    { id: 3, title: 'آزمایشات', code: 'SRV-003', count: 1, requestPrice: 35000000, confirmedPrice: 0 },
    { id: 4, title: 'دارو', code: 'DRG-001', count: 10, requestPrice: 20000000, confirmedPrice: 0 }
  ],
  diagnoses: [
    { code: 'J18.9', title: 'پنومونی' }
  ]
})

function formatCurrency(amount: number): string {
  return new Intl.NumberFormat('fa-IR').format(amount / 10) + ' تومان'
}
</script>

<template>
  <div class="claim-detail">
    <!-- Page Header -->
    <div class="page-header">
      <h1><i class="bi bi-file-earmark-medical"></i> جزئیات ادعا - {{ claim.trackingCode }}</h1>
      <div class="header-actions">
        <button class="btn btn-info">
          <i class="bi bi-pencil"></i>
          ویرایش
        </button>
        <button class="btn btn-success">
          <i class="bi bi-check-lg"></i>
          ارزیابی
        </button>
      </div>
    </div>

    <!-- Main Info -->
    <div class="row">
      <!-- Claim Info -->
      <div class="col-8">
        <div class="card">
          <div class="card-header">
            <i class="bi bi-info-circle"></i>
            اطلاعات ادعا
          </div>
          <div class="card-body">
            <div class="row">
              <div class="col-3">
                <div class="info-item">
                  <div class="info-label">کد پیگیری</div>
                  <div class="info-value">{{ claim.trackingCode }}</div>
                </div>
              </div>
              <div class="col-3">
                <div class="info-item">
                  <div class="info-label">وضعیت</div>
                  <div class="info-value"><span class="badge badge-warning">{{ claim.status }}</span></div>
                </div>
              </div>
              <div class="col-3">
                <div class="info-item">
                  <div class="info-label">نوع ادعا</div>
                  <div class="info-value">{{ claim.claimType }}</div>
                </div>
              </div>
              <div class="col-3">
                <div class="info-item">
                  <div class="info-label">نوع پذیرش</div>
                  <div class="info-value">{{ claim.admissionType }}</div>
                </div>
              </div>
            </div>
            <div class="row mt-3">
              <div class="col-3">
                <div class="info-item">
                  <div class="info-label">تاریخ پذیرش</div>
                  <div class="info-value">{{ claim.admissionDate }}</div>
                </div>
              </div>
              <div class="col-3">
                <div class="info-item">
                  <div class="info-label">تاریخ ترخیص</div>
                  <div class="info-value">{{ claim.dischargeDate }}</div>
                </div>
              </div>
              <div class="col-3">
                <div class="info-item">
                  <div class="info-label">تاریخ خدمت</div>
                  <div class="info-value">{{ claim.serviceDate }}</div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Items -->
        <div class="card">
          <div class="card-header">
            <i class="bi bi-list-check"></i>
            اقلام ادعا
          </div>
          <div class="card-body">
            <table>
              <thead>
                <tr>
                  <th>ردیف</th>
                  <th>کد</th>
                  <th>عنوان</th>
                  <th>تعداد</th>
                  <th>مبلغ درخواستی</th>
                  <th>مبلغ تایید</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(item, index) in claim.items" :key="item.id">
                  <td>{{ index + 1 }}</td>
                  <td>{{ item.code }}</td>
                  <td>{{ item.title }}</td>
                  <td>{{ item.count }}</td>
                  <td>{{ formatCurrency(item.requestPrice) }}</td>
                  <td>{{ item.confirmedPrice > 0 ? formatCurrency(item.confirmedPrice) : '-' }}</td>
                </tr>
              </tbody>
              <tfoot>
                <tr>
                  <td colspan="4"><strong>جمع کل</strong></td>
                  <td><strong>{{ formatCurrency(claim.requestAmount) }}</strong></td>
                  <td><strong>{{ claim.approvedAmount > 0 ? formatCurrency(claim.approvedAmount) : '-' }}</strong></td>
                </tr>
              </tfoot>
            </table>
          </div>
        </div>
      </div>

      <!-- Sidebar -->
      <div class="col-4">
        <!-- Member Info -->
        <div class="card">
          <div class="card-header">
            <i class="bi bi-person"></i>
            بیمه‌شده
          </div>
          <div class="card-body">
            <div class="info-item mb-3">
              <div class="info-label">نام</div>
              <div class="info-value">{{ claim.member.name }}</div>
            </div>
            <div class="info-item mb-3">
              <div class="info-label">کد ملی</div>
              <div class="info-value">{{ claim.member.nationalCode }}</div>
            </div>
            <div class="info-item mb-3">
              <div class="info-label">شماره بیمه‌نامه</div>
              <div class="info-value">{{ claim.member.policyNumber }}</div>
            </div>
            <div class="info-item">
              <div class="info-label">نسبت</div>
              <div class="info-value">{{ claim.member.relation }}</div>
            </div>
          </div>
        </div>

        <!-- Center Info -->
        <div class="card">
          <div class="card-header">
            <i class="bi bi-hospital"></i>
            مرکز درمانی
          </div>
          <div class="card-body">
            <div class="info-item mb-3">
              <div class="info-label">نام مرکز</div>
              <div class="info-value">{{ claim.center.name }}</div>
            </div>
            <div class="info-item mb-3">
              <div class="info-label">کد مرکز</div>
              <div class="info-value">{{ claim.center.code }}</div>
            </div>
            <div class="info-item">
              <div class="info-label">نوع مرکز</div>
              <div class="info-value">{{ claim.center.type }}</div>
            </div>
          </div>
        </div>

        <!-- Diagnosis -->
        <div class="card">
          <div class="card-header">
            <i class="bi bi-clipboard2-pulse"></i>
            تشخیص
          </div>
          <div class="card-body">
            <div v-for="diag in claim.diagnoses" :key="diag.code" class="diagnosis-item">
              <span class="diag-code">{{ diag.code }}</span>
              <span class="diag-title">{{ diag.title }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.header-actions {
  display: flex;
  gap: 12px;
}

.info-item {
  margin-bottom: 16px;
}

.info-label {
  font-size: 0.85rem;
  color: var(--text-muted);
  margin-bottom: 4px;
}

.info-value {
  font-weight: 600;
  color: var(--text-dark);
}

.diagnosis-item {
  display: flex;
  gap: 12px;
  padding: 10px;
  background: var(--bg-light);
  border-radius: 8px;
}

.diag-code {
  font-weight: 600;
  color: var(--primary);
}

.diag-title {
  color: var(--text-dark);
}

tfoot tr {
  background: var(--bg-light);
}
</style>
