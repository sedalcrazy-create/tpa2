import { Module } from '@nestjs/common';
import { ConfigModule, ConfigService } from '@nestjs/config';
import { TypeOrmModule } from '@nestjs/typeorm';
import { AppController } from './app.controller';
import { AppService } from './app.service';
import { AuthModule } from './auth/auth.module';
import { User } from './entities/user.entity';
import { Province } from './entities/province.entity';
import { InsuredPerson } from './entities/insured-person.entity';
import { Case } from './entities/case.entity';
import { Document } from './entities/document.entity';
import { ExpertOpinion } from './entities/expert-opinion.entity';
import { Meeting } from './entities/meeting.entity';
import { Verdict } from './entities/verdict.entity';
import { CaseTimeline } from './entities/case-timeline.entity';
import { AuditLog } from './entities/audit-log.entity';
import { Notification } from './entities/notification.entity';
import { CaseType } from './entities/case-type.entity';
import { VerdictTemplate } from './entities/verdict-template.entity';
import { SocialWorkCase } from './entities/social-work-case.entity';
import { ReferralLetter } from './entities/referral-letter.entity';
import { UsersModule } from './users/users.module';
import { ProvincesModule } from './provinces/provinces.module';
import { CasesModule } from './cases/cases.module';
import { InsuredPersonsModule } from './insured-persons/insured-persons.module';
import { ImportModule } from './import/import.module';
import { CaseTypesModule } from './case-types/case-types.module';
import { VerdictTemplatesModule } from './verdict-templates/verdict-templates.module';
import { DashboardModule } from './dashboard/dashboard.module';
import { SocialWorkModule } from './social-work/social-work.module';

@Module({
  imports: [
    ConfigModule.forRoot({
      isGlobal: true,
      envFilePath: '.env',
    }),
    TypeOrmModule.forRootAsync({
      inject: [ConfigService],
      useFactory: (configService: ConfigService) => {
        const dbType = configService.get('DB_TYPE') || 'postgres';
        const entities = [
          User,
          Province,
          InsuredPerson,
          Case,
          Document,
          ExpertOpinion,
          Meeting,
          Verdict,
          CaseTimeline,
          AuditLog,
          Notification,
          CaseType,
          VerdictTemplate,
          SocialWorkCase,
          ReferralLetter,
        ];

        // SQLite configuration
        if (dbType === 'sqlite') {
          return {
            type: 'sqlite' as const,
            database: (configService.get('DB_DATABASE') || 'medical_commission.sqlite') as string,
            entities,
            synchronize: true,
            logging: configService.get('NODE_ENV') === 'development',
          };
        }

        // PostgreSQL configuration (default)
        return {
          type: 'postgres' as const,
          host: (configService.get('DB_HOST') || 'localhost') as string,
          port: parseInt(configService.get('DB_PORT') || '5432', 10),
          username: (configService.get('DB_USERNAME') || 'postgres') as string,
          password: (configService.get('DB_PASSWORD') || '123') as string,
          database: (configService.get('DB_NAME') || 'medical_commission') as string,
          entities,
          synchronize: true, // Disable in production!
          logging: configService.get('NODE_ENV') === 'development',
        };
      },
    }),
    AuthModule,
    UsersModule,
    ProvincesModule,
    CasesModule,
    InsuredPersonsModule,
    ImportModule,
    CaseTypesModule,
    VerdictTemplatesModule,
    DashboardModule,
    SocialWorkModule,
  ],
  controllers: [AppController],
  providers: [AppService],
})
export class AppModule {}
