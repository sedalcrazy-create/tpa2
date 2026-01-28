import { IsEnum, IsNotEmpty, IsOptional, IsString, IsUUID } from 'class-validator';
import { SocialWorkCaseType } from '../../common/enums/social-work-case-type.enum';

export class CreateSocialWorkDto {
  @IsEnum(SocialWorkCaseType)
  @IsNotEmpty()
  caseType: SocialWorkCaseType;

  @IsUUID()
  @IsNotEmpty()
  insuredPersonId: string;

  @IsUUID()
  @IsOptional()
  medicalCaseId?: string;

  @IsOptional()
  requestDetails?: any; // JSON object با جزئیات درخواست
}
