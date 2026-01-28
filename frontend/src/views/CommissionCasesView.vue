<template>
  <div class="cases-view">
    <div class="mb-4">
      <h2>مدیریت پرونده‌ها</h2>
    </div>

    <!-- Search Section -->
    <Card class="mb-4">
      <template #title>
        <div class="flex align-items-center gap-2">
          <i class="pi pi-search"></i>
          <span>جستجوی کارمند</span>
        </div>
      </template>
      <template #content>
        <div class="flex flex-column gap-3">
          <div class="flex gap-3">
            <div class="flex-1">
              <label for="searchInput" class="block mb-2">کد پرسنلی</label>
              <InputText
                id="searchInput"
                v-model="searchQuery"
                placeholder="کد پرسنلی را وارد کنید"
                class="w-full"
                @keyup.enter="searchInsuredPerson"
              />
            </div>
            <div class="flex align-items-end">
              <Button
                label="جستجو"
                icon="pi pi-search"
                @click="searchInsuredPerson"
                :loading="searching"
              />
            </div>
          </div>

          <!-- Search Results -->
          <div v-if="searchResults.length > 0" class="mt-3">
            <label class="block mb-2">نتایج جستجو</label>
            <DataTable :value="searchResults" selectionMode="single" @row-select="onPersonSelect">
              <Column field="personnelCode" header="کد پرسنلی" />
              <Column field="nationalId" header="کد ملی" />
              <Column header="نام و نام خانوادگی">
                <template #body="{ data }">
                  {{ data.firstName }} {{ data.lastName }}
                </template>
              </Column>
              <Column field="phoneNumber" header="شماره تماس" />
            </DataTable>
          </div>

          <!-- No Results Message -->
          <div v-if="searchPerformed && searchResults.length === 0" class="text-center p-4">
            <i class="pi pi-info-circle text-4xl text-400 mb-3"></i>
            <p class="text-500">کارمندی با این کد پرسنلی یافت نشد</p>
          </div>
        </div>
      </template>
    </Card>

    <!-- Selected Person and Cases -->
    <Card v-if="selectedPerson">
      <template #title>
        <div class="flex justify-content-between align-items-center">
          <div>
            <div class="text-xl font-bold mb-1">
              {{ selectedPerson.firstName }} {{ selectedPerson.lastName }}
            </div>
            <div class="text-sm text-500">
              کد پرسنلی: {{ selectedPerson.personnelCode }} | کد ملی: {{ selectedPerson.nationalId }}
            </div>
          </div>
          <Button
            label="پرونده جدید"
            icon="pi pi-plus"
            @click="showDialog = true"
            severity="success"
          />
        </div>
      </template>
      <template #content>
        <DataTable :value="cases" :loading="loading" paginator :rows="10">
          <Column field="caseNumber" header="شماره پرونده" sortable />
          <Column header="نوع پرونده">
            <template #body="{ data }">
              {{ data.caseType?.name || '-' }}
            </template>
          </Column>
          <Column header="رای کمیسیون">
            <template #body="{ data }">
              {{ data.verdictTemplate?.title || '-' }}
            </template>
          </Column>
          <Column field="status" header="وضعیت" />
          <Column header="سطح کمیسیون">
            <template #body="{ data }">
              {{ data.caseType?.isCentralCommission ? 'مرکزی' : 'استانی' }}
            </template>
          </Column>
          <Column field="createdAt" header="تاریخ ثبت" sortable />
          <Column header="عملیات">
            <template #body="{ data }">
              <Button
                icon="pi pi-eye"
                severity="info"
                text
                rounded
                v-tooltip.top="'مشاهده جزئیات'"
                @click="viewCaseDetails(data.id)"
              />
              <Button
                icon="pi pi-pencil"
                severity="warning"
                text
                rounded
                v-tooltip.top="'ویرایش'"
              />
            </template>
          </Column>
        </DataTable>
      </template>
    </Card>

    <Dialog v-model:visible="showDialog" header="پرونده جدید" :style="{ width: '50vw' }" :modal="true">
      <div class="flex flex-column gap-3">
        <!-- Display selected person info -->
        <div v-if="selectedPerson" class="p-3 surface-100 border-round">
          <div class="text-sm text-500 mb-1">بیمه‌شده</div>
          <div class="font-bold">
            {{ selectedPerson.firstName }} {{ selectedPerson.lastName }}
          </div>
          <div class="text-sm text-500">
            کد پرسنلی: {{ selectedPerson.personnelCode }} | کد ملی: {{ selectedPerson.nationalId }}
          </div>
        </div>

        <div class="p-3 surface-50 border-round">
          <i class="pi pi-info-circle text-blue-500 ml-2"></i>
          <span class="text-sm text-600">شماره پرونده به صورت خودکار توسط سیستم ایجاد می‌شود</span>
        </div>
        <div class="flex flex-column gap-2">
          <label for="caseType">نوع پرونده</label>
          <Dropdown
            id="caseType"
            v-model="formData.caseTypeId"
            :options="caseTypeOptions"
            optionLabel="label"
            optionValue="value"
            placeholder="انتخاب کنید"
          />
        </div>
        <div class="flex flex-column gap-2" v-if="formData.caseTypeId">
          <label>سطح کمیسیون</label>
          <InputText :value="commissionLevelText" disabled />
        </div>
        <div class="flex flex-column gap-2" v-if="formData.caseTypeId">
          <label for="verdictTemplate">رای کمیسیون</label>
          <Dropdown
            id="verdictTemplate"
            v-model="formData.verdictTemplateId"
            :options="verdictTemplateOptions"
            optionLabel="label"
            optionValue="value"
            placeholder="انتخاب کنید"
            :loading="loadingVerdictTemplates"
          />
        </div>
        <div class="flex flex-column gap-2">
          <label for="province">استان</label>
          <Dropdown
            id="province"
            v-model="formData.provinceId"
            :options="provinceOptions"
            optionLabel="label"
            optionValue="value"
            placeholder="انتخاب استان"
          />
        </div>
        <div class="flex flex-column gap-2">
          <label for="description">توضیحات</label>
          <Textarea id="description" v-model="formData.description" rows="4" />
        </div>
        <div class="flex flex-column gap-2">
          <label for="files">پیوست فایل‌ها (عکس یا PDF)</label>
          <FileUpload
            id="files"
            mode="basic"
            :multiple="true"
            accept="image/*,.pdf"
            :maxFileSize="50000000"
            :auto="false"
            chooseLabel="انتخاب فایل‌ها"
            @select="onFileSelect"
          />
          <small class="text-muted">حداکثر 10 فایل، هر کدام حداکثر 50 مگابایت</small>
        </div>
      </div>
      <template #footer>
        <Button label="انصراف" severity="secondary" @click="showDialog = false" />
        <Button label="ذخیره" @click="createCase" :loading="saving" />
      </template>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue';
