# Observaciones del auditor: seguimiento

Matriz que relaciona cada observación preliminar del auditor (letra del obligatorio,
sección 1) con las tareas del tablero que la resuelven y la evidencia de su resolución.
Estado al 22/11; se completa en S5-15

Estados: ✅ Resuelta · 🟡 Parcial (resuelta en local, falta AWS o parte) · ⏳ Pendiente

## 1. Observaciones del auditor

| # | Observación | Tareas que la resuelven | Evidencia | Estado |
|---|---|---|---|---|
| 1 | Contraseñas de la base y credenciales del admin escritas en archivos del repositorio | S1-09, S1-14, S2-21, S1-11 (Gitleaks), S2-16 (historial), S3-03 (Secrets Manager), S5-14 (probar el gate) | PR [#88](https://github.com/Obligatorio-DevOps-ORT/Tero/pull/88) (`.env` + `.env.example`), PR [#107](https://github.com/Obligatorio-DevOps-ORT/Tero/pull/107) (admin sin valores por defecto), PR [#86](https://github.com/Obligatorio-DevOps-ORT/Tero/pull/86) (Gitleaks en el CI) | 🟡 Falta el historial de git (S2-16) y Secrets Manager en AWS (S3-03) |
| 2 | Imágenes construidas en un solo paso sobre bases completas, procesos como root | S2-01 a S2-05, S2-14, S2-15 | PRs de [#99](…/pull/99), [#101](…/pull/101), [#105](…/pull/105), [#106](…/pull/106), [#108](…/pull/108), [#109](…/pull/109); tabla de tamaños y usuarios en el issue [#31](…/issues/31) (−90 % en total, ninguno como root) | ✅ |
| 3 | Base de datos y Redis publicados en puertos del anfitrión | S1-10, S2-10 (subredes privadas y SG), S3-01, S3-02 | PR [#90](…/pull/90) (local) | 🟡 Falta AWS |
| 4 | Sin pruebas automatizadas ni análisis de código ni de dependencias | S2-06 (Sonar), S2-07 (Semgrep), S2-08 (Trivy fs), S3-06 (Trivy image), S3-15 y S3-16 (Postman/Newman), S4-15, S4-16, S5-14 | — | ⏳ |
| 5 | No se sabe qué versión corre ni se puede repetir un despliegue | S3-07 (tag por SHA), S1-12, S1-13, S2-11, S2-13, S4-18, S3-09, S3-19, S4-01 a S4-04, S4-10 (rollback) | PR [#103](…/pull/103) (imágenes con el SHA), [#85](…/pull/85), [#92](…/pull/92), [#102](…/pull/102), [#100](…/pull/100) | 🟡 Falta despliegue y rollback |
| 6 | Algunos servicios exponen métricas pero nadie las recolecta | S3-12, S3-13, S4-05, S4-06, S4-07, S4-13, S4-19, S5-02 | — | ⏳ |

## 2. Deudas técnicas y de seguridad que el auditor no mencionó

| # | Hallazgo | Dónde se detectó | Tareas | Evidencia | Estado |
|---|---|---|---|---|---|
| A | Admin con valores por defecto de secretos en `app.ts` | Revisión de S1-09 | S2-21 | PR [#107](…/pull/107) | ✅ |
| B | Los secretos viejos siguen en el historial público de git | Análisis de Gitleaks | S2-16 | — | ⏳ |
| C | carts no reintenta la conexión a Postgres al arrancar (en ECS no existe `depends_on`) | Issue [#7](…/issues/7) (S1-07) | S3-04 (health check + reinicio de ECS) | — | ⏳ |
| D | ui y admin sin `package-lock.json`, `yarn.lock` de checkout desincronizado (builds no reproducibles) | S2-04 y S2-05 | S2-04, S2-05 | PRs [#106](…/pull/106), [#108](…/pull/108) | ✅ |
| E | Admin entra a las 3 bases con el mismo usuario (sin mínimo privilegio) | Propuesta de plan | ADR (S5-04) | — | ⏳ |
| F | Credenciales del CI de larga duración (`tero-ci`) | PR [#102](…/pull/102) | S3-08, S4-17 (rotación) | — | ⏳ |
| G | Retención de 5 imágenes en ECR puede borrar la de Prod y achicar el rollback | PR [#102](…/pull/102) | S2-11 (pendiente), S4-10 | — | ⏳ |
| H | Imágenes base con tag de versión pero sin digest (build no 100 % reproducible) | Revisión de S2-01 | ADR (S5-04) | — | ⏳ |
