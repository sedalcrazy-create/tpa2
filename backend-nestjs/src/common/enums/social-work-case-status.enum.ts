export enum SocialWorkCaseStatus {
  DRAFT = 'DRAFT', // پیش‌نویس
  UNDER_ASSESSMENT = 'UNDER_ASSESSMENT', // در حال ارزیابی
  REFERRED = 'REFERRED', // ارجاع شده به حسابداری
  CLOSED = 'CLOSED', // بسته شده
}

export const SocialWorkCaseStatusLabels: Record<SocialWorkCaseStatus, string> = {
  [SocialWorkCaseStatus.DRAFT]: 'پیش‌نویس',
  [SocialWorkCaseStatus.UNDER_ASSESSMENT]: 'در حال ارزیابی',
  [SocialWorkCaseStatus.REFERRED]: 'ارجاع شده',
  [SocialWorkCaseStatus.CLOSED]: 'بسته شده',
};
