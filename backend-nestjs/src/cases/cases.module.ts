import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { CasesService } from './cases.service';
import { CasesController } from './cases.controller';
import { CaseNumberGeneratorService } from './case-number-generator.service';
import { Case } from '../entities/case.entity';
import { Document } from '../entities/document.entity';
import { InsuredPerson } from '../entities/insured-person.entity';
import { CaseTimeline } from '../entities/case-timeline.entity';

@Module({
  imports: [TypeOrmModule.forFeature([Case, Document, InsuredPerson, CaseTimeline])],
  controllers: [CasesController],
  providers: [CasesService, CaseNumberGeneratorService],
  exports: [CasesService],
})
export class CasesModule {}
