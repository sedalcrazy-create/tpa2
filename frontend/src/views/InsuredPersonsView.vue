<template>
  <div>
    <div class="flex justify-content-between align-items-center mb-4">
      <h2>مدیریت بیمه‌شدگان</h2>
      <Button label="بیمه‌شده جدید" icon="pi pi-plus" @click="showCreateDialog" />
    </div>
    <Card>
      <template #content>
        <div class="grid">
          <div class="col-3">
            <label class="block mb-2">جستجو بر اساس کد ملی</label>
            <InputText v-model="searchNationalId" placeholder="کد ملی را وارد کنید" class="w-full" @input="filterPersons" />
          </div>
          <div class="col-3">
            <label class="block mb-2">جستجو بر اساس کد پرسنلی</label>
            <InputText v-model="searchPersonnelCode" placeholder="کد پرسنلی را وارد کنید" class="w-full" @input="filterPersons" />
          </div>
          <div class="col-3">
            <label class="block mb-2">جستجو بر اساس نام</label>
            <InputText v-model="searchFirstName" placeholder="نام را وارد کنید" class="w-full" @input="filterPersons" />
          </div>
          <div class="col-3">
            <label class="block mb-2">جستجو بر اساس نام خانوادگی</label>
            <InputText v-model="searchLastName" placeholder="نام خانوادگی را وارد کنید" class="w-full" @input="filterPersons" />
          </div>
        </div>
        <DataTable :value="filteredPersons" :loading="loading" paginator :rows="10" stripedRows>
          <Column field="nationalId" header="کد ملی" sortable />
          <Column field="personnelCode" header="کد پرسنلی" sortable />
          <Column field="firstName" header="نام" sortable />
          <Column field="lastName" header="نام خانوادگی" sortable />
          <Column field="familyRelation" header="نسبت" sortable>
            <template #body="{ data }">
              {{ translateRelation(data.familyRelation) }}
            </template>
          </Column>
          <Column header="واحد کمیسیون پزشکی" sortable>
            <template #body="{ data }">
              {{ getProvinceName(data.provinceId) }}
            </template>
          </Column>
          <Column field="officeLocation" header="اداره امور" />
          <Column header="عملیات" style="width: 120px">
            <template #body="{ data }">
              <Button
                icon="pi pi-eye"
                severity="info"
                text
                rounded
                @click="viewPerson(data)"
                v-tooltip.top="'مشاهده جزئیات'"
              />
              <Button
                icon="pi pi-pencil"
                severity="warning"
                text
                rounded
                @click="editPerson(data)"
                v-tooltip.top="'ویرایش'"
              />
            </template>
          </Column>
        </DataTable>
      </template>
    </Card>

    <!-- View Dialog -->
    <Dialog v-model:visible="viewDialogVisible" header="جزئیات بیمه‌شده" :modal="true" :style="{ width: '50vw' }">
      <div v-if="selectedPerson" class="grid">
        <div class="col-6">
          <strong>کد ملی:</strong> {{ selectedPerson.nationalId }}
        </div>
        <div class="col-6">
          <strong>کد پرسنلی:</strong> {{ selectedPerson.personnelCode }}
        </div>
        <div class="col-6">
          <strong>نام:</strong> {{ selectedPerson.firstName }}
        </div>
        <div class="col-6">
          <strong>نام خانوادگی:</strong> {{ selectedPerson.lastName }}
        </div>
        <div class="col-6">
          <strong>تاریخ تولد:</strong> {{ formatDate(selectedPerson.birthDate) }}
        </div>
        <div class="col-6">
          <strong>نسبت:</strong> {{ translateRelation(selectedPerson.familyRelation) }}
        </div>
        <div class="col-6">
          <strong>شماره بیمه:</strong> {{ selectedPerson.insuranceNumber || '-' }}
        </div>
        <div class="col-6">
          <strong>تلفن:</strong> {{ selectedPerson.phone || '-' }}
        </div>
        <div class="col-6">
          <strong>وضعیت خدمت:</strong> {{ selectedPerson.employmentStatus || '-' }}
        </div>
        <div class="col-6">
          <strong>واحد کمیسیون پزشکی:</strong> {{ getProvinceName(selectedPerson.provinceId) }}
        </div>
        <div class="col-6">
          <strong>اداره امور:</strong> {{ selectedPerson.officeLocation || '-' }}
        </div>
        <div class="col-6">
          <strong>شهر محل خدمت:</strong> {{ selectedPerson.city || '-' }}
        </div>
        <div class="col-12">
          <strong>آدرس:</strong> {{ selectedPerson.address || '-' }}
        </div>
      </div>
      <template #footer>
        <Button label="بستن" icon="pi pi-times" @click="viewDialogVisible = false" />
      </template>
    </Dialog>

    <!-- Edit Dialog -->
    <Dialog v-model:visible="editDialogVisible" header="ویرایش بیمه‌شده" :modal="true" :style="{ width: '50vw' }">
      <div v-if="selectedPerson" class="grid">
        <div class="col-6">
          <label>نام</label>
          <InputText v-model="selectedPerson.firstName" class="w-full" />
        </div>
        <div class="col-6">
          <label>نام خانوادگی</label>
          <InputText v-model="selectedPerson.lastName" class="w-full" />
        </div>
        <div class="col-6">
          <label>شماره بیمه</label>
          <InputText v-model="selectedPerson.insuranceNumber" class="w-full" />
        </div>
        <div class="col-6">
          <label>تلفن</label>
          <InputText v-model="selectedPerson.phone" class="w-full" />
        </div>
        <div class="col-6">
          <label>وضعیت خدمت</label>
          <InputText v-model="selectedPerson.employmentStatus" class="w-full" />
        </div>
        <div class="col-6">
          <label>واحد کمیسیون پزشکی</label>
          <Dropdown
            v-model="selectedPerson.provinceId"
            :options="provinces"
            optionLabel="name"
            optionValue="id"
            placeholder="واحد کمیسیون پزشکی را انتخاب کنید"
            class="w-full"
          />
        </div>
        <div class="col-6">
          <label>اداره امور</label>
          <InputText v-model="selectedPerson.officeLocation" class="w-full" placeholder="اداره امور / محل خدمت" />
        </div>
        <div class="col-6">
          <label>شهر محل خدمت</label>
          <InputText v-model="selectedPerson.city" class="w-full" placeholder="برای شعب مستقل" />
        </div>
        <div class="col-12">
          <label>آدرس</label>
          <InputText v-model="selectedPerson.address" class="w-full" />
        </div>
      </div>
      <template #footer>
        <Button label="لغو" icon="pi pi-times" @click="editDialogVisible = false" text />
        <Button label="ذخیره" icon="pi pi-check" @click="savePerson" :loading="saving" />
      </template>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import Card from 'primevue/card';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import Button from 'primevue/button';
