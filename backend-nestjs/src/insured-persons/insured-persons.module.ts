import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { InsuredPersonsService } from './insured-persons.service';
import { InsuredPersonsController } from './insured-persons.controller';
import { InsuredPerson } from '../entities/insured-person.entity';

@Module({
  imports: [TypeOrmModule.forFeature([InsuredPerson])],
  controllers: [InsuredPersonsController],
  providers: [InsuredPersonsService],
  exports: [InsuredPersonsService],
})
export class InsuredPersonsModule {}
