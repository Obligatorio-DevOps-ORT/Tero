# Registro de decisiones de arquitectura (ADRs)

Un ADR (*Architecture Decision Record*) documenta **una** decisión importante: qué problema había, qué se decidió, qué otras opciones había y qué consecuencias aceptamos. La letra del obligatorio (sección 9) pide documentar las decisiones de diseño más relevantes, y en la defensa el auditor pregunta "¿por qué lo resolvieron así y qué pasaría si cambiaran las condiciones?". Cada ADR intenta responder exactamente esa pregunta.

Todos los archivos siguen el mismo patrón: `ADR-` + número correlativo de 4 cifras + nombre corto (por ejemplo, `ADR-0010-salida-internet-subredes-privadas.md`). Dentro de cada ADR figura la tarea del tablero de la que salió.

Un ADR con estado **"Aceptada, con un punto pendiente"** tiene la decisión principal tomada e implementada, pero una parte sin confirmar. Esa parte está marcada al principio del ADR con un recuadro **"Pendiente de decisión"**, que se borra cuando se confirma.

## Índice

### Gestión y código

| ADR | Decisión | Estado |
| --- | --- | --- |
| [ADR-0001](ADR-0001-registrar-decisiones-con-adrs.md) | Registrar las decisiones con ADRs | Aceptada |
| [ADR-0002](ADR-0002-kanban-en-github-projects.md) | Kanban con límites WIP en GitHub Projects | Aceptada |
| [ADR-0003](ADR-0003-github-flow-y-proteccion-de-main.md) | GitHub Flow, PRs con revisión y protección de `main` | Aceptada, con un punto pendiente |
| [ADR-0004](ADR-0004-repositorio-publico-en-organizacion-propia.md) | Repositorio público en una organización propia | Aceptada |

### Infraestructura (Terraform y AWS)

| ADR | Decisión | Estado |
| --- | --- | --- |
| [ADR-0005](ADR-0005-estado-remoto-s3.md) | Estado remoto de Terraform en S3, cifrado y con bloqueo nativo | Aceptada |
| [ADR-0006](ADR-0006-workspaces-por-ambiente.md) | Un estado por ambiente con workspaces y guard de ambiente | Aceptada |
| [ADR-0007](ADR-0007-version-terraform.md) | Versión mínima de Terraform 1.17 | Aceptada |
| [ADR-0008](ADR-0008-modulo-de-red.md) | Módulo de red propio, con subredes públicas y privadas en 2 zonas | Aceptada |
| [ADR-0009](ADR-0009-tfvars-por-ambiente.md) | Valores por ambiente con `.tfvars`, combinados con workspaces | Aceptada |
| [ADR-0010](ADR-0010-salida-internet-subredes-privadas.md) | Un NAT Gateway por ambiente, apagado por defecto | Aceptada |
| [ADR-0011](ADR-0011-security-groups.md) | Cuatro security groups con reglas entre SGs | Aceptada |
| [ADR-0012](ADR-0012-repositorios-ecr-compartidos.md) | Repositorios ECR compartidos en un root separado | Aceptada |
| [ADR-0013](ADR-0013-configuracion-repositorios-ecr.md) | Configuración de ECR: tags inmutables, escaneo y retención | Aceptada, con un punto pendiente |
| [ADR-0014](ADR-0014-cuenta-aws-propia.md) | Cuenta AWS propia en `us-east-2` en lugar del Learner Lab | Aceptada, con un punto pendiente |
| [ADR-0015](ADR-0015-credenciales-aws-ci.md) | Credenciales del CI: un usuario IAM `tero-ci` con claves de acceso y permisos de administrador | Aceptada |
| [ADR-0016](ADR-0016-computo-fargate.md) | Cómputo del cluster ECS con Fargate | Aceptada |
| [ADR-0017](ADR-0017-capacity-providers.md) | Fargate Spot en dev y test; en prod, 1 task on-demand fija | Aceptada |
| [ADR-0018](ADR-0018-container-insights.md) | Container Insights activado solo en prod | Aceptada |

