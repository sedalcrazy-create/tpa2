export enum SocialWorkCaseType {
  INSTALLMENT_PAYMENT = 'INSTALLMENT_PAYMENT', // تقسیط بدهی
  CONSULTATION = 'CONSULTATION', // مشاوره مددکاری
  TRANSPORTATION = 'TRANSPORTATION', // ایاب ذهاب
  DISABILITY = 'DISABILITY', // معلولیت
  FINANCIAL_AID = 'FINANCIAL_AID', // کمک مالی و مساعده
  LOAN = 'LOAN', // وام
  DEATH_SUPPORT = 'DEATH_SUPPORT', // فوت
  MEDICAL_EQUIPMENT = 'MEDICAL_EQUIPMENT', // تجهیزات پزشکی
  PATIENT_VISIT = 'PATIENT_VISIT', // عیادت از بیمه شدگان
}

export const SocialWorkCaseTypeLabels: Record<SocialWorkCaseType, string> = {
  [SocialWorkCaseType.INSTALLMENT_PAYMENT]: 'تقسیط بدهی',
  [SocialWorkCaseType.CONSULTATION]: 'مشاوره مددکاری',
  [SocialWorkCaseType.TRANSPORTATION]: 'ایاب ذهاب',
  [SocialWorkCaseType.DISABILITY]: 'معلولیت',
  [SocialWorkCaseType.FINANCIAL_AID]: 'کمک مالی و مساعده',
  [SocialWorkCaseType.LOAN]: 'وام',
  [SocialWorkCaseType.DEATH_SUPPORT]: 'کمک هزینه فوت',
  [SocialWorkCaseType.MEDICAL_EQUIPMENT]: 'تجهیزات پزشکی',
  [SocialWorkCaseType.PATIENT_VISIT]: 'عیادت از بیمه‌شدگان',
};
