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

provider "aws" {
  alias  = "us-east-2"
  region = "us-east-2"
}

provider "aws" {
  alias  = "ap-southeast-2"
  region = "ap-southeast-2"
}

resource "aws_vpc" "ipfs-vpc-1" {
  cidr_block = "10.1.0.0/16"
  provider = aws.us-east-2
  tags = {
    Name = "ipfs-pegasys"
  }
}

resource "aws_vpc" "ipfs-vpc-2" {
  cidr_block = "10.2.0.0/16"
  provider = aws.ap-southeast-2
  tags = {
    Name = "ipfs-pegasys"
  }
}

resource "aws_subnet" "ipfs-subnet-1" {
  cidr_block = "10.0.1.0/24"
  vpc_id     = aws_vpc.ipfs-vpc-1.id
  availability_zone = "us-east-2"
  tags = {
    Name = "ipfs-pegasys"
  }
}

resource "aws_subnet" "ipfs-subnet-2" {
  cidr_block = "10.0.2.0/24"
  vpc_id     = aws_vpc.ipfs-vpc-2.id
  availability_zone = "ap-southeast-2"
  tags = {
    Name = "ipfs-pegasys"
  }
}

resource "aws_internet_gateway" "gateway-1" {
  vpc_id = aws_vpc.ipfs-vpc-1.id
  tags = {
    Name = "ipfs-pegasys"
  }
}

resource "aws_internet_gateway" "gateway-2" {
  vpc_id = aws_vpc.ipfs-vpc-2.id
  tags = {
    Name = "ipfs-pegasys"
  }
}

resource "aws_route_table" "route-1" {
  vpc_id = aws_vpc.ipfs-vpc-1.id
  depends_on = [aws_internet_gateway.gateway-1, aws_internet_gateway.gateway-2]
  route = [
    {
      cidr_block =  "10.0.1.0/24"
      local_gateway_id = aws_internet_gateway.gateway-1.id
    }, {
      cidr_block =  "10.0.2.0/24"
      nat_gateway_id = aws_internet_gateway.gateway-2.id
    }
  ]
  tags = {
    Name = "ipfs-pegasys"
  }
}

resource "aws_route_table" "route-2" {
  vpc_id = aws_vpc.ipfs-vpc-2.id
  route = [
    {
      cidr_block =  "10.0.2.0/24"
      local_gateway_id = aws_internet_gateway.gateway-2.id
    }, {
      cidr_block =  "10.0.1.0/24"
      nat_gateway_id = aws_internet_gateway.gateway-1.id
    }
  ]
  tags = {
    Name = "ipfs-pegasys"
  }
}

resource "aws_route_table_association" "rta-1" {
  route_table_id = aws_route_table.route-1.id
  subnet_id = aws_subnet.ipfs-subnet-1.id
}

resource "aws_route_table_association" "rta-2" {
  route_table_id = aws_route_table.route-2.id
  subnet_id = aws_subnet.ipfs-subnet-2.id
}

resource "aws_vpc_peering_connection" "peer" {
  peer_vpc_id   = aws_vpc.ipfs-vpc-1.id
  vpc_id        = aws_vpc.ipfs-vpc-2.id
  auto_accept   = true

  accepter {
    allow_remote_vpc_dns_resolution = true
  }

  requester {
    allow_remote_vpc_dns_resolution = true
  }

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
  subnet_id = aws_subnet.ipfs-subnet-1.id
}

resource "aws_instance" "node2" {
  ami = "ami-07738c5c0ee584ed1"
  instance_type = "t2.small"
  tags = {
    Name = "ipfs-pegasys"
  }
  subnet_id = aws_subnet.ipfs-subnet-1.id
}

# nodes  in a separate network

resource "aws_instance" "node3" {
  ami = "ami-07738c5c0ee584ed1"
  instance_type = "t2.small"
  tags = {
    Name = "ipfs-pegasys"
  }
  subnet_id = aws_subnet.ipfs-subnet-2.id
}

resource "aws_instance" "node4" {
  ami = "ami-07738c5c0ee584ed1"
  instance_type = "t2.small"
  tags = {
    Name = "ipfs-pegasys"
  }
  subnet_id = aws_subnet.ipfs-subnet-2.id
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