import { IsString, IsNotEmpty, IsOptional, IsBoolean } from 'class-validator';

export class CreateCaseTypeDto {
  @IsString()
  @IsNotEmpty()
  name: string;

  @IsString()
  @IsOptional()
  description?: string;

  @IsBoolean()
  @IsOptional()
  isCentralCommission?: boolean;

  @IsBoolean()
  @IsOptional()
  isActive?: boolean;
}
