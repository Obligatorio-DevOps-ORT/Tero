# ADR-0004: Repositorio público en una organización propia

- **Estado:** Aceptada
- **Fecha:** 2026-11-04
- **Decisores:** @Mariocgf, @nicokuro, @RominaRamon10
- **Tareas / PRs:** S1-02, S1-03

## Contexto

El código de Tero está en un repo **privado** de la organización del curso (`ORT-ATI-CertificadoDevOps/Tero`). La forma obvia de trabajar sería un fork, pero:

- El fork de un repo privado no se puede hacer público.
- En repos privados de cuentas gratuitas, las reglas de protección de ramas y los GitHub Environments requieren un plan pago.
- SonarCloud es gratis solo para repos públicos.

Además, la rúbrica de "Gestión de código" evalúa un repositorio bien organizado con separación de componentes.

## Decisión

Creamos la organización **`Obligatorio-DevOps-ORT`** (los tres como Owners) y dentro un repo **público** `Tero`, copiado del original **con todo su historial** (`git clone --bare` + `git push --mirror`), no como fork. Es un **monorepo**: el código de los 6 servicios (`src/`), la infraestructura (`infra/`), los workflows (`.github/`) y la documentación (`docs/`) viven juntos.

## Alternativas consideradas

| Opción | A favor | En contra |
| --- | --- | --- |
| **Repo público nuevo en una organización** | Protección de ramas, Environments y SonarCloud sin costo; los tres con los mismos permisos; el tablero es de la organización | El código y su historial quedan públicos |
| Fork privado | Un clic; mantiene el vínculo con el original | Sin protección de ramas ni Environments sin plan pago (GitHub Pro del Student Pack, por persona); sin SonarCloud gratis |
| Repo privado en la cuenta de uno | Privado | Un solo dueño; mismos límites del plan gratuito |
| Varios repos (uno por servicio + uno de infra) | Separación fuerte | Más pipelines y permisos que mantener; para 3 personas y 5 semanas no compensa |

## Consecuencias

- **Lo bueno:** todas las funciones de GitHub que necesitamos están disponibles; el historial conserva el "antes" (útil como evidencia de lo que encontró el auditor).
- **Lo que aceptamos:**
  - Las contraseñas que venían fijas en el código original (base de datos, usuario del admin y secreto de los JWT) siguen en el historial, que ahora es público. Por eso esos valores tienen que cambiarse y no reutilizarse nunca en ningún ambiente (S2-16, ver [ADR-0021](ADR-0021-secretos-fuera-del-repo-en-local.md)).
  - En un repo público, los PRs que vienen de forks no reciben los secrets: el pipeline publica imágenes solo en push a `main` (ver [ADR-0019](ADR-0019-imagenes-etiquetadas-con-sha.md)).
  - No recibimos cambios del repo original automáticamente; si el curso lo actualiza, hay que traer los cambios a mano.

## Uso de IA

El análisis de opciones y la redacción de este ADR se hicieron con apoyo de Claude (Anthropic). La decisión la tomó el equipo, y el repo se copió con todo su historial (`git clone --bare` + `git push --mirror`).
