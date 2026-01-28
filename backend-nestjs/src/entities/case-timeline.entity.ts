import {
  Entity,
  PrimaryGeneratedColumn,
  Column,
  CreateDateColumn,
  ManyToOne,
  JoinColumn,
} from 'typeorm';
import { Case } from './case.entity';
import { User } from './user.entity';

@Entity('case_timeline')
export class CaseTimeline {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @ManyToOne(() => Case, (caseEntity) => caseEntity.timeline)
  @JoinColumn({ name: 'caseId' })
  case: Case;

  @Column({ type: 'uuid' })
  caseId: string;

  @Column({ length: 200 })
  action: string; // نوع اقدام (مثلاً "پرونده ثبت شد"، "ارجاع به پزشک")

  @Column({ type: 'text', nullable: true })
  description: string; // توضیحات

  @ManyToOne(() => User, { nullable: true })
  @JoinColumn({ name: 'userId' })
  user: User;

  @Column({ type: 'uuid', nullable: true })
  userId: string;

  @CreateDateColumn()
  createdAt: Date;
}
