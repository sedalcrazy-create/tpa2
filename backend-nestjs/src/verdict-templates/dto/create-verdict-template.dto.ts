import { IsString, IsNotEmpty, IsOptional, IsBoolean, IsUUID, IsInt } from 'class-validator';

export class CreateVerdictTemplateDto {
  @IsString()
  @IsNotEmpty()
  title: string;

  @IsString()
  @IsOptional()
  description?: string;

  @IsUUID()
  @IsNotEmpty()
  caseTypeId: string;

  @IsInt()
  @IsOptional()
  sortOrder?: number;

  @IsBoolean()
  @IsOptional()
  isActive?: boolean;
}
