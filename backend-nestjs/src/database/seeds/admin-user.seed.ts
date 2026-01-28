import { DataSource } from 'typeorm';
import * as bcrypt from 'bcryptjs';
import { User } from '../../entities/user.entity';
import { UserRole } from '../../common/enums/user-role.enum';

export async function seedAdminUser(dataSource: DataSource) {
  const userRepository = dataSource.getRepository(User);

  // Check if admin user already exists
  const existingAdmin = await userRepository.findOne({
    where: { email: 'admin@test.com' },
  });

  if (existingAdmin) {
    console.log('Admin user already exists');
    return;
  }

  // Create admin user
  const hashedPassword = await bcrypt.hash('admin123', 10);

  const adminUser = userRepository.create({
    nationalId: '1234567890',
    firstName: 'Admin',
    lastName: 'System',
    email: 'admin@test.com',
    phone: '09123456789',
    password: hashedPassword,
    roles: [UserRole.SYSTEM_ADMIN],
    isActive: true,
    isMfaEnabled: false,
  });

  await userRepository.save(adminUser);
  console.log('✅ Admin user created successfully: admin@test.com / admin123');
}
