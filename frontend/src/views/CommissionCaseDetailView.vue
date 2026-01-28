<template>
  <div class="case-details-wrapper">
    <!-- دکمه بازگشت -->
    <div class="mb-4">
      <Button
        label="بازگشت به لیست پرونده‌ها"
        icon="pi pi-arrow-right"
        severity="secondary"
        @click="goBack"
      />
    </div>

    <!-- حالت بارگذاری -->
    <div v-if="loading" class="flex justify-content-center align-items-center" style="min-height: 400px">
      <ProgressSpinner />
    </div>

    <!-- محتوای اصلی -->
    <div v-else-if="caseData">
      <!-- بخش برجسته - شماره پرونده -->
      <div class="highlight-box mb-4">
        <div class="flex justify-content-between align-items-center">
          <div class="flex align-items-center gap-3">
            <div style="font-size: 2rem; color: #1e40af;">📁</div>
            <div>
              <div class="text-sm text-muted mb-1">شماره پرونده</div>
              <div class="number-value">{{ caseData.caseNumber }}</div>
            </div>
          </div>
          <div>
            <Tag
              :value="getStatusLabel(caseData.status)"
              :severity="getStatusSeverity(caseData.status)"
              style="font-size: 1.1rem; padding: 0.5rem 1rem;"
            />
          </div>
        </div>
      </div>

      <!-- دکمه‌های عملیاتی -->
      <Card class="mb-4 actions-card">
        <template #content>
          <div class="flex flex-wrap gap-2">
            <Button label="چاپ پرونده" icon="pi pi-print" severity="secondary" outlined />
            <Button label="ویرایش اطلاعات" icon="pi pi-pencil" severity="secondary" outlined />
            <Button label="ارجاع به متخصص" icon="pi pi-user-plus" severity="info" @click="showAssignDialog = true" />
            <Button label="تغییر وضعیت" icon="pi pi-refresh" severity="warning" @click="showStatusDialog = true" />
            <Button label="ثبت یادداشت" icon="pi pi-comment" severity="secondary" outlined />
            <Button label="تأیید نهایی" icon="pi pi-check-circle" severity="success" />
          </div>
        </template>
      </Card>

      <!-- محتوای اصلی - استفاده از Grid PrimeFlex -->
      <div class="grid">
        <!-- ستون اصلی - 8 از 12 -->
        <div class="col-12 lg:col-8">
          <!-- اطلاعات بیمه‌شده -->
          <Card class="mb-4">
            <template #title>
              <div class="custom-card-title">
                <i class="pi pi-user ml-2"></i>
                اطلاعات بیمه‌شده
              </div>
            </template>
            <template #content>
              <div class="grid">
                <div class="col-12 md:col-6">
                  <div class="data-row">
                    <span class="data-label">نام و نام خانوادگی:</span>
                    <span class="data-value font-bold">
                      {{ caseData.insuredPerson?.firstName }} {{ caseData.insuredPerson?.lastName }}
                    </span>
                  </div>
                  <div class="data-row">
                    <span class="data-label">کد ملی:</span>
                    <span class="data-value">{{ caseData.insuredPerson?.nationalId }}</span>
                  </div>
                  <div class="data-row">
                    <span class="data-label">کد پرسنلی:</span>
                    <span class="data-value">{{ caseData.insuredPerson?.personnelCode }}</span>
                  </div>
                </div>
                <div class="col-12 md:col-6">
                  <div class="data-row">
                    <span class="data-label">شماره تماس:</span>
                    <span class="data-value">{{ caseData.insuredPerson?.phone || '-' }}</span>
                  </div>
                  <div class="data-row">
                    <span class="data-label">آدرس:</span>
                    <span class="data-value">{{ caseData.insuredPerson?.address || '-' }}</span>
                  </div>
                  <div class="data-row">
                    <span class="data-label">نسبت خانوادگی:</span>
                    <span class="data-value">{{ getFamilyRelation(caseData.insuredPerson?.familyRelation) }}</span>
                  </div>
                </div>
              </div>
            </template>
          </Card>

          <!-- جزئیات پرونده -->
          <Card class="mb-4">
            <template #title>
              <div class="custom-card-title">
                <i class="pi pi-file ml-2"></i>
                جزئیات پرونده
              </div>
            </template>
            <template #content>
              <div class="grid">
                <div class="col-12 md:col-6">
                  <div class="data-row">
                    <span class="data-label">نوع پرونده:</span>
                    <span class="data-value font-medium">{{ caseData.caseType?.name || '-' }}</span>
                  </div>
                  <div class="data-row">
                    <span class="data-label">رای کمیسیون:</span>
                    <span class="data-value font-medium">{{ caseData.verdictTemplate?.title || '-' }}</span>
                  </div>
                  <div class="data-row">
                    <span class="data-label">سطح کمیسیون:</span>
                    <span class="data-value">
                      <Tag
                        :value="caseData.commissionLevel === 'CENTRAL' ? 'مرکزی (تهران)' : 'استانی'"
                        :severity="caseData.commissionLevel === 'CENTRAL' ? 'info' : 'success'"
                      />
                    </span>
                  </div>
                </div>
                <div class="col-12 md:col-6">
                  <div class="data-row">
                    <span class="data-label">متخصص مسئول:</span>
                    <span class="data-value">
                      <span v-if="caseData.assignedTo">
                        {{ caseData.assignedTo.firstName }} {{ caseData.assignedTo.lastName }}
                        <Tag
                          v-if="caseData.assignedTo.specialty"
                          :value="caseData.assignedTo.specialty"
                          severity="info"
                          class="mr-2"
                        />
                      </span>
                      <span v-else class="text-muted">تخصیص داده نشده</span>
                    </span>
                  </div>
                  <div class="data-row">
                    <span class="data-label">تاریخ ثبت:</span>
                    <span class="data-value">{{ formatDate(caseData.createdAt) }}</span>
                  </div>
                  <div class="data-row">
                    <span class="data-label">آخرین بروزرسانی:</span>
                    <span class="data-value">{{ formatDate(caseData.updatedAt) }}</span>
                  </div>
                </div>
              </div>

              <div class="custom-divider"></div>

              <div class="data-row" v-if="caseData.description">
                <span class="data-label">توضیحات:</span>
                <span class="data-value">{{ caseData.description }}</span>
              </div>
            </template>
          </Card>

          <!-- اسناد پیوست -->
          <Card v-if="caseData.documents && caseData.documents.length > 0" class="mb-4">
            <template #title>
              <div class="custom-card-title">
                <i class="pi pi-paperclip ml-2"></i>
                اسناد پیوست ({{ caseData.documents.length }})
              </div>
            </template>
            <template #content>
              <div class="grid">
                <div
                  v-for="(doc, index) in caseData.documents"
                  :key="index"
                  class="col-12 md:col-6"
                >
                  <div class="doc-box">
                    <!-- Image preview thumbnail -->
                    <div v-if="isImageFile(doc.fileType)" class="file-preview-thumb">
                      <Image
                        :src="getFileUrl(doc.fileName)"
                        :alt="doc.title"
                        width="60"
                        height="60"
                        preview
                        imageStyle="object-fit: cover; border-radius: 0.5rem;"
                      />
                    </div>
                    <!-- PDF icon -->
                    <div v-else-if="isPdfFile(doc.fileType)" class="file-icon">
                      <i class="pi pi-file-pdf" style="color: #dc2626; font-size: 2rem;"></i>
                    </div>
                    <!-- Other file icon -->
                    <div v-else class="file-icon">
                      <i class="pi pi-file" style="color: #1e40af; font-size: 1.25rem;"></i>
                    </div>

                    <div class="flex-1 ml-3">
                      <div class="font-medium">{{ doc.title }}</div>
                      <div class="text-sm text-muted">{{ formatFileSize(doc.fileSize) }}</div>
                      <div class="text-xs text-muted">{{ doc.fileType }}</div>
                    </div>

                    <!-- Action buttons -->
                    <div class="flex gap-2">
                      <Button
                        v-if="isPdfFile(doc.fileType)"
                        icon="pi pi-eye"
                        label="مشاهده"
                        size="small"
                        severity="info"
                        outlined
                        @click="openPdfInNewTab(doc.fileName)"
                        v-tooltip.top="'باز کردن PDF'"
                      />
                      <Button
                        icon="pi pi-download"
                        size="small"
                        text
                        rounded
                        severity="secondary"
                        v-tooltip.top="'دانلود فایل'"
                      />
                    </div>
                  </div>
                </div>
              </div>
            </template>
          </Card>
        </div>

        <!-- ستون کناری - 4 از 12 -->
        <div class="col-12 lg:col-4">
          <Card class="timeline-box">
            <template #title>
              <div class="custom-card-title">
                <i class="pi pi-clock ml-2"></i>
                جدول زمانی رویدادها
              </div>
            </template>
            <template #content>
              <Timeline
                v-if="timelineEvents.length > 0"
                :value="timelineEvents"
                align="right"
                class="custom-timeline"
              >
                <template #marker="slotProps">
                  <span
                    class="timeline-dot"
                    :style="{ backgroundColor: slotProps.item.color }"
                  >
                    <i :class="slotProps.item.icon"></i>
                  </span>
                </template>
                <template #content="slotProps">
                  <div class="timeline-box-content">
                    <div class="timeline-event-title">{{ slotProps.item.title }}</div>
                    <div class="timeline-event-date">{{ slotProps.item.date }}</div>
                    <div v-if="slotProps.item.description" class="timeline-event-desc">
                      {{ slotProps.item.description }}
                    </div>
                  </div>
                </template>
              </Timeline>

              <div v-else class="text-center text-muted p-4">
                <i class="pi pi-info-circle text-4xl mb-3"></i>
                <p>هنوز رویدادی ثبت نشده است</p>
              </div>
            </template>
          </Card>
        </div>
      </div>
    </div>

    <!-- دیالوگ ارجاع به متخصص -->
    <Dialog
      v-model:visible="showAssignDialog"
      header="ارجاع پرونده به متخصص"
      :style="{ width: '500px' }"
      :modal="true"
    >
      <div class="flex flex-column gap-4">
        <div class="flex flex-column gap-2">
          <label for="specialist" class="font-medium">انتخاب متخصص</label>
          <Dropdown
            id="specialist"
            v-model="assignData.specialistId"
            :options="specialists"
            optionLabel="label"
            optionValue="value"
            placeholder="متخصص مورد نظر را انتخاب کنید"
            :loading="loadingSpecialists"
            filter
            class="w-full"
          />
        </div>
        <div class="flex flex-column gap-2">
          <label for="assignNote" class="font-medium">توضیحات ارجاع</label>
          <Textarea
            id="assignNote"
            v-model="assignData.note"
            rows="4"
            placeholder="توضیحات تکمیلی در مورد ارجاع را وارد کنید..."
          />
        </div>
      </div>
      <template #footer>
        <Button label="انصراف" severity="secondary" outlined @click="showAssignDialog = false" />
        <Button label="ثبت ارجاع" icon="pi pi-check" @click="assignToSpecialist" :loading="assigning" />
      </template>
    </Dialog>

    <!-- دیالوگ تغییر وضعیت -->
    <Dialog
      v-model:visible="showStatusDialog"
      header="تغییر وضعیت پرونده"
      :style="{ width: '500px' }"
      :modal="true"
    >
      <div class="flex flex-column gap-4">
        <div class="flex flex-column gap-2">
          <label for="newStatus" class="font-medium">وضعیت جدید</label>
          <Dropdown
            id="newStatus"
            v-model="statusData.newStatus"
            :options="statusOptions"
            optionLabel="label"
            optionValue="value"
            placeholder="وضعیت جدید را انتخاب کنید"
            class="w-full"
          />
        </div>
        <div class="flex flex-column gap-2">
          <label for="statusNote" class="font-medium">توضیحات</label>
          <Textarea
            id="statusNote"
            v-model="statusData.note"
            rows="4"
            placeholder="توضیحات مربوط به تغییر وضعیت را وارد کنید..."
          />
        </div>
      </div>
      <template #footer>
        <Button label="انصراف" severity="secondary" outlined @click="showStatusDialog = false" />
        <Button label="ثبت تغییر" icon="pi pi-check" @click="changeStatus" :loading="changingStatus" />
      </template>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import Card from 'primevue/card';
