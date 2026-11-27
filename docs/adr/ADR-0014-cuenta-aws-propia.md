# ADR-0014: Infraestructura en una cuenta AWS propia en lugar del Learner Lab

- **Estado:** Aceptada, con un punto pendiente
- **Fecha:** 2026-11-17
- **Tarea:** S2-11, repositorios de imágenes en ECR (issue #27, PR #102)
- **Área:** Infra (Terraform, `infra/backend.tf`, `infra/providers.tf`, `infra/global`)

> **Pendiente de decisión:** cómo acceden los otros dos integrantes a la consola de AWS. Hoy solo entra el dueño de la cuenta. Las opciones son una invitación a cada uno, con su usuario y MFA (recomendado), o compartir el correo y la contraseña de la cuenta. Decide el equipo antes de la defensa. El uso de la cuenta propia, la región y el backend están confirmados (ver "Pendientes" al final).

## Contexto

Hasta S2-11 la infraestructura se aplicaba en el **AWS Academy Learner Lab** (`us-east-1`), con el state en el bucket `obligatorio-devops-tfstate-26`.

Las credenciales del Learner Lab son **temporales**: pertenecen al rol `voclabs` y vencen cuando termina la sesión del laboratorio. Para que el pipeline use AWS (push a ECR y, más adelante, Terraform), habría que reemplazar los secrets de GitHub en cada sesión nueva. Si no, el CI falla en cada push o merge a main.

Además, el Learner Lab tiene restricciones propias:

- Solo las regiones `us-east-1` y `us-west-2`.
- Solo una lista de servicios permitidos, que puede cambiar.
- No se pueden crear roles de IAM. Los servicios usan el rol `LabRole`, que ya viene creado.

Ningún otro integrante del equipo estaba usando recursos de AWS.

## Opciones consideradas

### A. Seguir en el Learner Lab

- Sin costo propio y con el entorno que propone el curso.
- Hay que actualizar los 3 secrets (`AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_SESSION_TOKEN`) cada vez que se abre una sesión.

### B. Una cuenta AWS propia en el plan gratuito

- Credenciales que no vencen. El CI funciona sin intervención.
- Permite crear roles de IAM y elegir la región.
- El presupuesto pasa a ser el de la cuenta: créditos y vencimiento del plan gratuito.

## Decisión

**Opción B.** La infraestructura pasa a una cuenta AWS propia en el plan gratuito, región **`us-east-2`**.

- **Backend:** el state del root principal y del root global va al bucket `tero-tfstate-550842784335` de la cuenta nueva. Cada root mantiene su key.
- **Región como variable:** `aws_region` (default `us-east-2`) en lugar de un valor fijo en el provider.
- **Sin migración de state:** el state viejo describía recursos de otra cuenta. Se inicializó el backend nuevo con `terraform init -reconfigure` y la VPC de dev se volvió a crear en la cuenta nueva.
- El cambio se hizo en un **commit separado** dentro del PR de S2-11, para poder revertirlo solo.

**Presupuesto y créditos (al 18/11):**

- **Sin tarjeta vinculada.** El plan gratuito no factura: cuando se terminan los créditos o vence el período, AWS corta el acceso a los servicios.
- **USD 120 de crédito** (los USD 100 iniciales más USD 20 por completar una actividad de AWS), con vencimiento el **12/05/2027**. Hay hasta USD 80 más disponibles completando las otras 5 actividades (Lambda, Bedrock, RDS, DynamoDB, EC2).
- **AWS Budgets:** presupuesto mensual de USD 50 ("My Monthly Cost Budget", desde el 01/11/2026), con alerta por mail al llegar al monto. Excluye créditos y reembolsos, así que mide el consumo real aunque lo paguen los créditos.
- **Regla de uso:** los ambientes se prenden para trabajar o demostrar y se apagan después: NAT con `enable_nat_gateway` (ADR-0010) y services en `desired_count = 0` (ADR-0016).

## Consecuencias

**Positivas**

- El CI no depende de renovar credenciales.
- Se pueden crear los roles de IAM que necesita ECS, en lugar de usar `LabRole` para todo (sin mínimo privilegio).
- Se pueden usar mecanismos que el Learner Lab no permite, como OIDC para GitHub Actions.

**Negativas y riesgos aceptados**

- **El presupuesto es finito y tiene fecha.** El vencimiento (12/05/2027) queda lejos, pero los créditos se pueden agotar antes: si pasa antes de la entrega (07/12) o de la defensa (08 al 18/12), Tero deja de funcionar y no se puede mostrar. Funciona como una tarjeta prepaga: no genera deuda, pero sin saldo no hay servicio. El NAT Gateway es el mayor consumo (ver ADR-0010).
- **La cuenta depende de una sola persona.** Los demás integrantes no tienen acceso a la consola. Si el dueño no está disponible en la defensa, nadie más puede operar desde la consola.
- **Región `us-east-2`.** `cart` tiene `us-east-1` fijo en `dynamodb_service.py`. No afecta mientras se use Postgres, pero hay que tenerlo presente.
- **Volver al Learner Lab no es solo cambiar el bucket.** Habría que:
  - Cambiar la región a `us-east-1` o `us-west-2`.
  - Crear todo de nuevo allá, porque el state describe recursos de esta cuenta, y destruir lo de esta cuenta para no dejar recursos huérfanos.
  - Reemplazar la creación de roles de IAM por el ARN de `LabRole`. Para que eso no implique reescribir código, los módulos de ECS reciben el rol como variable.
- **Quedan recursos huérfanos en el Learner Lab** (la VPC de S1-13). No tienen costo.

## Cuándo revisar esta decisión

- Si se acaban los créditos o el plan vence antes de la defensa.
- Si la letra o los docentes exigen usar el AWS Academy.
- Si otro integrante necesita operar la infraestructura: darle acceso a la cuenta.

## Pendientes

- **Acceso del equipo a la consola**, antes de la defensa: cada integrante por invitación, con su propio usuario y MFA, en lugar de compartir el correo y la contraseña de la cuenta. Con la invitación, las acciones quedan registradas a nombre de quien las hizo y no se comparte la credencial raíz, que puede cerrar la cuenta o cambiar cualquier configuración.
- **Verificar el segundo mail del presupuesto:** la consola muestra "1/2 emails verified", así que uno de los destinatarios no recibe las alertas.
- Evaluar alertas más tempranas: hoy solo avisa al llegar a USD 50. Con avisos al 50 % y al 80 % (o una alerta por pronóstico), hay margen para reaccionar.
- Revisar los créditos restantes en cada reunión semanal, sobre todo cuando se levanten RDS, ElastiCache y el NAT.
- Documentar en el README el `terraform init -reconfigure` que hace falta cuando cambia el bucket del backend.
- Cerrar S1-08 (prueba en Learner Lab).

## Referencias

- AWS Academy Learner Lab, documentación de *Foundational Services* (regiones y servicios permitidos).
- Issue #27 y PR #102 del repositorio Tero.
- Prueba en Learner Lab de S1-08 (rol `LabRole` como rol de ejecución de ECS).
- ADR-0010: Salida a internet de las subredes privadas.

## Uso de IA

El análisis de opciones y la redacción de este ADR se hicieron con apoyo de Claude (Anthropic). La decisión la tomó Infra, y se verificó con `terraform plan` en la cuenta nueva (`No changes` en el workspace `dev`).
