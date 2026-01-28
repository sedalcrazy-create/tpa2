<template>
  <div class="create-case">
    <div class="page-header mb-4">
      <h2>{{ pageTitle }}</h2>
      <Button
        label="بازگشت"
        icon="pi pi-arrow-right"
        class="p-button-text"
        @click="goBack"
      />
    </div>

    <Card>
      <template #content>
        <form @submit.prevent="submitForm">
          <div class="p-fluid">
            <!-- نوع خدمت -->
            <div class="field mb-4">
              <label for="caseType">نوع خدمت مددکاری *</label>
              <Dropdown
                id="caseType"
                v-model="formData.caseType"
                :options="caseTypeOptions"
                optionLabel="label"
                optionValue="value"
                placeholder="انتخاب کنید"
                :class="{ 'p-invalid': submitted && !formData.caseType }"
              />
              <small v-if="submitted && !formData.caseType" class="p-error">
                نوع خدمت الزامی است
              </small>
            </div>

            <!-- بیمه شده -->
            <div class="field mb-4">
              <label for="insuredPerson">بیمه شده *</label>
              <Dropdown
                id="insuredPerson"
                v-model="formData.insuredPersonId"
                :options="insuredPersons"
                optionLabel="label"
                optionValue="value"
                placeholder="انتخاب کنید"
                :loading="loadingInsuredPersons"
                filter
                :class="{ 'p-invalid': submitted && !formData.insuredPersonId }"
              />
              <small v-if="submitted && !formData.insuredPersonId" class="p-error">
                انتخاب بیمه شده الزامی است
              </small>
            </div>

            <!-- پرونده پزشکی (اختیاری) -->
            <div class="field mb-4">
              <label for="medicalCase">پرونده پزشکی مرتبط (اختیاری)</label>
              <Dropdown
                id="medicalCase"
                v-model="formData.medicalCaseId"
                :options="medicalCases"
                optionLabel="label"
                optionValue="value"
                placeholder="انتخاب کنید"
                :loading="loadingMedicalCases"
                filter
                showClear
              />
            </div>

            <!-- جزئیات درخواست -->
            <div class="field mb-4">
              <label for="requestDetails">جزئیات درخواست</label>
              <Textarea
                id="requestDetails"
                v-model="requestDetailsText"
                rows="5"
                placeholder="توضیحات تکمیلی در مورد درخواست..."
              />
            </div>

            <!-- دکمه‌ها -->
            <div class="flex gap-2 justify-content-end">
              <Button
                label="انصراف"
                icon="pi pi-times"
                class="p-button-text"
                @click="goBack"
              />
              <Button
                type="submit"
                label="ثبت پرونده"
                icon="pi pi-check"
                :loading="submitting"
              />
            </div>
          </div>
        </form>
      </template>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useToast } from 'primevue/usetoast';
import Card from 'primevue/card';
import Button from 'primevue/button';
import Dropdown from 'primevue/dropdown';
import Textarea from 'primevue/textarea';
import { socialWorkService } from '../services/socialWorkService';
import { SocialWorkCaseType, SocialWorkCaseTypeLabels } from '../types/social-work';
import api from '../services/api';

const router = useRouter();
const route = useRoute();
const toast = useToast();

const formData = ref({
  caseType: null as SocialWorkCaseType | null,
  insuredPersonId: '',
  medicalCaseId: '',
});

const requestDetailsText = ref('');
const submitted = ref(false);
const submitting = ref(false);

const insuredPersons = ref<{ label: string; value: string }[]>([]);
const medicalCases = ref<{ label: string; value: string }[]>([]);
const loadingInsuredPersons = ref(false);
const loadingMedicalCases = ref(false);

// گزینه‌های نوع خدمت
const caseTypeOptions = Object.entries(SocialWorkCaseTypeLabels).map(([value, label]) => ({
  value,
  label,
}));

const loadInsuredPersons = async () => {
  try {
    loadingInsuredPersons.value = true;
    const response = await api.get('/insured-persons');
    insuredPersons.value = response.data.map((person: any) => ({
      label: `${person.firstName} ${person.lastName} (${person.nationalId})`,
      value: person.id,
    }));
  } catch (error) {
    console.error('Error loading insured persons:', error);
  } finally {
    loadingInsuredPersons.value = false;
  }
};

const loadMedicalCases = async () => {
  try {
    loadingMedicalCases.value = true;
    const response = await api.get('/cases');
    medicalCases.value = response.data.map((caseItem: any) => ({
      label: caseItem.caseNumber,
      value: caseItem.id,
    }));
  } catch (error) {
    console.error('Error loading medical cases:', error);
  } finally {
    loadingMedicalCases.value = false;
  }
};

const submitForm = async () => {
  submitted.value = true;

  if (!formData.value.caseType || !formData.value.insuredPersonId) {
    toast.add({
      severity: 'warn',
      summary: 'هشدار',
      detail: 'لطفاً فیلدهای الزامی را پر کنید',
      life: 3000,
    });
    return;
  }

  try {
    submitting.value = true;

    const data: any = {
      caseType: formData.value.caseType,
      insuredPersonId: formData.value.insuredPersonId,
    };

    if (formData.value.medicalCaseId) {
      data.medicalCaseId = formData.value.medicalCaseId;
    }

    if (requestDetailsText.value.trim()) {
      data.requestDetails = { description: requestDetailsText.value };
    }

    const newCase = await socialWorkService.createCase(data);

    toast.add({
      severity: 'success',
      summary: 'موفق',
      detail: 'پرونده مددکاری با موفقیت ثبت شد',
      life: 3000,
    });

    // هدایت به صفحه جزئیات
    router.push(`/social-work/${newCase.id}`);
  } catch (error: any) {
    toast.add({
      severity: 'error',
      summary: 'خطا',
      detail: error.response?.data?.message || 'خطا در ثبت پرونده',
      life: 3000,
    });
  } finally {
    submitting.value = false;
  }
};

const goBack = () => {
  router.push('/social-work');
};

// عنوان صفحه بر اساس نوع خدمت انتخاب شده
const pageTitle = computed(() => {
  if (formData.value.caseType) {
    return `ثبت پرونده ${SocialWorkCaseTypeLabels[formData.value.caseType]}`;
  }
  return 'ثبت پرونده مددکاری جدید';
});

onMounted(() => {
  // خواندن نوع خدمت از query parameter
  const typeFromQuery = route.query.type as SocialWorkCaseType;
  if (typeFromQuery && Object.values(SocialWorkCaseType).includes(typeFromQuery)) {
    formData.value.caseType = typeFromQuery;
  }

  loadInsuredPersons();
  loadMedicalCases();
});
</script>

<style scoped>
.create-case {
  padding: 20px;
  max-width: 800px;
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
</style>
