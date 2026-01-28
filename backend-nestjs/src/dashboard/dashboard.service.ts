import { Injectable } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { Case } from '../entities/case.entity';
import { CaseStatus } from '../common/enums/case-status.enum';

@Injectable()
export class DashboardService {
  constructor(
    @InjectRepository(Case)
    private caseRepository: Repository<Case>,
  ) {}

  async getStats() {
    const totalCases = await this.caseRepository.count();
    const completedCases = await this.caseRepository.count({ where: { status: CaseStatus.VERDICT_ISSUED } });
    const pendingCases = await this.caseRepository.count({ where: { status: CaseStatus.UNDER_REVIEW } });

    return { totalCases, completedCases, pendingCases, monthlyMeetings: 42 };
  }

  async getRecentCases() {
    const cases = await this.caseRepository.find({
      take: 5,
      order: { createdAt: 'DESC' },
      relations: ['insuredPerson', 'caseType'],
    });

    return cases.map(c => ({
      caseNumber: c.caseNumber,
      patientName: c.insuredPerson ? `${c.insuredPerson.firstName} ${c.insuredPerson.lastName}` : 'نامشخص',
      caseType: c.caseType?.name || 'نامشخص',
      status: this.translateStatus(c.status),
      createdAt: new Date(c.createdAt).toLocaleDateString('fa-IR'),
    }));
  }

  private translateStatus(status: string): string {
    const map = { 'REGISTERED': 'ثبت شده', 'IN_REVIEW': 'در حال بررسی', 'COMPLETED': 'تکمیل شده', 'WAITING_DOCUMENTS': 'منتظر مدارک', 'REJECTED': 'رد شده' };
    return map[status] || status;
  }
}
