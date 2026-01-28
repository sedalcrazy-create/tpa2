import { PartialType } from '@nestjs/mapped-types';
import { CreateVerdictTemplateDto } from './create-verdict-template.dto';

export class UpdateVerdictTemplateDto extends PartialType(CreateVerdictTemplateDto) {}
