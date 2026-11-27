# ADR-0022: Seguridad del pipeline: Gitleaks, Semgrep y actions fijadas a SHA

- **Estado:** Aceptada (Semgrep y las actions fijadas a SHA entraron con el PR #112; Semgrep es check obligatorio del ruleset)
- **Fecha:** 2026-11-09 (Gitleaks) · 2026-11-22 (Semgrep y actions fijadas a SHA)
- **Decisores:** @nicokuro (propuesta), revisado por @Mariocgf y @RominaRamon10
- **Tareas / PRs:** S1-11 (#86), S2-07 (#112), S2-16, S3-06

## Contexto

Dos observaciones del auditor apuntan al pipeline: hay contraseñas escritas en archivos del repo y no hay ningún análisis de código. La letra (sección 5) pide un detector de secretos que **falle el build** si encuentra credenciales, análisis estático de seguridad (SAST) con hallazgos reportados y cero secretos en el repositorio.

Hay tres condiciones que pesan en la decisión:

- **Tero usa tres lenguajes** (Go, Python y TypeScript), además de Terraform y Dockerfiles. Una herramienta por lenguaje multiplica la configuración y los umbrales a mantener.
- **El historial ya tiene secretos.** Las contraseñas que venían fijas en el código original (base de datos, usuario del admin y secreto de los JWT) están en commits públicos (ver [ADR-0004](ADR-0004-repositorio-publico-en-organizacion-propia.md)). Un escáner que revise todo el historial falla en cada PR por algo que ningún PR nuevo puede arreglar.
- **El pipeline tiene acceso a AWS.** Desde el PR #103, el job de build usa las claves de `tero-ci`, que tiene permisos de administrador (ver [ADR-0015](ADR-0015-credenciales-aws-ci.md)). Cualquier código de terceros que corra en el pipeline, como las actions de GitHub, puede leer esos secretos.

## Decisión

### 1. Detector de secretos: Gitleaks

- Corre en cada PR hacia `main` y en cada push a `main`, como job propio ("Detección de secretos (Gitleaks)"). Es un **check obligatorio** del ruleset.
- Se ejecuta con la **imagen de Docker** de Gitleaks (`docker run ... gitleaks git /repo`), con el mismo comando que en local.
- Baja el repo con todo el historial (`fetch-depth: 0`), porque Gitleaks revisa commits, no solo archivos.
- **Escanea solo el rango de commits nuevos** (`--log-opts`):
  - En un PR: de la base del PR al último commit del PR.
  - En un push a `main`: del commit anterior al nuevo.
- `--redact`: si encuentra un secreto, lo tapa en el log. El repo es público y los logs de Actions también.
- **Regla propia** en `.gitleaks.toml`, además de las reglas por defecto (`useDefault = true`):

  ```toml
  [[rules]]
  id = "password-en-config"
  regex = '''(?i)(?:password|passwd|secret)[a-z_]*[ \t]*[:=][ \t]*["']?([^\s"'$<{][^\s"']{3,})'''
  secretGroup = 1
  path = '''(?i)(\.ya?ml|\.env(\.[a-z]+)?|\.tfvars|\.properties|\.ini|\.conf)$'''
  ```

  Busca `password`, `passwd` o `secret` seguidos de un valor en archivos de configuración. Ignora los valores que empiezan con `$`, `<` o `{`, que son referencias a variables (`${DB_PASSWORD}`) o placeholders, no secretos. Usa `[ \t]*` en vez de `\s*` para que la regla no salte de línea y dé falsos positivos.

### 2. SAST: Semgrep

- Job "Análisis estático (Semgrep)", en cada PR y en cada push a `main`.
- Imagen **fijada**: `semgrep/semgrep:1.186.0`.
- `semgrep scan --config p/default --metrics=off --error`:
  - `p/default` es el ruleset general (693 reglas). Cubre Go, Python, TypeScript, Terraform, Dockerfiles y YAML de una sola pasada.
  - `--error` hace que el job falle con cualquier hallazgo. Sin ese flag, Semgrep termina con éxito aunque encuentre problemas.
  - `--metrics=off`: no se envían estadísticas de uso.
- **Umbral: cero hallazgos**, de cualquier severidad. La primera corrida encontró 9 y se resolvieron todos antes de activar el job:

  | Hallazgo | Dónde | Resolución |
  | --- | --- | --- |
  | Actions con tag mutable (6) | `ci.yml` | Arreglado: fijadas a commit SHA (punto 3) |
  | `missing-integrity` (2) | `index.html` de ui y admin | Falso positivo: el favicon está embebido como `data:` URI, no es un recurso externo. Marcado con `nosemgrep` y el ID de la regla, con la justificación en la línea anterior |
  | Yarn sin *minimal age gate* (1) | `src/checkout/.yarnrc.yml` | Arreglado: `npmMinimalAgeGate: "7d"`, para que Yarn no instale versiones publicadas hace menos de 7 días |

- Un falso positivo se marca **siempre** con el ID de la regla y una justificación al lado, y se acepta en la revisión del PR. Nunca se desactiva una regla para todo el repo.

### 3. Actions de terceros fijadas a commit SHA

Todas las actions del workflow se referencian por el **SHA completo del commit**, con la versión en un comentario:

```yaml
uses: actions/checkout@11d5960a326750d5838078e36cf38b85af677262 # v4.4.0
```

| Action | Versión |
| --- | --- |
| `actions/checkout` | v4.4.0 |
| `aws-actions/configure-aws-credentials` | v4.3.1 |
| `aws-actions/amazon-ecr-login` | v2.2.0 |
| `hashicorp/setup-terraform` | v3.1.2 |

Los SHA se tomaron de los releases oficiales de cada action (`git ls-remote --tags`).

## Alternativas consideradas

### Detector de secretos

| Opción | A favor | En contra |
| --- | --- | --- |
| Gitleaks con `docker run` | El mismo comando en local y en el CI; reglas propias en un archivo versionado | Hay que armar a mano el rango de commits |
| `gitleaks/gitleaks-action` | Más corto de escribir | En repos de una organización pide una licencia; es una action de terceros más con acceso a los secrets |
| TruffleHog | Verifica si los secretos encontrados siguen activos | Otra herramienta para aprender; las reglas propias son menos directas |
| Escanear todo el historial en cada PR | Detecta cualquier secreto, viejo o nuevo | Falla siempre por los secretos del commit inicial, que ningún PR puede sacar (ver S2-16) |
| Solo el secret scanning de GitHub | Gratis en repos públicos, sin configurar nada | No detecta contraseñas genéricas como las del compose; no es un check del pipeline |

### SAST

| Opción | A favor | En contra |
| --- | --- | --- |
| Semgrep | Una herramienta para los tres lenguajes, Terraform y Dockerfiles; rápida; corre igual en local | Las reglas se descargan del registro de Semgrep en cada ejecución |
| Una herramienta por lenguaje (gosec, Bandit, eslint-plugin-security) | Reglas especializadas | Tres configuraciones, tres umbrales y tres reportes que mantener |
| CodeQL | Gratis en repos públicos; análisis profundo | Más lento; no cubre Terraform ni Dockerfiles; no se corre igual en local |
| Solo SonarCloud | Ya lo vamos a usar (S2-06) | Está orientado a calidad; su cobertura de seguridad es más limitada |

### Referencias a las actions

| Opción | A favor | En contra |
| --- | --- | --- |
| Commit SHA | Inmutable: identifica un contenido exacto | Las actualizaciones no llegan solas; menos legible (lo compensa el comentario) |
| Tag de versión mayor (`@v4`) | Recibe parches automáticamente | El dueño de la action puede mover el tag. Si la action se compromete, el pipeline corre código malicioso con acceso a los secrets de AWS (pasó en 2025 con `tj-actions/changed-files`) |
| Tag exacto (`@v4.4.0`) | Más legible | También se puede mover: un tag es solo un nombre |

## Consecuencias

**Lo bueno:**

- Ningún secreto nuevo puede entrar a `main` sin que el pipeline lo frene, y el log no lo expone.
- El código de los tres lenguajes, la infraestructura y los Dockerfiles se analiza en cada PR, con una sola herramienta.
- El pipeline queda protegido de que una action comprometida corra código distinto al revisado.
- Cada hallazgo resuelto o marcado como falso positivo queda documentado en su PR: es la evidencia de "hallazgos reportados" que pide la letra.

**Lo que aceptamos:**

- **Gitleaks no revisa el historial viejo.** Los secretos del commit inicial siguen ahí. Se tratan aparte en S2-16 (baseline o excepción documentada). En cualquier caso, esos valores no se usan en ningún ambiente ([ADR-0021](ADR-0021-secretos-fuera-del-repo-en-local.md)).
- **Las actions fijadas no se actualizan solas.** Hay que subir el SHA a mano o con una herramienta (ver pendientes).
- **Semgrep depende del registro de Semgrep:** si `semgrep.dev` no responde, el job falla aunque el código esté bien. Se resuelve con un re-run.
- Un falso positivo frena el PR hasta que alguien lo justifique y lo marque. Es el precio de tener el umbral en cero.

**Qué haría cambiar esta decisión:** que el docente exija una herramienta concreta, que Semgrep dé demasiados falsos positivos en un lenguaje (se sumaría una herramienta específica para ese lenguaje), o que GitHub permita fijar las actions de forma centralizada para toda la organización.

## Pendientes

- **Fijar también la imagen de Gitleaks:** hoy usa `zricethezav/gitleaks:latest`, que va contra este mismo ADR. Pasar a una versión fija.
- Agregar `needs: gitleaks` al job de build, para no publicar imágenes si Gitleaks falla en `main` (sugerido en el #103; también figura en [ADR-0019](ADR-0019-imagenes-etiquetadas-con-sha.md)).
- Sumar Dependabot (ecosistema `github-actions`) para que las actions fijadas se actualicen con PRs revisables.
- Decidir qué hacer con los secretos del historial (S2-16).
- Activar el *secret scanning* y la *push protection* de GitHub como capa extra: son gratis en repos públicos y frenan el push antes de que llegue al pipeline.
- SCA con Trivy en modo filesystem (S2-08) y escaneo de imágenes con Trivy (S3-06) van en su propio ADR, junto con la interpretación de "sin parche disponible" que se consultó al docente.

## Uso de IA

El análisis de opciones y la redacción de este ADR se hicieron con apoyo de Claude (Anthropic). La decisión la tomó CI/CD. En la revisión del PR #112 se verificaron los SHA de las actions con `git ls-remote` y se corrió Semgrep en local sin hallazgos.
