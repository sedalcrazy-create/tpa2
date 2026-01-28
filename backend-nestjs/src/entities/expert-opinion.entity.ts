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

@Entity('expert_opinions')
export class ExpertOpinion {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @ManyToOne(() => Case, (caseEntity) => caseEntity.expertOpinions)
  @JoinColumn({ name: 'caseId' })
  case: Case;

  @Column({ type: 'uuid' })
  caseId: string;

  @ManyToOne(() => User, (user) => user.expertOpinions)
  @JoinColumn({ name: 'expertId' })
  expert: User;

  @Column({ type: 'uuid' })
  expertId: string;

  @Column({ type: 'text' })
  opinion: string; // نظر کارشناسی

  @Column({ type: 'int', nullable: true })
  disabilityPercentage: number; // درصد ازکارافتادگی (if applicable)

  @Column({ type: 'text', nullable: true })
  recommendations: string; // توصیه‌ها

  @Column({ default: false })
  isApproved: boolean; // تأیید شده توسط رئیس

  @CreateDateColumn()
  createdAt: Date;

  @UpdateDateColumn()
  updatedAt: Date;
}
