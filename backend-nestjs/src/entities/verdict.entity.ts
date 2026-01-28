import {
  Entity,
  PrimaryGeneratedColumn,
  Column,
  CreateDateColumn,
  UpdateDateColumn,
  ManyToOne,
  JoinColumn,
} from 'typeorm';
import { Case } from './case.entity';
import { User } from './user.entity';
import { Meeting } from './meeting.entity';

@Entity('verdicts')
export class Verdict {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column({ unique: true, length: 50 })
  verdictNumber: string; // شماره رأی

  @ManyToOne(() => Case, (caseEntity) => caseEntity.verdicts)
  @JoinColumn({ name: 'caseId' })
  case: Case;

  @Column({ type: 'uuid' })
  caseId: string;

  @ManyToOne(() => Meeting, { nullable: true })
  @JoinColumn({ name: 'meetingId' })
  meeting: Meeting;

  @Column({ type: 'uuid', nullable: true })
  meetingId: string;

  @Column({ type: 'text' })
  content: string; // محتوای رأی

  @Column({ type: 'int', nullable: true })
  disabilityPercentage: number; // درصد ازکارافتادگی نهایی

  @Column({ default: false })
  isApproved: boolean;

  @ManyToOne(() => User)
  @JoinColumn({ name: 'approvedById' })
  approvedBy: User;

  @Column({ type: 'uuid', nullable: true })
  approvedById: string;

  @Column({ type: 'timestamp', nullable: true })
  approvedAt: Date;

  @Column({ type: 'text', nullable: true })
  digitalSignature: string; // امضای دیجیتال

  @Column({ default: false })
  isDispatched: boolean; // ارسال شده به واحدهای مرتبط

  @Column({ type: 'timestamp', nullable: true })
  dispatchedAt: Date;

  @CreateDateColumn()
  createdAt: Date;

  @UpdateDateColumn()
  updatedAt: Date;
}
