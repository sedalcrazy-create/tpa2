<script setup lang="ts">
// Placeholder
</script>

<template>
  <div class="settings-view">
    <div class="page-header">
      <h1><i class="bi bi-gear"></i> تنظیمات</h1>
    </div>

    <div class="row">
      <div class="col-6">
        <div class="card">
          <div class="card-header">
            <i class="bi bi-sliders"></i>
            تنظیمات عمومی
          </div>
          <div class="card-body">
            <div class="form-group">
              <label class="form-label">نام سامانه</label>
              <input type="text" class="form-control" value="سامانه TPA" />
            </div>
            <div class="form-group">
              <label class="form-label">تعداد آیتم در صفحه</label>
              <select class="form-control">
                <option>10</option>
                <option selected>20</option>
                <option>50</option>
              </select>
            </div>
            <button class="btn btn-primary">ذخیره تغییرات</button>
          </div>
        </div>
      </div>
      <div class="col-6">
        <div class="card">
          <div class="card-header">
            <i class="bi bi-database"></i>
            اطلاعات سیستم
          </div>
          <div class="card-body">
            <div class="info-row">
              <span>نسخه سامانه</span>
              <span>۱.۰.۰</span>
            </div>
            <div class="info-row">
              <span>آخرین به‌روزرسانی</span>
              <span>۱۴۰۳/۱۰/۲۶</span>
            </div>
            <div class="info-row">
              <span>تعداد کاربران</span>
              <span>۱۵</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.info-row {
  display: flex;
  justify-content: space-between;
  padding: 12px 0;
  border-bottom: 1px solid var(--border-color);
}

.info-row:last-child {
  border-bottom: none;
}
</style>
