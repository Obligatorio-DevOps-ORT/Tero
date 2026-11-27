import json
import logging
import sys
from datetime import datetime, timezone

# Nombre del servicio que va en cada línea de log.
# Usamos "carts" porque así se llama el servicio en docker-compose y en ECS.
SERVICE = "carts"

# Python llama "WARNING" a lo que Go (slog) llama "WARN": lo unificamos
# para que en CloudWatch se pueda filtrar igual en todos los servicios.
NIVELES = {"WARNING": "WARN", "CRITICAL": "ERROR"}


class JsonFormatter(logging.Formatter):
    """Convierte cada mensaje de log en una línea JSON con el esquema común."""

    def format(self, record):
        linea = {
            "time": datetime.fromtimestamp(record.created, tz=timezone.utc)
            .isoformat(timespec="milliseconds")
            .replace("+00:00", "Z"),
            "level": NIVELES.get(record.levelname, record.levelname),
            "msg": record.getMessage(),
            "service": SERVICE,
        }
        # Campos extra (method, path, status...) que se pasan con extra={"campos": {...}}
        linea.update(getattr(record, "campos", {}))
        # Si el log viene de una excepción, se agrega el detalle (traceback) en un campo.
        if record.exc_info:
            linea["error"] = self.formatException(record.exc_info)
        return json.dumps(linea, ensure_ascii=False, separators=(",", ":"))


def configurar_logs():
    """Hace que todo lo que se loguee (nuestro código y uvicorn) salga en JSON por stdout."""
    handler = logging.StreamHandler(sys.stdout)  # stdout, igual que catalog y orders
    handler.setFormatter(JsonFormatter())

    # Logger raíz: el que usa nuestro código (logging.getLogger("carts")).
    raiz = logging.getLogger()
    raiz.handlers = [handler]
    raiz.setLevel(logging.INFO)

    # uvicorn configura sus propios loggers con formato de texto al arrancar.
    # Les sacamos ese formato y los mandamos al logger raíz (que ya es JSON).
    for nombre in ("uvicorn", "uvicorn.error"):
        logger = logging.getLogger(nombre)
        logger.handlers = []
        logger.propagate = True

    # El access log de uvicorn no mide la duración: lo apagamos y lo
    # reemplaza nuestro middleware (en main.py), igual que en catalog y orders.
    logging.getLogger("uvicorn.access").disabled = True
