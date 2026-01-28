<template>
  <div class="case-detail" v-if="caseData">
    <div class="page-header mb-4">
      <h2>جزئیات پرونده مددکاری</h2>
      <Button
        label="بازگشت"
        icon="pi pi-arrow-right"
        class="p-button-text"
        @click="goBack"
      />
    </div>

    <!-- اطلاعات پرونده -->
    <Card class="mb-4">
      <template #title>
        <div class="flex justify-content-between align-items-center">
          <span>{{ caseData.caseNumber }}</span>
          <Tag
            :value="getStatusLabel(caseData.status)"
            :severity="getStatusSeverity(caseData.status)"
          />
        </div>
      </template>
      <template #content>
        <div class="grid">
          <div class="col-12 md:col-6">
            <strong>نوع خدمت:</strong>
            <p>{{ getCaseTypeLabel(caseData.caseType) }}</p>
          </div>
          <div class="col-12 md:col-6">
            <strong>بیمه شده:</strong>
            <p v-if="caseData.insuredPerson">
              {{ caseData.insuredPerson.firstName }} {{ caseData.insuredPerson.lastName }}
            </p>
          </div>
          <div class="col-12 md:col-6">
            <strong>مددکار مسئول:</strong>
            <p v-if="caseData.socialWorker">
              {{ caseData.socialWorker.firstName }} {{ caseData.socialWorker.lastName }}
            </p>
          </div>
          <div class="col-12 md:col-6">
            <strong>تاریخ ثبت:</strong>
            <p>{{ formatDate(caseData.createdAt) }}</p>
          </div>
        </div>
      </template>
    </Card>

    <!-- ثبت ارزیابی -->
    <Card class="mb-4" v-if="caseData.status === 'DRAFT' || caseData.status === 'UNDER_ASSESSMENT'">
      <template #title>گزارش ارزیابی</template>
      <template #content>
        <form @submit.prevent="submitAssessment">
          <div class="p-fluid">
            <div class="field">
              <label for="assessment">گزارش ارزیابی مددکار *</label>
              <Textarea
                id="assessment"
                v-model="assessmentReport"
                rows="8"
                placeholder="گزارش کامل ارزیابی وضعیت اجتماعی-اقتصادی بیمه شده را وارد کنید..."
                :class="{ 'p-invalid': assessmentSubmitted && !assessmentReport }"
              />
              <small v-if="assessmentSubmitted && !assessmentReport" class="p-error">
                گزارش ارزیابی الزامی است
              </small>
            </div>

            <Button
              type="submit"
              label="ثبت ارزیابی"
              icon="pi pi-save"
              :loading="submittingAssessment"
              class="mt-2"
            />
          </div>
        </form>
      </template>
    </Card>

    <!-- گزارش ارزیابی ثبت شده -->
    <Card class="mb-4" v-else-if="caseData.assessmentReport">
      <template #title>گزارش ارزیابی</template>
      <template #content>
        <div class="assessment-content">
          {{ caseData.assessmentReport }}
        </div>
        <Divider />
        <p class="text-sm text-color-secondary" v-if="caseData.assessedAt">
          <i class="pi pi-calendar mr-2"></i>
          تاریخ ارزیابی: {{ formatDate(caseData.assessedAt) }}
        </p>
      </template>
    </Card>

    <!-- صدور معرفی‌نامه -->
    <Card class="mb-4" v-if="caseData.status === 'UNDER_ASSESSMENT' && caseData.assessmentReport">
      <template #title>صدور معرفی‌نامه</template>
      <template #content>
        <form @submit.prevent="submitReferral">
          <div class="p-fluid">
            <div class="field mb-3">
              <label for="referredTo">ارجاع به</label>
              <InputText
                id="referredTo"
                v-model="referralData.referredTo"
                placeholder="حسابداری"
              />
            </div>

            <div class="field mb-3">
              <label for="additionalNotes">یادداشت‌های تکمیلی</label>
              <Textarea
                id="additionalNotes"
                v-model="referralData.additionalNotes"
                rows="4"
                placeholder="توضیحات اضافی در صورت نیاز..."
              />
            </div>

            <Button
              type="submit"
              label="صدور معرفی‌نامه"
              icon="pi pi-send"
              severity="success"
              :loading="submittingReferral"
            />
          </div>
        </form>
      </template>
    </Card>

    <!-- معرفی‌نامه‌های صادر شده -->
    <Card v-if="caseData.referralLetters && caseData.referralLetters.length > 0">
      <template #title>معرفی‌نامه‌های صادر شده</template>
      <template #content>
        <DataTable :value="caseData.referralLetters" responsiveLayout="scroll">
          <Column field="letterNumber" header="شماره" />
          <Column field="referredTo" header="ارجاع به" />
          <Column field="generatedAt" header="تاریخ صدور">
            <template #body="slotProps">
              {{ formatDate(slotProps.data.generatedAt) }}
            </template>
          </Column>
          <Column header="عملیات">
            <template #body="slotProps">
              <Button
                icon="pi pi-eye"
                class="p-button-rounded p-button-text"
                @click="viewReferralLetter(slotProps.data)"
                v-tooltip="'مشاهده'"
              />
            </template>
          </Column>
        </DataTable>
      </template>
    </Card>

    <!-- Dialog نمایش معرفی‌نامه -->
    <Dialog
      v-model:visible="referralDialog"
      :style="{ width: '700px' }"
      header="معرفی‌نامه"
      :modal="true"
    >
      <div v-if="selectedReferral" class="referral-content">
        <pre style="white-space: pre-wrap; font-family: Vazirmatn;">{{ selectedReferral.content }}</pre>
      </div>
      <template #footer>
        <Button label="بستن" icon="pi pi-times" @click="referralDialog = false" />
      </template>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useToast } from 'primevue/usetoast';
