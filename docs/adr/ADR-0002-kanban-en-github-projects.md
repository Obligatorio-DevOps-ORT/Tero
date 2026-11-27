# ADR-0002: Kanban con límites WIP en GitHub Projects

- **Estado:** Aceptada
- **Fecha:** 2026-11-04
- **Decisores:** @Mariocgf, @nicokuro, @RominaRamon10
- **Tareas / PRs:** S1-02, S1-04, S1-05 (#84)

## Contexto

La letra pide usar una metodología ágil (Kanban o Scrum) con un tablero que incluya tareas de infraestructura, seguridad, testing y observabilidad, y documentar su estado en al menos tres momentos (inicio, mitad y cierre).

El equipo son tres estudiantes con horarios distintos, sin un Product Owner real, con unas cinco semanas hasta la entrega y tareas de tamaños muy distintos (desde "quitar un puerto del compose" hasta "desplegar en ECS"). La mayor parte del trabajo es de infraestructura y operaciones, donde el teórico ubica a Kanban como el enfoque más natural.

## Decisión

Usamos **Kanban** en **GitHub Projects** (tablero "Tareas Obligatorio DevOps", proyecto de la organización):

- Columnas: Backlog → Por hacer → En curso → En revisión → Hecho.
- Límites WIP: **En curso 6** (3 personas × 2 tarjetas) y **En revisión 3**. GitHub no bloquea al pasarse, así que la regla es nuestra: si la columna está llena, primero se ayuda a terminar o a revisar ("terminar antes de empezar").
- Cada tarea es un issue (`S1-01 …`), con etiqueta de área (`area: infra`, `area: ci-cd`, `area: cont-obs`, `area: todos`), etiqueta de tema (una por sección de la letra) y los campos Prioridad y Semana.
- Vistas: semanal y una por área.
- Cada PR lleva `Closes #N`: al mergear, la issue se cierra y la tarjeta pasa sola a Hecho.
- Tres capturas del tablero como tareas propias (S1-05, S3-14, S5-10).

**Reglas de trabajo (confirmadas por el equipo el 20/11):**

| Regla | Por qué |
| --- | --- |
| **Máximo 2 tarjetas por persona en En curso.** Antes de empezar una tarea, cada uno se la asigna. Si la columna está llena (o ya tenés 2), primero se ayuda a terminar o a revisar lo que hay. | Es la idea central de Kanban: terminar antes de empezar. Tener muchas tareas abiertas a la vez hace que todas avancen lento y que nada llegue a Hecho. Como GitHub cuenta el límite por columna y no por persona, asignarse la tarea es lo que permite controlarlo. |
| **Las tareas sin código** (configurar GitHub, reuniones, capturas que no se suben al repo) **pasan a Hecho cuando otro integrante verifica que están bien.** Quien verifica deja un comentario en la issue y la cierra. | No tienen un PR que las cierre con `Closes #N`, así que no se mueven solas. Pedir que otro las verifique mantiene el mismo control que tienen las tareas con código: nadie se aprueba su propio trabajo. |

## Alternativas consideradas

| Opción | A favor | En contra |
| --- | --- | --- |
| **Kanban** | Flujo continuo, sin roles ni ceremonias obligatorias; encaja con tareas de tamaños distintos | Sin sprints, el riesgo es perder el ritmo: lo compensa una reunión semanal |
| Scrum | Sprint Goal y ritmo fijo | Exige PO y Scrum Master que no tenemos; Planning, Daily, Review y Retro pesan mucho para 3 personas en 5 semanas |
| Jira / Trello | Más funciones de gestión | Quedan fuera de GitHub: se pierde el vínculo automático tarea → PR → código |

## Consecuencias

- **Lo bueno:** trazabilidad completa de tarea a código, que el auditor puede seguir desde la tarjeta hasta el commit; el estado del tablero se actualiza solo al mergear.
- **Lo que aceptamos:**
  - GitHub no bloquea al pasarse del límite: solo lo marca en la columna. La regla se cumple por acuerdo, no por la herramienta.
  - Una tarea sin código puede quedar esperando en En revisión hasta que alguien la verifique. Se repasa en la reunión semanal.
  - Sumamos una reunión semanal corta para repriorizar (Kanban con una práctica de Scrum, a veces llamado "Scrumban").
- **Qué haría cambiar esta decisión:** si el equipo creciera o hubiera un cliente que priorizara por iteraciones.
- **Pendientes:** fijar el día de la reunión semanal (S1-16).

## Uso de IA

El análisis de opciones y la redacción de este ADR se hicieron con apoyo de Claude (Anthropic). La metodología la eligió el equipo al empezar el obligatorio, y las reglas de trabajo (WIP por persona y tareas sin código) se confirmaron en equipo el 04/11. El tablero, las etiquetas y los campos se crearon con scripts de `gh` revisados por el equipo.
