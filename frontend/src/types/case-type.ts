export interface CaseType {
  id: string;
  name: string;
  description?: string;
  isActive: boolean;
  isCentralCommission: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface CreateCaseTypeDto {
  name: string;
  description?: string;
  isActive?: boolean;
  isCentralCommission?: boolean;
}

export interface UpdateCaseTypeDto {
  name?: string;
  description?: string;
  isActive?: boolean;
  isCentralCommission?: boolean;
}
