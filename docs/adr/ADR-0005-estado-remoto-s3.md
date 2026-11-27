# ADR-0005: Estado remoto de Terraform en S3 con cifrado y lock nativo

- **Estado:** Aceptada (el bucket y la región fueron reemplazados por ADR-0014; el mecanismo sigue igual)
- **Fecha:** 2026-11-09
- **Tarea:** S1-12, bucket S3 del estado y backend remoto cifrado (issue #12, PR #85)
- **Área:** Infra (Terraform, `infra/backend.tf`)

## Contexto

Terraform guarda en el **estado** (`tfstate`) qué recursos creó y con qué ID. Si cada integrante lo tiene en su máquina, dos personas pueden crear el mismo recurso dos veces o pisarse los cambios.

La letra (sección 7) exige un **backend remoto** y no guardar valores sensibles en el estado sin cifrar. También permite que el bucket del estado sea la única pieza creada a mano: Terraform no puede guardar su estado en algo que todavía no creó.

## Opciones consideradas

### A. Estado local

- Cero configuración.
- Un estado por persona: no hay fuente de verdad compartida ni bloqueo. No cumple la letra.

### B. S3 con lock en DynamoDB

- El patrón clásico.
- Requiere crear y mantener una tabla de DynamoDB solo para el bloqueo.

### C. S3 con lock nativo (`use_lockfile`)

- Desde Terraform 1.10, S3 bloquea el estado con un archivo `.tflock` en el mismo bucket.
- Sin DynamoDB: una pieza menos.

### D. Terraform Cloud

- Estado y bloqueo administrados.
- Otra cuenta y otra herramienta fuera de lo visto en el curso.

## Decisión

**Opción C.** Un bucket S3 creado a mano, con versionado, y este backend:

```hcl
backend "s3" {
  bucket       = "obligatorio-devops-tfstate-26"
  key          = "terraform.tfstate"
  region       = "us-east-1"
  encrypt      = true
  use_lockfile = true
}
```

- `encrypt = true`: el estado se guarda cifrado.
- `use_lockfile = true`: mientras alguien corre `apply`, nadie más puede escribir el estado.
- El bloque `backend` no acepta variables, así que los valores van fijos.

## Consecuencias

**Positivas**

- Una sola fuente de verdad para el equipo y, más adelante, para el pipeline.
- Cumple los requisitos de estado remoto y cifrado de la letra.
- El versionado del bucket permite recuperar un estado anterior si se corrompe.

**Negativas y riesgos aceptados**

- El bucket queda fuera de Terraform: si se borra, hay que recrearlo a mano.
- El bucket vive en una sola cuenta. Quien no tenga credenciales de esa cuenta no puede correr `plan`; puede validar con `terraform init -backend=false` y `terraform validate`.
- Cambiar de bucket o de cuenta requiere `terraform init -reconfigure` (lo que pasó en ADR-0014).

## Cuándo revisar esta decisión

- Si el equipo pasa a otra cuenta o región (ya ocurrió: ADR-0014).
- Si se adopta Terraform Cloud o otro backend administrado.

## Referencias

- Letra del obligatorio, sección 7 (IaC).
- Documentación de Terraform: backend S3, `use_lockfile`.
- PR #85 del repositorio Tero.

## Uso de IA

El análisis de opciones y la redacción de este ADR se hicieron con apoyo de Claude (Anthropic). La decisión la tomó Infra y se verificó con `terraform init` y `terraform plan` contra el bucket.