import { useRouter } from 'vue-router';
import Card from 'primevue/card';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import Button from 'primevue/button';
import Dialog from 'primevue/dialog';
import InputText from 'primevue/inputtext';
import Dropdown from 'primevue/dropdown';
import Textarea from 'primevue/textarea';
import FileUpload from 'primevue/fileupload';
import api from '../services/api';
import type { CaseType } from '../types/case-type';
import type { VerdictTemplate } from '../types/verdict-template';

const router = useRouter();

// Search state
const searchQuery = ref('');
const searching = ref(false);
const searchResults = ref<any[]>([]);
const searchPerformed = ref(false);
const selectedPerson = ref<any>(null);

// Cases state
const cases = ref([]);
const loading = ref(false);
const showDialog = ref(false);
const saving = ref(false);

// Form state
const caseTypes = ref<CaseType[]>([]);
const verdictTemplates = ref<VerdictTemplate[]>([]);
const provinces = ref<any[]>([]);
const loadingVerdictTemplates = ref(false);

const formData = ref({
  insuredPersonId: '',
  caseTypeId: '',
  verdictTemplateId: '',
  provinceId: '',
  description: ''
});

const uploadedFiles = ref([]);

const caseTypeOptions = computed(() =>
  caseTypes.value.filter(ct => ct.isActive).map(ct => ({
    label: ct.name,
    value: ct.id
  }))
);

const verdictTemplateOptions = computed(() =>
  verdictTemplates.value.filter(vt => vt.isActive).map(vt => ({
    label: vt.title,
    value: vt.id
  }))
);

const provinceOptions = computed(() =>
  provinces.value.map(p => ({
    label: p.name,
    value: p.id
  }))
);

const selectedCaseType = computed(() =>
  caseTypes.value.find(ct => ct.id === formData.value.caseTypeId)
);

const commissionLevelText = computed(() => {
  if (!selectedCaseType.value) return '';
  return selectedCaseType.value.isCentralCommission ? 'مرکزی (تهران)' : 'استانی';
});

