import {
  Entity,
  PrimaryGeneratedColumn,
  Column,
  CreateDateColumn,
  UpdateDateColumn,
  ManyToOne,
  ManyToMany,
  JoinTable,
  JoinColumn,
} from 'typeorm';
import { CommissionLevel } from '../common/enums/commission-level.enum';
import { Case } from './case.entity';
import { User } from './user.entity';
import { Province } from './province.entity';

@Entity('meetings')
export class Meeting {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column({ unique: true, length: 50 })
  meetingNumber: string; // شماره جلسه

  @Column({ type: 'timestamp' })
  scheduledAt: Date;

  @Column({ length: 200 })
  location: string; // محل برگزاری

  @Column({ type: 'varchar' })
  commissionLevel: CommissionLevel;

  @ManyToOne(() => Province, { nullable: true })
  @JoinColumn({ name: 'provinceId' })
  province: Province;

  @Column({ type: 'uuid', nullable: true })
  provinceId: string;

  @Column({ type: 'text', nullable: true })
  agenda: string; // دستور جلسه

  @Column({ type: 'text', nullable: true })
  minutes: string; // صورت‌جلسه

  @Column({ default: false })
  isCompleted: boolean;

  @Column({ type: 'timestamp', nullable: true })
  completedAt: Date;

  @CreateDateColumn()
  createdAt: Date;

  @UpdateDateColumn()
  updatedAt: Date;

  // Relations
  @ManyToOne(() => Case)
  @JoinColumn({ name: 'caseId' })
  case: Case;

  @Column({ type: 'uuid' })
  caseId: string;

  @ManyToMany(() => User)
  @JoinTable({
    name: 'meeting_attendees',
    joinColumn: { name: 'meetingId', referencedColumnName: 'id' },
    inverseJoinColumn: { name: 'userId', referencedColumnName: 'id' },
  })
  attendees: User[];
}
