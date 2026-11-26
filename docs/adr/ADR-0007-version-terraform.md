# ADR-0007: Versión mínima de Terraform 1.17

- **Estado:** Aceptada
- **Fecha:** 2026-11-09
- **Tarea:** S1-12, bucket S3 del estado y backend remoto cifrado (issue #12, PR #85)
- **Área:** Infra (Terraform, `infra/backend.tf`, `README.md`)

## Contexto

Tres personas y el pipeline van a correr Terraform sobre el mismo estado. Con versiones distintas, un integrante puede usar una función que otro no tiene, o un Terraform más nuevo puede escribir el estado en un formato que uno más viejo no entiende.

El backend elegido (ADR-0005) usa `use_lockfile`, que existe **desde Terraform 1.10**. Ese es el mínimo técnico.

## Opciones consideradas

### A. `required_version = ">=1.10"`

- El mínimo técnico. Acepta versiones muy distintas entre integrantes.

### B. `required_version = ">=1.17"`

- Una versión reciente, la que tenía instalada el equipo al empezar.
- Todos trabajan con un comportamiento parecido.

### C. Versión exacta (`= 1.17.x`)

- Máxima reproducibilidad.
- Obliga a todos a la misma versión de parche; cualquier actualización rompe al resto.

## Decisión

**Opción B.**

```hcl
terraform {
  required_version = ">=1.17"
}
```

El README lista "Terraform v1.17+" como prerrequisito. La versión del **provider** se fija aparte, con `.terraform.lock.hcl` (ADR-0008), uno por root (`infra/` e `infra/global`).

**En el pipeline la versión es exacta:** el job `Terraform (fmt y validate)` fija Terraform `1.17.1` (PR #100), corre `init -backend=false` (no necesita credenciales), `fmt -check` y `validate` en cada PR, y es un check obligatorio del ruleset de `main` (ver [ADR-0003](ADR-0003-github-flow-y-proteccion-de-main.md)).

## Consecuencias

**Positivas**

- Si alguien tiene una versión vieja, Terraform falla al iniciar con un mensaje claro, en lugar de fallar a mitad de camino.
- Cubre el requisito técnico de `use_lockfile`.

**Negativas y riesgos aceptados**

- Es un mínimo, no una versión exacta: dos integrantes pueden tener 1.17 y 1.18.
- Quien tenga una versión anterior tiene que actualizar Terraform.

## Cuándo revisar esta decisión

- Si una versión nueva trae algo que el equipo necesita: subir el mínimo y avisar.
- Si aparecen diferencias de comportamiento entre integrantes: fijar la versión del pipeline.

## Referencias

- Documentación de Terraform: `required_version`, backend S3 (`use_lockfile` desde 1.10).
- PR #85 del repositorio Tero.

## Uso de IA

La redacción de este ADR se hizo con apoyo de Claude (Anthropic). La versión la eligió el equipo.
