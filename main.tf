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
  user_data = <<-EOF
    #!/bin/sh
    sudo apt-get update
    sudo apt install -y unzip
    wget https://github.com/grafana/loki/releases/download/v2.3.0/loki-linux-amd64.zip
    wget https://dl.grafana.com/oss/release/grafana-8.1.5.linux-amd64.tar.gz
    wget https://raw.githubusercontent.com/grafana/loki/v2.3.0/cmd/loki/loki-local-config.yaml
    wget https://raw.githubusercontent.com/ConsenSys/ipfs-lookup-measurement/main/monitor/grafana-datasources.yml
    wget https://raw.githubusercontent.com/ConsenSys/ipfs-lookup-measurement/main/monitor/grafana-dashboards.yml
    wget https://raw.githubusercontent.com/ConsenSys/ipfs-lookup-measurement/main/monitor/ipfs-dashboard.json
    unzip loki-linux-amd64.zip -y
    tar -zxvf grafana-8.1.5.linux-amd64.tar.gz
    mv grafana-datasources.yml ./grafana-8.1.5/conf/provisioning/datasources/datasources.yml
    mv grafana-dashboards.yml ./grafana-8.1.5/conf/provisioning/dashboards/dashboards.yml
    sudo mkdir --parents /var/lib/grafana/dashboards
    mv ipfs-dashboard.json /var/lib/grafana/dashboards/
    nohup ./loki-linux-amd64 -config.file=loki-local-config.yaml &
    cd ./grafana-8.1.5/bin
    nohup ./grafana-server &
  EOF
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