# ADR-0009: Valores por ambiente con `.tfvars`, combinados con workspaces

- **Estado:** Aceptada
- **Fecha:** 2026-11-12
- **Tarea:** S1-13, primer módulo de red en Terraform: VPC y subredes (issue #13, PR #92). Cubre también S2-13 (issue #29).
- **Área:** Infra (Terraform, `infra/envs/*.tfvars`, `infra/variables.tf`)

## Contexto

En S1-12 el equipo eligió workspaces para separar el estado de cada ambiente (ADR-0006), sin `.tfvars`: `plan` pedía el valor de `env` por consola.

Al leer la letra, la sección 7.2 exige: *"usar archivos `.tfvars` diferenciados para Dev, Test y Prod"*. Con el módulo de red aparecen valores que cambian por ambiente (los rangos de IP), y escribirlos por consola no es reproducible: el pipeline no puede responder preguntas.

## Opciones consideradas

### A. Seguir sin `.tfvars` y pasar valores por consola o con `-var`

- No cumple la letra.
- Las listas hay que tipearlas en sintaxis HCL: un error de tipeo crea otra red, sin dejar registro.

### B. `terraform.tfvars` o `*.auto.tfvars`

- Terraform los carga solo.
- Es uno por carpeta: no sirve para tres ambientes en el mismo root.

### C. Un `.tfvars` por ambiente en `infra/envs/`, pasado con `-var-file`

- Valores versionados y revisados en un PR.
- Funciona igual por consola y desde el pipeline.

## Decisión

**Opción C, combinada con workspaces.** El workspace decide **dónde** se guarda el estado; el `.tfvars` decide **con qué valores** se crea la infraestructura.

```
terraform workspace select dev
terraform plan -var-file="envs/dev.tfvars"
```

- `env` pasa a estar dentro de cada `.tfvars`. Si alguien selecciona el workspace `prod` y pasa `dev.tfvars`, el guard de ADR-0006 frena el `plan`.
- Un rango de IP distinto por ambiente. El segundo número indica el ambiente, y las privadas empiezan en el tercer número 11:

| Ambiente | VPC | Públicas | Privadas |
|---|---|---|---|
| dev | `10.0.0.0/16` | `10.0.1.0/24`, `10.0.2.0/24` | `10.0.11.0/24`, `10.0.12.0/24` |
| test | `10.1.0.0/16` | `10.1.1.0/24`, `10.1.2.0/24` | `10.1.11.0/24`, `10.1.12.0/24` |
| prod | `10.2.0.0/16` | `10.2.1.0/24`, `10.2.2.0/24` | `10.2.11.0/24`, `10.2.12.0/24` |

- El nombre de los recursos (`tero-${var.env}`) se arma en el root y no en el `.tfvars`, para que siempre coincida con el ambiente.
- Los `.tfvars` **se commitean**: no tienen secretos. Las contraseñas van a Secrets Manager, nunca a un `.tfvars`.

## Consecuencias

**Positivas**

- Cumple el requisito de la letra y vale en la rúbrica ("IaC modular y parametrizada").
- Despliegues reproducibles: el mismo archivo da siempre el mismo resultado.
- El guard pasa a proteger contra la combinación equivocada de workspace y `.tfvars`.
- Los rangos no se superponen si algún día se conectan las VPCs.

**Negativas y riesgos aceptados**

- Cada comando lleva `-var-file`; olvidarlo hace que Terraform pida los valores por consola.
- Hay que mantener los tres `.tfvars` alineados cuando se agrega una variable nueva.

## Cuándo revisar esta decisión

- Si aparece un valor sensible: va a Secrets Manager o a una variable `sensitive`, nunca al `.tfvars`.
- Si los ambientes pasan a cuentas distintas.

## Referencias

- Letra del obligatorio, sección 7.2 (variables parametrizadas por ambiente).
- ADR-0006: Un estado por ambiente con workspaces y guard de ambiente.
- PR #92 del repositorio Tero.

## Uso de IA

El análisis de opciones y la redacción de este ADR se hicieron con apoyo de Claude (Anthropic). El cambio de criterio salió de releer la letra; el equipo lo aceptó y se verificó con `plan` y `apply` en dev.
