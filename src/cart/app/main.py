import logging
import time

from fastapi import FastAPI, Request
from prometheus_fastapi_instrumentator import Instrumentator
from app.config import settings
from app.logging_config import configurar_logs
from app.router import router

# Primero los logs: así cualquier mensaje del arranque ya sale en JSON.
configurar_logs()
logger = logging.getLogger("carts")

app = FastAPI(title="Cart Service")


@app.middleware("http")
async def registrar_request(request: Request, call_next):
    """Registra cada request HTTP como una línea JSON (salvo /health)."""
    path = request.url.path
    if path == "/health":
        return await call_next(request)

    inicio = time.perf_counter()
    try:
        response = await call_next(request)
    except Exception:
        # Si el handler explota, igual dejamos registrado el 500 y el motivo.
        logger.exception("request", extra={"campos": {
            "method": request.method,
            "path": path,
            "status": 500,
            "duration_ms": round((time.perf_counter() - inicio) * 1000, 3),
        }})
        raise

    status = response.status_code
    nivel = logging.INFO
    if status >= 500:
        nivel = logging.ERROR
    elif status >= 400:
        nivel = logging.WARNING

    logger.log(nivel, "request", extra={"campos": {
        "method": request.method,
        "path": path,
        "status": status,
        "duration_ms": round((time.perf_counter() - inicio) * 1000, 3),
    }})
    return response


@app.on_event("startup")
def startup():
    if settings.cart_persistence_provider == "postgres":
        from app.postgres_service import PostgresCartService
        app.state.cart_service = PostgresCartService()
        logger.info("usando base de datos postgres", extra={"campos": {
            "host": settings.postgres_host,
            "db": settings.postgres_db,
        }})
    else:
        from app.service import InMemoryCartService
        app.state.cart_service = InMemoryCartService()
        logger.info("usando base de datos en memoria")


@app.get("/topology")
def topology():
    return {"provider": settings.cart_persistence_provider}


@app.get("/health")
def health():
    return {"status": "UP"}


app.include_router(router)
Instrumentator().instrument(app).expose(app, endpoint="/metrics")
