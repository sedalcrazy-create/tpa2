import { IsString, IsNotEmpty, IsBoolean, IsOptional } from 'class-validator';

export class CreateProvinceDto {
  @IsString()
  @IsNotEmpty()
  name: string;

  @IsString()
  @IsNotEmpty()
  code: string;

  @IsBoolean()
  @IsOptional()
  isActive?: boolean;
}
