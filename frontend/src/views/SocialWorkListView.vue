<template>
  <div class="social-work-list">
    <div class="page-header mb-4">
      <h2>پرونده‌های مددکاری</h2>
      <Button
        label="ثبت پرونده جدید"
        icon="pi pi-plus"
        @click="navigateToCreate"
        severity="success"
      />
    </div>

    <DataTable
      :value="cases"
      :loading="loading"
      paginator
      :rows="10"
      :rowsPerPageOptions="[5, 10, 20, 50]"
      stripedRows
      responsiveLayout="scroll"
      class="p-datatable-sm"
    >
      <Column field="caseNumber" header="شماره پرونده" sortable />

      <Column field="caseType" header="نوع خدمت" sortable>
        <template #body="slotProps">
          <Tag :value="getCaseTypeLabel(slotProps.data.caseType)" severity="info" />
        </template>
      </Column>

      <Column field="insuredPerson" header="بیمه شده" sortable>
        <template #body="slotProps">
          <span v-if="slotProps.data.insuredPerson">
            {{ slotProps.data.insuredPerson.firstName }} {{ slotProps.data.insuredPerson.lastName }}
          </span>
        </template>
      </Column>

      <Column field="socialWorker" header="مددکار" sortable>
        <template #body="slotProps">
          <span v-if="slotProps.data.socialWorker">
            {{ slotProps.data.socialWorker.firstName }} {{ slotProps.data.socialWorker.lastName }}
          </span>
        </template>
      </Column>

      <Column field="status" header="وضعیت" sortable>
        <template #body="slotProps">
          <Tag :value="getStatusLabel(slotProps.data.status)" :severity="getStatusSeverity(slotProps.data.status)" />
        </template>
      </Column>

      <Column field="createdAt" header="تاریخ ثبت" sortable>
        <template #body="slotProps">
          {{ formatDate(slotProps.data.createdAt) }}
        </template>
      </Column>

      <Column header="عملیات" style="width: 200px">
        <template #body="slotProps">
          <Button
            icon="pi pi-eye"
            class="p-button-rounded p-button-text p-button-info"
            @click="viewCase(slotProps.data.id)"
            v-tooltip="'مشاهده جزئیات'"
          />
          <Button
            icon="pi pi-trash"
            class="p-button-rounded p-button-text p-button-danger"
            @click="confirmDelete(slotProps.data)"
            v-tooltip="'حذف'"
          />
        </template>
      </Column>
    </DataTable>

    <Dialog
      v-model:visible="deleteDialog"
      :style="{ width: '450px' }"
      header="تأیید حذف"
      :modal="true"
    >
      <div class="confirmation-content">
        <i class="pi pi-exclamation-triangle mr-3" style="font-size: 2rem" />
        <span v-if="selectedCase">
          آیا از حذف پرونده <b>{{ selectedCase.caseNumber }}</b> اطمینان دارید؟
        </span>
      </div>
      <template #footer>
        <Button label="خیر" icon="pi pi-times" class="p-button-text" @click="deleteDialog = false" />
        <Button label="بله" icon="pi pi-check" class="p-button-text" @click="deleteCase" />
      </template>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useToast } from 'primevue/usetoast';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import Button from 'primevue/button';
import Tag from 'primevue/tag';
import Dialog from 'primevue/dialog';
import { socialWorkService } from '../services/socialWorkService';
import type { SocialWorkCase } from '../types/social-work';
import {
  SocialWorkCaseTypeLabels,
  SocialWorkCaseStatusLabels,
  SocialWorkCaseStatus,
} from '../types/social-work';

const router = useRouter();
const toast = useToast();

const cases = ref<SocialWorkCase[]>([]);
const loading = ref(false);
const deleteDialog = ref(false);
const selectedCase = ref<SocialWorkCase | null>(null);

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

const loadCases = async () => {
  try {
    loading.value = true;
    cases.value = await socialWorkService.getAllCases();
  } catch (error: any) {
    toast.add({
      severity: 'error',
      summary: 'خطا',
      detail: 'خطا در بارگذاری پرونده‌ها',
      life: 3000,
    });
  } finally {
    loading.value = false;
  }
};

const navigateToCreate = () => {
  router.push('/social-work/create');
};

const viewCase = (id: string) => {
  router.push(`/social-work/${id}`);
};

const confirmDelete = (caseItem: SocialWorkCase) => {
  selectedCase.value = caseItem;
  deleteDialog.value = true;
};

const deleteCase = async () => {
  if (!selectedCase.value) return;

  try {
    await socialWorkService.deleteCase(selectedCase.value.id);
    toast.add({
      severity: 'success',
      summary: 'موفق',
      detail: 'پرونده با موفقیت حذف شد',
      life: 3000,
    });
    deleteDialog.value = false;
    loadCases();
  } catch (error: any) {
    toast.add({
      severity: 'error',
      summary: 'خطا',
      detail: 'خطا در حذف پرونده',
      life: 3000,
    });
  }
};

onMounted(() => {
  loadCases();
});
</script>

<style scoped>
.social-work-list {
  padding: 20px;
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

.confirmation-content {
  display: flex;
  align-items: center;
  gap: 10px;
}
</style>
