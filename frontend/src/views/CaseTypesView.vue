<template>
  <div>
    <div class="flex justify-content-between align-items-center mb-4">
      <h2>مدیریت انواع پرونده</h2>
      <Button label="نوع جدید" icon="pi pi-plus" @click="openCreateDialog" />
    </div>
    <Card>
      <template #content>
        <DataTable :value="caseTypes" :loading="loading" paginator :rows="10">
          <Column field="name" header="نام نوع پرونده" />
          <Column field="description" header="توضیحات" />
          <Column field="isCentralCommission" header="نوع کمیسیون">
            <template #body="{ data }">
              <Badge
                :value="data.isCentralCommission ? 'مرکزی (تهران)' : 'استانی'"
                :severity="data.isCentralCommission ? 'info' : 'success'"
              />
            </template>
          </Column>
          <Column field="isActive" header="وضعیت">
            <template #body="{ data }">
              <Badge :value="data.isActive ? 'فعال' : 'غیرفعال'" :severity="data.isActive ? 'success' : 'danger'" />
            </template>
          </Column>
          <Column header="عملیات">
            <template #body="{ data }">
              <Button
                icon="pi pi-list"
                severity="info"
                text
                rounded
                v-tooltip.top="'مدیریت رای‌های کمیسیون'"
                @click="goToVerdictTemplates(data.id)"
              />
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
      :header="editMode ? 'ویرایش نوع پرونده' : 'افزودن نوع پرونده جدید'"
      :modal="true"
      :style="{ width: '500px' }"
    >
      <div class="flex flex-column gap-3">
        <div class="flex flex-column gap-2">
          <label for="name">نام نوع پرونده</label>
          <InputText id="name" v-model="formData.name" />
        </div>
        <div class="flex flex-column gap-2">
          <label for="description">توضیحات</label>
          <Textarea id="description" v-model="formData.description" rows="3" />
        </div>
        <div class="flex align-items-center gap-2">
          <Checkbox v-model="formData.isCentralCommission" :binary="true" inputId="isCentralCommission" />
          <label for="isCentralCommission">ارجاع به کمیسیون مرکزی (تهران)</label>
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
        <span>آیا مطمئن هستید که می‌خواهید این نوع پرونده را حذف کنید؟</span>
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
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import Card from 'primevue/card';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import Button from 'primevue/button';
import Badge from 'primevue/badge';
import Dialog from 'primevue/dialog';
import InputText from 'primevue/inputtext';
import Textarea from 'primevue/textarea';
import Checkbox from 'primevue/checkbox';
import api from '../services/api';
import type { CaseType, CreateCaseTypeDto, UpdateCaseTypeDto } from '../types/case-type';

const router = useRouter();

const caseTypes = ref<CaseType[]>([]);
const loading = ref(false);
const dialogVisible = ref(false);
const deleteDialogVisible = ref(false);
const editMode = ref(false);
const saving = ref(false);
const deleting = ref(false);

const formData = ref<CreateCaseTypeDto>({
  name: '',
  description: '',
  isActive: true,
  isCentralCommission: false,
});

const selectedCaseType = ref<CaseType | null>(null);

const fetchCaseTypes = async () => {
  loading.value = true;
  try {
    const response = await api.get<CaseType[]>('/case-types');
    caseTypes.value = response.data;
  } finally {
    loading.value = false;
  }
};

const openCreateDialog = () => {
  editMode.value = false;
  formData.value = {
    name: '',
    description: '',
    isActive: true,
    isCentralCommission: false,
  };
  dialogVisible.value = true;
};

const openEditDialog = (caseType: CaseType) => {
  editMode.value = true;
  selectedCaseType.value = caseType;
  formData.value = {
    name: caseType.name,
    description: caseType.description,
    isActive: caseType.isActive,
    isCentralCommission: caseType.isCentralCommission,
  };
  dialogVisible.value = true;
};

const handleSave = async () => {
  saving.value = true;
  try {
    if (editMode.value && selectedCaseType.value) {
      await api.patch<CaseType>(`/case-types/${selectedCaseType.value.id}`, formData.value as UpdateCaseTypeDto);
    } else {
      await api.post<CaseType>('/case-types', formData.value);
    }
    dialogVisible.value = false;
    await fetchCaseTypes();
  } finally {
    saving.value = false;
  }
};

const confirmDelete = (caseType: CaseType) => {
  selectedCaseType.value = caseType;
  deleteDialogVisible.value = true;
};

const handleDelete = async () => {
  if (!selectedCaseType.value) return;

  deleting.value = true;
  try {
    await api.delete(`/case-types/${selectedCaseType.value.id}`);
    deleteDialogVisible.value = false;
    await fetchCaseTypes();
  } finally {
    deleting.value = false;
  }
};

const goToVerdictTemplates = (caseTypeId: string) => {
  router.push(`/case-types/${caseTypeId}/verdict-templates`);
};

onMounted(() => fetchCaseTypes());
</script>
