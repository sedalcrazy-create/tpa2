import { Injectable, NotFoundException } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { CaseType } from '../entities/case-type.entity';
import { CreateCaseTypeDto } from './dto/create-case-type.dto';
import { UpdateCaseTypeDto } from './dto/update-case-type.dto';

@Injectable()
export class CaseTypesService {
  constructor(
    @InjectRepository(CaseType)
    private caseTypeRepository: Repository<CaseType>,
  ) {}

  async create(createCaseTypeDto: CreateCaseTypeDto): Promise<CaseType> {
    const caseType = this.caseTypeRepository.create(createCaseTypeDto);
    return await this.caseTypeRepository.save(caseType);
  }

  async findAll(): Promise<CaseType[]> {
    return await this.caseTypeRepository.find({
      order: { name: 'ASC' },
    });
  }

  async findOne(id: string): Promise<CaseType> {
    const caseType = await this.caseTypeRepository.findOne({ where: { id } });
    if (!caseType) {
      throw new NotFoundException(`نوع پرونده با شناسه ${id} یافت نشد`);
    }
    return caseType;
  }

  async update(id: string, updateCaseTypeDto: UpdateCaseTypeDto): Promise<CaseType> {
    const caseType = await this.findOne(id);
    Object.assign(caseType, updateCaseTypeDto);
    return await this.caseTypeRepository.save(caseType);
  }

  async remove(id: string): Promise<void> {
    const caseType = await this.findOne(id);
    await this.caseTypeRepository.remove(caseType);
  }
}
