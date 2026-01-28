import {
  Entity,
  PrimaryGeneratedColumn,
  Column,
  CreateDateColumn,
  ManyToOne,
  JoinColumn,
} from 'typeorm';
import { User } from './user.entity';

@Entity('commission_audit_logs')
export class AuditLog {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @ManyToOne(() => User, { nullable: true })
  @JoinColumn({ name: 'user_id' })
  user: User;

  @Column({ name: 'user_id', type: 'uuid', nullable: true })
  userId: string;

  @Column({ length: 100 })
  action: string; // نوع عملیات (مثلاً "LOGIN", "CREATE_CASE", "UPDATE_VERDICT")

  @Column({ name: 'entity_type', length: 100 })
  entityType: string; // نوع موجودیت (مثلاً "Case", "User")

  @Column({ name: 'entity_id', type: 'uuid', nullable: true })
  entityId: string; // شناسه موجودیت

  @Column({ name: 'old_values', type: 'json', nullable: true })
  oldValues: any; // مقادیر قبلی (برای UPDATE)

  @Column({ name: 'new_values', type: 'json', nullable: true })
  newValues: any; // مقادیر جدید (برای UPDATE/CREATE)

  @Column({ name: 'ip_address', length: 50, nullable: true })
  ipAddress: string;

  @Column({ name: 'user_agent', length: 500, nullable: true })
  userAgent: string;

  @CreateDateColumn({ name: 'created_at' })
  createdAt: Date;
}
