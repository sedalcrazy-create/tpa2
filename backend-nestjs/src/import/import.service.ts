import { Injectable } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import * as XLSX from 'xlsx';
import { InsuredPerson } from '../entities/insured-person.entity';
import { FamilyRelation } from '../common/enums/family-relation.enum';

@Injectable()
export class ImportService {
  constructor(
    @InjectRepository(InsuredPerson)
    private insuredPersonRepository: Repository<InsuredPerson>,
  ) {}

  async importExcel(buffer: Buffer): Promise<{
    success: boolean;
    message: string;
    imported: number;
    updated: number;
    errors: any[];
  }> {
    const workbook = XLSX.read(buffer, { type: 'buffer' });
    const sheetName = workbook.SheetNames[0];
    const sheet = workbook.Sheets[sheetName];
    const data = XLSX.utils.sheet_to_json(sheet);

    let imported = 0;
    let updated = 0;
    const errors = [];

    for (const row of data as any[]) {
      try {
        // Map Excel columns to database fields
        const nationalId = this.cleanNationalId(row['کد ملی']);
        const personnelCode = row['استخدامی']?.toString();
        const firstName = row['نام'];
        const lastName = row['نام خانوادگی'];
        const birthYear = row['سال تولد'];
        const birthMonth = row['ماه تولد'];
        const birthDay = row['روز تولد'];
        const familyRelation = this.mapFamilyRelation(row['نسبت']);
        const phone = row['شماره حساب']?.toString(); // Using account as phone placeholder
        const address = row['اداره امور'] || row['تفکیک محل'] || '';
        const employmentStatus = row['وضعیت خدمت همکار'] || '';

        // Validate required fields
        if (!nationalId || !personnelCode || !firstName || !lastName) {
          errors.push({
            row: row,
            error: 'Missing required fields',
          });
          continue;
        }

        // Convert Persian date to Gregorian (approximate)
        const birthDate = this.convertPersianToGregorian(
          birthYear,
          birthMonth,
          birthDay,
        );

        // Check if person already exists
        let person = await this.insuredPersonRepository.findOne({
          where: { nationalId },
        });

        if (person) {
          // Update existing person
          person.personnelCode = personnelCode;
          person.firstName = firstName;
          person.lastName = lastName;
          person.birthDate = birthDate;
          person.familyRelation = familyRelation;
          person.phone = phone;
          person.address = address;
          person.employmentStatus = employmentStatus;
          await this.insuredPersonRepository.save(person);
          updated++;
        } else {
          // Create new person
          person = this.insuredPersonRepository.create({
            nationalId,
            personnelCode,
            firstName,
            lastName,
            birthDate,
            familyRelation,
            phone,
            address,
            employmentStatus,
          });
          await this.insuredPersonRepository.save(person);
          imported++;
        }
      } catch (error) {
        errors.push({
          row: row,
          error: error.message,
        });
      }
    }

    return {
      success: true,
      message: `Import completed: ${imported} new records, ${updated} updated records`,
      imported,
      updated,
      errors: errors.slice(0, 100), // Return first 100 errors only
    };
  }

  private cleanNationalId(nationalId: string): string {
    if (!nationalId) return '';
    // Remove dashes and spaces
    return nationalId.toString().replace(/[-\s]/g, '');
  }

  private mapFamilyRelation(relation: string): FamilyRelation {
    const mapping: { [key: string]: FamilyRelation } = {
      'خود': FamilyRelation.SELF,
      'پسر': FamilyRelation.CHILD,
      'دختر': FamilyRelation.CHILD,
      'همسر': FamilyRelation.SPOUSE,
      'پدر': FamilyRelation.FATHER,
      'مادر': FamilyRelation.MOTHER,
    };
    return mapping[relation] || FamilyRelation.SELF;
  }

  private convertPersianToGregorian(
    year: number,
    month: number,
    day: number,
  ): Date {
    // Simple approximation: Persian year 1400 ≈ Gregorian 2021
    // For accurate conversion, use a proper library like moment-jalaali
    const gregorianYear = year + 621;

    // Approximate month conversion
    let gregorianMonth = month;
    let gregorianDay = day;

    // Adjust for Persian calendar differences
    if (month > 9) {
      gregorianMonth = month - 9;
      gregorianDay = day;
    } else {
      gregorianMonth = month + 3;
    }

    return new Date(gregorianYear, gregorianMonth - 1, gregorianDay);
  }
}
