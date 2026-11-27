# ADR-0021: Secretos fuera del repo, *fail fast* y sin puertos de datos en local

- **Estado:** Aceptada
- **Fecha:** 2026-11-11 (*fail fast* del admin el 2026-11-19)
- **Decisores:** @RominaRamon10, revisado por @Mariocgf y @nicokuro
- **Tareas / PRs:** S1-09 (#88), S1-10 (#90), S1-14, S2-21 (#107), S1-11 (#86)

## Contexto

El auditor encontró contraseñas escritas en el repo: la de la base de datos, repetida en cuatro servicios del `docker-compose.yml`; el usuario y la contraseña del admin, también publicados en el README; y el secreto con el que se firman los JWT. Además, `admin/src/app.ts` tenía esos mismos valores como *default*, así que aunque se sacaran del compose, la app arrancaba igual con las contraseñas conocidas. Por último, Postgres (5433) y Redis (6379) estaban publicados en el host.

Esta decisión cubre el entorno **local**. Los secretos en AWS (Secrets Manager o Parameter Store) se deciden en un ADR posterior.

## Decisión

- **Secretos en un `.env`** que no se versiona (`.env` en el `.gitignore`). En el repo queda un **`.env.example`** con los nombres de las variables y sin valores reales.
- **Credenciales fuera del README.**
- **Gitleaks** en el pipeline: falla si detecta un secreto. Es un check obligatorio para mergear.
- ***Fail fast* en el admin**: se quitaron los valores por defecto de `DB_PASSWORD`, `ADMIN_USERNAME`, `ADMIN_PASSWORD` y `ADMIN_JWT_SECRET`. Si falta alguno o está vacío, el admin muestra cuáles faltan y termina con código 1. `DB_HOST`, `DB_PORT`, `DB_USER` y `PORT` conservan su valor por defecto porque no son secretos.
- **Sin puertos publicados para Postgres y Redis**: los servicios los alcanzan por la red interna de compose. Para entrar a mano: `docker compose exec db psql …` o `docker compose exec redis redis-cli`.

## Alternativas consideradas

| Opción | A favor | En contra |
| --- | --- | --- |
| **`.env` + `.env.example` + *fail fast*** | Estándar en compose; cada uno elige sus valores; un error de configuración se ve al arrancar | Cada integrante tiene que mantener su `.env` |
| Dejar valores por defecto "de desarrollo" en el código | Arranca sin configurar nada | Es justo lo que señaló el auditor: si alguien olvida configurar Prod, corre con la contraseña conocida |
| Docker secrets en compose | Los secretos no quedan como variables de entorno | La app lee variables de entorno: habría que cambiar el código de los 6 servicios |
| Cifrar el `.env` en el repo (SOPS, git-crypt) | Un solo archivo para todos | Otra herramienta y otra clave que compartir |

## Consecuencias

- **Lo bueno:** el repo deja de tener secretos nuevos; un admin mal configurado no arranca con valores conocidos, falla y dice qué falta (en compose se ve como un contenedor que se reinicia en bucle; el motivo está en `docker compose logs admin`).
- **Lo que aceptamos:**
  - Los valores viejos siguen en el **historial público** de git (ver [ADR-0004](ADR-0004-repositorio-publico-en-organizacion-propia.md)). No se pueden usar en ningún ambiente. Qué hacer con los hallazgos de Gitleaks en el historial (baseline o excepción documentada) se decide en S2-16.
  - Si alguien cambia `DB_PASSWORD`, tiene que recrear la base con `docker compose down -v`, porque Postgres guarda la contraseña la primera vez que crea el volumen.
- **Pendientes:**
  - En AWS, los secretos se inyectan con el bloque `secrets` de la task definition (S3-03). Ojo: `ADMIN_USERNAME` también hay que pasarlo, y `DB_HOST`/`DB_USER` tienen que apuntar a RDS, porque sus valores por defecto son los de compose.
  - Documentar en el README cómo armar el `.env` (S2-20).

## Uso de IA

El análisis de opciones y la redacción de este ADR se hicieron con apoyo de Claude (Anthropic). Los cambios los hizo y probó Cont./Obs.: el admin falla al arrancar con un secreto vacío y los puertos de Postgres y Redis ya no responden desde el host.
