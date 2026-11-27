# ADR-0006: Un estado por ambiente con workspaces y guard de ambiente

- **Estado:** Aceptada
- **Fecha:** 2026-11-09
- **Tarea:** S1-12, bucket S3 del estado y backend remoto cifrado (issue #12, PR #85)
- **Área:** Infra (Terraform, `infra/backend.tf`, `infra/main.tf`)

## Contexto

Tero tiene tres ambientes: dev, test y prod. Cada uno necesita su propio estado, para que un `apply` en dev nunca toque prod.

El equipo vio **workspaces** en el curso. Con workspaces, el mismo código y el mismo backend guardan un estado distinto por ambiente.

El riesgo de los workspaces es operativo: estar parado en `prod` creyendo que estás en `dev`.

## Opciones consideradas

### A. Workspaces del mismo root

- Un solo código. El estado de cada ambiente queda en `env:/<workspace>/terraform.tfstate` del mismo bucket.
- Es lo visto en el curso.
- Depende de que la persona seleccione el workspace correcto.

### B. Una carpeta por ambiente (`envs/dev`, `envs/test`, `envs/prod`), cada una con su backend

- Separación explícita en el sistema de archivos.
- Más archivos repetidos y otro patrón que el del curso.

## Decisión

**Opción A**, con un **guard** que frena el error de ambiente:

```hcl
resource "terraform_data" "env_guard" {
  lifecycle {
    precondition {
      condition     = var.env == terraform.workspace
      error_message = "El workspace '${terraform.workspace}' no coincide con el env '${var.env}'"
    }
  }
}
```

- La variable `env` no tiene `default`: siempre hay que indicarla (desde S1-13, viene del `.tfvars`; ver ADR-0009).
- Si `env` y el workspace no coinciden, el `plan` falla antes de tocar AWS.
- `terraform_data` no crea nada en AWS; solo sirve para colgar la validación.

Cada workspace se crea una vez con `terraform workspace new <env>` y se usa con `terraform workspace select <env>`.

## Consecuencias

**Positivas**

- Un solo código para los tres ambientes.
- Estados separados: un error en dev no afecta a prod.
- El guard convierte el error más común de los workspaces en un error de `plan`.

**Negativas y riesgos aceptados**

- El workspace seleccionado vive en `.terraform/environment`, que es local: cada persona y el pipeline tienen que seleccionarlo.
- Todo lo que está en el root se aplica en cada ambiente. Lo que tiene que existir una sola vez (por ejemplo, ECR) necesita otro root (ver ADR-0012).

## Cuándo revisar esta decisión

- Si los ambientes pasan a cuentas AWS distintas: conviene un backend por cuenta.
- Si los ambientes empiezan a diferir en estructura y no solo en valores.

## Referencias

- Material del curso: Terraform workspaces.
- Documentación de Terraform: `terraform_data`, `precondition`.
- PR #85 del repositorio Tero.

## Uso de IA

El análisis de opciones y la redacción de este ADR se hicieron con apoyo de Claude (Anthropic). La decisión la tomó el equipo y se verificó creando los workspaces y probando el guard con un `env` que no coincide.
