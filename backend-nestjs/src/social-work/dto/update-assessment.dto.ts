import { IsNotEmpty, IsString } from 'class-validator';

export class UpdateAssessmentDto {
  @IsString()
  @IsNotEmpty()
  assessmentReport: string; // گزارش ارزیابی مددکار
}
