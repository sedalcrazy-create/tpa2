import { Injectable, NotFoundException } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { VerdictTemplate } from '../entities/verdict-template.entity';
import { CreateVerdictTemplateDto } from './dto/create-verdict-template.dto';
import { UpdateVerdictTemplateDto } from './dto/update-verdict-template.dto';

@Injectable()
export class VerdictTemplatesService {
  constructor(
    @InjectRepository(VerdictTemplate)
    private verdictTemplateRepository: Repository<VerdictTemplate>,
  ) {}

  async create(createVerdictTemplateDto: CreateVerdictTemplateDto): Promise<VerdictTemplate> {
    const template = this.verdictTemplateRepository.create(createVerdictTemplateDto);
    return await this.verdictTemplateRepository.save(template);
  }

  async findAll(caseTypeId?: string): Promise<VerdictTemplate[]> {
    const where = caseTypeId ? { caseTypeId } : {};
    return await this.verdictTemplateRepository.find({
      where,
      order: { sortOrder: 'ASC', title: 'ASC' },
    });
  }

  async findOne(id: string): Promise<VerdictTemplate> {
    const template = await this.verdictTemplateRepository.findOne({ where: { id } });
    if (!template) {
      throw new NotFoundException(`رای کمیسیون با شناسه ${id} یافت نشد`);
    }
    return template;
  }

  async update(id: string, updateVerdictTemplateDto: UpdateVerdictTemplateDto): Promise<VerdictTemplate> {
    const template = await this.findOne(id);
    Object.assign(template, updateVerdictTemplateDto);
    return await this.verdictTemplateRepository.save(template);
  }

  async remove(id: string): Promise<void> {
    const template = await this.findOne(id);
    await this.verdictTemplateRepository.remove(template);
  }
}
