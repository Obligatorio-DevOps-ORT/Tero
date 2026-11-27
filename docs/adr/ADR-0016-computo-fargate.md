# ADR-0016: Cómputo del cluster ECS con Fargate en lugar de EC2

- **Estado:** Aceptada
- **Fecha:** 2026-11-23
- **Tarea:** S2-12, cluster ECS (issue #28)
- **Área:** Infra (Terraform, `infra/modules/ecs-cluster`)

## Contexto

Tero tiene seis microservicios que hay que desplegar en un orquestador en la nube. La letra (sección 3.3) pide desplegar en *"AWS ECS, EKS u equivalente"*, pero **no define el tipo de cómputo**. El plan del equipo ya eligió ECS.

Un cluster ECS es solo una agrupación lógica y no tiene costo propio. Lo que se cobra es el cómputo donde corren las tasks, y ese cómputo puede ser de dos tipos:

- **EC2**: instancias propias que se registran en el cluster. Las tasks corren dentro de ellas.
- **Fargate**: AWS asigna cómputo aislado a cada task, sin instancias que administrar.

Restricciones:

- Cuenta AWS personal en el plan gratuito, con **créditos limitados**, región `us-east-2` (ver ADR-0014).
- Tres ambientes (dev, test, prod) que se prenden para trabajar o demostrar y se apagan después.
- Equipo chico, sin tiempo para operar servidores.

## Opciones consideradas

Costos aproximados de `us-east-2` al momento de decidir. Hay que verificarlos en la página de precios de AWS.

### A. ECS sobre EC2

- Hay que administrar la AMI optimizada para ECS, los parches, el Auto Scaling Group y la capacidad de las instancias.
- **Costo:** se paga la instancia **las 24 h, haya o no tasks corriendo**. Una `t3.micro` cuesta ~USD 0.0104/h (~USD 7.6/mes) por instancia y por ambiente.
- Bajar las tasks a cero no reduce el gasto. Para dejar de pagar hay que bajar también el Auto Scaling Group.

### B. ECS sobre Fargate

- Sin instancias: AWS administra el cómputo y el sistema operativo.
- **Costo:** solo mientras la task está en estado RUNNING. Se cobra por vCPU/hora y GB/hora: una task de 0.25 vCPU y 0.5 GB cuesta ~USD 0.012/h (~USD 9/mes si corre 24/7).
- Fargate **no tiene capa gratuita**: el uso se descuenta de los créditos del plan.
- Es lo que usan el teórico *T06.1 Contenedores en AWS* y el laboratorio *T06 › ECS*.

## Decisión

**Opción B.** El cluster ECS usa **Fargate** en los tres ambientes.

Para dejar de consumir créditos, se bajan los services a `desired_count = 0`. El cluster no se destruye, porque no tiene costo.

## Consecuencias

**Positivas**

- El gasto es proporcional al uso real. Con las tasks en cero, el consumo de cómputo es nulo, lo que encaja con prender y apagar los ambientes.
- No hay servidores que parchear ni capacidad de instancias que dimensionar.
- Sigue lo que enseña el curso, así que el equipo y quien opere después lo pueden entender y mantener.

**Negativas y riesgos aceptados**

- **Una task 24/7 sale más cara en Fargate** que la misma carga bien empaquetada en EC2. Se acepta porque los ambientes no van a estar prendidos todo el tiempo.
- **Menos control:** no hay acceso al host, y las combinaciones de CPU y memoria son las que permite Fargate.
- **Los créditos se consumen sin factura visible.** Hay que acordarse de bajar las tasks al terminar.

## Cuándo revisar esta decisión

- Si Tero pasa a producción real con carga constante: evaluar EC2 con instancias reservadas o Savings Plans.
- Si algún servicio necesita acceso al host, GPU o configuraciones que Fargate no soporta.

## Referencias

- Letra del obligatorio, sección 3 (Containerización y Despliegue).
- Teórico *T06.1 Contenedores en AWS*, sección "ECS vs Fargate".
- Laboratorio del curso T06 › Containers › ECS.
- ADR-0014: Infraestructura en una cuenta AWS propia.

## Uso de IA

El análisis de opciones y costos, y la redacción de este ADR, se hicieron con apoyo de Claude (Anthropic). El equipo verificó la decisión contra la letra y el material del curso.
