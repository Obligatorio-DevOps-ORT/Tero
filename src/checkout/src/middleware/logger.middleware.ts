import { Injectable, Logger, NestMiddleware } from '@nestjs/common';
import { Request, Response, NextFunction } from 'express';

// Registra cada request HTTP (salvo /health) cuando termina, con los mismos
// campos que los otros servicios: method, path, status y duration_ms.
// No registra la IP ni el navegador del cliente, ni el body ni los headers.
@Injectable()
export class LoggerMiddleware implements NestMiddleware {
  private logger = new Logger('HTTP');

  use(request: Request, response: Response, next: NextFunction) {
    // originalUrl incluye la query (?algo=...): nos quedamos solo con la ruta.
    const path = request.originalUrl.split('?')[0];
    if (path === '/health') {
      next();
      return;
    }

    const inicio = process.hrtime.bigint();

    response.on('finish', () => {
      const status = response.statusCode;
      const duracionMs = Number(process.hrtime.bigint() - inicio) / 1_000_000;
      const datos = {
        msg: 'request',
        method: request.method,
        path,
        status,
        duration_ms: Math.round(duracionMs * 1000) / 1000,
      };

      if (status >= 500) {
        this.logger.error(datos);
      } else if (status >= 400) {
        this.logger.warn(datos);
      } else {
        this.logger.log(datos);
      }
    });

    next();
  }
}
