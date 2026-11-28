import { ValidationPipe } from '@nestjs/common';
import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import { JsonLogger, registrarEventosDelProceso } from './json.logger';

registrarEventosDelProceso();

async function bootstrap() {
  // logger: todo lo que escribe Nest (arranque, rutas, errores) sale en JSON.
  const app = await NestFactory.create(AppModule, { logger: new JsonLogger() });
  app.useGlobalPipes(new ValidationPipe());
  app.enableShutdownHooks();

  const port = process.env.PORT || 8080;
  await app.listen(port);
}
bootstrap();
