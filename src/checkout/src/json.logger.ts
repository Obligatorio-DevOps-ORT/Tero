import { LoggerService } from '@nestjs/common';

// Nombre del servicio que va en cada línea de log.
const SERVICE = 'checkout';

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

// Logger para Nest: Nest no escribe directo en la consola, sino que le pasa
// cada mensaje a un "LoggerService". Este lo escribe en JSON con el esquema común.
// Se registra en main.ts con NestFactory.create(AppModule, { logger: new JsonLogger() }).
export class JsonLogger implements LoggerService {
  log(message: unknown, ...params: unknown[]): void {
    this.escribir('INFO', message, params);
  }

  warn(message: unknown, ...params: unknown[]): void {
    this.escribir('WARN', message, params);
  }

  error(message: unknown, ...params: unknown[]): void {
    this.escribir('ERROR', message, params);
  }

  fatal(message: unknown, ...params: unknown[]): void {
    this.escribir('ERROR', message, params);
  }

  // debug y verbose no se escriben: igual que en los otros servicios,
  // registramos desde INFO para arriba.
  debug(): void {}
  verbose(): void {}

  private escribir(level: Nivel, message: unknown, params: unknown[]): void {
    const campos: Record<string, unknown> = {};
    const resto = [...params];

    // Nest pasa como último parámetro el "contexto": quién escribe el mensaje
    // (por ejemplo "HTTP", "NestFactory" o "ExceptionsHandler").
    if (typeof resto[resto.length - 1] === 'string') {
      campos.context = resto.pop();
    }
    // En logger.error(mensaje, stack, contexto), lo que queda es el stack del error.
    if (level === 'ERROR' && typeof resto[0] === 'string') {
      campos.error = resto[0];
    }

    let msg: string;
    if (message instanceof Error) {
      // Nest a veces pasa el error entero (por ejemplo cuando un endpoint falla).
      msg = message.message;
      campos.error = message.stack ?? message.message;
    } else if (typeof message === 'object' && message !== null) {
      // Objeto con campos, por ejemplo { msg: 'request', method, path, status }.
      const { msg: texto, ...otros } = message as Record<string, unknown>;
      msg = String(texto ?? '');
      Object.assign(campos, otros);
    } else {
      msg = String(message);
    }

    log(level, msg, campos);
  }
}

// Mensajes que Node escribe por su cuenta, en texto (igual que en ui y admin):
// - Si el proceso se cae por un error que nadie atrapó, queda en una sola línea
//   JSON y el proceso termina igual que antes (código 1).
// - Los avisos (warnings) de Node pasan a JSON.
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
