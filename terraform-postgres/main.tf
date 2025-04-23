provider "aws" {
  region = "eu-north-1" # change as needed
}

resource "aws_db_subnet_group" "default" {
  name       = "default-subnet-group"
  subnet_ids = data.aws_subnets.default.ids
  tags = {
    Name = "Default subnet group"
  }
}

data "aws_vpc" "default" {
  default = true
}

data "aws_subnets" "default" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.default.id]
  }
}

resource "aws_security_group" "rds_sg" {
  name        = "rds-sg"
  description = "Allow PostgreSQL access"
  vpc_id      = data.aws_vpc.default.id

  ingress {
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"] # can limit to your IP for better security
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "rds-sg"
  }
}

resource "aws_db_instance" "postgres" {
  identifier        = "my-free-postgres"
  allocated_storage = 20
  engine            = "postgres"
  engine_version    = "17.2"
  instance_class    = "db.t3.micro"
  username          = "postgres"
  password          = "postgrespassword123" # Replace with a secure secret
  db_subnet_group_name = aws_db_subnet_group.default.name
  vpc_security_group_ids = [aws_security_group.rds_sg.id]
  publicly_accessible    = true
  skip_final_snapshot    = true
  delete_automated_backups = true

  tags = {
    Name = "My Free Tier PostgreSQL DB"
  }
}
