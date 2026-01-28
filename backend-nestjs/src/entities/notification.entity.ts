import {
  Entity,
  PrimaryGeneratedColumn,
  Column,
  CreateDateColumn,
  UpdateDateColumn,
  ManyToOne,
  JoinColumn,
} from 'typeorm';
import { User } from './user.entity';

@Entity('commission_notifications')
export class Notification {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @ManyToOne(() => User)
  @JoinColumn({ name: 'user_id' })
  user: User;

  @Column({ name: 'user_id', type: 'uuid' })
  userId: string;

  @Column({ length: 200 })
  title: string;

  @Column({ type: 'text' })
  message: string;

  @Column({ length: 100 })
  type: string; // نوع اعلان (مثلاً "MEETING", "CASE_ASSIGNED", "VERDICT_ISSUED")

  @Column({ name: 'is_read', default: false })
  isRead: boolean;

  @Column({ type: 'json', nullable: true })
  metadata: any; // داده‌های اضافی (مثلاً لینک، شناسه پرونده)

  @CreateDateColumn({ name: 'created_at' })
  createdAt: Date;

  @UpdateDateColumn({ name: 'updated_at' })
  updatedAt: Date;
}
