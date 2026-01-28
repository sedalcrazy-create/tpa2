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

@Entity('documents')
export class Document {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @ManyToOne(() => Case, (caseEntity) => caseEntity.documents)
  @JoinColumn({ name: 'caseId' })
  case: Case;

  @Column({ type: 'uuid' })
  caseId: string;

  @Column({ length: 200 })
  title: string; // عنوان مدرک

  @Column({ length: 500 })
  filePath: string; // مسیر فایل

  @Column({ length: 100 })
  fileName: string;

  @Column({ length: 50 })
  fileType: string; // MIME type

  @Column({ type: 'int' })
  fileSize: number; // در بایت

  @Column({ type: 'text', nullable: true })
  description: string;

  @ManyToOne(() => User)
  @JoinColumn({ name: 'uploadedById' })
  uploadedBy: User;

  @Column({ type: 'uuid', nullable: true })
  uploadedById: string;

  @CreateDateColumn()
  createdAt: Date;

  @UpdateDateColumn()
  updatedAt: Date;
}
