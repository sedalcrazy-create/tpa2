<template>
  <div>
    <div class="flex justify-content-between align-items-center mb-4">
      <div>
        <Button
          icon="pi pi-arrow-right"
          label="بازگشت"
          text
          @click="goBack"
          class="mb-2"
        />
        <h2>مدیریت رای‌های کمیسیون - {{ caseTypeName }}</h2>
      </div>
      <Button label="رای جدید" icon="pi pi-plus" @click="openCreateDialog" />
    </div>
    <Card>
      <template #content>
        <DataTable :value="verdictTemplates" :loading="loading" paginator :rows="10">
          <Column field="sortOrder" header="ترتیب" sortable />
          <Column field="title" header="عنوان رای" />
          <Column field="description" header="توضیحات" />
          <Column field="isActive" header="وضعیت">
            <template #body="{ data }">
              <Badge :value="data.isActive ? 'فعال' : 'غیرفعال'" :severity="data.isActive ? 'success' : 'danger'" />
            </template>
          </Column>
          <Column header="عملیات">
            <template #body="{ data }">
              <Button
                icon="pi pi-pencil"
                severity="warning"
                text
                rounded
                @click="openEditDialog(data)"
              />
              <Button
                icon="pi pi-trash"
                severity="danger"
                text
                rounded
                @click="confirmDelete(data)"
              />
            </template>
          </Column>
        </DataTable>
      </template>
    </Card>

    <!-- Create/Edit Dialog -->
    <Dialog
      v-model:visible="dialogVisible"
      :header="editMode ? 'ویرایش رای کمیسیون' : 'افزودن رای کمیسیون جدید'"
      :modal="true"
      :style="{ width: '500px' }"
    >
      <div class="flex flex-column gap-3">
        <div class="flex flex-column gap-2">
          <label for="title">عنوان رای</label>
          <InputText id="title" v-model="formData.title" />
        </div>
        <div class="flex flex-column gap-2">
          <label for="description">توضیحات</label>
          <Textarea id="description" v-model="formData.description" rows="3" />
        </div>
        <div class="flex flex-column gap-2">
          <label for="sortOrder">ترتیب نمایش</label>
          <InputNumber id="sortOrder" v-model="formData.sortOrder" :min="0" />
        </div>
        <div class="flex align-items-center gap-2">
          <Checkbox v-model="formData.isActive" :binary="true" inputId="isActive" />
          <label for="isActive">فعال</label>
        </div>
      </div>
      <template #footer>
        <Button label="انصراف" icon="pi pi-times" @click="dialogVisible = false" text />
        <Button
          :label="editMode ? 'ویرایش' : 'ذخیره'"
          icon="pi pi-check"
          @click="handleSave"
          :loading="saving"
        />
      </template>
    </Dialog>

    <!-- Delete Confirmation Dialog -->
    <Dialog
      v-model:visible="deleteDialogVisible"
      header="تأیید حذف"
      :modal="true"
      :style="{ width: '400px' }"
    >
      <div class="flex align-items-center gap-3">
        <i class="pi pi-exclamation-triangle" style="font-size: 2rem; color: var(--red-500)" />
        <span>آیا مطمئن هستید که می‌خواهید این رای کمیسیون را حذف کنید؟</span>
      </div>
      <template #footer>
        <Button label="انصراف" icon="pi pi-times" @click="deleteDialogVisible = false" text />
        <Button
          label="حذف"
          icon="pi pi-check"
          severity="danger"
          @click="handleDelete"
          :loading="deleting"
        />
      </template>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import Card from 'primevue/card';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import Button from 'primevue/button';
import Badge from 'primevue/badge';
import Dialog from 'primevue/dialog';
import InputText from 'primevue/inputtext';
import Textarea from 'primevue/textarea';
import InputNumber from 'primevue/inputnumber';
import Checkbox from 'primevue/checkbox';
import api from '../services/api';
import type { VerdictTemplate, CreateVerdictTemplateDto, UpdateVerdictTemplateDto } from '../types/verdict-template';
import type { CaseType } from '../types/case-type';

const router = useRouter();
const route = useRoute();

const caseTypeId = computed(() => route.params.caseTypeId as string);

const verdictTemplates = ref<VerdictTemplate[]>([]);
const caseType = ref<CaseType | null>(null);
const loading = ref(false);
const dialogVisible = ref(false);
const deleteDialogVisible = ref(false);
const editMode = ref(false);
const saving = ref(false);
const deleting = ref(false);

const formData = ref<CreateVerdictTemplateDto>({
  title: '',
  description: '',
  caseTypeId: caseTypeId.value,
  sortOrder: 0,
  isActive: true,
});

const selectedVerdictTemplate = ref<VerdictTemplate | null>(null);

const caseTypeName = computed(() => caseType.value?.name || 'در حال بارگذاری...');

const fetchCaseType = async () => {
  try {
    const response = await api.get<CaseType>(`/case-types/${caseTypeId.value}`);
    caseType.value = response.data;
  } catch (error) {
    console.error('Error fetching case type:', error);
  }
};

const fetchVerdictTemplates = async () => {
  loading.value = true;
  try {
    const response = await api.get<VerdictTemplate[]>(`/verdict-templates?caseTypeId=${caseTypeId.value}`);
    verdictTemplates.value = response.data;
  } finally {
    loading.value = false;
  }
};

const openCreateDialog = () => {
  editMode.value = false;
  formData.value = {
    title: '',
    description: '',
    caseTypeId: caseTypeId.value,
    sortOrder: 0,
    isActive: true,
  };
  dialogVisible.value = true;
};

const openEditDialog = (verdictTemplate: VerdictTemplate) => {
  editMode.value = true;
  selectedVerdictTemplate.value = verdictTemplate;
  formData.value = {
    title: verdictTemplate.title,
    description: verdictTemplate.description,
    caseTypeId: caseTypeId.value,
    sortOrder: verdictTemplate.sortOrder,
    isActive: verdictTemplate.isActive,
  };
  dialogVisible.value = true;
};

const handleSave = async () => {
  saving.value = true;
  try {
    if (editMode.value && selectedVerdictTemplate.value) {
      await api.patch<VerdictTemplate>(`/verdict-templates/${selectedVerdictTemplate.value.id}`, formData.value as UpdateVerdictTemplateDto);
    } else {
      await api.post<VerdictTemplate>('/verdict-templates', formData.value);
    }
    dialogVisible.value = false;
    await fetchVerdictTemplates();
  } finally {
    saving.value = false;
  }
};

const confirmDelete = (verdictTemplate: VerdictTemplate) => {
  selectedVerdictTemplate.value = verdictTemplate;
  deleteDialogVisible.value = true;
};

const handleDelete = async () => {
  if (!selectedVerdictTemplate.value) return;

  deleting.value = true;
  try {
    await api.delete(`/verdict-templates/${selectedVerdictTemplate.value.id}`);
    deleteDialogVisible.value = false;
    await fetchVerdictTemplates();
  } finally {
    deleting.value = false;
  }
};

const goBack = () => {
  router.push('/case-types');
};

onMounted(() => {
  fetchCaseType();
  fetchVerdictTemplates();
});
</script>