// Search for insured person by personnel code
const searchInsuredPerson = async () => {
  if (!searchQuery.value.trim()) {
    return;
  }

  searching.value = true;
  searchPerformed.value = true;
  try {
    // Search by personnel code only
    const response = await api.get('/insured-persons', {
      params: {
        personnelCode: searchQuery.value
      }
    });
    searchResults.value = response.data;
  } catch (error) {
    console.error('Error searching insured person:', error);
    searchResults.value = [];
  } finally {
    searching.value = false;
  }
};

// Handle person selection from search results
const onPersonSelect = async (event: any) => {
  selectedPerson.value = event.data;
  formData.value.insuredPersonId = event.data.id;

  // Fetch cases for this person
  await fetchCases(event.data.id);
};

// Fetch cases for a specific insured person
const fetchCases = async (insuredPersonId: string) => {
  loading.value = true;
  try {
    const response = await api.get('/cases', {
      params: { insuredPersonId }
    });
    cases.value = response.data;
  } catch (error) {
    console.error('Error fetching cases:', error);
    cases.value = [];
  } finally {
    loading.value = false;
  }
};

const fetchCaseTypes = async () => {
  try {
    const response = await api.get<CaseType[]>('/case-types');
    caseTypes.value = response.data;
  } catch (error) {
    console.error('Error fetching case types:', error);
  }
};

const fetchProvinces = async () => {
  try {
    const response = await api.get('/provinces');
    provinces.value = response.data;
  } catch (error) {
    console.error('Error fetching provinces:', error);
  }
};

const fetchVerdictTemplates = async (caseTypeId: string) => {
  if (!caseTypeId) {
    verdictTemplates.value = [];
    return;
  }

  loadingVerdictTemplates.value = true;
  try {
    const response = await api.get<VerdictTemplate[]>(`/verdict-templates?caseTypeId=${caseTypeId}`);
    verdictTemplates.value = response.data;
  } catch (error) {
    console.error('Error fetching verdict templates:', error);
  } finally {
    loadingVerdictTemplates.value = false;
  }
};

// Watch for caseTypeId changes to fetch verdict templates
watch(() => formData.value.caseTypeId, (newCaseTypeId) => {
  formData.value.verdictTemplateId = ''; // Reset verdict template when case type changes
  fetchVerdictTemplates(newCaseTypeId);
});

const onFileSelect = (event: any) => {
  uploadedFiles.value = event.files;
};

const createCase = async () => {
  if (!selectedPerson.value) {
    return;
  }

  saving.value = true;
  try {
    const formDataToSend = new FormData();

    // Add form fields (caseNumber will be auto-generated by backend)
    formDataToSend.append('insuredPersonId', selectedPerson.value.id);
    if (formData.value.caseTypeId) {
      formDataToSend.append('caseTypeId', formData.value.caseTypeId);
    }
    if (formData.value.verdictTemplateId) {
      formDataToSend.append('verdictTemplateId', formData.value.verdictTemplateId);
    }
    if (formData.value.provinceId) {
      formDataToSend.append('provinceId', formData.value.provinceId);
    }
    // Determine commissionLevel based on caseType
    const commissionLevel = selectedCaseType.value?.isCentralCommission ? 'CENTRAL' : 'PROVINCIAL';
    formDataToSend.append('commissionLevel', commissionLevel);
    if (formData.value.description) {
      formDataToSend.append('description', formData.value.description);
    }

    // Add files
    if (uploadedFiles.value && uploadedFiles.value.length > 0) {
      for (const file of uploadedFiles.value) {
        formDataToSend.append('files', file);
      }
    }

    await api.post('/cases', formDataToSend, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    });

    showDialog.value = false;
    formData.value = {
      insuredPersonId: selectedPerson.value.id,
      caseTypeId: '',
      verdictTemplateId: '',
      provinceId: '',
      description: ''
    };
    uploadedFiles.value = [];

    // Reload cases for the selected person
    await fetchCases(selectedPerson.value.id);
  } catch (error) {
    console.error('Error creating case:', error);
  } finally {
    saving.value = false;
  }
};

// Navigate to case details
const viewCaseDetails = (caseId: string) => {
  router.push(`/commission/${caseId}`);
};

onMounted(() => {
  fetchCaseTypes();
  fetchProvinces();
});
</script>

<style scoped>
h2 {
  margin: 0;
  color: #495057;
}
</style>
