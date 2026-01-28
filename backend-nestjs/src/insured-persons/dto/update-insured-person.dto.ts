import { PartialType } from '@nestjs/mapped-types';
import { CreateInsuredPersonDto } from './create-insured-person.dto';

export class UpdateInsuredPersonDto extends PartialType(CreateInsuredPersonDto) {}
