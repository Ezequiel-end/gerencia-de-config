import express from 'express';
import router from './router.js';
import db from './database.js';

const app = express();
const port = Number(process.env.PORT) || 3000;

app.use(express.json());
app.use('/', router);

try {
  await db.authenticate();
  console.log('Eba!');
  if (process.env.RUN_SYNC === 'true') {
    await db.sync();
    console.log('Schema sincronizado (RUN_SYNC).');
  }
} catch (error) {
  console.log('Não eba...', error?.message || error);
}

app.listen(port, '0.0.0.0', () => {
  console.log(`Rodando na porta ${port}.`);
});
