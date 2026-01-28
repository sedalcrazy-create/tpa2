import { PartialType } from '@nestjs/mapped-types';
import { IsOptional, IsDate } from 'class-validator';
import { Type } from 'class-transformer';
import { CreateCaseDto } from './create-case.dto';

export class UpdateCaseDto extends PartialType(CreateCaseDto) {
  @IsOptional()
  @IsDate()
  @Type(() => Date)
  assignedAt?: Date;
}
