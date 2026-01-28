import api from './api';
import type {
  SocialWorkCase,
  CreateSocialWorkDto,
  UpdateAssessmentDto,
  GenerateReferralDto,
  ReferralLetter,
} from '../types/social-work';

export const socialWorkService = {
  // ثبت پرونده مددکاری جدید
  async createCase(data: CreateSocialWorkDto): Promise<SocialWorkCase> {
    const response = await api.post('/social-work', data);
    return response.data;
  },

  // لیست پرونده‌های مددکاری
  async getAllCases(): Promise<SocialWorkCase[]> {
    const response = await api.get('/social-work');
    return response.data;
  },

  // جزئیات یک پرونده
  async getCaseById(id: string): Promise<SocialWorkCase> {
    const response = await api.get(`/social-work/${id}`);
    return response.data;
  },

  // ثبت ارزیابی
  async updateAssessment(
    id: string,
    data: UpdateAssessmentDto
  ): Promise<SocialWorkCase> {
    const response = await api.patch(`/social-work/${id}/assessment`, data);
    return response.data;
  },

  // صدور معرفی‌نامه
  async generateReferralLetter(
    id: string,
    data: GenerateReferralDto
  ): Promise<ReferralLetter> {
    const response = await api.post(`/social-work/${id}/referral`, data);
    return response.data;
  },

  // حذف پرونده
  async deleteCase(id: string): Promise<void> {
    await api.delete(`/social-work/${id}`);
  },
};
