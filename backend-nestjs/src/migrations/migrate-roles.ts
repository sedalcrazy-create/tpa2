import { DataSource } from 'typeorm';
import * as dotenv from 'dotenv';

// Load environment variables
dotenv.config();

async function migrateRoles() {
  const dataSource = new DataSource({
    type: 'postgres',
    host: process.env.DB_HOST || 'localhost',
    port: parseInt(process.env.DB_PORT || '5432', 10),
    username: process.env.DB_USERNAME || 'postgres',
    password: process.env.DB_PASSWORD || '123',
    database: process.env.DB_NAME || 'medical_commission',
  });

  try {
    await dataSource.initialize();
    console.log('✓ Database connection established');

    // Check if 'role' column exists (old schema)
    const checkColumnQuery = `
      SELECT column_name
      FROM information_schema.columns
      WHERE table_name = 'users' AND column_name = 'role';
    `;
    const roleColumnExists = await dataSource.query(checkColumnQuery);

    if (roleColumnExists && roleColumnExists.length > 0) {
      console.log('✓ Old "role" column found, starting migration...');

      // Get all users with their current role
      const users = await dataSource.query(
        'SELECT id, role FROM users WHERE role IS NOT NULL;'
      );

      console.log(`✓ Found ${users.length} users to migrate`);

      // Add new 'roles' column if it doesn't exist
      const checkRolesColumnQuery = `
        SELECT column_name
        FROM information_schema.columns
        WHERE table_name = 'users' AND column_name = 'roles';
      `;
      const rolesColumnExists = await dataSource.query(checkRolesColumnQuery);

      if (!rolesColumnExists || rolesColumnExists.length === 0) {
        await dataSource.query(
          `ALTER TABLE users ADD COLUMN roles TEXT;`
        );
        console.log('✓ Added new "roles" column');
      }

      // Migrate each user's role to roles array
      for (const user of users) {
        if (user.role) {
          // Convert single role to array format (simple-array uses comma separation)
          await dataSource.query(
            'UPDATE users SET roles = $1 WHERE id = $2;',
            [user.role, user.id]
          );
        }
      }

      console.log('✓ Migrated all users successfully');

      // Drop old 'role' column
      await dataSource.query('ALTER TABLE users DROP COLUMN role;');
      console.log('✓ Dropped old "role" column');

      console.log('✅ Migration completed successfully!');
    } else {
      console.log('ℹ No migration needed - "role" column does not exist');
    }

    await dataSource.destroy();
  } catch (error) {
    console.error('❌ Migration failed:', error);
    process.exit(1);
  }
}

// Run migration
migrateRoles();
