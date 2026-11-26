# ADR-0015: Credenciales de AWS para GitHub Actions

- **Estado:** Aceptada (permisos de administrador confirmados por Infra el 2026-11-25)
- **Fecha:** 2026-11-18
- **Tarea:** S2-11, repositorios de imágenes en ECR (issue #27, PR #102)
- **Área:** Infra (IAM de la cuenta AWS, secrets y variables de GitHub Actions)

> **Historial de los permisos de `tero-ci`:**
>
> 1. **18/11, PR #102:** `tero-ci` se creó con `AdministratorAccess`. En la revisión del PR quedó un comentario de Infra diciendo que ese permiso era **provisorio** ("para probar") y que se iba a cambiar por permisos equivalentes a `LabRole`.
> 2. **25/11:** Infra lo revisó con el equipo y confirmó que **`tero-ci` queda como administrador**. Ese comentario del PR quedó desactualizado: la decisión vigente es la de este ADR.
>
> **Por qué se mantuvo:** Terraform necesita crear roles de IAM, y quien puede crear roles en la práctica ya tiene permisos de administrador. Además, `LabRole` también es muy amplio, así que igualarlo casi no reducía el riesgo.

## Contexto

El pipeline necesita credenciales de AWS para:

- **Ahora:** subir las imágenes a ECR.
- **Más adelante:** correr `terraform plan` en cada PR (S4-18, issue #70) y `terraform apply` al mergear a main (S3-19, issue #95).

Para Terraform, la identidad necesita permisos muy amplios: crear y borrar VPC, ECS, RDS, ElastiCache y **roles de IAM**. Quien puede crear roles puede darse cualquier permiso, así que en la práctica es un administrador.

La infraestructura está en una cuenta AWS propia (ADR-0014), pero el equipo quiere poder volver al Learner Lab si hace falta. En el Learner Lab **no se pueden crear usuarios ni proveedores de identidad**: hay una sola credencial temporal (rol `voclabs`) con access key, secret y session token.

## Opciones consideradas

### A. OIDC entre GitHub y AWS

- GitHub se identifica ante AWS y recibe credenciales temporales que vencen al terminar el job. En GitHub no se guarda ninguna clave, solo el ARN del rol.
- Se puede restringir por repo y por rama.
- **No existe en el Learner Lab.** Volver al Lab implica cambiar la autenticación del pipeline.

### B. Dos usuarios IAM con access keys

- `tero-ci-ecr`, con permiso solo para pushear a los repos `tero-*`.
- `tero-ci-terraform`, con permisos amplios, solo para los jobs de Terraform.
- Si se filtran las claves de push, el daño se limita a subir imágenes.
- No tiene equivalente en el Learner Lab, donde hay una sola credencial.

### C. Un usuario IAM con access keys

- Un usuario `tero-ci` para todo el pipeline.
- Mismo formato que la credencial del Learner Lab: volver al Lab es reemplazar las claves y agregar el session token.

## Decisión

**Opción C.**

1. **Usuario IAM `tero-ci`**, creado a mano desde la consola, **sin acceso a la consola** (solo claves).
2. **Policy `AdministratorAccess`.** Achicarla no protegería de verdad, porque Terraform tiene que crear roles de IAM.
3. **Access key creada a mano**, nunca con Terraform. Si se crea con Terraform, el secreto queda guardado en el state.
4. **En GitHub (Settings → Secrets and variables → Actions):**

| Tipo | Nombre | Contenido |
|---|---|---|
| Secret | `AWS_ACCESS_KEY_ID` | Access key de `tero-ci` |
| Secret | `AWS_SECRET_ACCESS_KEY` | Secret key de `tero-ci` |
| Variable | `AWS_REGION` | `us-east-2` |
| Variable | `ECR_REGISTRY` | URL del registry de la cuenta |

Lo secreto va en **Secrets**, que GitHub oculta en los logs. Lo que no es secreto va en **Variables**, para que el equipo lo vea. La imagen completa se arma como `${ECR_REGISTRY}/tero-<servicio>:<sha>`.

5. **Rotación sin cortar el CI:** AWS permite 2 claves activas por usuario. Se crea la nueva, se actualizan los dos secrets, se verifica que el pipeline pase y se desactiva y borra la vieja. Se ensaya en S4-17, porque "rotar una credencial" es una tarea típica de la parte "Hacer" de la defensa.

## Consecuencias

**Positivas**

- El pipeline funciona sin renovar credenciales.
- Volver al Learner Lab es reemplazar dos secrets y agregar `AWS_SESSION_TOKEN`.
- El registry no está escrito en el workflow: la variable evita fijar el número de cuenta en el código.

**Negativas y riesgos aceptados**

- **Claves de administrador que no vencen.** Si se filtran, dan control total de la cuenta hasta que se revoquen. En el plan gratuito el riesgo no es una factura, sino que se consuman los créditos y la cuenta se cierre antes de la defensa. Mitigaciones:
  - Las claves solo existen en los secrets de GitHub. No se copian en archivos ni en chats.
  - El pipeline corre Gitleaks para detectar secretos en el código.
  - La cuenta tiene una alerta de presupuesto configurada.
- **Todos los jobs que usen las claves tienen permisos de administrador**, aunque solo pusheen imágenes.
- **Las credenciales son más poderosas que las del Learner Lab.** El formato es el mismo, pero `voclabs` no puede crear roles de IAM, usa solo 2 regiones y una lista acotada de servicios. Algo que funcione con `tero-ci` puede fallar en el Lab (ver ADR-0014).

## Cuándo revisar esta decisión

- Si se confirma que no se vuelve al Learner Lab: migrar a OIDC.
- Cuando se implementen S4-18 y S3-19: proteger el `apply` con un *environment* de GitHub que pida aprobación manual.
- **Después de la entrega y la defensa: eliminar las claves y el usuario `tero-ci`.**

## Pendientes

- Ensayar la rotación de las claves (S4-17).

## Referencias

- Issues #70 (S4-18) y #95 (S3-19) del repositorio Tero.
- AWS Academy Learner Lab, documentación de *Foundational Services*.
- ADR-0014: Infraestructura en una cuenta AWS propia en lugar del Learner Lab.

## Uso de IA

El análisis de opciones y la redacción de este ADR se hicieron con apoyo de Claude (Anthropic). La decisión la tomó Infra, priorizando la compatibilidad con el Learner Lab sobre el mínimo privilegio.
