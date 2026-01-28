import {
  Entity,
  PrimaryGeneratedColumn,
  Column,
  CreateDateColumn,
  UpdateDateColumn,
  OneToMany,
  ManyToOne,
  JoinColumn,
} from 'typeorm';
import { FamilyRelation } from '../common/enums/family-relation.enum';
import { Case } from './case.entity';
import { Province } from './province.entity';

@Entity('insured_persons')
export class InsuredPerson {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column({ unique: true, length: 10 })
  nationalId: string; // کد ملی

  @Column({ length: 50 })
  personnelCode: string; // کد پرسنلی (shared across family)

  @Column({ length: 100 })
  firstName: string;

  @Column({ length: 100 })
  lastName: string;

  @Column({ type: 'date' })
  birthDate: Date;

  @Column({ type: 'varchar', default: FamilyRelation.SELF })
  familyRelation: FamilyRelation; // نسبت با بیمه‌شده اصلی

  @Column({ length: 50, nullable: true })
  insuranceNumber: string; // شماره بیمه‌نامه

  @Column({ length: 20, nullable: true })
  phone: string;

  @Column({ length: 200, nullable: true })
  address: string;

  @Column({ length: 50, nullable: true })
  employmentStatus: string; // وضعیت خدمت

  @Column({ length: 200, nullable: true })
  officeLocation: string; // اداره امور / محل خدمت

  @Column({ length: 100, nullable: true })
  city: string; // شهر محل خدمت (برای شعب مستقل)

  // Province association (واحد کمیسیون پزشکی)
  @ManyToOne(() => Province, { nullable: true })
  @JoinColumn({ name: 'provinceId' })
  province: Province;

  @Column({ type: 'uuid', nullable: true })
  provinceId: string;

  @CreateDateColumn()
  createdAt: Date;

  @UpdateDateColumn()
  updatedAt: Date;

  // Relations
  @OneToMany(() => Case, (caseEntity) => caseEntity.insuredPerson)
  cases: Case[];
}
