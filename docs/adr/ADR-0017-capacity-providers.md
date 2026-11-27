# ADR-0017: Capacity providers FARGATE y FARGATE_SPOT por ambiente

- **Estado:** Aceptada
- **Fecha:** 2026-11-23
- **Tarea:** S2-12, cluster ECS (issue #28)
- **Área:** Infra (Terraform, `infra/modules/ecs-cluster`, `infra/envs/*.tfvars`)

## Contexto

El cluster usa Fargate (ver ADR-0016). Fargate ofrece dos capacity providers:

| Provider | Precio | Disponibilidad |
|---|---|---|
| `FARGATE` | Completo (on-demand) | AWS no reclama la capacidad |
| `FARGATE_SPOT` | Con un descuento grande (hasta ~70%) | AWS puede reclamar la task con 2 minutos de aviso |

Cuando AWS reclama una task Spot, ECS le manda `SIGTERM`, la saca del ALB y le da hasta `stopTimeout` (máximo 120 s en Fargate) para terminar lo que está procesando. Después la corta con `SIGKILL`.

**AWS puede reclamar todas las tasks Spot al mismo tiempo.** Si un servicio corre solo en Spot, puede quedar sin ninguna task. Si además no hay capacidad Spot libre, ECS no tiene dónde volver a lanzarlas.

La estrategia por defecto del cluster reparte las tasks entre los providers con dos valores:

- `base`: cuántas tasks van sí o sí a un provider antes de repartir el resto. Solo un provider puede tener `base`.
- `weight`: en qué proporción se reparten las tasks que quedan.

Restricciones: créditos limitados del plan gratuito, y los servicios críticos de prod (por ejemplo, checkout) no pueden quedar caídos.

## Opciones consideradas

### A. Solo FARGATE en todos los ambientes

- Máxima disponibilidad y precio completo en todos lados, incluso en dev y test, donde una interrupción no afecta a nadie.

### B. Solo FARGATE_SPOT en todos los ambientes

- Mínimo costo, pero prod puede quedar sin ninguna task si AWS reclama la capacidad.

### C. Estrategia según el ambiente

- dev y test: todo en `FARGATE_SPOT`.
- prod: una task fija en `FARGATE` y el resto en `FARGATE_SPOT`.

## Decisión

**Opción C.** El módulo habilita los dos providers y define la estrategia por defecto así:

```hcl
FARGATE      → base = var.fargate_base, weight = 0
FARGATE_SPOT → weight = 1
```

Con **una sola variable** (`fargate_base`) se cubren los dos casos:

| Ambiente | `fargate_base` | Resultado |
|---|---|---|
| dev, test | `0` (default) | Todas las tasks en Spot |
| prod | `1` | 1 task en on-demand, el resto en Spot |

Es la estrategia **por defecto** del cluster. Un service puede definir su propia estrategia si lo necesita.

## Consecuencias

**Positivas**

- Dev y test cuestan lo mínimo posible.
- En prod siempre queda al menos una task on-demand que sigue atendiendo si AWS reclama el Spot.
- La variable es simple y tiene un default seguro para el costo. Sigue el mismo patrón que `enable_nat_gateway`.

**Negativas y riesgos aceptados**

- **En dev y test los servicios se pueden interrumpir.** Se acepta: son ambientes de trabajo.
- **En prod, una sola task on-demand puede no alcanzar** si AWS reclama todo el Spot en un pico de carga. Funciona, pero degradado.
- **La app tiene que manejar `SIGTERM`** para terminar las requests en curso (un pago, por ejemplo) antes del corte. Eso es responsabilidad del código de cada servicio, no del cluster.
- `fargate_base` vale para todo el cluster. Si un servicio necesita otra distribución, la define en su propio service.

## Cuándo revisar esta decisión

- Si se ven interrupciones de Spot frecuentes que afectan las demos o las pruebas.
- Si Tero pasa a producción real: subir `fargate_base`, o usar estrategias por service para los servicios críticos.

## Referencias

- Documentación de AWS: *Amazon ECS capacity providers for Fargate* y *Fargate Spot termination notices*.
- Laboratorio del curso T06 › Containers › ECS.
- ADR-0016: Cómputo del cluster ECS con Fargate.
- ADR-0010: Salida a internet de las subredes privadas (patrón de toggle por variable).

## Uso de IA

El análisis de opciones y la redacción de este ADR se hicieron con apoyo de Claude (Anthropic). La decisión la tomó Infra.
