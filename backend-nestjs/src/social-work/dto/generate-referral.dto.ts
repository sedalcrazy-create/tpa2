import { IsNotEmpty, IsOptional, IsString } from 'class-validator';

export class GenerateReferralDto {
  @IsString()
  @IsOptional()
  referredTo?: string; // به کجا ارجاع شود (پیش‌فرض: حسابداری)

  @IsString()
  @IsOptional()
  additionalNotes?: string; // یادداشت‌های اضافی
}
