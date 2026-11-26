resource "aws_security_group" "alb" {
  name        = "${var.name}-alb-sg"
  description = "ALB: recibe HTTP desde internet"
  vpc_id      = aws_vpc.this.id

  tags = {
    Name = "${var.name}-alb-sg"
  }
}

resource "aws_security_group" "ecs" {
  name        = "${var.name}-ecs-sg"
  description = "Servicios ECS: Reciben trafico del ALB y entre ellos"
  vpc_id      = aws_vpc.this.id

  tags = {
    Name = "${var.name}-ecs-sg"
  }
}

resource "aws_security_group" "rds" {
  name        = "${var.name}-rds-sg"
  description = "Postgres: solo acepta conexiones de los servicios ECS"
  vpc_id      = aws_vpc.this.id

  tags = {
    Name = "${var.name}-rds-sg"
  }
}

resource "aws_security_group" "redis" {
  name        = "${var.name}-redis-sg"
  description = "Redis: solo acepta conexiones de los servicios ECS"
  vpc_id      = aws_vpc.this.id

  tags = {
    Name = "${var.name}-redis-sg"
  }
}

# Reglas del ALB

resource "aws_vpc_security_group_ingress_rule" "alb_http" {
  security_group_id = aws_security_group.alb.id
  description       = "HTTP desde internet"
  ip_protocol       = "tcp"
  from_port         = 80
  to_port           = 80
  cidr_ipv4         = "0.0.0.0/0"
}

resource "aws_vpc_security_group_egress_rule" "alb_to_ecs" {
  security_group_id            = aws_security_group.alb.id
  description                  = "Hacia los servicios ECS"
  ip_protocol                  = "tcp"
  from_port                    = var.app_port
  to_port                      = var.app_port
  referenced_security_group_id = aws_security_group.ecs.id
}

# Reglas de ECS

resource "aws_vpc_security_group_ingress_rule" "ecs_from_alb" {
  security_group_id            = aws_security_group.ecs.id
  description                  = "Trafico desde el ALB"
  ip_protocol                  = "tcp"
  from_port                    = var.app_port
  to_port                      = var.app_port
  referenced_security_group_id = aws_security_group.alb.id
}

resource "aws_vpc_security_group_ingress_rule" "ecs_from_ecs" {
  security_group_id            = aws_security_group.ecs.id
  description                  = "Llamadas entre servicios"
  ip_protocol                  = "tcp"
  from_port                    = var.app_port
  to_port                      = var.app_port
  referenced_security_group_id = aws_security_group.ecs.id
}

resource "aws_vpc_security_group_egress_rule" "ecs_all" {
  security_group_id = aws_security_group.ecs.id
  description       = "Salida: ECR, logs, secretos, RDS, Redis y otros servicios"
  ip_protocol       = "-1"
  cidr_ipv4         = "0.0.0.0/0"
}

resource "aws_vpc_security_group_ingress_rule" "rds_from_ecs" {
  security_group_id            = aws_security_group.rds.id
  description                  = "Postgres desde los servicios ECS"
  ip_protocol                  = "tcp"
  from_port                    = var.db_port
  to_port                      = var.db_port
  referenced_security_group_id = aws_security_group.ecs.id
}

resource "aws_vpc_security_group_ingress_rule" "redis_from_ecs" {
  security_group_id            = aws_security_group.redis.id
  description                  = "Redis desde los servicios ECS"
  ip_protocol                  = "tcp"
  from_port                    = var.redis_port
  to_port                      = var.redis_port
  referenced_security_group_id = aws_security_group.ecs.id
}