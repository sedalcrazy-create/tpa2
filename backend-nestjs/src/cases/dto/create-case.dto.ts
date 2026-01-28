import { IsString, IsNotEmpty, IsEnum, IsOptional, IsUUID } from 'class-validator';
import { CaseStatus } from '../../common/enums/case-status.enum';
import { CommissionLevel } from '../../common/enums/commission-level.enum';

export class CreateCaseDto {
  // caseNumber حذف شد - سیستم خودکار آن را تولید می‌کند

  @IsUUID()
  @IsNotEmpty()
  insuredPersonId: string;

  @IsUUID()
  @IsOptional()
  caseTypeId?: string;

  @IsUUID()
  @IsOptional()
  verdictTemplateId?: string;

  @IsEnum(CaseStatus)
  @IsOptional()
  status?: CaseStatus;

  @IsEnum(CommissionLevel)
  @IsNotEmpty()
  commissionLevel: CommissionLevel;

  @IsUUID()
  @IsNotEmpty()
  provinceId: string; // اجباری برای تولید شماره پرونده

  @IsString()
  @IsOptional()
  description?: string;

  @IsString()
  @IsOptional()
  medicalHistory?: string;

  @IsUUID()
  @IsOptional()
  assignedToId?: string;
}
