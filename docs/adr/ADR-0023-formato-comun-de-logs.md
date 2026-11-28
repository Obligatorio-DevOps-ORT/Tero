# ADR-0023: Logs en JSON con un esquema común en los 6 servicios

- **Estado:** Aceptada
- **Fecha:** 2026-11-27
- **Decisores:** @RominaRamon10
- **Tarea / PRs:** S3-13 (#47), PR #119 (revisado por @nicokuro)

## Contexto

La letra (sección 8) pide centralizar los logs de todos los servicios "con formato estructurado que permita búsquedas y filtros". Cada servicio de Tero escribía sus logs en texto y a su manera:

- **catalog y orders (Go + Gin):** el logger de texto de Gin y `log.Fatal`/`fmt.Println`. orders registraba además cada `/health`. GORM escribía sus propios avisos en varias líneas (por ejemplo, un `record not found` con la consulta SQL en cada 404).
- **cart (Python + FastAPI):** solo el *access log* de uvicorn, en texto y sin la duración de cada request.
- **checkout (NestJS):** el logger de Nest, en texto con códigos de color. Registraba la IP y el navegador del cliente.
- **ui y admin (Express):** no registraban ningún request; solo una línea al arrancar.

En CloudWatch (S3-12) cada línea es un evento. Con texto libre, encontrar "qué servicio está devolviendo errores 500" obliga a escribir una expresión distinta por servicio, y un error de varias líneas (un *traceback* de Python, un *stack* de Node) se parte en muchos eventos sueltos.

## Decisión

Todos los servicios escriben **una línea JSON por evento en stdout**, con los mismos nombres de campos:

| Campo | Siempre | Ejemplo |
| --- | --- | --- |
| `time` | Sí | `2026-11-27T22:33:55.354Z` (UTC) |
| `level` | Sí | `INFO`, `WARN` o `ERROR` (en mayúsculas, iguales en los 4 lenguajes) |
| `msg` | Sí | `request`, `servidor escuchando`, `error de conexión con Redis` |
| `service` | Sí | `catalog`, `orders`, `carts`, `checkout`, `ui`, `admin` (el nombre del servicio en compose y ECS) |
| `method`, `path`, `status`, `duration_ms` | En cada request | `POST`, `/checkout/user-x/submit`, `201`, `202.044` |
| `error` | Si hay un error | El mensaje o el *stack* completo, dentro de una sola línea |

Reglas:

- **Nivel según el status:** 5xx → `ERROR`, 4xx → `WARN`, el resto → `INFO`. No se escriben `DEBUG` ni `verbose`.
- **`/health` no se registra:** los health checks de ECS y del ALB llegan cada pocos segundos y no aportan información.
- **No se registran datos sensibles:** ni el body (el login del admin manda la contraseña), ni cookies (el JWT del admin), ni headers, ni la IP o el navegador del cliente. El `path` va sin la *query* (`?page=2`), y las consultas SQL se registran sin sus valores (con `$1`, `$2`…), porque pueden traer datos de los clientes.
- **Los campos comunes no se pueden pisar:** si un mensaje trae un campo extra con el mismo nombre que `time`, `level`, `msg` o `service`, se descarta el extra.
- **Sin librerías nuevas:** se usa lo que trae cada lenguaje.

| Servicio | Cómo |
| --- | --- |
| catalog, orders | `log/slog` (librería estándar de Go) con `JSONHandler`; un middleware de Gin propio; un reemplazo de `gin.Recovery()` que, si un handler entra en pánico, responde 500 y lo registra en JSON (`gin.Recovery()` lo escribe en texto, en varias líneas y con colores); y un logger propio para GORM: `ERROR` para errores de SQL o de conexión, `WARN` para consultas de más de 200 ms, sin `record not found` y sin los valores de las consultas |
| cart | `logging` de Python con un `JsonFormatter` propio; los logs de uvicorn pasan por el mismo formato y su *access log* se reemplaza por un middleware de FastAPI que mide la duración |
| ui, admin | `logger.ts` propio con `JSON.stringify`; middleware que registra cada request al terminar; los errores del proxy de ui (`[HPM]`) y los avisos de Node también pasan a JSON |
| checkout | `JsonLogger` propio que implementa el `LoggerService` de Nest; el campo extra `context` indica qué parte de Nest escribió el mensaje (`HTTP`, `ExceptionsHandler`, `Redis`) |

## Alternativas consideradas

| Opción | A favor | En contra |
| --- | --- | --- |
| **Lo que trae cada lenguaje + un esquema común (elegida)** | Sin dependencias nuevas: no hay que regenerar lockfiles ni suma paquetes que Trivy tenga que escanear. Los mismos campos en los 6 servicios | Hay que escribir y mantener un poco de código en cada servicio (entre 30 y 90 líneas) |
| Librerías de logging (zap en Go, structlog en Python, pino en Node) | Muy usadas en la industria, con más funciones (muestreo, *redaction*) | 4 librerías distintas, cada una con sus propios nombres de campos; más dependencias que actualizar y escanear |
| `ConsoleLogger({ json: true })` de Nest 11 en checkout | Es una línea de código | Usa otros nombres de campos (`timestamp` en milisegundos, `message`, `level: "log"`) y no tiene `service`: un filtro `level = "ERROR"` en CloudWatch no encontraría los errores de checkout |
| Formato de texto con una expresión regular en CloudWatch | No cambia el código | Una expresión distinta por servicio; los errores de varias líneas igual se parten en muchos eventos |
| Escribir los logs en archivos dentro del contenedor | — | ECS (driver `awslogs`) solo toma stdout y stderr; además los contenedores corren como usuario no-root ([ADR-0020](ADR-0020-imagenes-multi-stage-minimas-no-root.md)) |

## Consecuencias

- **Lo bueno:** en CloudWatch Logs Insights se puede consultar todo el sistema con una sola consulta, por ejemplo `filter level = "ERROR" | stats count() by service`, o comparar la latencia de un mismo request en ui y en el servicio que lo atiende. Un error de varias líneas llega como un solo evento con su causa. Un login fallido del admin queda como `WARN 401` en `/auth/login`, sin la contraseña.
- **Lo que aceptamos:**
  - El código de logging está repetido en cada servicio: no se puede compartir porque cada servicio es una imagen separada y los lenguajes son distintos.
  - El `path` incluye ids (del carrito, del producto), así que para agrupar por ruta hay que normalizarlo en la consulta.
  - checkout escribe al arrancar unas 15 líneas `INFO` de Nest (módulos y rutas).
  - Un mismo request aparece en ui (con el prefijo `/api`) y en el servicio que lo atiende.
- **Qué haría cambiar esta decisión:** si hiciera falta correlacionar un request entre servicios (un *trace id* o *request id*), conviene pasar a OpenTelemetry o a librerías que lo soporten, en lugar de seguir agregando código propio.
- **Pendientes:**
  - Enviar los logs a CloudWatch Logs (S3-12).
  - `GIN_MODE=release` en las task definitions de catalog y orders; si no, Gin escribe líneas `[GIN-debug]` en texto (S3-04).
  - Bug de la app encontrado al probar: si la base de datos no responde, admin no devuelve un 500, sino que el proceso entero se cae (los handlers `async` no tienen `try/catch` y Express 4 no los atrapa). Ahora queda registrado en JSON (`promesa rechazada sin manejar`), pero hay que corregirlo en una tarea aparte.

## Uso de IA

El análisis de opciones, el código de logging de cada servicio y la redacción de este ADR se hicieron con apoyo de Claude (Anthropic). Los cambios los aplicó y probó Cont./Obs. con `docker compose` en los 6 servicios: tráfico normal (pedido completo en la tienda; login, productos y órdenes en el admin), errores 404 y 500, arranque sin base de datos, servicios caídos detrás del proxy de ui, Redis caído, y búsqueda de la contraseña del admin en los logs.
