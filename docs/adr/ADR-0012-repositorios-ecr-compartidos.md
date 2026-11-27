# ADR-0012: Repositorios ECR compartidos en un root de Terraform separado

- **Estado:** Aceptada
- **Fecha:** 2026-11-13
- **Tarea:** S2-11, repositorios de imágenes en ECR (issue #27, PR #102)
- **Área:** Infra (Terraform, `infra/modules/ecr`, `infra/global`)

## Contexto

El pipeline de CI ya construye las imágenes de los 6 servicios con una matrix, pero no tiene dónde publicarlas. ECR quedó confirmado como registry en S1-02.

El root principal de Terraform (`infra/`) maneja los tres ambientes (dev, test, prod) con **un workspace por ambiente**, y los tres viven en **la misma cuenta de AWS**. Cada workspace aplica todo lo que está en `infra/main.tf`.

Los nombres de los repositorios ECR son **únicos por cuenta y región**. Con la VPC eso no pasa: puede haber varias VPCs con el mismo tag `Name`.

El CI etiqueta las imágenes como `tero-<servicio>:<sha del commit>`, con el nombre de la carpeta en `src/`.

## Opciones consideradas

### A. Módulo ECR dentro del root principal

- Los repos se crean en cada workspace, como la VPC.
- **No funciona:** el workspace `dev` crea `tero-cart`, y el `apply` del workspace `test` falla con `RepositoryAlreadyExistsException`.

### B. Un juego de repos por ambiente (`tero-dev-cart`, `tero-test-cart`, ...)

- Entra en la estructura actual sin cambios.
- Son 18 repos, y el CI tiene que construir o copiar la imagen una vez por ambiente.
- Rompe el principio *build once, deploy many*: lo que se probó en test no es exactamente el mismo artefacto que llega a prod.

### C. Repos compartidos en un root separado (`infra/global`)

- 6 repos `tero-<servicio>` que usan los tres ambientes.
- La imagen se construye una sola vez y se promueve de dev a test y a prod sin tocarla.
- Requiere un segundo root de Terraform, con su propio state y su propio `apply`.

## Decisión

**Opción C.**

1. **Módulo reutilizable `infra/modules/ecr`**, que crea un repo por nombre recibido.
2. **Root nuevo `infra/global`**, que llama al módulo con los 6 servicios. Usa el mismo bucket que el root principal, pero **otra key de state** (`global/terraform.tfstate`), y trabaja en el workspace `default`, sin workspaces por ambiente.
3. **Nombres `tero-<servicio>`**, iguales a las imágenes que construye el CI: `tero-catalog`, `tero-cart`, `tero-orders`, `tero-checkout`, `tero-ui` y `tero-admin`.
4. **`for_each` sobre un `set(string)`**, no `count`. Cada repo queda identificado en el state por su nombre (`["tero-cart"]`), no por su posición. Si se saca un servicio de la lista, solo se borra ese repo, y no se recrean los que estaban después.

## Consecuencias

**Positivas**

- *Build once, deploy many:* el mismo artefacto (mismo SHA) pasa por los tres ambientes.
- El root global se aplica una vez y no se vuelve a tocar salvo que cambien los repos. Las imágenes no son recursos de Terraform: el CI las sube y AWS borra las viejas, y el `plan` del root global sigue dando `No changes`.
- Destruir un ambiente para ahorrar créditos no borra las imágenes.

**Negativas y riesgos aceptados**

- **Dos roots para operar.** Regla para no confundirse: en `infra/` se elige workspace y tfvars del ambiente; en `infra/global` no.
- **La key del state es crítica.** Si `infra/global` usara la misma key que el root principal, cada root intentaría pisar el state del otro.
- **La lista de servicios está en dos lugares:** la matrix del CI y `infra/global/main.tf`. Si se agrega un servicio hay que tocar ambos.
- **El CI valida solo `infra/`.** El job de Terraform corre `fmt -recursive`, que sí cubre `infra/global`, pero el `init` y el `validate` se ejecutan solo en `infra/`.
- **`cart` y `carts` no son lo mismo.** El repo y la imagen se llaman `tero-cart` (por la carpeta), pero el servicio se encuentra en la red como `carts` (`docker-compose.yml` y la UI llaman a `http://carts:8080`). En ECS, el nombre de descubrimiento del servicio tiene que ser `carts`.

## Cuándo revisar esta decisión

- Si cada ambiente pasa a una cuenta de AWS distinta: los repos irían en una cuenta central, con permisos para que las otras cuentas descarguen las imágenes.
- Si se agrega o se quita un servicio: actualizar la lista del root global y la matrix del CI juntas.

## Referencias

- Letra del obligatorio, sección 3.2.
- Issue #27 y PR #102 del repositorio Tero.
- `.github/workflows/ci.yml` (matrix de build y formato del tag).
- `docker-compose.yml` (nombres de red de los servicios).

## Uso de IA

El análisis de opciones y la redacción de este ADR se hicieron con apoyo de Claude (Anthropic). El equipo verificó la decisión con `terraform plan` (12 recursos a crear) y contra la cuenta de AWS después del `apply`.