import Button from 'primevue/button';
import Tag from 'primevue/tag';
import Timeline from 'primevue/timeline';
import Dialog from 'primevue/dialog';
import Dropdown from 'primevue/dropdown';
import Textarea from 'primevue/textarea';
import ProgressSpinner from 'primevue/progressspinner';
import Image from 'primevue/image';
import api from '../services/api';

const route = useRoute();
const router = useRouter();

const caseData = ref<any>(null);
const loading = ref(false);
const showAssignDialog = ref(false);
const showStatusDialog = ref(false);
const specialists = ref<any[]>([]);
const loadingSpecialists = ref(false);
const assigning = ref(false);
const changingStatus = ref(false);

const assignData = ref({
  specialistId: '',
  note: ''
});

const statusData = ref({
  newStatus: '',
  note: ''
});

const statusOptions = [
  { label: 'در انتظار بررسی دبیرخانه', value: 'PENDING_SECRETARIAT' },
  { label: 'ارجاع به متخصص', value: 'ASSIGNED_TO_SPECIALIST' },
  { label: 'در حال بررسی', value: 'UNDER_REVIEW' },
  { label: 'در انتظار جلسه', value: 'PENDING_MEETING' },
  { label: 'جلسه زمان‌بندی شده', value: 'MEETING_SCHEDULED' },
  { label: 'در انتظار رأی', value: 'PENDING_VERDICT' },
  { label: 'رأی صادر شد', value: 'VERDICT_ISSUED' },
  { label: 'بایگانی شده', value: 'ARCHIVED' },
  { label: 'رد شده', value: 'REJECTED' }
];

