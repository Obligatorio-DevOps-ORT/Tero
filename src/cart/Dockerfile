# syntax=docker/dockerfile:1

# ---------- Etapa 1: build ----------
# Imagen slim de Python con pip. Acá se instalan las dependencias
# dentro de un entorno virtual (venv) que después copiamos entero.
FROM python:3.12-slim AS build

# No generar archivos .pyc ni guardar la caché de pip (imagen más chica).
ENV PYTHONDONTWRITEBYTECODE=1 \
    PIP_NO_CACHE_DIR=1 \
    PIP_DISABLE_PIP_VERSION_CHECK=1

# Crear el venv y ponerlo primero en el PATH.
RUN python -m venv /opt/venv
ENV PATH="/opt/venv/bin:$PATH"

# 1) Copiar primero solo requirements.txt e instalar.
#    Docker cachea esta capa: si solo cambia el código, no reinstala todo.
COPY requirements.txt .
RUN pip install -r requirements.txt

# ---------- Etapa 2: runtime ----------
FROM python:3.12-slim

ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1 \
    PATH="/opt/venv/bin:$PATH" \
    CART_PERSISTENCE_PROVIDER=in-memory \
    PORT=8080

# Usuario sin privilegios, sin shell de login y sin carpeta home.
RUN groupadd --system --gid 10001 app \
 && useradd --system --uid 10001 --gid app --no-create-home --shell /usr/sbin/nologin app

WORKDIR /app

# Del build solo nos llevamos el venv con las dependencias.
COPY --from=build /opt/venv /opt/venv

# Código de la app (el dueño sigue siendo root: el usuario app puede leerlo pero no modificarlo).
COPY app/ ./app/

# No correr nunca como root.
USER app:app

EXPOSE 8080

# Forma "exec" con sh para poder leer $PORT; exec hace que uvicorn sea el PID 1 y reciba las señales de parada.
CMD ["sh", "-c", "exec uvicorn app.main:app --host 0.0.0.0 --port ${PORT}"]