import Card from 'primevue/card';
import Button from 'primevue/button';
import Tag from 'primevue/tag';
import Textarea from 'primevue/textarea';
import InputText from 'primevue/inputtext';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import Dialog from 'primevue/dialog';
import Divider from 'primevue/divider';
import { socialWorkService } from '../services/socialWorkService';
import type { SocialWorkCase, ReferralLetter } from '../types/social-work';
import {
  SocialWorkCaseTypeLabels,
  SocialWorkCaseStatusLabels,
  SocialWorkCaseStatus,
} from '../types/social-work';

const route = useRoute();
const router = useRouter();
const toast = useToast();

const caseData = ref<SocialWorkCase | null>(null);
const assessmentReport = ref('');
const assessmentSubmitted = ref(false);
const submittingAssessment = ref(false);

const referralData = ref({
  referredTo: 'حسابداری',
  additionalNotes: '',
});
const submittingReferral = ref(false);

const referralDialog = ref(false);
const selectedReferral = ref<ReferralLetter | null>(null);

const getCaseTypeLabel = (type: string) => {
  return SocialWorkCaseTypeLabels[type as keyof typeof SocialWorkCaseTypeLabels] || type;
};

const getStatusLabel = (status: string) => {
  return SocialWorkCaseStatusLabels[status as keyof typeof SocialWorkCaseStatusLabels] || status;
};

const getStatusSeverity = (status: string) => {
  switch (status) {
    case SocialWorkCaseStatus.DRAFT:
      return 'secondary';
    case SocialWorkCaseStatus.UNDER_ASSESSMENT:
      return 'info';
    case SocialWorkCaseStatus.REFERRED:
      return 'success';
    case SocialWorkCaseStatus.CLOSED:
      return 'danger';
    default:
      return 'secondary';
  }
};

const formatDate = (date: string) => {
  return new Date(date).toLocaleDateString('fa-IR');
};

const loadCase = async () => {
  try {
    const id = route.params.id as string;
    caseData.value = await socialWorkService.getCaseById(id);

    if (caseData.value.assessmentReport) {
      assessmentReport.value = caseData.value.assessmentReport;
    }
  } catch (error: any) {
    toast.add({
      severity: 'error',
      summary: 'خطا',
      detail: 'خطا در بارگذاری اطلاعات پرونده',
      life: 3000,
    });
  }
};

const submitAssessment = async () => {
  assessmentSubmitted.value = true;

  if (!assessmentReport.value.trim()) {
    toast.add({
      severity: 'warn',
      summary: 'هشدار',
      detail: 'لطفاً گزارش ارزیابی را وارد کنید',
      life: 3000,
    });
    return;
  }

  try {
    submittingAssessment.value = true;
    const id = route.params.id as string;
    await socialWorkService.updateAssessment(id, {
      assessmentReport: assessmentReport.value,
    });

    toast.add({
      severity: 'success',
      summary: 'موفق',
      detail: 'گزارش ارزیابی با موفقیت ثبت شد',
      life: 3000,
    });

    loadCase();
  } catch (error: any) {
    toast.add({
      severity: 'error',
      summary: 'خطا',
      detail: 'خطا در ثبت ارزیابی',
      life: 3000,
    });
  } finally {
    submittingAssessment.value = false;
  }
};

const submitReferral = async () => {
  try {
    submittingReferral.value = true;
    const id = route.params.id as string;
    await socialWorkService.generateReferralLetter(id, referralData.value);

    toast.add({
      severity: 'success',
      summary: 'موفق',
      detail: 'معرفی‌نامه با موفقیت صادر شد',
      life: 3000,
    });

    loadCase();
  } catch (error: any) {
    toast.add({
      severity: 'error',
      summary: 'خطا',
      detail: error.response?.data?.message || 'خطا در صدور معرفی‌نامه',
      life: 3000,
    });
  } finally {
    submittingReferral.value = false;
  }
};

const viewReferralLetter = (letter: ReferralLetter) => {
  selectedReferral.value = letter;
  referralDialog.value = true;
};

const goBack = () => {
  router.push('/social-work');
};

onMounted(() => {
  loadCase();
});
</script>

<style scoped>
.case-detail {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-header h2 {
  margin: 0;
  color: #333;
}

.assessment-content {
  white-space: pre-wrap;
  line-height: 1.8;
  padding: 15px;
  background: #f8f9fa;
  border-radius: 8px;
}

.referral-content {
  padding: 20px;
  background: #f8f9fa;
  border-radius: 8px;
  max-height: 500px;
  overflow-y: auto;
}
</style>
