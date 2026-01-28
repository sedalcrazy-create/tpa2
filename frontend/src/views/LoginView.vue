<script setup lang="ts">
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const username = ref('')
const password = ref('')
const isLoading = ref(false)
const error = ref('')

async function handleSubmit() {
  if (!username.value || !password.value) {
    error.value = 'لطفا نام کاربری و رمز عبور را وارد کنید'
    return
  }

  isLoading.value = true
  error.value = ''

  const success = await authStore.login(username.value, password.value)

  if (success) {
    const redirect = route.query.redirect as string || '/'
    router.push(redirect)
  } else {
    error.value = 'نام کاربری یا رمز عبور اشتباه است'
  }

  isLoading.value = false
}
</script>

<template>
  <div class="login-page">
    <div class="login-container">
      <div class="login-card">
        <div class="login-header">
          <div class="login-logo">
            <i class="bi bi-heart-pulse"></i>
          </div>
          <h1>سامانه TPA</h1>
          <p>مدیریت اسناد درمانی</p>
        </div>

        <form @submit.prevent="handleSubmit" class="login-form">
          <div v-if="error" class="alert alert-danger">
            <i class="bi bi-exclamation-circle"></i>
            {{ error }}
          </div>

          <div class="form-group">
            <label class="form-label">نام کاربری</label>
            <div class="input-wrapper">
              <i class="bi bi-person"></i>
              <input
                v-model="username"
                type="text"
                class="form-control"
                placeholder="نام کاربری را وارد کنید"
                :disabled="isLoading"
              />
            </div>
          </div>

          <div class="form-group">
            <label class="form-label">رمز عبور</label>
            <div class="input-wrapper">
              <i class="bi bi-lock"></i>
              <input
                v-model="password"
                type="password"
                class="form-control"
                placeholder="رمز عبور را وارد کنید"
                :disabled="isLoading"
              />
            </div>
          </div>

          <button type="submit" class="btn btn-primary w-100" :disabled="isLoading">
            <span v-if="isLoading" class="spinner"></span>
            <span v-else>
              <i class="bi bi-box-arrow-in-right"></i>
              ورود به سامانه
            </span>
          </button>
        </form>

        <div class="login-footer">
          <p>نسخه ۱.۰.۰ - تمامی حقوق محفوظ است</p>
        </div>
      </div>
    </div>

    <div class="login-side">
      <div class="side-content">
        <div class="side-icon">
          <i class="bi bi-clipboard2-pulse"></i>
        </div>
        <h2>سامانه جامع مدیریت اسناد درمانی</h2>
        <ul class="features-list">
          <li><i class="bi bi-check-circle-fill"></i> ثبت و پیگیری ادعاهای درمانی</li>
          <li><i class="bi bi-check-circle-fill"></i> ارزیابی و تایید اسناد</li>
          <li><i class="bi bi-check-circle-fill"></i> مدیریت مراکز درمانی طرف قرارداد</li>
          <li><i class="bi bi-check-circle-fill"></i> تسویه حساب و پرداخت</li>
          <li><i class="bi bi-check-circle-fill"></i> گزارشات تحلیلی و مدیریتی</li>
        </ul>
      </div>
    </div>
  </div>
</template>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  direction: rtl;
}

.login-container {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px;
  background: var(--bg-light);
}

.login-card {
  width: 100%;
  max-width: 420px;
  background: var(--bg-white);
  border-radius: 20px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.08);
  padding: 40px;
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.login-logo {
  width: 80px;
  height: 80px;
  margin: 0 auto 20px;
  background: linear-gradient(135deg, var(--primary) 0%, var(--accent) 100%);
  border-radius: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 10px 30px rgba(255, 107, 107, 0.35);
}

.login-logo i {
  font-size: 2.5rem;
  color: #fff;
}

.login-header h1 {
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--text-dark);
  margin-bottom: 8px;
}

.login-header p {
  color: var(--text-muted);
  font-size: 0.9rem;
}

.login-form {
  margin-bottom: 24px;
}

.input-wrapper {
  position: relative;
}

.input-wrapper i {
  position: absolute;
  right: 14px;
  top: 50%;
  transform: translateY(-50%);
  color: var(--text-muted);
  font-size: 1.1rem;
}

.input-wrapper .form-control {
  padding-right: 44px;
}

.login-footer {
  text-align: center;
  color: var(--text-muted);
  font-size: 0.8rem;
}

.login-side {
  width: 45%;
  background: linear-gradient(135deg, var(--primary) 0%, var(--primary-dark) 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px;
  position: relative;
  overflow: hidden;
}

.login-side::before {
  content: '';
  position: absolute;
  width: 300px;
  height: 300px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 50%;
  top: -100px;
  right: -100px;
}

.login-side::after {
  content: '';
  position: absolute;
  width: 400px;
  height: 400px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 50%;
  bottom: -150px;
  left: -150px;
}

.side-content {
  position: relative;
  z-index: 1;
  color: #fff;
  text-align: center;
  max-width: 400px;
}

.side-icon {
  width: 100px;
  height: 100px;
  margin: 0 auto 24px;
  background: rgba(255, 255, 255, 0.15);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.side-icon i {
  font-size: 3rem;
}

.side-content h2 {
  font-size: 1.5rem;
  font-weight: 700;
  margin-bottom: 32px;
}

.features-list {
  list-style: none;
  padding: 0;
  text-align: right;
}

.features-list li {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 0;
  font-size: 1rem;
  opacity: 0.9;
}

.features-list i {
  color: var(--accent);
}

.spinner {
  display: inline-block;
  width: 20px;
  height: 20px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

@media (max-width: 992px) {
  .login-side {
    display: none;
  }
}
</style>
