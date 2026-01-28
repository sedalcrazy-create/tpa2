import {
  Entity,
  PrimaryGeneratedColumn,
  Column,
  CreateDateColumn,
  UpdateDateColumn,
  OneToMany,
} from 'typeorm';
import { Case } from './case.entity';
import { VerdictTemplate } from './verdict-template.entity';

@Entity('case_types')
export class CaseType {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column({ unique: true, length: 200 })
  name: string; // نام نوع پرونده (مثلاً: معاینه جهت تعیین درصد نقص عضو)

  @Column({ type: 'text', nullable: true })
  description: string; // توضیحات

  @Column({ default: true })
  isActive: boolean;

  // نوع ارجاع: true = کمیسیون مرکزی (تهران), false = کمیسیون استانی
  @Column({ default: false })
  isCentralCommission: boolean;

  @CreateDateColumn()
  createdAt: Date;

  @UpdateDateColumn()
  updatedAt: Date;

  // Relations
  @OneToMany(() => Case, (caseEntity) => caseEntity.caseType)
  cases: Case[];

  @OneToMany(() => VerdictTemplate, (template) => template.caseType)
  verdictTemplates: VerdictTemplate[];
}
