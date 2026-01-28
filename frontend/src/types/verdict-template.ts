export interface VerdictTemplate {
  id: string;
  title: string;
  description?: string;
  isActive: boolean;
  sortOrder: number;
  caseTypeId: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateVerdictTemplateDto {
  title: string;
  description?: string;
  caseTypeId: string;
  sortOrder?: number;
  isActive?: boolean;
}

export interface UpdateVerdictTemplateDto {
  title?: string;
  description?: string;
  sortOrder?: number;
  isActive?: boolean;
}
