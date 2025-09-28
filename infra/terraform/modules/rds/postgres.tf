resource "aws_db_subnet_group" "postgres_subnet_group" {
  name       = "${var.project_name}-postgres-subnet-group"
  subnet_ids = [var.private_subnet_1a, var.private_subnet_1b]

  tags = {
    Name = "${var.project_name}-postgres-subnet-group"
  }
}

resource "aws_security_group" "rds_sg" {
  name_prefix = "${var.project_name}-rds-"
  vpc_id      = var.vpc_id

  ingress {
    description     = "PostgreSQL access from EKS"
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = [var.cluster_security_group_id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "${var.project_name}-rds-security-group"
  }
}

resource "aws_db_instance" "postgres_instance" {
  identifier     = "${var.project_name}-postgres"
  engine         = "postgres"
  engine_version = "17.4"
  instance_class = "db.t4g.micro"

  allocated_storage     = 20
  max_allocated_storage = 100
  storage_type          = "gp3"
  storage_encrypted     = true

  db_name  = var.db_name
  username = var.db_username
  password = var.db_password

  db_subnet_group_name   = aws_db_subnet_group.postgres_subnet_group.name
  vpc_security_group_ids = [aws_security_group.rds_sg.id]
  publicly_accessible    = false

  skip_final_snapshot = true # Change to false in production

  tags = merge(
    {
      Name = "${var.project_name}-postgres"
    },
    var.tags
  )

  lifecycle {
    prevent_destroy = false # Change to true in production
  }
}
