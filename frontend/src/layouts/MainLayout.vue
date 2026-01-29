<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink, RouterView, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const authStore = useAuthStore()

const currentDate = computed(() => {
  return new Intl.DateTimeFormat('fa-IR', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  }).format(new Date())
})

const pageTitle = computed(() => {
  const titles: Record<string, string> = {
    dashboard: 'داشبورد',
    claims: 'مدیریت ادعاها',
    packages: 'بسته‌های اسناد',
    prescriptions: 'نسخه‌های پزشکی',
    centers: 'مراکز درمانی',
    settlements: 'تسویه حساب',
    members: 'استعلام بیمه‌شدگان',
    reports: 'گزارشات',
    users: 'مدیریت کاربران',
    settings: 'تنظیمات',
    // Commission
    'commission-cases': 'پرونده‌های کمیسیون',
    'commission-case-detail': 'جزئیات پرونده کمیسیون',
    'case-types': 'انواع پرونده',
    'verdict-templates': 'قالب آرای کمیسیون',
    'insured-persons': 'بیمه‌شدگان',
    // Social Work
    'social-work-list': 'پرونده‌های مددکاری',
    'social-work-create': 'ایجاد پرونده مددکاری',
    'social-work-detail': 'جزئیات پرونده مددکاری',
    // Base Data - Drugs & Services
    'drugs': 'بانک دارویی',
    'services': 'خدمات درمانی',
    'drug-prices': 'قیمت دارو',
    'service-prices': 'قیمت خدمات',
    // Base Data - Centers & Contracts
    'doctors': 'پزشکان',
    'contracts': 'قراردادها',
    // Base Data - Pricing & Rules
    'price-conditions': 'شرایط قیمت‌گذاری',
    'insurance-rules': 'قوانین بیمه',
    'tariffs': 'تعرفه‌ها',
    // Base Data - Personnel
    'employees': 'کارمندان',
    'employee-sync': 'به‌روزرسانی کارمندان',
    // Base Data - Geographic
    'provinces': 'استان‌ها',
    'cities': 'شهرها',
    // Base Data - Medical
    'diagnoses': 'تشخیص‌ها (ICD-10)'
  }
  return titles[route.name as string] || 'داشبورد'
})

const menuItems = [
  { section: 'منوی اصلی', items: [
    { name: 'dashboard', title: 'داشبورد', icon: '🏠' }
  ]},
  { section: 'عملیات اسناد', items: [
    { name: 'claims', title: 'ادعاهای درمانی', icon: '🏥' },
    { name: 'packages', title: 'بسته‌های اسناد', icon: '📦' },
    { name: 'prescriptions', title: 'نسخه‌های پزشکی', icon: '💊' }
  ]},
  { section: 'مراکز و مالی', items: [
    { name: 'centers', title: 'مراکز درمانی', icon: '🏥' },
    { name: 'settlements', title: 'تسویه حساب', icon: '💰' }
  ]},
  { section: 'اطلاعات پایه', items: [
    { name: 'price-conditions', title: 'شرایط قیمت‌گذاری', icon: '💰' },
    { name: 'insurance-rules', title: 'قوانین بیمه', icon: '🛡️' },
    { name: 'contracts', title: 'قراردادها', icon: '📄' },
    { name: 'prescriptions', title: 'نسخه‌ها', icon: '💊' },
    { name: 'employee-sync', title: 'به‌روزرسانی کارمندان', icon: '🔄', roles: ['system_admin', 'insurer_admin'] }
  ]},
  { section: 'کمیسیون پزشکی', items: [
    { name: 'commission-cases', title: 'پرونده‌های کمیسیون', icon: '📋' },
    { name: 'insured-persons', title: 'بیمه‌شدگان', icon: '👥' },
    { name: 'case-types', title: 'انواع پرونده', icon: '📑', roles: ['system_admin'] },
    { name: 'verdict-templates', title: 'قالب آرا', icon: '📝', roles: ['system_admin'] }
  ]},
  { section: 'مددکاری', items: [
    { name: 'social-work-list', title: 'پرونده‌های مددکاری', icon: '❤️' },
    { name: 'social-work-create', title: 'ایجاد پرونده جدید', icon: '➕' }
  ]},
  { section: 'گزارشات', items: [
    { name: 'members', title: 'استعلام بیمه‌شده', icon: '🔍' },
    { name: 'reports', title: 'گزارشات', icon: '📊' }
  ]},
  { section: 'مدیریت', items: [
    { name: 'users', title: 'کاربران', icon: '👤', roles: ['system_admin', 'insurer_admin'] },
    { name: 'settings', title: 'تنظیمات', icon: '⚙️', roles: ['system_admin', 'insurer_admin'] }
  ]}
]

function canShowItem(item: { roles?: string[] }): boolean {
  if (!item.roles) return true
  return authStore.hasRole(item.roles)
}

async function handleLogout() {
  await authStore.logout()
}
</script>

<template>
  <div class="app-layout">
    <!-- Sidebar -->
    <nav class="sidebar">
      <div class="sidebar-brand">
        <div class="brand-logo">
          <img src="/logo.jpg" alt="TPA Logo" />
        </div>
        <div class="brand-text">
          <div class="brand-title">سامانه TPA</div>
          <div class="brand-subtitle">مدیریت اسناد درمانی</div>
        </div>
      </div>

      <div class="sidebar-nav">
        <template v-for="section in menuItems" :key="section.section">
          <div class="nav-section" v-if="section.items.some(canShowItem)">
            <div class="nav-section-title">{{ section.section }}</div>
            <div class="nav-item" v-for="item in section.items" :key="item.name">
              <RouterLink
                v-if="canShowItem(item)"
                :to="{ name: item.name }"
                class="nav-link"
                :class="{ active: route.name === item.name }"
              >
                <span class="nav-icon">{{ item.icon }}</span>
                <span>{{ item.title }}</span>
              </RouterLink>
            </div>
          </div>
        </template>
      </div>

      <div class="sidebar-footer">
        <div class="user-card">
          <div class="user-avatar">{{ authStore.initials }}</div>
          <div class="user-info">
            <div class="user-name">{{ authStore.fullName }}</div>
            <div class="user-role">{{ authStore.user?.roleTitle }}</div>
          </div>
        </div>
        <button class="btn-logout" @click="handleLogout">
          <i class="bi bi-box-arrow-right"></i>
          <span>خروج از حساب</span>
        </button>
      </div>
    </nav>

    <!-- Main Wrapper -->
    <div class="main-wrapper">
      <!-- Top Header -->
      <header class="top-header">
        <div class="header-right">
          <div class="breadcrumb">
            <RouterLink :to="{ name: 'dashboard' }">
              <i class="bi bi-house"></i>
            </RouterLink>
            <span class="breadcrumb-separator">/</span>
            <span class="breadcrumb-current">{{ pageTitle }}</span>
          </div>
        </div>
        <div class="header-left">
          <div class="header-date">
            <i class="bi bi-calendar3"></i>
            <span>{{ currentDate }}</span>
          </div>
        </div>
      </header>

      <!-- Main Content -->
      <main class="main-content">
        <RouterView />
      </main>
    </div>
  </div>
</template>

<style scoped>
.app-layout {
  display: flex;
  min-height: 100vh;
}
</style>
