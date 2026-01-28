export enum SocialWorkCaseType {
  INSTALLMENT_PAYMENT = 'INSTALLMENT_PAYMENT',
  CONSULTATION = 'CONSULTATION',
  TRANSPORTATION = 'TRANSPORTATION',
  DISABILITY = 'DISABILITY',
  FINANCIAL_AID = 'FINANCIAL_AID',
  LOAN = 'LOAN',
  DEATH_SUPPORT = 'DEATH_SUPPORT',
  MEDICAL_EQUIPMENT = 'MEDICAL_EQUIPMENT',
  PATIENT_VISIT = 'PATIENT_VISIT',
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

export enum SocialWorkCaseStatus {
  DRAFT = 'DRAFT',
  UNDER_ASSESSMENT = 'UNDER_ASSESSMENT',
  REFERRED = 'REFERRED',
  CLOSED = 'CLOSED',
}

export const SocialWorkCaseStatusLabels: Record<SocialWorkCaseStatus, string> = {
  [SocialWorkCaseStatus.DRAFT]: 'پیش‌نویس',
  [SocialWorkCaseStatus.UNDER_ASSESSMENT]: 'در حال ارزیابی',
  [SocialWorkCaseStatus.REFERRED]: 'ارجاع شده',
  [SocialWorkCaseStatus.CLOSED]: 'بسته شده',
};

export interface SocialWorkCase {
  id: string;
  caseNumber: string;
  caseType: SocialWorkCaseType;
  status: SocialWorkCaseStatus;
  insuredPersonId: string;
  insuredPerson?: any;
  socialWorkerId: string;
  socialWorker?: any;
  medicalCaseId?: string;
  medicalCase?: any;
  requestDetails?: any;
  assessmentReport?: string;
  assessedAt?: string;
  referredAt?: string;
  createdAt: string;
  updatedAt: string;
  closedAt?: string;
  referralLetters?: ReferralLetter[];
}

export interface ReferralLetter {
  id: string;
  letterNumber: string;
  socialWorkCaseId: string;
  content: string;
  pdfPath?: string;
  referredTo: string;
  generatedAt: string;
}

export interface CreateSocialWorkDto {
  caseType: SocialWorkCaseType;
  insuredPersonId: string;
  medicalCaseId?: string;
  requestDetails?: any;
}

export interface UpdateAssessmentDto {
  assessmentReport: string;
}

export interface GenerateReferralDto {
  referredTo?: string;
  additionalNotes?: string;
}
