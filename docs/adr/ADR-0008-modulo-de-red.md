# ADR-0008: Módulo de red reutilizable con subredes públicas y privadas en 2 zonas

- **Estado:** Aceptada (completado por ADR-0010 y ADR-0011)
- **Fecha:** 2026-11-12
- **Tarea:** S1-13, primer módulo de red en Terraform: VPC y subredes (issue #13, PR #92)
- **Área:** Infra (Terraform, `infra/modules/network`, `infra/providers.tf`)

## Contexto

La letra (sección 7.2) exige **al menos un módulo reutilizable**, y da como ejemplo un módulo de red. Además, el auditor observó que la base de datos y Redis están publicados en puertos de la máquina anfitriona: en AWS tienen que quedar en subredes privadas.

El tablero divide la red en dos tareas: S1-13 (VPC y subredes) y S2-10 (red completa con ruteo y security groups).

Dos servicios que vamos a usar exigen subredes en **al menos 2 zonas de disponibilidad**:

- El **Application Load Balancer**.
- El **DB subnet group** de RDS, aunque la base corra en una sola instancia.

## Opciones consideradas

### A. Recursos de red directo en el root

- Menos archivos.
- No cumple el requisito de módulo y no se reutiliza entre ambientes.

### B. Módulo comunitario `terraform-aws-modules/vpc`

- Completo y probado.
- Esconde la lógica que el equipo tiene que poder explicar en la defensa.

### C. Módulo propio `infra/modules/network`

- Solo lo que Tero necesita, escrito y entendido por el equipo.

Para las zonas:

- **Una sola zona:** más simple, pero obliga a rehacer la red al crear el ALB o la RDS.
- **Dos zonas:** cumple lo que exigen ALB y RDS, y tolera la caída de una zona.

## Decisión

**Opción C, con 2 zonas de disponibilidad.**

- El módulo recibe `name`, `vpc_cidr`, `public_subnet_cidrs` y `private_subnet_cidrs`, todas **sin `default`**: los valores los decide quien lo llama.
- Crea una VPC (con `enable_dns_support` y `enable_dns_hostnames`) y una subred pública y una privada por zona, con `count = length(lista)`.
- Las zonas salen de `data "aws_availability_zones"`, no se escriben a mano.
- Devuelve `vpc_id`, `vpc_cidr`, `public_subnet_ids` y `private_subnet_ids`, todos con `description`.
- El provider AWS se fija en `~> 6.0` y su versión exacta queda en `.terraform.lock.hcl`, que se commitea. Los recursos llevan `default_tags` con el proyecto y el ambiente.
- El Internet Gateway, el ruteo, la salida de las privadas y los security groups quedan **fuera de alcance**: son S2-10.

## Consecuencias

**Positivas**

- Cumple el requisito de módulo reutilizable: el mismo módulo crea la red de dev, test y prod.
- El ALB y RDS ya tienen las subredes que exigen.
- Agregar una tercera zona es agregar un CIDR a cada lista, sin tocar el módulo.

**Negativas y riesgos aceptados**

- Hasta S2-10, "pública" y "privada" son solo el diseño: sin ruteo, ninguna tiene salida a internet.
- Si una lista tiene más CIDR que zonas disponibles en la región, el `plan` falla.

## Cuándo revisar esta decisión

- Si se necesita tolerancia a la caída de dos zonas: pasar a 3.
- Si el módulo crece mucho: separar security groups u otros recursos en módulos propios.

## Referencias

- Letra del obligatorio, sección 7.2 (modularización, outputs documentados).
- Observación del auditor sobre puertos publicados de la base y Redis.
- Documentación de AWS: requisitos de subredes del ALB y de los DB subnet groups de RDS.
- PR #92 del repositorio Tero.

## Uso de IA

El análisis de opciones y la redacción de este ADR se hicieron con apoyo de Claude (Anthropic). Infra escribió el código y lo verificó con `terraform apply` en dev (`5 to add`: 1 VPC y 4 subredes) y en la consola de AWS.
