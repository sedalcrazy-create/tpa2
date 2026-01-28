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
import { UserRole } from '../common/enums/user-role.enum';
import { Province } from './province.entity';
import { Case } from './case.entity';
import { ExpertOpinion } from './expert-opinion.entity';
import { AuditLog } from './audit-log.entity';

@Entity('commission_users')
export class User {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column({ name: 'national_id', unique: true, length: 10 })
  nationalId: string; // کد ملی

  @Column({ name: 'first_name', length: 100 })
  firstName: string;

  @Column({ name: 'last_name', length: 100 })
  lastName: string;

  @Column({ unique: true, length: 100 })
  email: string;

  @Column({ length: 20, nullable: true })
  phone: string;

  @Column({ type: 'text' })
  password: string; // hashed password

  @Column({ type: 'simple-array' })
  roles: UserRole[];

  @Column({ name: 'is_active', default: true })
  isActive: boolean;

  @Column({ name: 'is_mfa_enabled', default: false })
  isMfaEnabled: boolean;

  @Column({ name: 'mfa_secret', type: 'text', nullable: true })
  mfaSecret: string; // For TOTP (Time-based One-Time Password)

  @Column({ name: 'refresh_token', type: 'text', nullable: true })
  refreshToken: string;

  @Column({ name: 'last_login_at', type: 'timestamp', nullable: true })
  lastLoginAt: Date;

  // Province association (for provincial users)
  @ManyToOne(() => Province, { nullable: true })
  @JoinColumn({ name: 'province_id' })
  province: Province;

  @Column({ name: 'province_id', type: 'uuid', nullable: true })
  provinceId: string;

  // Specialty for commission members
  @Column({ length: 200, nullable: true })
  specialty: string; // تخصص پزشکی

  @CreateDateColumn({ name: 'created_at' })
  createdAt: Date;

  @UpdateDateColumn({ name: 'updated_at' })
  updatedAt: Date;

  // Relations
  @OneToMany(() => Case, (caseEntity) => caseEntity.assignedTo)
  assignedCases: Case[];

  @OneToMany(() => ExpertOpinion, (opinion) => opinion.expert)
  expertOpinions: ExpertOpinion[];

  @OneToMany(() => AuditLog, (log) => log.user)
  auditLogs: AuditLog[];
}
