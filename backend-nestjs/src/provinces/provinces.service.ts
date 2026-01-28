import { Injectable, NotFoundException } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { Province } from '../entities/province.entity';
import { CreateProvinceDto } from './dto/create-province.dto';
import { UpdateProvinceDto } from './dto/update-province.dto';

@Injectable()
export class ProvincesService {
  constructor(
    @InjectRepository(Province)
    private provinceRepository: Repository<Province>,
  ) {}

  async create(createProvinceDto: CreateProvinceDto) {
    const province = this.provinceRepository.create(createProvinceDto);
    return this.provinceRepository.save(province);
  }

  async findAll() {
    return this.provinceRepository.find({ where: { isActive: true } });
  }

  async findOne(id: string) {
    const province = await this.provinceRepository.findOne({ where: { id } });
    if (!province) {
      throw new NotFoundException('استان یافت نشد');
    }
    return province;
  }

  async update(id: string, updateProvinceDto: UpdateProvinceDto) {
    await this.findOne(id);
    await this.provinceRepository.update(id, updateProvinceDto);
    return this.findOne(id);
  }

  async remove(id: string) {
    const province = await this.findOne(id);
    await this.provinceRepository.remove(province);
    return { message: 'استان با موفقیت حذف شد' };
  }
}