const timelineEvents = computed(() => {
  if (!caseData.value) return [];

  if (caseData.value.timeline && caseData.value.timeline.length > 0) {
    return caseData.value.timeline.map((event: any) => ({
      title: event.action,
      date: formatDate(event.createdAt),
      description: event.description,
      icon: getTimelineIcon(event.action),
      color: getTimelineColor(event.action)
    }));
  }

  const events = [];

  events.push({
    title: 'ثبت پرونده',
    date: formatDate(caseData.value.createdAt),
    icon: 'pi pi-plus',
    color: '#22C55E'
  });

  if (caseData.value.assignedTo) {
    events.push({
      title: 'ارجاع به متخصص',
      date: formatDate(caseData.value.assignedAt),
      description: `${caseData.value.assignedTo.firstName} ${caseData.value.assignedTo.lastName}`,
      icon: 'pi pi-user',
      color: '#3B82F6'
    });
  }

  events.push({
    title: getStatusLabel(caseData.value.status),
    date: formatDate(caseData.value.updatedAt),
    icon: 'pi pi-flag',
    color: '#F59E0B'
  });

  return events;
});

const fetchCaseDetails = async () => {
  loading.value = true;
  try {
    const response = await api.get(`/cases/${route.params.id}`);
    caseData.value = response.data;
  } catch (error) {
    console.error('Error fetching case details:', error);
  } finally {
    loading.value = false;
  }
};

