# ADR-0011: Diseño de los security groups

- **Estado:** Aceptada
- **Fecha:** 2026-11-22
- **Tarea:** S2-10, VPC completa con subredes públicas, privadas y security groups (issue #26)
- **Área:** Infra (Terraform, `infra/modules/network/security_groups.tf`)

## Contexto

El auditor observó que la base de datos y Redis están publicados en puertos de la máquina anfitriona. En AWS el control equivalente son los **security groups**: definen qué tráfico puede llegar a cada componente.

Flujos reales de la aplicación, tomados de `docker-compose.yml`. Todos los servicios escuchan en el puerto **8080** dentro del contenedor:

| Origen | Destino | Puerto |
|---|---|---|
| Internet | ALB | 80 |
| ALB | `ui`, `admin` | 8080 |
| `ui` | `catalog`, `carts`, `checkout`, `orders` | 8080 |
| `checkout` | `orders` | 8080 |
| `checkout` | Redis | 6379 |
| `catalog`, `carts`, `orders`, `admin` | Postgres | 5432 |

Las tareas de ECS Fargate cambian de IP cada vez que se crean, así que las reglas por IP no sirven para el tráfico interno.

## Decisión

### 1. Cuatro security groups, con reglas que referencian a otros SGs

| SG | Ingress | Origen | Egress |
|---|---|---|---|
| `alb` | 80/tcp | `0.0.0.0/0` | 8080/tcp hacia el SG `ecs` |
| `ecs` | 8080/tcp | SG `alb` y el propio SG `ecs` | Todo (`0.0.0.0/0`) |
| `rds` | 5432/tcp | SG `ecs` | Ninguno |
| `redis` | 6379/tcp | SG `ecs` | Ninguno |

Las reglas internas apuntan a **SGs, no a CIDRs**. Es el mismo patrón del laboratorio *T04 › 02-networking › 03*, donde el SG App solo acepta tráfico desde el SG del Bastion.

### 2. Un SG compartido para los 6 servicios de ECS

No se crea un SG por servicio.

### 3. Reglas como recursos separados

Se usan `aws_vpc_security_group_ingress_rule` y `aws_vpc_security_group_egress_rule`, no bloques `ingress`/`egress` dentro del SG. Los SGs `alb` y `ecs` se referencian mutuamente. Con bloques inline eso es una dependencia circular que Terraform no puede resolver. Con recursos separados, primero se crean los SGs y después las reglas. Además, es lo que recomienda la documentación del provider AWS.

### 4. Puertos como variables del módulo

`app_port`, `db_port` y `redis_port` tienen valores por defecto (8080, 5432 y 6379).

## Consecuencias

**Positivas**

- RDS y Redis solo son alcanzables desde los servicios. No tienen exposición a internet ni a otros recursos de la VPC.
- RDS y Redis no pueden iniciar conexiones salientes. Las respuestas igual llegan, porque los SGs son *stateful*: el tráfico de retorno de una conexión permitida se permite automáticamente.
- Las reglas siguen funcionando aunque las tareas cambien de IP o escalen.

**Negativas y riesgos aceptados**

- **SG compartido en ECS:** cualquier servicio puede llamar a cualquier otro en el 8080, aunque no lo necesite (por ejemplo, `catalog` → `orders`). Se acepta por simplicidad. La alternativa es un SG por servicio con reglas explícitas por flujo.
- **Egress abierto en ECS:** las tareas pueden salir a cualquier destino. Hace falta para ECR, CloudWatch Logs, Secrets Manager, RDS y Redis. Restringirlo por puerto o destino es una mejora posible.
- **Solo HTTP en el ALB:** HTTPS requiere un dominio y un certificado (ACM), fuera del alcance por ahora.
- **El panel `admin` queda expuesto a través del ALB.** El plan del equipo prevé restringirlo (SG por IP o API Gateway con autenticación); está pendiente.

## Cuándo revisar esta decisión

- Si se agregan servicios con requisitos de aislamiento distintos: SG por servicio.
- Cuando se defina el acceso al `admin`: nueva regla o componente delante del ALB.
- Si se consigue un dominio: agregar 443 en el ALB y redirigir el 80.

## Referencias

- `docker-compose.yml` del repositorio Tero (puertos y dependencias entre servicios).
- Teórico T04.1 Cloud Computing, sección Security Groups (stateful).
- Laboratorio del curso T04 › 02-networking › 03-Subnets Públicas y Privadas.
- ADR-0010: Salida a internet de las subredes privadas.

## Uso de IA

El diseño de las reglas y la redacción de este ADR se hicieron con apoyo de Claude (Anthropic). El equipo verificó los flujos contra `docker-compose.yml` y el material del curso.
