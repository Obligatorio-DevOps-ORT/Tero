# ADR-0019: Imágenes etiquetadas con el SHA del commit, publicadas solo desde `main`

- **Estado:** Aceptada
- **Fecha:** 2026-11-19
- **Decisores:** @nicokuro, revisado por @RominaRamon10
- **Tareas / PRs:** S1-11 (#86), S2-09 (#87), S3-07 (#103)

## Contexto

Una de las observaciones del auditor es que "nadie sabe qué versión está corriendo" y que no hay despliegues repetibles. En la defensa, una tarea típica es volver a una versión anterior. Necesitamos poder ir de una imagen al commit exacto que la generó, y al revés.

## Decisión

- El pipeline construye las 6 imágenes con una **matrix** de GitHub Actions.
- Cada imagen se publica en ECR con el **SHA completo del commit** como tag: `<registry>/tero-<servicio>:<sha>`.
- Se publica **la misma imagen que se construyó y probó** en el job (`docker tag` + `docker push`), sin reconstruirla.
- **Solo se publica en `push` a `main`**, es decir, después de un PR aprobado. En los PRs, los pasos de publicación aparecen como *skipped*.
- Al promover a Prod se agrega un tag semántico (`v1.0.0`) a la misma imagen.
- El workflow tiene `permissions: contents: read`: cualquier permiso que no se nombre queda en "ninguno". Un job que necesite más lo pide en su propio bloque (mínimo privilegio).

## Alternativas consideradas

| Opción | A favor | En contra |
| --- | --- | --- |
| **SHA del commit** | Único, automático, lleva directo al código; compatible con tags inmutables | No dice nada para una persona; se complementa con `v*` |
| `latest` | Simple | No dice qué versión corre; incompatible con `IMMUTABLE`; rollback imposible |
| Número de ejecución del pipeline | Corto y ordenado | No lleva al código sin pasar por GitHub Actions |
| Solo versión semántica | Legible | Requiere decidir la versión en cada merge |
| Publicar también desde los PRs | Se pueden probar imágenes antes del merge | Sube código no aprobado; en un repo público los PRs desde forks no reciben los secrets |

## Consecuencias

- **Lo bueno:** el rollback es redesplegar el tag anterior; desde la consola de ECS se sabe exactamente qué commit corre.
- **Lo que aceptamos:**
  - Cada merge publica las 6 imágenes aunque cambie un solo servicio. Esto consume la retención de ECR más rápido (ver [ADR-0013](ADR-0013-configuracion-repositorios-ecr.md)).
  - Con tags inmutables, si un run falla a mitad de camino y se re-ejecuta, la imagen reconstruida tiene otro digest y ECR rechaza el push del mismo tag en los servicios que ya se habían publicado. Hay que tenerlo en cuenta al diagnosticar el pipeline.
- **Pendientes:**
  - Agregar `needs: gitleaks` al job de build para no publicar si Gitleaks falla en `main` (sugerido en el #103).
  - Trivy en modo imagen antes de publicar (S3-06).
  - Verificar la primera ejecución en `main` y que estén las 6 imágenes en ECR con el SHA.

## Uso de IA

La redacción de este ADR se hizo con apoyo de Claude (Anthropic). La decisión la tomó CI/CD y se revisó en el PR #103, donde los pasos de publicación aparecen como *skipped* en los PRs y solo corren en push a `main`.
