# use github.com/fussion-suite/ipfs-cluster-aws instead
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
  region = "us-east-2"
}

resource "aws_vpc" "ipfs-vpc" {
  cidr_block = "10.0.0.0/16"
  tags = {
    Name = "ipfs-pegasys"
  }
}

resource "aws_subnet" "ipfs-subnet" {
  cidr_block = "10.0.1.0/24"
  vpc_id     = aws_vpc.ipfs-vpc.id
  tags = {
    Name = "ipfs-pegasys"
  }
}

resource "aws_instance" "node1" {
  ami = "ami-07738c5c0ee584ed1"
  instance_type = "t2.small"
  tags = {
    Name = "ipfs-pegasys"
  }
  subnet_id = aws_subnet.ipfs-subnet.id
}

resource "aws_instance" "node2" {
  ami = "ami-07738c5c0ee584ed1"
  instance_type = "t2.small"
  tags = {
    Name = "ipfs-pegasys"
  }
  subnet_id = aws_subnet.ipfs-subnet.id
}

resource "aws_instance" "node3" {
  ami = "ami-07738c5c0ee584ed1"
  instance_type = "t2.small"
  tags = {
    Name = "ipfs-pegasys"
  }
  subnet_id = aws_subnet.ipfs-subnet.id
}



## don't have permission
#resource "aws_budgets_budget" "ec2" {
#  name              = "budget-ec2-monthly"
#  budget_type       = "COST"
#  limit_amount      = "500"
#  limit_unit        = "USD"
#  time_period_end   = "2021-10-30_00:00"
#  time_period_start = "2021-09-28_00:00"
#  time_unit         = "MONTHLY"
#
#  notification {
#    comparison_operator        = "GREATER_THAN"
#    threshold                  = 100
#    threshold_type             = "PERCENTAGE"
#    notification_type          = "FORECASTED"
#    subscriber_email_addresses = ["igor.zenyuk@gmail.com"]
#  }
#}