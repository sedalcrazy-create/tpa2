import {
  Entity,
  PrimaryGeneratedColumn,
  Column,
  CreateDateColumn,
  ManyToOne,
  JoinColumn,
} from 'typeorm';
import { SocialWorkCase } from './social-work-case.entity';

@Entity('referral_letters')
export class ReferralLetter {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column({ unique: true, length: 50 })
  letterNumber: string; // شماره معرفی‌نامه (REF-2025-XXXX)

  // پرونده مددکاری مرتبط
  @ManyToOne(() => SocialWorkCase, (swCase) => swCase.referralLetters, {
    nullable: false,
  })
  @JoinColumn({ name: 'socialWorkCaseId' })
  socialWorkCase: SocialWorkCase;

  @Column({ type: 'uuid' })
  socialWorkCaseId: string;

  // محتوای معرفی‌نامه
  @Column({ type: 'text' })
  content: string;

  // مسیر فایل PDF (اگر تولید شده باشد)
  @Column({ type: 'varchar', length: 500, nullable: true })
  pdfPath: string;

  // ارجاع به کجا (حسابداری، متخصص و غیره)
  @Column({ type: 'varchar', length: 100, default: 'حسابداری' })
  referredTo: string;

  @CreateDateColumn()
  generatedAt: Date;
}
