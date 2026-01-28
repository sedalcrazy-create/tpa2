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
import { CaseStatus } from '../common/enums/case-status.enum';
import { CaseType } from './case-type.entity';
import { VerdictTemplate } from './verdict-template.entity';
import { CommissionLevel } from '../common/enums/commission-level.enum';
import { InsuredPerson } from './insured-person.entity';
import { Province } from './province.entity';
import { User } from './user.entity';
import { Document } from './document.entity';
import { ExpertOpinion } from './expert-opinion.entity';
import { Meeting } from './meeting.entity';
import { Verdict } from './verdict.entity';
import { CaseTimeline } from './case-timeline.entity';

@Entity('cases')
export class Case {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column({ unique: true, length: 50 })
  caseNumber: string; // شماره پرونده

  @ManyToOne(() => InsuredPerson, (person) => person.cases)
  @JoinColumn({ name: 'insuredPersonId' })
  insuredPerson: InsuredPerson;

  @Column({ type: 'uuid' })
  insuredPersonId: string;

  // نوع پرونده (مثل معاینه نقص عضو، بیماری خاص و غیره)
  @ManyToOne(() => CaseType, (caseType) => caseType.cases, { nullable: true })
  @JoinColumn({ name: 'caseTypeId' })
  caseType: CaseType;

  @Column({ type: 'uuid', nullable: true })
  caseTypeId: string;

  // رای کمیسیون انتخابی
  @ManyToOne(() => VerdictTemplate, { nullable: true })
  @JoinColumn({ name: 'verdictTemplateId' })
  verdictTemplate: VerdictTemplate;

  @Column({ type: 'uuid', nullable: true })
  verdictTemplateId: string;

  @Column({ type: 'varchar', default: CaseStatus.PENDING_SECRETARIAT })
  status: CaseStatus;

  @Column({ type: 'varchar' })
  commissionLevel: CommissionLevel;

  @ManyToOne(() => Province, { nullable: true })
  @JoinColumn({ name: 'provinceId' })
  province: Province;

  @Column({ type: 'uuid', nullable: true })
  provinceId: string;

  @Column({ type: 'text', nullable: true })
  description: string; // توضیحات اولیه

  @Column({ type: 'text', nullable: true })
  medicalHistory: string; // سابقه پزشکی

  // Assigned specialist/expert
  @ManyToOne(() => User, { nullable: true })
  @JoinColumn({ name: 'assignedToId' })
  assignedTo: User;

  @Column({ type: 'uuid', nullable: true })
  assignedToId: string;

  @Column({ type: 'timestamp', nullable: true })
  assignedAt: Date;

  @CreateDateColumn()
  createdAt: Date;

  @UpdateDateColumn()
  updatedAt: Date;

  @Column({ type: 'timestamp', nullable: true })
  closedAt: Date;

  // Relations
  @OneToMany(() => Document, (doc) => doc.case)
  documents: Document[];

  @OneToMany(() => ExpertOpinion, (opinion) => opinion.case)
  expertOpinions: ExpertOpinion[];

  @OneToMany(() => CaseTimeline, (timeline) => timeline.case)
  timeline: CaseTimeline[];

  @OneToMany(() => Meeting, (meeting) => meeting.case)
  meetings: Meeting[];

  @OneToMany(() => Verdict, (verdict) => verdict.case)
  verdicts: Verdict[];
}
