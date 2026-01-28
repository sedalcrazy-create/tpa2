import { Injectable, NotFoundException } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { Case } from '../entities/case.entity';
import { Document } from '../entities/document.entity';
import { CaseTimeline } from '../entities/case-timeline.entity';
import { CreateCaseDto } from './dto/create-case.dto';
import { UpdateCaseDto } from './dto/update-case.dto';
import { CaseNumberGeneratorService } from './case-number-generator.service';

@Injectable()
export class CasesService {
  constructor(
    @InjectRepository(Case)
    private caseRepository: Repository<Case>,
    @InjectRepository(Document)
    private documentRepository: Repository<Document>,
    @InjectRepository(CaseTimeline)
    private timelineRepository: Repository<CaseTimeline>,
    private caseNumberGenerator: CaseNumberGeneratorService,
  ) {}

  async create(createCaseDto: CreateCaseDto, files?: Express.Multer.File[]) {
    // تولید خودکار شماره پرونده با استفاده از کد استخدامی
    const caseNumber = await this.caseNumberGenerator.generateCaseNumber(
      createCaseDto.insuredPersonId,
    );

    const caseEntity = this.caseRepository.create({
      ...createCaseDto,
      caseNumber,
    });
    const savedCase = await this.caseRepository.save(caseEntity);

    // Create timeline entry for case creation
    await this.createTimelineEntry(
      savedCase.id,
      'ثبت پرونده',
      `پرونده با شماره ${caseNumber} ثبت شد`,
      null, // TODO: Replace with actual authenticated user ID
    );

    // Create document records for uploaded files
    if (files && files.length > 0) {
      for (const file of files) {
        const document = this.documentRepository.create({
          caseId: savedCase.id,
          title: file.originalname,
          fileName: file.filename,
          filePath: file.path,
          fileType: file.mimetype,
          fileSize: file.size,
          uploadedById: null, // TODO: Replace with actual authenticated user ID
        });
        await this.documentRepository.save(document);
      }
    }

    return this.findOne(savedCase.id);
  }

  async findAll(insuredPersonId?: string) {
    const queryOptions: any = {
      relations: ['insuredPerson', 'province', 'assignedTo', 'caseType', 'verdictTemplate'],
      order: { createdAt: 'DESC' },
    };

    if (insuredPersonId) {
      queryOptions.where = { insuredPersonId };
    }

    return this.caseRepository.find(queryOptions);
  }

  async findOne(id: string) {
    const caseEntity = await this.caseRepository.findOne({
      where: { id },
      relations: ['insuredPerson', 'province', 'assignedTo', 'caseType', 'verdictTemplate', 'documents', 'timeline'],
    });
    if (!caseEntity) {
      throw new NotFoundException('پرونده یافت نشد');
    }
    return caseEntity;
  }

  async update(id: string, updateCaseDto: UpdateCaseDto, files?: Express.Multer.File[]) {
    const existingCase = await this.findOne(id);

    // Track what changed for timeline entries
    const statusChanged = updateCaseDto.status && updateCaseDto.status !== existingCase.status;
    const assignmentChanged = updateCaseDto.assignedToId && updateCaseDto.assignedToId !== existingCase.assignedToId;

    // Update assignedAt if assignment is changing
    if (assignmentChanged) {
      updateCaseDto.assignedAt = new Date();
    }

    await this.caseRepository.update(id, updateCaseDto);

    // Create timeline entries for changes
    if (statusChanged) {
      await this.createTimelineEntry(
        id,
        'تغییر وضعیت',
        `وضعیت پرونده به "${this.getStatusLabel(updateCaseDto.status)}" تغییر یافت`,
        null, // TODO: Replace with actual authenticated user ID
      );
    }

    if (assignmentChanged) {
      await this.createTimelineEntry(
        id,
        'ارجاع به متخصص',
        `پرونده به متخصص ارجاع داده شد`,
        null, // TODO: Replace with actual authenticated user ID
      );
    }

    // Create document records for uploaded files
    if (files && files.length > 0) {
      for (const file of files) {
        const document = this.documentRepository.create({
          caseId: existingCase.id,
          title: file.originalname,
          fileName: file.filename,
          filePath: file.path,
          fileType: file.mimetype,
          fileSize: file.size,
          uploadedById: null, // TODO: Replace with actual authenticated user ID
        });
        await this.documentRepository.save(document);
      }
    }

    return this.findOne(id);
  }

  async remove(id: string) {
    const caseEntity = await this.findOne(id);
    await this.caseRepository.remove(caseEntity);
    return { message: 'پرونده با موفقیت حذف شد' };
  }

  // Helper method to create timeline entries
  private async createTimelineEntry(
    caseId: string,
    action: string,
    description: string,
    userId: string | null,
  ): Promise<void> {
    const timelineEntry = this.timelineRepository.create({
      caseId,
      action,
      description,
      userId,
    });
    await this.timelineRepository.save(timelineEntry);
  }

  // Helper method to get Persian label for status
  private getStatusLabel(status: string): string {
    const statusLabels: Record<string, string> = {
      PENDING_SECRETARIAT: 'در انتظار بررسی دبیرخانه',
      ASSIGNED_TO_SPECIALIST: 'ارجاع به متخصص',
      UNDER_REVIEW: 'در حال بررسی',
      PENDING_MEETING: 'در انتظار جلسه',
      MEETING_SCHEDULED: 'جلسه زمان‌بندی شده',
      PENDING_VERDICT: 'در انتظار رأی',
      VERDICT_ISSUED: 'رأی صادر شد',
      ARCHIVED: 'بایگانی شده',
      REJECTED: 'رد شده',
    };
    return statusLabels[status] || status;
  }
}
