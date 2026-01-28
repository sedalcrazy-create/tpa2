import { Controller, Get, Post, Body, Patch, Param, Delete, Query } from '@nestjs/common';
import { InsuredPersonsService } from './insured-persons.service';
import { CreateInsuredPersonDto } from './dto/create-insured-person.dto';
import { UpdateInsuredPersonDto } from './dto/update-insured-person.dto';

@Controller('insured-persons')
export class InsuredPersonsController {
  constructor(private readonly insuredPersonsService: InsuredPersonsService) {}

  @Post()
  create(@Body() createInsuredPersonDto: CreateInsuredPersonDto) {
    return this.insuredPersonsService.create(createInsuredPersonDto);
  }

  @Get()
  findAll(
    @Query('nationalId') nationalId?: string,
    @Query('personnelCode') personnelCode?: string,
    @Query('firstName') firstName?: string,
    @Query('lastName') lastName?: string,
  ) {
    return this.insuredPersonsService.findAll({
      nationalId,
      personnelCode,
      firstName,
      lastName,
    });
  }

  @Get(':id')
  findOne(@Param('id') id: string) {
    return this.insuredPersonsService.findOne(id);
  }

  @Patch(':id')
  update(@Param('id') id: string, @Body() updateInsuredPersonDto: UpdateInsuredPersonDto) {
    return this.insuredPersonsService.update(id, updateInsuredPersonDto);
  }

  @Delete(':id')
  remove(@Param('id') id: string) {
    return this.insuredPersonsService.remove(id);
  }
}
