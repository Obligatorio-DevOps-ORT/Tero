# ADR-0013: Configuración de los repositorios ECR

- **Estado:** Aceptada, con un punto pendiente
- **Fecha:** 2026-11-12
- **Tarea:** S2-11, repositorios de imágenes en ECR (issue #27, PR #102)
- **Área:** Infra (Terraform, `infra/modules/ecr`)

> **Pendiente de decisión:** la retención de imágenes. Hoy se guardan 5, y como cada merge a `main` publica las 6 imágenes, 5 merges sin promover pueden borrar la imagen que corre en prod. Falta elegir entre subir el número (por ejemplo, a 20) o conservar siempre los tags `v*` (ver "Negativas y riesgos aceptados"). Decide Infra antes del primer deploy a prod. Tags inmutables, escaneo y `force_delete` están confirmados.

## Contexto

Cada repo ECR tiene opciones que definen cómo se guardan y protegen las imágenes. Las que importan para Tero:

- El CI etiqueta cada imagen con el **SHA del commit**, que es único por build.
- ECS **no guarda copia de la imagen**: cada vez que levanta una tarea (deploy, escalado o reinicio) la vuelve a bajar de ECR.
- ECR cobra el **almacenamiento**, unos USD 0.10 por GB al mes. Sin limpieza, cada build suma una imagen y el repo crece para siempre.
- El root de ECR se crea y se destruye en un entorno de práctica, no en producción real.

## Decisión

### 1. Tags inmutables (`image_tag_mutability = "IMMUTABLE"`)

Un tag que ya existe no se puede volver a pushear. Cada tag apunta siempre a la misma imagen, así que lo que dice el tag es exactamente lo que corre, y un rollback a un tag anterior es confiable.

Es compatible con el CI, porque el SHA nunca se repite. A cambio, **no se puede usar `latest`**.

### 2. Escaneo de vulnerabilidades al subir (`scan_on_push = true`)

ECR escanea cada imagen en busca de CVEs conocidos apenas se sube. Es el escaneo básico, sin costo extra.

El escaneo **no bloquea** el push. Decidir si el pipeline frena un despliegue cuando aparecen vulnerabilidades críticas le corresponde a CI/CD. Infra solo deja el escaneo activado, porque es un atributo del repo.

### 3. Lifecycle policy: se guardan las últimas 5 imágenes

Una regla de `aws_ecr_lifecycle_policy` borra las imágenes más viejas cuando hay más de 5 (`imageCountMoreThan`, `tagStatus = any`).

Se evaluó guardar 3. La diferencia de costo es de centavos y 5 da más margen para rollback. Con los Dockerfiles de S2 (ver [ADR-0020](ADR-0020-imagenes-multi-stage-minimas-no-root.md)) las imágenes ocupan entre 6 y 73 MB comprimidas, así que 6 repos × 5 imágenes son menos de 2 GB: unos USD 0,20 al mes.

La retención es una variable del módulo (`image_retention_count`, default 5) con una validación que no permite menos de 1.

### 4. `force_delete` como variable, con default `false`

Terraform no puede borrar un repo que tiene imágenes. Con `force_delete = true`, el `destroy` borra el repo y sus imágenes.

En el módulo vale `false`, que es lo seguro: quien lo use tiene que pedir el comportamiento peligroso explícitamente. El root `infra/global` lo pone en `true` porque es un entorno de práctica que se crea y se destruye.

### 5. Lo que queda fijo en el módulo

Los tags inmutables y el escaneo **no son variables**. Son decisiones de seguridad del equipo, no cambian según el ambiente, y hacerlas configurables invita a apagarlas.

El cifrado no se declara: ECR cifra las imágenes por defecto con AES-256.

## Consecuencias

**Positivas**

- Cada despliegue es trazable a un commit exacto.
- El costo de almacenamiento queda acotado sin intervención manual.
- Las vulnerabilidades quedan registradas desde el primer push.

**Negativas y riesgos aceptados**

- **La lifecycle policy puede borrar la imagen que está corriendo.** Si el CI pusheara una imagen en cada PR, 5 builds de PR alcanzarían para desplazar la imagen de producción. La próxima vez que ECS reinicie una tarea, fallaría con `CannotPullContainerError`. Por eso se acordó con CI/CD que **el push a ECR se haga solo al mergear a main**; en los PR solo se construye (implementado en el PR #103, ver [ADR-0019](ADR-0019-imagenes-etiquetadas-con-sha.md)).
  - **El riesgo sigue con los merges:** cada merge a `main` publica las 6 imágenes, aunque cambie un solo servicio. Después de 5 merges sin promover a prod, ECR puede borrar la imagen que corre en prod, y también se achica la ventana de rollback (S4-10). Como la retención es una variable, hay dos salidas baratas: subir `image_retention_count` en `infra/global` (por ejemplo a 20, que son menos de 6 GB en total), o agregar una regla con `rulePriority` 1 que conserve los tags `v*` (las versiones promovidas) y pasar la regla actual a prioridad 2. **Pendiente de decidir antes del primer deploy a prod.**
- **Sin `latest`**, cualquier herramienta o script que lo use falla. Hay que referenciar siempre el SHA.
- **`force_delete = true` en el root global** permite borrar todas las imágenes con un `destroy`. Se acepta en un entorno de práctica; en una cuenta productiva iría en `false`.

## Cuándo revisar esta decisión

- Si el CI empieza a publicar imágenes de PR: subir la retención o separar esas imágenes con una regla por prefijo de tag.
- Antes del primer deploy a prod: decidir entre subir la retención o conservar los tags `v*` (ver riesgos).
- Si el equipo necesita bloquear despliegues por vulnerabilidades: CI/CD lee el resultado del escaneo, o se evalúa el escaneo mejorado (Amazon Inspector).
- Si Tero pasa a producción: `force_delete = false`.

## Referencias

- Letra del obligatorio, sección 3.2.
- Issue #27 y PR #102 del repositorio Tero.
- `.github/workflows/ci.yml` (tag por `github.sha`).
- ADR-0012: Repositorios ECR compartidos en un root de Terraform separado.

## Uso de IA

El análisis de opciones y la redacción de este ADR se hicieron con apoyo de Claude (Anthropic). El equipo verificó la configuración en AWS con `aws ecr describe-repositories` y `aws ecr get-lifecycle-policy`.
