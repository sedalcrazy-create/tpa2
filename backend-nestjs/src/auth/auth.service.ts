import { Injectable, UnauthorizedException, ConflictException } from '@nestjs/common';
import { JwtService } from '@nestjs/jwt';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import * as bcrypt from 'bcryptjs';
import { User } from '../entities/user.entity';
import { RegisterDto } from './dto/register.dto';
import { LoginDto } from './dto/login.dto';

@Injectable()
export class AuthService {
  constructor(
    @InjectRepository(User)
    private userRepository: Repository<User>,
    private jwtService: JwtService,
  ) {}

  async register(registerDto: RegisterDto) {
    const { email, nationalId, password, ...userData } = registerDto;

    // Check if user exists
    const existingUser = await this.userRepository.findOne({
      where: [{ email }, { nationalId }],
    });

    if (existingUser) {
      throw new ConflictException('کاربر با این ایمیل یا کد ملی قبلاً ثبت شده است');
    }

    // Hash password
    const hashedPassword = await bcrypt.hash(password, 10);

    // Create user
    const user = this.userRepository.create({
      email,
      nationalId,
      password: hashedPassword,
      ...userData,
    });

    await this.userRepository.save(user);

    const { password: _, ...result } = user;
    return {
      user: result,
      token: this.generateToken(user),
    };
  }

  async login(loginDto: LoginDto) {
    const { email, password } = loginDto;

    const user = await this.userRepository.findOne({
      where: { email },
      relations: ['province'],
    });

    if (!user) {
      throw new UnauthorizedException('ایمیل یا رمز عبور اشتباه است');
    }

    if (!user.isActive) {
      throw new UnauthorizedException('حساب کاربری غیرفعال است');
    }

    const isPasswordValid = await bcrypt.compare(password, user.password);
    if (!isPasswordValid) {
      throw new UnauthorizedException('ایمیل یا رمز عبور اشتباه است');
    }

    // Update last login
    user.lastLoginAt = new Date();
    await this.userRepository.save(user);

    const { password: _, ...result } = user;
    return {
      user: result,
      token: this.generateToken(user),
    };
  }

  async getProfile(userId: string) {
    const user = await this.userRepository.findOne({
      where: { id: userId },
      relations: ['province'],
    });

    if (!user) {
      throw new UnauthorizedException('کاربر یافت نشد');
    }

    const { password, ...result } = user;
    return result;
  }

  private generateToken(user: User): string {
    const payload = {
      id: user.id,
      email: user.email,
      roles: user.roles,
      provinceId: user.provinceId,
    };

    return this.jwtService.sign(payload);
  }
}