import Dialog from 'primevue/dialog';
import InputText from 'primevue/inputtext';
import Dropdown from 'primevue/dropdown';
import { useToast } from 'primevue/usetoast';
import api from '../services/api';
import type { InsuredPerson } from '../types/insured-person';
import type { Province } from '../types/province';

const toast = useToast();
const persons = ref<InsuredPerson[]>([]);
const filteredPersons = ref<InsuredPerson[]>([]);
const provinces = ref<Province[]>([]);
const loading = ref(false);
const saving = ref(false);
const viewDialogVisible = ref(false);
const editDialogVisible = ref(false);
const selectedPerson = ref<InsuredPerson | null>(null);
const searchNationalId = ref('');
const searchPersonnelCode = ref('');
const searchFirstName = ref('');
const searchLastName = ref('');

const fetchPersons = async () => {
  loading.value = true;
  try {
    const params: any = {};

    if (searchNationalId.value.trim()) {
      params.nationalId = searchNationalId.value.trim();
    }
    if (searchPersonnelCode.value.trim()) {
      params.personnelCode = searchPersonnelCode.value.trim();
    }
    if (searchFirstName.value.trim()) {
      params.firstName = searchFirstName.value.trim();
    }
    if (searchLastName.value.trim()) {
      params.lastName = searchLastName.value.trim();
    }

    const response = await api.get('/insured-persons', { params });
    persons.value = response.data;
    filteredPersons.value = response.data;
  } catch (error: any) {
    toast.add({
      severity: 'error',
      summary: 'خطا',
      detail: error.response?.data?.message || 'خطا در دریافت اطلاعات',
      life: 3000,
    });
  } finally {
    loading.value = false;
  }
};

const filterPersons = () => {
  // Now search is done on server side
  fetchPersons();
};

const viewPerson = (person: any) => {
  selectedPerson.value = { ...person };
  viewDialogVisible.value = true;
};

const editPerson = (person: any) => {
  selectedPerson.value = { ...person };
  editDialogVisible.value = true;
};

const savePerson = async () => {
  if (!selectedPerson.value) return;

  saving.value = true;
  try {
    await api.patch(`/insured-persons/${selectedPerson.value.id}`, selectedPerson.value);
    toast.add({
      severity: 'success',
      summary: 'موفقیت',
      detail: 'اطلاعات با موفقیت ذخیره شد',
      life: 3000,
    });
    editDialogVisible.value = false;
    await fetchPersons();
  } catch (error: any) {
    toast.add({
      severity: 'error',
      summary: 'خطا',
      detail: error.response?.data?.message || 'خطا در ذخیره اطلاعات',
      life: 3000,
    });
  } finally {
    saving.value = false;
  }
};

const showCreateDialog = () => {
  toast.add({
    severity: 'info',
    summary: 'اطلاع',
    detail: 'این قابلیت به زودی اضافه می‌شود',
    life: 3000,
  });
};

const formatDate = (date: string) => {
  if (!date) return '-';
  return new Date(date).toLocaleDateString('fa-IR');
};

const translateRelation = (relation: string) => {
  const relations: { [key: string]: string } = {
    'SELF': 'خود',
    'SPOUSE': 'همسر',
    'CHILD': 'فرزند',
    'FATHER': 'پدر',
    'MOTHER': 'مادر',
    'OTHER': 'سایر',
  };
  return relations[relation] || relation;
};

const fetchProvinces = async () => {
  try {
    const response = await api.get('/provinces');
    provinces.value = response.data;
  } catch (error: any) {
    console.error('Error fetching provinces:', error);
  }
};

const getProvinceName = (provinceId: string | undefined) => {
  if (!provinceId) return '-';
  const province = provinces.value.find((p) => p.id === provinceId);
  return province ? province.name : '-';
};

onMounted(() => {
  fetchPersons();
  fetchProvinces();
});
</script>
