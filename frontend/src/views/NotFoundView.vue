<script setup lang="ts">
import { RouterLink } from 'vue-router'
</script>

<template>
  <div class="not-found">
    <div class="not-found-content">
      <div class="not-found-code">۴۰۴</div>
      <h1>صفحه مورد نظر یافت نشد</h1>
      <p>صفحه‌ای که به دنبال آن هستید وجود ندارد یا منتقل شده است.</p>
      <RouterLink :to="{ name: 'dashboard' }" class="btn btn-primary">
        <i class="bi bi-house"></i>
        بازگشت به داشبورد
      </RouterLink>
    </div>
  </div>
</template>

<style scoped>
.not-found {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-light);
}

.not-found-content {
  text-align: center;
}

.not-found-code {
  font-size: 8rem;
  font-weight: 800;
  background: linear-gradient(135deg, var(--primary) 0%, var(--accent) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin-bottom: 20px;
}

.not-found h1 {
  font-size: 1.5rem;
  margin-bottom: 12px;
}

.not-found p {
  color: var(--text-muted);
  margin-bottom: 30px;
}
</style>
