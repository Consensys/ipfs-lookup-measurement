terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.5.0"
    }
  }
  required_version = ">= 0.13"
}

provider "aws" {
  region = "ap-southeast-2"
}

resource "aws_instance" "ipfs-testing-monitor" {
  ami           = "ami-0567f647e75c7bc05"
  instance_type = "t2.small"
  tags = {
    Name = "ipfs-testing-monitor"
  }
  security_groups = ["security_ipfs_testing_monitor"]
}

resource "aws_security_group" "security_ipfs_testing_monitor" {
  name        = "security_ipfs_testing_monitor"
  description = "security group for ipfs testing monitor"

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }


  ingress {
    from_port   = 3000
    to_port     = 3000
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 3100
    to_port     = 3100
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "security_ipfs_testing_monitor"
  }
}