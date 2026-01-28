import {
  Controller,
  Get,
  Post,
  Body,
  Patch,
  Param,
  Delete,
  UseGuards,
  Request,
} from '@nestjs/common';
import { JwtAuthGuard } from '../auth/guards/jwt-auth.guard';
import { SocialWorkService } from './social-work.service';
import { CreateSocialWorkDto } from './dto/create-social-work.dto';
import { UpdateAssessmentDto } from './dto/update-assessment.dto';
import { GenerateReferralDto } from './dto/generate-referral.dto';

@Controller('social-work')
@UseGuards(JwtAuthGuard)
export class SocialWorkController {
  constructor(private readonly socialWorkService: SocialWorkService) {}

  @Post()
  create(@Body() createSocialWorkDto: CreateSocialWorkDto, @Request() req) {
    return this.socialWorkService.create(createSocialWorkDto, req.user.userId);
  }

  @Get()
  findAll() {
    return this.socialWorkService.findAll();
  }

  @Get(':id')
  findOne(@Param('id') id: string) {
    return this.socialWorkService.findOne(id);
  }

  @Patch(':id/assessment')
  updateAssessment(
    @Param('id') id: string,
    @Body() updateAssessmentDto: UpdateAssessmentDto,
  ) {
    return this.socialWorkService.updateAssessment(id, updateAssessmentDto);
  }

  @Post(':id/referral')
  generateReferralLetter(
    @Param('id') id: string,
    @Body() generateReferralDto: GenerateReferralDto,
  ) {
    return this.socialWorkService.generateReferralLetter(
      id,
      generateReferralDto,
    );
  }

  @Delete(':id')
  remove(@Param('id') id: string) {
    return this.socialWorkService.remove(id);
  }
}
