import { IsString, IsNotEmpty, IsEnum, IsOptional, IsDateString, IsUUID } from 'class-validator';
import { FamilyRelation } from '../../common/enums/family-relation.enum';

export class CreateInsuredPersonDto {
  @IsString()
  @IsNotEmpty()
  nationalId: string;

  @IsString()
  @IsNotEmpty()
  personnelCode: string;

  @IsString()
  @IsNotEmpty()
  firstName: string;

  @IsString()
  @IsNotEmpty()
  lastName: string;

  @IsDateString()
  @IsNotEmpty()
  birthDate: Date;

  @IsEnum(FamilyRelation)
  @IsOptional()
  familyRelation?: FamilyRelation;

  @IsString()
  @IsOptional()
  insuranceNumber?: string;

  @IsString()
  @IsOptional()
  phone?: string;

  @IsString()
  @IsOptional()
  address?: string;

  @IsString()
  @IsOptional()
  employmentStatus?: string;

  @IsString()
  @IsOptional()
  officeLocation?: string; // اداره امور / محل خدمت

  @IsString()
  @IsOptional()
  city?: string; // شهر محل خدمت

  @IsUUID()
  @IsOptional()
  provinceId?: string; // واحد کمیسیون پزشکی
}
