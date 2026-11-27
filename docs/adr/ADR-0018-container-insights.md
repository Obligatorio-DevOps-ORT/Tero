# ADR-0018: Container Insights activado solo en prod

- **Estado:** Aceptada
- **Fecha:** 2026-11-23
- **Tarea:** S2-12, cluster ECS (issue #28)
- **Área:** Infra (Terraform, `infra/modules/ecs-cluster`, `infra/envs/prod.tfvars`)

## Contexto

La letra (sección 8, Observabilidad) pide, **para producción**:

- Métricas de infraestructura (CPU, memoria, red) o de aplicación de los microservicios.
- Logs centralizados con formato estructurado.
- Al menos 1 dashboard y al menos 2 alertas, cada una con su condición de disparo y su procedimiento de respuesta.

ECS publica en CloudWatch, sin costo, métricas básicas de **CPU y memoria a nivel de service**. Esas métricas no incluyen la red, y el promedio del service no muestra qué task tiene problemas: una task al 100% y dos al 10% dan un promedio de 40%.

**Container Insights** agrega métricas por task y por contenedor (CPU, memoria, red, disco, cantidad de tasks) y dashboards en CloudWatch. Se cobra: genera datos en CloudWatch Logs y métricas.

Restricción: créditos limitados del plan gratuito.

## Opciones consideradas

### A. Apagado en todos los ambientes

- Sin costo, pero no cubre la red ni la visibilidad por task que pide la sección 8.

### B. Activado en todos los ambientes

- Cubre la letra, pero gasta créditos en dev y test, donde casi no se va a mirar.

### C. Activado solo en prod

- Cubre la letra donde la pide ("en producción") y no gasta en dev y test.

## Decisión

**Opción C.** El módulo recibe `enable_container_insights` (`bool`, default `false`) y lo traduce al setting `containerInsights` (`enabled` o `disabled`) del cluster. Solo `prod.tfvars` lo pone en `true`.

## Consecuencias

**Positivas**

- Cubre las métricas de CPU, memoria y red de la sección 8, con visibilidad por task.
- Da un dashboard en CloudWatch que sirve de base para el dashboard que pide la letra.
- En dev y test no hay costo de observabilidad.

**Negativas y riesgos aceptados**

- **En dev y test no hay métricas por task.** Para diagnosticar un problema ahí, se activa temporalmente con `-var="enable_container_insights=true"`.
- **El costo en prod crece con la cantidad de tasks y servicios.** Hay que tenerlo presente cuando prod esté prendido.
- **El cluster no resuelve toda la sección 8.** Los logs centralizados (driver `awslogs` en cada task definition), el dashboard propio y las 2 alertas son tareas aparte de observabilidad.

## Cuándo revisar esta decisión

- Si el costo de Container Insights en prod se vuelve relevante frente a los créditos que quedan.
- Si el equipo pasa a otra solución de métricas (por ejemplo, Prometheus y Grafana).

## Referencias

- Letra del obligatorio, sección 8 (Observabilidad).
- Documentación de AWS: *Amazon ECS Container Insights* y precios de CloudWatch.
- ADR-0016: Cómputo del cluster ECS con Fargate.

## Uso de IA

El análisis de opciones y la redacción de este ADR se hicieron con apoyo de Claude (Anthropic). La decisión la tomó Infra y se verificó contra la letra del obligatorio.
