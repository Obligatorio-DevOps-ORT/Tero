# ADR-0010: Salida a internet de las subredes privadas

- **Estado:** Aceptada
- **Fecha:** 2026-11-19
- **Tarea:** S2-10, VPC completa con subredes públicas, privadas y security groups (issue #26)
- **Área:** Infra (Terraform, `infra/modules/network`)

## Contexto

El auditor observó que la base de datos y Redis están publicados en puertos de la máquina anfitriona. El plan del equipo resuelve eso así: **solo el ALB es público**. Los servicios (ECS Fargate), RDS y ElastiCache van en subredes privadas.

Un recurso en una subred privada no tiene IP pública y no puede salir a internet por el Internet Gateway. Pero las tareas de ECS necesitan llegar a servicios de AWS para arrancar y operar:

| Necesidad | Servicio |
|---|---|
| Bajar la imagen | ECR (API + registry) y S3 (las capas de las imágenes) |
| Enviar logs | CloudWatch Logs |
| Leer credenciales | Secrets Manager |

RDS y ElastiCache **no necesitan salida a internet**. Solo reciben conexiones de los servicios.

Restricciones:

- Cuenta AWS personal en el plan gratuito, con **créditos limitados**, región `us-east-2`.
- Tres ambientes (dev, test, prod) con workspaces del mismo root.
- La letra del obligatorio no exige un mecanismo de salida concreto. Exige IaC con Terraform y despliegue en un orquestador (ECS/EKS).

## Opciones consideradas

Costos aproximados de `us-east-2` al momento de decidir. Hay que verificarlos en la página de precios de AWS.

### A. Un NAT Gateway por ambiente, en una subred pública

- Se crea una Elastic IP, el NAT Gateway en la subred pública 1, y una route table privada con `0.0.0.0/0 → NAT`.
- **Costo:** unos USD 0.045/h (~USD 33/mes) por NAT mientras exista, más USD 0.045/GB procesado. **Cobra aunque no haya tráfico.**
- Es el patrón del laboratorio del curso *T04 Cloud › 02-networking › 03-Subnets Públicas y Privadas*.

### B. Interface VPC endpoints (PrivateLink)

- Un endpoint por servicio (`ecr.api`, `ecr.dkr`, `logs`, `secretsmanager`) en cada AZ, más un endpoint gateway de S3.
- **Costo:** USD 0.01/h por endpoint por AZ. Con 4 endpoints × 2 AZ son ~USD 58/mes, **más caro que un NAT**.
- No da salida a internet en general, solo a los servicios de AWS listados.
- No se ve en el curso.

### C. Tareas de ECS en subredes públicas con IP pública

- Las tareas se exponen con `assign_public_ip = true`, y el security group solo acepta tráfico desde el ALB. RDS y Redis siguen en subredes privadas.
- **Costo:** ~USD 0.005/h por IP pública (~USD 3.6/mes por tarea).
- Es lo que hace el laboratorio del curso *T06 › ECS › Load Balancing*, por simplicidad.
- Contradice la decisión del plan de que "solo el ALB es público". Además, cada tarea queda direccionable desde internet y la única barrera es el security group.

## Decisión

**Opción A**, con dos ajustes por costo:

1. **Un solo NAT Gateway por ambiente**, no uno por AZ.
2. **El NAT es activable por variable** (`enable_nat_gateway`). Vale `false` por defecto en el módulo y en los tfvars de los 3 ambientes, así que el estado seguro es no gastar. Se prende cuando el ambiente se va a usar (`-var="enable_nat_gateway=true"`) y se apaga con un `apply` normal. Solo nacen y mueren el NAT, su EIP y la ruta privada; la VPC, las subredes y las route tables quedan.

Las tareas de ECS, RDS y ElastiCache van en las subredes privadas.

## Consecuencias

**Positivas**

- Las tareas no tienen IP pública. Para llegar desde internet hay que pasar por el ALB, que es defensa en profundidad además de los security groups.
- Sigue el patrón que enseña el curso, así que el equipo y quien opere después lo pueden entender y mantener.
- El toggle evita gastar créditos en NATs ociosos. Con los 3 ambientes prendidos todo el mes serían unos USD 100.

**Negativas y riesgos aceptados**

- **Punto único de falla en la salida:** si cae la AZ del NAT, las tareas de la otra AZ pierden salida. Se acepta para un proyecto académico. En producción real iría un NAT por AZ, con una route table privada por AZ.
- **Si el NAT está apagado o se elimina:** las tareas nuevas no pueden bajar la imagen de ECR y quedan fallando al arrancar. Las que ya corren siguen atendiendo tráfico, pero fallan al enviar logs o leer secretos. Esto hay que tenerlo presente al prender o apagar un ambiente.
- Hay costo por GB procesado. Mejora posible: agregar el endpoint **gateway** de S3, que es gratis, para que las capas de las imágenes no pasen por el NAT.

## Cuándo revisar esta decisión

- Si Tero pasa a producción real: un NAT por AZ.
- Si el costo de los datos procesados por el NAT se vuelve relevante: endpoint gateway de S3, y evaluar los interface endpoints.
- Si algún servicio necesita llamar APIs externas: el NAT ya lo cubre. Con la opción B no sería posible.

## Referencias

- Letra del obligatorio, secciones 3 (Containerización y Despliegue) y 7 (IaC).
- Plan del equipo: *Obligatorio DevOps — Propuesta de plan (Tero)*, sección de red.
- Laboratorio del curso T04 › 02-networking › 03-Subnets Públicas y Privadas.
- Laboratorio del curso T05 › 03-02 Meta-Argument `depends_on`.
- Laboratorio del curso T06 › ECS › 04 Load Balancing.

## Uso de IA

El análisis de opciones y costos, y la redacción de este ADR, se hicieron con apoyo de Claude (Anthropic). El equipo verificó la decisión contra la letra, el plan y los laboratorios del curso.
