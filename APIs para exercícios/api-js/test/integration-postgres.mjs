/**
 * Integração: autentica Sequelize contra Postgres (CI com service postgres).
 * Define env ANTES do dynamic import de database.js.
 */
process.env.DB_HOST ??= 'localhost';
process.env.DB_NAME ??= 'apijs_test';
process.env.DB_USER ??= 'postgres';
process.env.DB_PASSWORD ??= 'postgres';
process.env.DB_DIALECT = 'postgres';

const { default: db } = await import('../database.js');

try {
  await db.authenticate();
  console.log('integration-postgres: ok');
} finally {
  await db.close();
}
