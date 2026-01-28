import { Injectable, NotFoundException } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { Case } from '../entities/case.entity';
import { InsuredPerson } from '../entities/insured-person.entity';
import * as moment from 'moment-jalaali';

@Injectable()
export class CaseNumberGeneratorService {
  constructor(
    @InjectRepository(Case)
    private caseRepository: Repository<Case>,
    @InjectRepository(InsuredPerson)
    private insuredPersonRepository: Repository<InsuredPerson>,
  ) {}

  /**
   * تولید شماره پرونده به فرمت: سال-کد استخدامی-شماره سریال
   * مثال: 1403-12345-00001
   *
   * @param insuredPersonId - شناسه بیمه‌شده
   * @returns شماره پرونده یکتا
   */
  async generateCaseNumber(insuredPersonId: string): Promise<string> {
    // دریافت سال شمسی جاری
    const persianYear = moment().format('jYYYY');

    // دریافت کد استخدامی بیمه‌شده
    const insuredPerson = await this.insuredPersonRepository.findOne({
      where: { id: insuredPersonId },
    });

    if (!insuredPerson) {
      throw new NotFoundException('بیمه‌شده یافت نشد');
    }

    const personnelCode = insuredPerson.personnelCode;

    // پیدا کردن آخرین شماره سریال برای این سال و کد استخدامی
    const serialNumber = await this.getNextSerialNumber(persianYear, personnelCode);

    // ساخت شماره پرونده نهایی
    return `${persianYear}-${personnelCode}-${serialNumber}`;
  }

  /**
   * دریافت شماره سریال بعدی برای سال و کد استخدامی مشخص
   */
  private async getNextSerialNumber(year: string, personnelCode: string): Promise<string> {
    // الگوی شماره پرونده برای جستجو
    const pattern = `${year}-${personnelCode}-%`;

    // پیدا کردن آخرین شماره پرونده با این الگو
    const lastCase = await this.caseRepository
      .createQueryBuilder('case')
      .where('case.caseNumber LIKE :pattern', { pattern })
      .orderBy('case.caseNumber', 'DESC')
      .getOne();

    let nextSerial = 1;

    if (lastCase && lastCase.caseNumber) {
      // استخراج شماره سریال از شماره پرونده
      const parts = lastCase.caseNumber.split('-');
      if (parts.length === 3) {
        const lastSerial = parseInt(parts[2], 10);
        nextSerial = lastSerial + 1;
      }
    }

    // تبدیل به فرمت 5 رقمی (00001)
    return nextSerial.toString().padStart(5, '0');
  }

  /**
   * بررسی یکتا بودن شماره پرونده
   */
  async isCaseNumberUnique(caseNumber: string): Promise<boolean> {
    const existingCase = await this.caseRepository.findOne({
      where: { caseNumber },
    });
    return !existingCase;
  }
}
