export interface InsuredPerson {
  id: string;
  nationalId: string;
  personnelCode: string;
  firstName: string;
  lastName: string;
  birthDate: string;
  familyRelation: string;
  insuranceNumber?: string;
  phone?: string;
  address?: string;
  employmentStatus?: string;
  officeLocation?: string;
  city?: string;
  provinceId?: string;
  createdAt: string;
  updatedAt: string;
}
