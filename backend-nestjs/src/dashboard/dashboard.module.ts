import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { DashboardController } from './dashboard.controller';
import { DashboardService } from './dashboard.service';
import { Case } from '../entities/case.entity';

@Module({
  imports: [TypeOrmModule.forFeature([Case])],
  controllers: [DashboardController],
  providers: [DashboardService]
})
export class DashboardModule {}
