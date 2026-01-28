import {
  Entity,
  PrimaryGeneratedColumn,
  Column,
  CreateDateColumn,
  UpdateDateColumn,
  ManyToOne,
  JoinColumn,
} from 'typeorm';
import { CaseType } from './case-type.entity';

@Entity('verdict_templates')
export class VerdictTemplate {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column({ length: 500 })
  title: string; // عنوان رای (مثلاً: دارای 25% نقص عضو)

  @Column({ type: 'text', nullable: true })
  description: string; // متن کامل رای / توضیحات

  @Column({ default: true })
  isActive: boolean;

  @Column({ type: 'int', default: 0 })
  sortOrder: number; // ترتیب نمایش

  @CreateDateColumn()
  createdAt: Date;

  @UpdateDateColumn()
  updatedAt: Date;

  // Relations
  @ManyToOne(() => CaseType, (caseType) => caseType.verdictTemplates, { onDelete: 'CASCADE' })
  @JoinColumn({ name: 'caseTypeId' })
  caseType: CaseType;

  @Column({ type: 'uuid' })
  caseTypeId: string;
}
