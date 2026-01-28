import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { ImportController } from './import.controller';
import { ImportService } from './import.service';
import { InsuredPerson } from '../entities/insured-person.entity';

@Module({
  imports: [TypeOrmModule.forFeature([InsuredPerson])],
  controllers: [ImportController],
  providers: [ImportService],
})
export class ImportModule {}
