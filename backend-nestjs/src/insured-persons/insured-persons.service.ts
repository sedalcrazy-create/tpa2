import { Injectable, NotFoundException } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository, Like, Or } from 'typeorm';
import { InsuredPerson } from '../entities/insured-person.entity';
import { CreateInsuredPersonDto } from './dto/create-insured-person.dto';
import { UpdateInsuredPersonDto } from './dto/update-insured-person.dto';
import { TextNormalizer } from '../common/utils/text-normalizer';

@Injectable()
export class InsuredPersonsService {
  constructor(
    @InjectRepository(InsuredPerson)
    private insuredPersonRepository: Repository<InsuredPerson>,
  ) {}

  async create(createInsuredPersonDto: CreateInsuredPersonDto) {
    // Normalize text fields to convert Arabic characters to Persian
    const normalizedDto = TextNormalizer.normalizeObject(createInsuredPersonDto, [
      'firstName',
      'lastName',
      'address',
      'city',
      'officeLocation',
    ]);

    const person = this.insuredPersonRepository.create(normalizedDto);
    return this.insuredPersonRepository.save(person);
  }

  async findAll(filters?: {
    nationalId?: string;
    personnelCode?: string;
    firstName?: string;
    lastName?: string;
  }) {
    let where: any = {};

    // When searching by nationalId and personnelCode with same value (for search),
    // use OR logic to find by either field
    if (filters?.nationalId && filters?.personnelCode &&
        filters.nationalId === filters.personnelCode) {
      where = [
        { nationalId: Like(`%${filters.nationalId}%`) },
        { personnelCode: Like(`%${filters.personnelCode}%`) }
      ];
    } else {
      // Otherwise use AND logic for separate filters
      if (filters?.nationalId) {
        where.nationalId = Like(`%${filters.nationalId}%`);
      }
      if (filters?.personnelCode) {
        where.personnelCode = Like(`%${filters.personnelCode}%`);
      }
      if (filters?.firstName) {
        where.firstName = Like(`%${filters.firstName}%`);
      }
      if (filters?.lastName) {
        where.lastName = Like(`%${filters.lastName}%`);
      }
    }

    // If no filters provided, limit to 100 records
    const hasFilters = (Array.isArray(where) && where.length > 0) ||
                       (!Array.isArray(where) && Object.keys(where).length > 0);
    const take = hasFilters ? 1000 : 100;

    return this.insuredPersonRepository.find({
      where,
      order: { createdAt: 'DESC' },
      take,
    });
  }

  async findOne(id: string) {
    const person = await this.insuredPersonRepository.findOne({ where: { id }, relations: ['cases'] });
    if (!person) {
      throw new NotFoundException('بیمه‌شده یافت نشد');
    }
    return person;
  }

  async update(id: string, updateInsuredPersonDto: UpdateInsuredPersonDto) {
    await this.findOne(id);

    // Normalize text fields to convert Arabic characters to Persian
    const normalizedDto = TextNormalizer.normalizeObject(updateInsuredPersonDto, [
      'firstName',
      'lastName',
      'address',
      'city',
      'officeLocation',
    ]);

    await this.insuredPersonRepository.update(id, normalizedDto);
    return this.findOne(id);
  }

  async remove(id: string) {
    const person = await this.findOne(id);
    await this.insuredPersonRepository.remove(person);
    return { message: 'بیمه‌شده با موفقیت حذف شد' };
  }
}