const fetchSpecialists = async () => {
  loadingSpecialists.value = true;
  try {
    const response = await api.get('/users', {
      params: { role: 'COMMISSION_MEMBER' }
    });
    specialists.value = response.data.map((user: any) => ({
      label: `${user.firstName} ${user.lastName}${user.specialty ? ` (${user.specialty})` : ''}`,
      value: user.id
    }));
  } catch (error) {
    console.error('Error fetching specialists:', error);
  } finally {
    loadingSpecialists.value = false;
  }
};

const assignToSpecialist = async () => {
  if (!assignData.value.specialistId) return;

  assigning.value = true;
  try {
    await api.patch(`/cases/${route.params.id}`, {
      assignedToId: assignData.value.specialistId,
      status: 'ASSIGNED_TO_SPECIALIST'
    });

    showAssignDialog.value = false;
    assignData.value = { specialistId: '', note: '' };
    await fetchCaseDetails();
  } catch (error) {
    console.error('Error assigning specialist:', error);
  } finally {
    assigning.value = false;
  }
};

const changeStatus = async () => {
  if (!statusData.value.newStatus) return;

  changingStatus.value = true;
  try {
    await api.patch(`/cases/${route.params.id}`, {
      status: statusData.value.newStatus
    });

    showStatusDialog.value = false;
    statusData.value = { newStatus: '', note: '' };
    await fetchCaseDetails();
  } catch (error) {
    console.error('Error changing status:', error);
  } finally {
    changingStatus.value = false;
  }
};

const getStatusLabel = (status: string) => {
  const option = statusOptions.find((o) => o.value === status);
  return option?.label || status;
};

const getStatusSeverity = (status: string) => {
  switch (status) {
    case 'PENDING_SECRETARIAT':
      return 'info';
    case 'ASSIGNED_TO_SPECIALIST':
    case 'UNDER_REVIEW':
      return 'warning';
    case 'PENDING_MEETING':
    case 'MEETING_SCHEDULED':
      return 'info';
    case 'PENDING_VERDICT':
      return 'warning';
    case 'VERDICT_ISSUED':
      return 'success';
    case 'ARCHIVED':
      return 'secondary';
    case 'REJECTED':
      return 'danger';
    default:
      return 'info';
  }
};

const getFamilyRelation = (relation: string) => {
  const relations: Record<string, string> = {
    SELF: 'خود بیمه‌شده',
    SPOUSE: 'همسر',
    CHILD: 'فرزند',
    FATHER: 'پدر',
    MOTHER: 'مادر'
  };
  return relations[relation] || relation;
};

const formatDate = (date: string) => {
  if (!date) return '-';
  return new Date(date).toLocaleDateString('fa-IR', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  });
};

const formatFileSize = (bytes: number) => {
  if (bytes < 1024) return bytes + ' بایت';
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' کیلوبایت';
  return (bytes / (1024 * 1024)).toFixed(1) + ' مگابایت';
};

