import {
  Controller,
  Get,
  Post,
  Body,
  Patch,
  Param,
  Delete,
  UseGuards,
  Query,
} from '@nestjs/common';
import { VerdictTemplatesService } from './verdict-templates.service';
import { CreateVerdictTemplateDto } from './dto/create-verdict-template.dto';
import { UpdateVerdictTemplateDto } from './dto/update-verdict-template.dto';
import { JwtAuthGuard } from '../auth/guards/jwt-auth.guard';

@Controller('verdict-templates')
@UseGuards(JwtAuthGuard)
export class VerdictTemplatesController {
  constructor(private readonly verdictTemplatesService: VerdictTemplatesService) {}

  @Post()
  create(@Body() createVerdictTemplateDto: CreateVerdictTemplateDto) {
    return this.verdictTemplatesService.create(createVerdictTemplateDto);
  }

  @Get()
  findAll(@Query('caseTypeId') caseTypeId?: string) {
    return this.verdictTemplatesService.findAll(caseTypeId);
  }

  @Get(':id')
  findOne(@Param('id') id: string) {
    return this.verdictTemplatesService.findOne(id);
  }

  @Patch(':id')
  update(@Param('id') id: string, @Body() updateVerdictTemplateDto: UpdateVerdictTemplateDto) {
    return this.verdictTemplatesService.update(id, updateVerdictTemplateDto);
  }

  @Delete(':id')
  remove(@Param('id') id: string) {
    return this.verdictTemplatesService.remove(id);
  }
}
