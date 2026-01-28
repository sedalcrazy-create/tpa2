import { PartialType } from '@nestjs/mapped-types';
import { CreateSocialWorkDto } from './create-social-work.dto';

export class UpdateSocialWorkDto extends PartialType(CreateSocialWorkDto) {}