const goBack = () => {
  router.push('/commission');
};

const getTimelineIcon = (action: string): string => {
  if (action.includes('ثبت')) return 'pi pi-plus';
  if (action.includes('ارجاع')) return 'pi pi-user';
  if (action.includes('تغییر وضعیت')) return 'pi pi-refresh';
  if (action.includes('جلسه')) return 'pi pi-calendar';
  if (action.includes('رأی') || action.includes('رای')) return 'pi pi-check';
  return 'pi pi-flag';
};

const getTimelineColor = (action: string): string => {
  if (action.includes('ثبت')) return '#22C55E';
  if (action.includes('ارجاع')) return '#3B82F6';
  if (action.includes('تغییر وضعیت')) return '#F59E0B';
  if (action.includes('جلسه')) return '#8B5CF6';
  if (action.includes('رأی') || action.includes('رای')) return '#10B981';
  if (action.includes('رد')) return '#EF4444';
  return '#6B7280';
};

// File preview helpers
const isImageFile = (fileType: string): boolean => {
  return Boolean(fileType && fileType.startsWith('image/'));
};

const isPdfFile = (fileType: string): boolean => {
  return fileType === 'application/pdf';
};

const getFileUrl = (fileName: string): string => {
  return `${import.meta.env.VITE_API_BASE_URL}/cases/files/${fileName}`;
};

const openPdfInNewTab = (fileName: string): void => {
  const url = getFileUrl(fileName);
  window.open(url, '_blank');
};

onMounted(() => {
  fetchCaseDetails();
  fetchSpecialists();
});
</script>

<style scoped>
.case-details-wrapper {
  padding: 2rem;
  max-width: 100%;
  width: 100%;
  margin: 0;
}

/* Highlight Box */
.highlight-box {
  background: linear-gradient(135deg, #eff6ff 0%, #dbeafe 100%);
  border: 2px solid #bfdbfe;
  border-radius: 0.75rem;
  padding: 1.5rem;
}

.number-value {
  font-size: 1.5rem;
  font-weight: bold;
  color: #1e40af;
}

/* Action Card */
.actions-card {
  background: linear-gradient(135deg, #f9fafb 0%, white 100%);
}

/* Card Title */
.custom-card-title {
  font-size: 1.125rem;
  font-weight: bold;
  color: #1f2937;
  margin-bottom: 1rem;
  padding-bottom: 0.5rem;
  border-bottom: 2px solid #1e40af;
  display: flex;
  align-items: center;
}

/* Data Row - HORIZONTAL LAYOUT */
.data-row {
  display: flex;
  align-items: flex-start;
  margin-bottom: 1rem;
  gap: 0.5rem;
}

.data-label {
  font-weight: 500;
  color: #374151;
  min-width: 140px;
  flex-shrink: 0;
}

.data-value {
  color: #111827;
  flex: 1;
}

/* Divider */
.custom-divider {
  height: 1px;
  background: linear-gradient(to left, transparent, #d1d5db, transparent);
  margin: 1.5rem 0;
}

/* Document Box */
.doc-box {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 1rem;
  background: #f9fafb;
  border-radius: 0.5rem;
  border: 1px solid #e5e7eb;
  transition: all 0.2s ease;
  margin-bottom: 0.5rem;
}

.doc-box:hover {
  background: white;
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
  transform: translateY(-2px);
}

/* Timeline */
.timeline-box {
  position: sticky;
  top: 1.5rem;
}

.custom-timeline :deep(.p-timeline-event-opposite) {
  display: none;
}

.timeline-dot {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2.5rem;
  height: 2.5rem;
  color: white;
  border-radius: 50%;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  font-size: 1.1rem;
}

.timeline-box-content {
  background: white;
  padding: 1rem;
  border-radius: 0.5rem;
  border: 1px solid #e5e7eb;
  margin-bottom: 1rem;
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
}

.timeline-event-title {
  font-weight: bold;
  color: #1f2937;
  margin-bottom: 0.25rem;
}

.timeline-event-date {
  font-size: 0.75rem;
  color: #6b7280;
  margin-bottom: 0.5rem;
}

.timeline-event-desc {
  font-size: 0.875rem;
  color: #4b5563;
}

/* Responsive */
@media (max-width: 992px) {
  .case-details-wrapper {
    padding: 1rem;
  }

  .timeline-box {
    position: relative;
    top: 0;
  }

  .data-row {
    flex-direction: column;
  }

  .data-label {
    min-width: auto;
  }
}
</style>
