export enum CaseStatus {
  PENDING_SECRETARIAT = 'PENDING_SECRETARIAT', // در انتظار بررسی دبیرخانه
  ASSIGNED_TO_SPECIALIST = 'ASSIGNED_TO_SPECIALIST', // ارجاع به متخصص
  UNDER_REVIEW = 'UNDER_REVIEW', // در حال بررسی
  PENDING_MEETING = 'PENDING_MEETING', // در انتظار جلسه
  MEETING_SCHEDULED = 'MEETING_SCHEDULED', // جلسه زمان‌بندی شده
  PENDING_VERDICT = 'PENDING_VERDICT', // در انتظار رأی
  VERDICT_ISSUED = 'VERDICT_ISSUED', // رأی صادر شد
  ARCHIVED = 'ARCHIVED', // بایگانی شده
  REJECTED = 'REJECTED', // رد شده
}
