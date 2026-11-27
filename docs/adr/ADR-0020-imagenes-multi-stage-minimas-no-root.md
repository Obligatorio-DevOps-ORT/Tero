# ADR-0020: Imágenes multi-stage, base mínima según el lenguaje y usuario no-root

- **Estado:** Aceptada
- **Fecha:** 2026-11-16 (aplicada a los 6 servicios hasta el 2026-11-20)
- **Decisores:** @RominaRamon10, revisado por @Mariocgf y @nicokuro
- **Tareas / PRs:** S2-01 (#17), S2-02 (#18), S2-03 (#105), S2-04 (#106), S2-05 (#108), S2-14, S2-15

## Contexto

Otra observación del auditor: las imágenes se construyen en un solo paso, traen todo (compiladores, herramientas de desarrollo, código fuente) y corren como `root`. Eso las hace enormes (de 426 MB a 499 MB comprimidas) y amplía la superficie de ataque: si alguien compromete el contenedor, tiene permisos de administrador y herramientas a mano.

Tero usa tres lenguajes, y cada uno necesita cosas distintas en tiempo de ejecución.

## Decisión

Todos los Dockerfiles son **multi-stage**: una etapa compila o instala, y la imagen final copia solo lo necesario para correr.

| Servicios | Construcción | Imagen final | Usuario |
| --- | --- | --- | --- |
| catalog, orders (Go) | `golang:1.24-alpine`, binario estático | `distroless/static-debian12:nonroot` | `nonroot` |
| cart (Python) | `python:3.12-slim`, dependencias en un venv | `python:3.12-slim` con el venv | `app` (UID 10001) |
| ui, admin (Node) | `node:22-alpine`: `npm ci` + `tsc`, y otra etapa con `npm ci --omit=dev` | `node:22-alpine` con `dist/` y dependencias de producción | `node` (UID 1000) |
| checkout (NestJS, Yarn 4) | `node:22-alpine`: `yarn install --immutable` + build, y `yarn workspaces focus --production` | `node:22-alpine` | `node` (UID 1000) |

Además:

- Cada servicio tiene su `.dockerignore`.
- El proceso arranca directo (`node dist/app.js`, `exec uvicorn`, el binario), sin `npm start`, para que reciba las señales de parada de ECS.
- Instalación reproducible: `npm ci` y `yarn install --immutable` con el lock file commiteado.
- Imágenes base con tag de versión, nunca `latest`.

Resultados (comprimida, antes → después): catalog 570 MB → 9,33 MB · orders 499 → 6,5 MB · cart 447 → 62,7 MB · ui 426 → 62,8 MB · admin 426 → 62,3 MB · checkout 466 → 72,4 MB. Go baja ~98 %; Python y Node, ~84–86 %.

## Alternativas consideradas

| Opción | A favor | En contra |
| --- | --- | --- |
| **Base mínima según el lenguaje** | Lo más chico posible para cada caso | Tres patrones distintos para mantener |
| Distroless para todos | Mínima superficie | Python y Node necesitan el intérprete; las variantes distroless de esos lenguajes complican depurar y no siempre tienen la versión que usamos |
| `scratch` para Go | Aún más chica | Sin certificados raíz, sin zona horaria y sin usuario `nonroot` ya definido |
| Alpine para Python | Más chica que slim | musl en lugar de glibc: muchas dependencias no tienen versión precompilada y hay que compilarlas |

## Consecuencias

- **Lo bueno:** imágenes más chicas (descargas y despliegues más rápidos, menos espacio en ECR) y con menos paquetes donde puedan aparecer vulnerabilidades; un atacante dentro del contenedor no es `root`.
- **Lo que aceptamos:**
  - Las imágenes distroless no tienen shell: `docker compose exec catalog sh` no funciona y el *health check* de ECS no puede usar `curl` dentro del contenedor. Para depurar existe la variante `:debug-nonroot`.
  - Las de Python y Node tienen shell pero no `curl`: el *health check* tiene que usar `wget` (Node) o el propio Python (cart).
  - catalog compilado sin CGO no puede usar el modo `in-memory` (SQLite). Tero siempre usa Postgres, así que no afecta.
  - Python queda fijo en 3.12: `psycopg2-binary` 2.9.9 no tiene versión precompilada para 3.13.
  - Los tags de versión reciben parches, así que dos builds del mismo commit pueden dar imágenes distintas.
- **Pendientes:** decidir si se fijan las imágenes base por digest (`@sha256:…`) y cómo se actualizan (a mano, Dependabot o Renovate). Va en un ADR aparte.

## Uso de IA

El análisis de opciones y la redacción de este ADR se hicieron con apoyo de Claude (Anthropic). Los Dockerfiles los escribió y probó Cont./Obs.: cada imagen se verificó corriendo Tero completo (tienda, pedido y admin), comprobando que el contenedor no corre como `root` y midiendo el tamaño antes y después.
