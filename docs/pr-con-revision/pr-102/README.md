# Evidencia: PR con revisión entre integrantes (PR #102)

Evidencia para la sección 6.3 de la letra del obligatorio: colaboración mediante
Pull Requests con revisiones de código entre integrantes. Tarea S2-17 (#33).

## Datos del PR

| Campo | Valor |
| --- | --- |
| PR | [#102 – feat(infra): repositorios ECR compartidos (S2-11)](https://github.com/Obligatorio-DevOps-ORT/Tero/pull/102) |
| Tarea | S2-11 (cierra #27) |
| Autor | @Mariocgf (Infra) |
| Revisores | @nicokuro (CI/CD) y @RominaRamon10 (Cont./Obs.) |
| Rama | `feature/s2-11-ecr` → `main` (2 commits, 12 archivos) |
| Checks | 8 en verde: Gitleaks, build de los 6 servicios y Terraform (fmt y validate) |
| Resultado | Aprobado por los dos revisores y mergeado el 18/11/2026 (commit `d89af24`) |

Los dos revisores son de áreas distintas a la del autor. La regla de protección
de `main` (`proteccion-main`) no deja mergear sin PR, sin al menos una aprobación
de otra persona y sin los checks en verde.

## Qué propuso el autor

El PR crea los 6 repositorios de ECR con un módulo reutilizable
(`infra/modules/ecr`) en un root aparte (`infra/global`) y explica cada decisión:
repos compartidos entre ambientes, tags inmutables, lifecycle de 5 imágenes,
`scan_on_push` y `force_delete` como variable. También pasa el backend de
Terraform a una cuenta AWS propia (`us-east-2`), en un commit separado por si
hay que revertirlo. Incluye cómo lo probó y qué queda pendiente.

## Qué se discutió en la revisión

**@nicokuro (CI/CD)** validó `infra/global` en local (`init -backend=false` +
`validate`), se ofreció a sumarlo al job de Terraform del CI y dejó tres
comentarios no bloqueantes:

- **Lifecycle de 5 imágenes:** después de 5 merges sin promover a prod, ECR
  borraría la imagen que corre en prod y ECS no podría volver a bajarla si
  reinicia la tarea. También achica la ventana de rollback. Propuso subir el
  número o conservar las imágenes con tag de versión (`v*`) antes de desplegar.
- **Cambio de cuenta:** llevarlo a la reunión y a un ADR: de quién es la cuenta,
  cómo se manejan los costos fuera del free tier y qué pasa si el dueño no está.
- **Credenciales del CI:** a futuro, usar OIDC en lugar de claves guardadas (lo
  recomienda el teórico de IAM). Mientras tanto, limitar los permisos de
  `tero-ci` a ECR y a lo que necesite Terraform.

**@RominaRamon10 (Cont./Obs.)** sumó tres puntos:

- Después del merge, cada integrante tiene que correr `terraform init -reconfigure`
  en `infra/`, porque cambió el bucket del backend.
- Crear una alerta en AWS Budgets (mail al 50 % y al 80 % de un monto acordado)
  mientras se define lo de la cuenta.
- Se comprometió a escribir el ADR del cambio de cuenta.

## Respuestas del autor

Mario respondió primero por WhatsApp. A pedido de @nicokuro (en la revisión del
PR #110), dejó las respuestas como comentario en el PR #102 el 23/11/2026, para
que la ida y vuelta quede registrada en GitHub:

- **Lifecycle de 5 imágenes:** se eligió 5 solo como margen para ahorrar en
  almacenamiento; se puede subir sin problema.
- **Cuenta AWS:** es una cuenta del plan gratuito, sin tarjeta vinculada. No
  cobra: cuando se terminan los créditos o el período gratuito, la cuenta se
  cierra. Hoy tiene 120 USD de créditos y 177 días restantes (hasta el
  12/05/2027), con 0 USD gastados en el mes. El acceso es con un correo que
  Mario creó para el proyecto; ofreció compartirlo o probar las invitaciones de
  AWS para que el resto del equipo tenga acceso.
- **Credenciales del CI:** se dejaron claves de acceso, igual que en el Learner
  Lab, para poder volver ahí si hace falta. Hoy `tero-ci` tiene permisos de
  administrador porque se usó para probar; **hay que reducirlos** a los
  equivalentes al LabRole. Se mantiene un solo usuario IAM porque el Lab solo
  permite uno y así es más fácil volver.
- **`terraform init -reconfigure`:** reconoció que faltaba documentarlo.
- **AWS Budgets:** la cuenta ya tenía una alerta, con un presupuesto mensual de
  50 USD.

## Qué salió de la revisión

- **PR #103** (de nicokuro): publica las imágenes en ECR con el SHA del commit al
  mergear a `main`, usando las credenciales que dejó este PR.
- **Issue #70** (S4-18): sumar `infra/global` al job de Terraform del CI.
- **ADR del cambio de cuenta AWS** (S1-15), a cargo de Romina, con los datos de
  las respuestas del autor.
- **Retención de imágenes:** se puede subir de 5; queda para decidir antes de
  desplegar.
- **Permisos de `tero-ci`:** pasar de administrador a lo mínimo necesario
  (S3-08).

## Capturas

### 1. Encabezado, revisores y decisiones del autor
![Encabezado del PR #102](01-encabezado-y-decisiones-2026-11-22.png)

### 2. Cambio de cuenta AWS y credenciales del CI
![Cambio de cuenta y credenciales](02-cambio-de-cuenta-y-credenciales-2026-11-22.png)

### 3. Cómo lo probó el autor y pendientes
![Pruebas y pendientes](03-pruebas-y-pendientes-2026-11-22.png)

### 4. Revisión de Nicolás
![Revisión de Nicolas](04-revision-nicolas-2026-11-22.png)

### 5. Revisión de Romina
![Revisión de Romina](05-revision-romina-2026-11-22.png)

### 6. Checks en verde y merge
![Checks y merge](06-checks-y-merge-2026-11-22.png)

### 7. Respuestas del autor
![Respuestas del autor](07-respuesta-autor-2026-11-23.png)

### 8. Créditos y costo de la cuenta AWS
![Créditos de la cuenta AWS](08-creditos-cuenta-aws-2026-11-23.png)

### 9. Alerta de AWS Budgets (50 USD por mes)
![Alerta de AWS Budgets](09-alerta-budgets-2026-11-23.png)
