import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { SocialWorkService } from './social-work.service';
import { SocialWorkController } from './social-work.controller';
import { SocialWorkCase } from '../entities/social-work-case.entity';
import { ReferralLetter } from '../entities/referral-letter.entity';

@Module({
  imports: [TypeOrmModule.forFeature([SocialWorkCase, ReferralLetter])],
  controllers: [SocialWorkController],
  providers: [SocialWorkService],
  exports: [SocialWorkService],
})
export class SocialWorkModule {}