### CI/CD, contenedores y seguridad

| ADR | Decisión | Estado |
| --- | --- | --- |
| [ADR-0019](ADR-0019-imagenes-etiquetadas-con-sha.md) | Imágenes etiquetadas con el SHA del commit, publicadas solo desde `main` | Aceptada |
| [ADR-0020](ADR-0020-imagenes-multi-stage-minimas-no-root.md) | Imágenes multi-stage, base mínima según el lenguaje y usuario no-root | Aceptada |
| [ADR-0021](ADR-0021-secretos-fuera-del-repo-en-local.md) | Secretos fuera del repo, *fail fast* y sin puertos de datos en local | Aceptada |
| [ADR-0022](ADR-0022-seguridad-del-pipeline.md) | Seguridad del pipeline: Gitleaks, Semgrep y actions fijadas a SHA | Aceptada |

### Decisiones pendientes (se documentan cuando se tomen)

**Bloquean el hito S3-10 (catalog desplegado en Dev):**

- Base de datos y Redis: RDS + ElastiCache o tareas de ECS. Los security groups (ADR-0011) ya asumen RDS y ElastiCache.
- Descubrimiento de servicios en ECS: Service Connect o Cloud Map. El nombre de `cart` en la red tiene que ser `carts` (ADR-0012).
- Secretos en AWS: Secrets Manager o SSM Parameter Store. Varios ADRs de Infra ya dan por hecho Secrets Manager.

**Puntos pendientes dentro de ADRs ya aceptados:**

- Convenciones de ramas (ADR-0003).
- Retención de imágenes en ECR (ADR-0013).
- Acceso del equipo a la consola de AWS (ADR-0014).

**Se pueden decidir más adelante:**

- Acceso al admin: restricción por IP o API Gateway (ADR-0011).
- Servicio serverless y su propósito.
- Permisos del admin sobre las tres bases (mínimo privilegio).
- Promoción entre ambientes (GitHub Environments, aprobación manual en prod).
- SCA y escaneo de imágenes con Trivy, y criterio de "sin parche disponible" (S2-08, S3-06).
- Calidad con SonarCloud y sus umbrales (S2-06).
- Imágenes base fijadas por digest y cómo se actualizan (ver [ADR-0020](ADR-0020-imagenes-multi-stage-minimas-no-root.md)).
- Qué hacer con los secretos del historial de git (S2-16).

## Cómo agregar un ADR

1. Copiá [plantilla.md](plantilla.md) con el número siguiente y un nombre corto: `ADR-0023-base-de-datos-rds-elasticache.md`. Poné la tarea del tablero en el campo **Tarea / PRs**.
2. Completalo y sumalo al índice de arriba.
3. Abrilo por PR como cualquier otro cambio, con los otros dos como reviewers. La discusión del PR es parte de la decisión.

Reglas:

- **Un ADR no se edita para cambiar la decisión.** Si la decisión cambia, se escribe uno nuevo que la reemplaza y en el viejo solo se cambia el estado a "Reemplazada por …". Así queda la historia de por qué cambiamos de idea, que es justo lo que el auditor puede preguntar.
- Se pueden corregir errores de redacción, agregar datos o referencias sin escribir uno nuevo.
- Estados posibles: **Propuesta** (en discusión), **Aceptada**, **Aceptada, con un punto pendiente**, **Rechazada**, **Reemplazada por …**.
- Cuando se confirma un punto pendiente, se completa en el mismo ADR (no es cambiar la decisión, es cerrar lo que el ADR dejaba abierto): se borra el recuadro, se pasa el estado a **Aceptada** y se actualiza el índice.

## Uso de IA

Estos ADRs se redactaron con apoyo de Claude (Anthropic), a partir de la propuesta de plan, los PRs y los documentos del equipo. El contenido fue revisado por el equipo.
