# ADR-0003: GitHub Flow, PRs con revisión y protección de `main`

- **Estado:** Aceptada, con un punto pendiente
- **Fecha:** 2026-11-04 (protección de `main` completada el 2026-11-16; reglas de revisión y merge confirmadas el 2026-11-04)
- **Decisores:** @Mariocgf, @nicokuro, @RominaRamon10
- **Tareas / PRs:** S1-02, S1-06, S2-18 (#104), #99, #100

> **Pendiente de decisión:** unificar las convenciones. Hoy conviven los prefijos de rama `feat/…` y `feature/…`.

## Contexto

La letra pide una estrategia de ramificación con evidencia de trabajo colaborativo y PRs revisados (rúbrica: "Estrategias de ramificación"). Somos tres, cada PR es chico (una tarea del tablero) y la promoción a Dev, Test y Prod la va a hacer el pipeline.

## Decisión

Usamos **GitHub Flow**: una sola rama de larga vida (`main`) y una rama corta por tarea, que entra a `main` solo por un PR aprobado.

**Ruleset `proteccion-main`** (sin bypass, ni siquiera para los Owners):

- PR obligatorio, con al menos **1 aprobación** de alguien distinto al autor.
- Las aprobaciones se descartan si se hace un push nuevo.
- Todas las conversaciones resueltas antes del merge.
- Checks obligatorios del CI: Gitleaks, Semgrep (desde el PR #112, ver [ADR-0022](ADR-0022-seguridad-del-pipeline.md)), el build de los 6 servicios y `Terraform (fmt y validate)`.
- No se puede borrar `main` ni hacer *force push*.
- **"Require branches to be up to date" desactivado** a propósito (ver consecuencias).

**Convenciones:**

- Ramas: `tipo/descripcion-corta`, por ejemplo `feat/s2-02-dockerfile-orders` o `fix/quitar-puertos-db-redis`.
- Commits y títulos de PR: `tipo: descripción` (estilo Conventional Commits; tipos `feat`, `fix`, `docs`, `ci`, `chore`, `style`, `refactor`, `test`).
- Descripción del PR con la plantilla `.github/pull_request_template.md`: Qué cambia, Cómo probarlo, Checklist y `Closes #N`.
- Reviewers: los otros dos integrantes.

**Revisión y merge (confirmado por el equipo el 20/11):**

| Regla | Por qué | Alternativa descartada |
| --- | --- | --- |
| **Revisa alguien que no trabaja en esa área.** Como hay una persona por área, en la práctica revisa cualquiera de los otros dos. En las tareas compartidas (por ejemplo, CI/CD + Infra), la aprobación tiene que venir del integrante que no participó. | En la defensa, el auditor elige quién responde sobre cada parte. Revisar el trabajo de otra área es la forma más barata de que los tres aprendamos todo. Además, alguien de afuera detecta lo que no está explicado en el PR. | Que revise el que más sabe del tema: revisiones más rápidas, pero el conocimiento queda concentrado en una persona. |
| **El autor hace el merge, después de la aprobación**, y borra la rama. | El autor sabe cuándo su cambio está listo (por ejemplo, si falta un commit o actualizar la rama con `main`) y es quien resuelve un conflicto si aparece. | Que mergee quien aprueba: el PR puede entrar antes de que el autor termine, y quien aprueba tendría que resolver conflictos de un código que no escribió. |

## Alternativas consideradas

| Opción | A favor | En contra |
| --- | --- | --- |
| **GitHub Flow** | Simple; PRs chicos y frecuentes; la promoción entre ambientes la hace el pipeline, así que la imagen que llega a Prod es la misma que se probó | No hay una rama que "muestre" cada ambiente |
| GitFlow (`develop`, `release/*`, `hotfix/*`) | Fácil de dibujar, una rama por ambiente | Pensado para releases planificados y versiones en paralelo, que no tenemos; duplica merges |
| Trunk-Based | Integración continua real | Exige feature flags y una batería de tests madura, que todavía no existe |

## Consecuencias

- **Lo bueno:** nada entra a `main` sin revisión y sin el CI en verde (probado: el push directo se rechaza). Cada PR queda vinculado a su tarea, con comentarios reales de revisión como evidencia (S2-17, S3-17).
- **Lo que aceptamos:**
  - Con "up to date" desactivado, dos PRs que pasan el CI por separado podrían romper al combinarse. Si se activara, junto con "descartar aprobaciones viejas" habría que volver a aprobar cada PR cada vez que entra algo a `main`. Con 3 personas y pocos PRs simultáneos, el riesgo es bajo.
  - Un check marcado como obligatorio antes de que su job exista en `main` deja trabados los PRs en "Expected — Waiting for status" (pasó el 16/11 con el job de Terraform). Regla: se marca como obligatorio **después** de mergear el job, y si se renombra un job hay que actualizar el ruleset.
- **Qué haría cambiar esta decisión:** muchos PRs simultáneos rompiendo `main` (activar "up to date" o una *merge queue*).
- **Pendientes:**
  - Unificar el prefijo de las ramas: hoy conviven `feature/…` y `feat/…`.
  - Diagrama de la estrategia de ramas (S3-18).

## Uso de IA

El análisis de opciones y la redacción de este ADR se hicieron con apoyo de Claude (Anthropic). La estrategia de ramas y las reglas de revisión y merge las decidió el equipo. La protección de `main` se verificó intentando un push directo, que el ruleset rechazó.
