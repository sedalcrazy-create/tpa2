import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { VerdictTemplatesController } from './verdict-templates.controller';
import { VerdictTemplatesService } from './verdict-templates.service';
import { VerdictTemplate } from '../entities/verdict-template.entity';

@Module({
  imports: [TypeOrmModule.forFeature([VerdictTemplate])],
  controllers: [VerdictTemplatesController],
  providers: [VerdictTemplatesService],
  exports: [VerdictTemplatesService],
})
export class VerdictTemplatesModule {}
