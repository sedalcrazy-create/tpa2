import {
  Entity,
  PrimaryGeneratedColumn,
  Column,
  CreateDateColumn,
  UpdateDateColumn,
  ManyToOne,
  OneToMany,
  JoinColumn,
} from 'typeorm';
import { SocialWorkCaseType } from '../common/enums/social-work-case-type.enum';
import { SocialWorkCaseStatus } from '../common/enums/social-work-case-status.enum';
import { InsuredPerson } from './insured-person.entity';
import { User } from './user.entity';
import { Case } from './case.entity';
import { ReferralLetter } from './referral-letter.entity';

@Entity('social_work_cases')
export class SocialWorkCase {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column({ unique: true, length: 50 })
  caseNumber: string; // شماره پرونده (MC-SW-2025-XXXX)

  @Column({ type: 'varchar' })
  caseType: SocialWorkCaseType; // نوع خدمت مددکاری

  @Column({ type: 'varchar', default: SocialWorkCaseStatus.DRAFT })
  status: SocialWorkCaseStatus; // وضعیت پرونده

  // بیمه شده مرتبط
  @ManyToOne(() => InsuredPerson, { nullable: false })
  @JoinColumn({ name: 'insuredPersonId' })
  insuredPerson: InsuredPerson;

  @Column({ type: 'uuid' })
  insuredPersonId: string;

  // مددکار مسئول
  @ManyToOne(() => User, { nullable: false })
  @JoinColumn({ name: 'socialWorkerId' })
  socialWorker: User;

  @Column({ type: 'uuid' })
  socialWorkerId: string;

  // ارتباط اختیاری با پرونده پزشکی
  @ManyToOne(() => Case, { nullable: true })
  @JoinColumn({ name: 'medicalCaseId' })
  medicalCase: Case;

  @Column({ type: 'uuid', nullable: true })
  medicalCaseId: string;

  // جزئیات درخواست (JSON)
  @Column({ type: 'simple-json', nullable: true })
  requestDetails: any; // اطلاعات فرم درخواست

  // گزارش ارزیابی مددکار
  @Column({ type: 'text', nullable: true })
  assessmentReport: string;

  // تاریخ ارزیابی
  @Column({ type: 'timestamp', nullable: true })
  assessedAt: Date;

  // تاریخ ارجاع
  @Column({ type: 'timestamp', nullable: true })
  referredAt: Date;

  @CreateDateColumn()
  createdAt: Date;

  @UpdateDateColumn()
  updatedAt: Date;

  @Column({ type: 'timestamp', nullable: true })
  closedAt: Date;

  // Relations
  @OneToMany(() => ReferralLetter, (letter) => letter.socialWorkCase)
  referralLetters: ReferralLetter[];
}
