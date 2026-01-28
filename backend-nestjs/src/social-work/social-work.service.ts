import { Injectable, NotFoundException } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { CreateSocialWorkDto } from './dto/create-social-work.dto';
import { UpdateAssessmentDto } from './dto/update-assessment.dto';
import { GenerateReferralDto } from './dto/generate-referral.dto';
import { SocialWorkCase } from '../entities/social-work-case.entity';
import { ReferralLetter } from '../entities/referral-letter.entity';
import { SocialWorkCaseStatus } from '../common/enums/social-work-case-status.enum';
import { SocialWorkCaseTypeLabels } from '../common/enums/social-work-case-type.enum';

@Injectable()
export class SocialWorkService {
  constructor(
    @InjectRepository(SocialWorkCase)
    private socialWorkCaseRepository: Repository<SocialWorkCase>,
    @InjectRepository(ReferralLetter)
    private referralLetterRepository: Repository<ReferralLetter>,
  ) {}

  async create(
    createSocialWorkDto: CreateSocialWorkDto,
    userId: string,
  ): Promise<SocialWorkCase> {
    // تولید شماره پرونده
    const count = await this.socialWorkCaseRepository.count();
    const caseNumber = `MC-SW-${new Date().getFullYear()}-${String(count + 1).padStart(4, '0')}`;

    const socialWorkCase = this.socialWorkCaseRepository.create({
      ...createSocialWorkDto,
      caseNumber,
      socialWorkerId: userId,
      status: SocialWorkCaseStatus.DRAFT,
    });

    return this.socialWorkCaseRepository.save(socialWorkCase);
  }

  async findAll(): Promise<SocialWorkCase[]> {
    return this.socialWorkCaseRepository.find({
      relations: ['insuredPerson', 'socialWorker', 'medicalCase', 'referralLetters'],
      order: { createdAt: 'DESC' },
    });
  }

  async findOne(id: string): Promise<SocialWorkCase> {
    const socialWorkCase = await this.socialWorkCaseRepository.findOne({
      where: { id },
      relations: ['insuredPerson', 'socialWorker', 'medicalCase', 'referralLetters'],
    });

    if (!socialWorkCase) {
      throw new NotFoundException(`پرونده مددکاری با شناسه ${id} یافت نشد`);
    }

    return socialWorkCase;
  }

  async updateAssessment(
    id: string,
    updateAssessmentDto: UpdateAssessmentDto,
  ): Promise<SocialWorkCase> {
    const socialWorkCase = await this.findOne(id);

    socialWorkCase.assessmentReport = updateAssessmentDto.assessmentReport;
    socialWorkCase.assessedAt = new Date();
    socialWorkCase.status = SocialWorkCaseStatus.UNDER_ASSESSMENT;

    return this.socialWorkCaseRepository.save(socialWorkCase);
  }

  async generateReferralLetter(
    id: string,
    generateReferralDto: GenerateReferralDto,
  ): Promise<ReferralLetter> {
    const socialWorkCase = await this.findOne(id);

    if (!socialWorkCase.assessmentReport) {
      throw new Error('ابتدا باید گزارش ارزیابی ثبت شود');
    }

    // تولید شماره معرفی‌نامه
    const count = await this.referralLetterRepository.count();
    const letterNumber = `REF-${new Date().getFullYear()}-${String(count + 1).padStart(5, '0')}`;

    // تولید محتوای معرفی‌نامه
    const caseTypeLabel = SocialWorkCaseTypeLabels[socialWorkCase.caseType];
    const content = this.generateLetterContent(
      socialWorkCase,
      caseTypeLabel,
      generateReferralDto.additionalNotes,
    );

    const referralLetter = this.referralLetterRepository.create({
      letterNumber,
      socialWorkCaseId: socialWorkCase.id,
      content,
      referredTo: generateReferralDto.referredTo || 'حسابداری',
    });

    // به‌روزرسانی وضعیت پرونده
    socialWorkCase.status = SocialWorkCaseStatus.REFERRED;
    socialWorkCase.referredAt = new Date();
    await this.socialWorkCaseRepository.save(socialWorkCase);

    return this.referralLetterRepository.save(referralLetter);
  }

  private generateLetterContent(
    socialWorkCase: SocialWorkCase,
    caseTypeLabel: string,
    additionalNotes?: string,
  ): string {
    return `
بسمه تعالی

معرفی‌نامه به حسابداری

شماره: ${socialWorkCase.caseNumber}
تاریخ: ${new Date().toLocaleDateString('fa-IR')}

احتراماً،

بیمه شده محترم با مشخصات زیر:
- نام و نام خانوادگی: [نام بیمه شده]
- کد پرسنلی: [کد پرسنلی]
- نوع درخواست: ${caseTypeLabel}

با توجه به بررسی‌های انجام شده توسط واحد مددکاری و با عنایت به شرایط اجتماعی-اقتصادی بیمه شده، نیاز به حمایت مالی تأیید می‌گردد.

گزارش ارزیابی:
${socialWorkCase.assessmentReport}

${additionalNotes ? `\nیادداشت‌های تکمیلی:\n${additionalNotes}` : ''}

لذا خواهشمند است بررسی‌های لازم را جهت اقدامات حمایتی مقتضی معمول فرمایید.

با تشکر
واحد مددکاری
    `.trim();
  }

  async remove(id: string): Promise<void> {
    const result = await this.socialWorkCaseRepository.delete(id);
    if (result.affected === 0) {
      throw new NotFoundException(`پرونده مددکاری با شناسه ${id} یافت نشد`);
    }
  }
}
