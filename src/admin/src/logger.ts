import { Request, Response, NextFunction } from 'express';

// Nombre del servicio que va en cada línea de log.
const SERVICE = 'admin';

type Nivel = 'INFO' | 'WARN' | 'ERROR';

// Escribe una línea JSON por stdout con el esquema común de los 6 servicios:
// time, level, msg, service y los campos extra que se pasen.
export function log(level: Nivel, msg: string, campos: Record<string, unknown> = {}): void {
  const linea = {
    time: new Date().toISOString(),
    level,
    msg,
    service: SERVICE,
    ...campos,
  };
  console.log(JSON.stringify(linea));
}

// Middleware que registra cada request HTTP (salvo /health) cuando termina.
// Nunca registra el body, las cookies ni los headers: ahí viajan la contraseña
// del login y el JWT del admin.
export function requestLogger(req: Request, res: Response, next: NextFunction): void {
  // Guardamos el path ANTES de seguir: los proxies de ui lo reescriben después.
  // req.path no incluye la query (?page=2...), solo la ruta.
  const path = req.path;
  if (path === '/health') {
    next();
    return;
  }

  const inicio = process.hrtime.bigint();

  // 'finish' se dispara cuando la respuesta terminó de enviarse.
  res.on('finish', () => {
    const status = res.statusCode;
    const level: Nivel = status >= 500 ? 'ERROR' : status >= 400 ? 'WARN' : 'INFO';
    const duracionMs = Number(process.hrtime.bigint() - inicio) / 1_000_000;
    log(level, 'request', {
      method: req.method,
      path,
      status,
      duration_ms: Math.round(duracionMs * 1000) / 1000,
    });
  });

  next();
}

// Mensajes que Node escribe por su cuenta, en texto:
// - Si el proceso se cae por un error que nadie atrapó, escribe el error en varias
//   líneas. Con esto queda en una sola línea JSON y el proceso termina igual que
//   antes (código 1), para que Docker o ECS lo reinicien.
// - Los avisos (warnings), por ejemplo cuando una librería usa una función vieja.
//   Se quita el aviso en texto que trae Node y se escribe en JSON.
export function registrarEventosDelProceso(): void {
  process.removeAllListeners('warning');
  process.on('warning', (aviso) => {
    log('WARN', aviso.message, { tipo: aviso.name });
  });

  process.on('uncaughtException', (err) => {
    log('ERROR', 'error no controlado', { error: err.stack ?? String(err) });
    process.exit(1);
  });
  process.on('unhandledRejection', (razon) => {
    log('ERROR', 'promesa rechazada sin manejar', {
      error: razon instanceof Error ? razon.stack ?? razon.message : String(razon),
    });
    process.exit(1);
  });
}
