# ADR-0001: Registrar las decisiones con ADRs

- **Estado:** Aceptada
- **Fecha:** 2026-11-20
- **Decisores:** @Mariocgf, @nicokuro, @RominaRamon10
- **Tareas / PRs:** S1-15 (#15), S5-04

## Contexto

La letra (sección 9) pide documentar las decisiones de diseño más relevantes, en formato ADR o libre. En la defensa, la parte "Explicar" consiste en que el auditor elige dos puntos y pregunta por qué se resolvieron así y qué pasaría si cambiaran las condiciones; además, cualquiera de los tres puede tener que responder.

Hasta ahora las decisiones quedaron repartidas en la propuesta de plan, en documentos sueltos y en las conversaciones de los PRs. Eso sirve mientras uno se acuerda, pero no para que otra persona (el auditor o un compañero de otra área) entienda el porqué.

## Decisión

Cada decisión relevante se registra como un ADR en Markdown, en `docs/adr/`, **un archivo por decisión** y con una estructura común: Contexto → Decisión → Alternativas consideradas → Consecuencias. Los ADRs entran por PR como cualquier otro cambio.

Todos se nombran igual: `ADR-` + un número correlativo de 4 cifras + un nombre corto (`ADR-0001-registrar-decisiones-con-adrs.md`), y el título del archivo repite el número (`# ADR-0001: …`). La tarea del tablero de la que salió cada decisión figura dentro del ADR. El índice (`README.md`) los reúne a todos.

## Alternativas consideradas

| Opción | A favor | En contra |
| --- | --- | --- |
| **Un archivo por ADR en `docs/adr/`** | Cada PR toca un archivo distinto, así que no hay conflictos de merge entre los tres; se puede enlazar a una decisión puntual; es el formato más usado (Michael Nygard, `adr-tools`) | Hay que mantener el índice |
| Un único archivo `docs/ADRs.md` | Todo en un lugar | Si dos personas agregan un ADR a la vez, conflicto seguro; el archivo crece sin límite |
| Wiki de GitHub | Fácil de editar | Queda fuera del repo: no se revisa por PR, no va en el `.zip` de la entrega sin exportarla |
| Formato libre dentro del README | Nada nuevo que aprender | Mezcla instrucciones de uso con historia de decisiones; difícil de mantener |

## Consecuencias

- **Lo bueno:** la documentación vive al lado del código y se versiona con él; el PR de cada ADR deja evidencia de que la decisión se discutió (suma a "Estrategias de ramificación" y "Gestión de proyecto"); sirve como guion para el simulacro "Explicar" (S5-08).
- **Lo que aceptamos:** escribir un ADR lleva tiempo. Para no frenar el trabajo, una decisión puede implementarse primero y documentarse después (como pasó con las de Infra), pero antes de la entrega.
- **Pendientes:** sumar `docs/` a la "Estructura del repositorio" del README (S2-20).

## Uso de IA

La redacción de este ADR se hizo con apoyo de Claude (Anthropic). La decisión de usar ADRs, un archivo por decisión, la tomó el equipo en S1-15.
