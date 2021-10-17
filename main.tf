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

variable "KEY" {
  type = string
}

resource "aws_instance" "ipfs-testing-monitor" {
  ami           = "ami-0567f647e75c7bc05"
  instance_type = "t2.small"
  tags = {
    Name = "ipfs-testing-monitor"
  }
  security_groups = ["security_ipfs_testing_monitor"]
  user_data       = <<-EOF
    #!/bin/sh
    cd /home/ubuntu/
    sudo apt-get update
    sudo apt install -y unzip
    wget https://github.com/grafana/loki/releases/download/v2.3.0/loki-linux-amd64.zip
    wget https://dl.grafana.com/oss/release/grafana-8.1.5.linux-amd64.tar.gz
    wget https://raw.githubusercontent.com/grafana/loki/v2.3.0/cmd/loki/loki-local-config.yaml
    wget https://raw.githubusercontent.com/ConsenSys/ipfs-lookup-measurement/main/monitor/grafana-datasources.yml
    wget https://raw.githubusercontent.com/ConsenSys/ipfs-lookup-measurement/main/monitor/grafana-dashboards.yml
    wget https://raw.githubusercontent.com/ConsenSys/ipfs-lookup-measurement/main/monitor/ipfs-dashboard.json
    unzip loki-linux-amd64.zip
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

resource "aws_instance" "ipfs-testing-node-1" {
  ami           = "ami-0567f647e75c7bc05"
  instance_type = "t2.small"
  tags = {
    Name = "ipfs-testing-node-1"
  }
  security_groups = ["security_ipfs_testing_node"]
  user_data       = <<-EOF
    #!/bin/sh
    cd /home/ubuntu/
    sudo apt-get update
    sudo apt install -y unzip git make build-essential
    wget https://github.com/grafana/loki/releases/download/v2.3.0/promtail-linux-amd64.zip
    wget https://golang.org/dl/go1.17.1.linux-amd64.tar.gz
    wget https://raw.githubusercontent.com/ConsenSys/ipfs-lookup-measurement/main/node/promtail-cloud-config.yaml
    unzip ./promtail-linux-amd64.zip
    sudo tar -C /usr/local -xzf go1.17.1.linux-amd64.tar.gz
    mkdir /home/ubuntu/go
    export HOME=/home/ubuntu
    export GOPATH=/home/ubuntu/go
    export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
    git clone https://github.com/ConsenSys/ipfs-lookup-measurement.git
    cd ipfs-lookup-measurement
    cd controller
    make agent
    cd ..
    cd ..
    git clone https://github.com/wcgcyx/go-libp2p-kad-dht.git
    cd go-libp2p-kad-dht
    git checkout more-logging
    cd ..
    git clone https://github.com/wcgcyx/go-ipfs.git
    cd go-ipfs
    git checkout more-logging
    echo "replace github.com/libp2p/go-libp2p-kad-dht => ../go-libp2p-kad-dht" >> go.mod
    go mod tidy
    make build > buildLog.txt 2>&1
    cd ..
    mkdir ./ipfs-tests/
    export KEY="${var.KEY}"
    export IP="${aws_instance.ipfs-testing-monitor.public_ip}"
    export PERFORMANCE_TEST_DIR=/home/ubuntu/ipfs-tests/
    export IPFS_PATH=/home/ubuntu/.ipfs
    export IPFS=/home/ubuntu/go-ipfs/cmd/ipfs/ipfs
    echo "$KEY" > ./ipfs-tests/.key
    echo "      host: node1" >> ./promtail-cloud-config.yaml
    echo "clients:" >> ./promtail-cloud-config.yaml
    echo "  - url: http://$IP:3100/loki/api/v1/push" >> ./promtail-cloud-config.yaml
    nohup ./promtail-linux-amd64 -config.file=promtail-cloud-config.yaml &
    ./go-ipfs/cmd/ipfs/ipfs init
    nohup ./go-ipfs/cmd/ipfs/ipfs daemon > /home/ubuntu/all.log 2>&1 &
    IPFS_LOGGING=INFO nohup ./ipfs-lookup-measurement/controller/agent > /home/ubuntu/agent.log 2>&1 &
  EOF
}

resource "aws_instance" "ipfs-testing-node-2" {
  ami           = "ami-0567f647e75c7bc05"
  instance_type = "t2.small"
  tags = {
    Name = "ipfs-testing-node-2"
  }
  security_groups = ["security_ipfs_testing_node"]
  user_data       = <<-EOF
    #!/bin/sh
    cd /home/ubuntu/
    sudo apt-get update
    sudo apt install -y unzip git make build-essential
    wget https://github.com/grafana/loki/releases/download/v2.3.0/promtail-linux-amd64.zip
    wget https://golang.org/dl/go1.17.1.linux-amd64.tar.gz
    wget https://raw.githubusercontent.com/ConsenSys/ipfs-lookup-measurement/main/node/promtail-cloud-config.yaml
    unzip ./promtail-linux-amd64.zip
    sudo tar -C /usr/local -xzf go1.17.1.linux-amd64.tar.gz
    mkdir /home/ubuntu/go
    export HOME=/home/ubuntu
    export GOPATH=/home/ubuntu/go
    export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
    git clone https://github.com/ConsenSys/ipfs-lookup-measurement.git
    cd ipfs-lookup-measurement
    cd controller
    make agent
    cd ..
    cd ..
    git clone https://github.com/wcgcyx/go-libp2p-kad-dht.git
    cd go-libp2p-kad-dht
    git checkout more-logging
    cd ..
    git clone https://github.com/wcgcyx/go-ipfs.git
    cd go-ipfs
    git checkout more-logging
    echo "replace github.com/libp2p/go-libp2p-kad-dht => ../go-libp2p-kad-dht" >> go.mod
    go mod tidy
    make build > buildLog.txt 2>&1
    cd ..
    mkdir ./ipfs-tests/
    export KEY="${var.KEY}"
    export IP="${aws_instance.ipfs-testing-monitor.public_ip}"
    export PERFORMANCE_TEST_DIR=/home/ubuntu/ipfs-tests/
    export IPFS_PATH=/home/ubuntu/.ipfs
    export IPFS=/home/ubuntu/go-ipfs/cmd/ipfs/ipfs
    echo "$KEY" > ./ipfs-tests/.key
    echo "      host: node2" >> ./promtail-cloud-config.yaml
    echo "clients:" >> ./promtail-cloud-config.yaml
    echo "  - url: http://$IP:3100/loki/api/v1/push" >> ./promtail-cloud-config.yaml
    nohup ./promtail-linux-amd64 -config.file=promtail-cloud-config.yaml &
    ./go-ipfs/cmd/ipfs/ipfs init
    nohup ./go-ipfs/cmd/ipfs/ipfs daemon > /home/ubuntu/all.log 2>&1 &
    IPFS_LOGGING=INFO nohup ./ipfs-lookup-measurement/controller/agent > /home/ubuntu/agent.log 2>&1 &
  EOF
}

resource "aws_instance" "ipfs-testing-node-3" {
  ami           = "ami-0567f647e75c7bc05"
  instance_type = "t2.small"
  tags = {
    Name = "ipfs-testing-node-3"
  }
  security_groups = ["security_ipfs_testing_node"]
  user_data       = <<-EOF
    #!/bin/sh
    cd /home/ubuntu/
    sudo apt-get update
    sudo apt install -y unzip git make build-essential
    wget https://github.com/grafana/loki/releases/download/v2.3.0/promtail-linux-amd64.zip
    wget https://golang.org/dl/go1.17.1.linux-amd64.tar.gz
    wget https://raw.githubusercontent.com/ConsenSys/ipfs-lookup-measurement/main/node/promtail-cloud-config.yaml
    unzip ./promtail-linux-amd64.zip
    sudo tar -C /usr/local -xzf go1.17.1.linux-amd64.tar.gz
    mkdir /home/ubuntu/go
    export HOME=/home/ubuntu
    export GOPATH=/home/ubuntu/go
    export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
    git clone https://github.com/ConsenSys/ipfs-lookup-measurement.git
    cd ipfs-lookup-measurement
    cd controller
    make agent
    cd ..
    cd ..
    git clone https://github.com/wcgcyx/go-libp2p-kad-dht.git
    cd go-libp2p-kad-dht
    git checkout more-logging
    cd ..
    git clone https://github.com/wcgcyx/go-ipfs.git
    cd go-ipfs
    git checkout more-logging
    echo "replace github.com/libp2p/go-libp2p-kad-dht => ../go-libp2p-kad-dht" >> go.mod
    go mod tidy
    make build > buildLog.txt 2>&1
    cd ..
    mkdir ./ipfs-tests/
    export KEY="${var.KEY}"
    export IP="${aws_instance.ipfs-testing-monitor.public_ip}"
    export PERFORMANCE_TEST_DIR=/home/ubuntu/ipfs-tests/
    export IPFS_PATH=/home/ubuntu/.ipfs
    export IPFS=/home/ubuntu/go-ipfs/cmd/ipfs/ipfs
    echo "$KEY" > ./ipfs-tests/.key
    echo "      host: node3" >> ./promtail-cloud-config.yaml
    echo "clients:" >> ./promtail-cloud-config.yaml
    echo "  - url: http://$IP:3100/loki/api/v1/push" >> ./promtail-cloud-config.yaml
    nohup ./promtail-linux-amd64 -config.file=promtail-cloud-config.yaml &
    ./go-ipfs/cmd/ipfs/ipfs init
    nohup ./go-ipfs/cmd/ipfs/ipfs daemon > /home/ubuntu/all.log 2>&1 &
    IPFS_LOGGING=INFO nohup ./ipfs-lookup-measurement/controller/agent > /home/ubuntu/agent.log 2>&1 &
  EOF
}

resource "aws_instance" "ipfs-testing-node-4" {
  ami           = "ami-0567f647e75c7bc05"
  instance_type = "t2.small"
  tags = {
    Name = "ipfs-testing-node-4"
  }
  security_groups = ["security_ipfs_testing_node"]
  user_data       = <<-EOF
    #!/bin/sh
    cd /home/ubuntu/
    sudo apt-get update
    sudo apt install -y unzip git make build-essential
    wget https://github.com/grafana/loki/releases/download/v2.3.0/promtail-linux-amd64.zip
    wget https://golang.org/dl/go1.17.1.linux-amd64.tar.gz
    wget https://raw.githubusercontent.com/ConsenSys/ipfs-lookup-measurement/main/node/promtail-cloud-config.yaml
    unzip ./promtail-linux-amd64.zip
    sudo tar -C /usr/local -xzf go1.17.1.linux-amd64.tar.gz
    mkdir /home/ubuntu/go
    export HOME=/home/ubuntu
    export GOPATH=/home/ubuntu/go
    export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
    git clone https://github.com/ConsenSys/ipfs-lookup-measurement.git
    cd ipfs-lookup-measurement
    cd controller
    make agent
    cd ..
    cd ..
    git clone https://github.com/wcgcyx/go-libp2p-kad-dht.git
    cd go-libp2p-kad-dht
    git checkout more-logging
    cd ..
    git clone https://github.com/wcgcyx/go-ipfs.git
    cd go-ipfs
    git checkout more-logging
    echo "replace github.com/libp2p/go-libp2p-kad-dht => ../go-libp2p-kad-dht" >> go.mod
    go mod tidy
    make build > buildLog.txt 2>&1
    cd ..
    mkdir ./ipfs-tests/
    export KEY="${var.KEY}"
    export IP="${aws_instance.ipfs-testing-monitor.public_ip}"
    export PERFORMANCE_TEST_DIR=/home/ubuntu/ipfs-tests/
    export IPFS_PATH=/home/ubuntu/.ipfs
    export IPFS=/home/ubuntu/go-ipfs/cmd/ipfs/ipfs
    echo "$KEY" > ./ipfs-tests/.key
    echo "      host: node4" >> ./promtail-cloud-config.yaml
    echo "clients:" >> ./promtail-cloud-config.yaml
    echo "  - url: http://$IP:3100/loki/api/v1/push" >> ./promtail-cloud-config.yaml
    nohup ./promtail-linux-amd64 -config.file=promtail-cloud-config.yaml &
    ./go-ipfs/cmd/ipfs/ipfs init
    nohup ./go-ipfs/cmd/ipfs/ipfs daemon > /home/ubuntu/all.log 2>&1 &
    IPFS_LOGGING=INFO nohup ./ipfs-lookup-measurement/controller/agent > /home/ubuntu/agent.log 2>&1 &
  EOF
}

resource "aws_instance" "ipfs-testing-node-5" {
  ami           = "ami-0567f647e75c7bc05"
  instance_type = "t2.small"
  tags = {
    Name = "ipfs-testing-node-5"
  }
  security_groups = ["security_ipfs_testing_node"]
  user_data       = <<-EOF
    #!/bin/sh
    cd /home/ubuntu/
    sudo apt-get update
    sudo apt install -y unzip git make build-essential
    wget https://github.com/grafana/loki/releases/download/v2.3.0/promtail-linux-amd64.zip
    wget https://golang.org/dl/go1.17.1.linux-amd64.tar.gz
    wget https://raw.githubusercontent.com/ConsenSys/ipfs-lookup-measurement/main/node/promtail-cloud-config.yaml
    unzip ./promtail-linux-amd64.zip
    sudo tar -C /usr/local -xzf go1.17.1.linux-amd64.tar.gz
    mkdir /home/ubuntu/go
    export HOME=/home/ubuntu
    export GOPATH=/home/ubuntu/go
    export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
    git clone https://github.com/ConsenSys/ipfs-lookup-measurement.git
    cd ipfs-lookup-measurement
    cd controller
    make agent
    cd ..
    cd ..
    git clone https://github.com/wcgcyx/go-libp2p-kad-dht.git
    cd go-libp2p-kad-dht
    git checkout more-logging
    cd ..
    git clone https://github.com/wcgcyx/go-ipfs.git
    cd go-ipfs
    git checkout more-logging
    echo "replace github.com/libp2p/go-libp2p-kad-dht => ../go-libp2p-kad-dht" >> go.mod
    go mod tidy
    make build > buildLog.txt 2>&1
    cd ..
    mkdir ./ipfs-tests/
    export KEY="${var.KEY}"
    export IP="${aws_instance.ipfs-testing-monitor.public_ip}"
    export PERFORMANCE_TEST_DIR=/home/ubuntu/ipfs-tests/
    export IPFS_PATH=/home/ubuntu/.ipfs
    export IPFS=/home/ubuntu/go-ipfs/cmd/ipfs/ipfs
    echo "$KEY" > ./ipfs-tests/.key
    echo "      host: node5" >> ./promtail-cloud-config.yaml
    echo "clients:" >> ./promtail-cloud-config.yaml
    echo "  - url: http://$IP:3100/loki/api/v1/push" >> ./promtail-cloud-config.yaml
    nohup ./promtail-linux-amd64 -config.file=promtail-cloud-config.yaml &
    ./go-ipfs/cmd/ipfs/ipfs init
    nohup ./go-ipfs/cmd/ipfs/ipfs daemon > /home/ubuntu/all.log 2>&1 &
    IPFS_LOGGING=INFO nohup ./ipfs-lookup-measurement/controller/agent > /home/ubuntu/agent.log 2>&1 &
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

resource "aws_security_group" "security_ipfs_testing_node" {
  name        = "security_ipfs_testing_node"
  description = "security group for ipfs testing node"

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
    from_port   = 4001
    to_port     = 4001
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 4001
    to_port     = 4001
    protocol    = "udp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 3030
    to_port     = 3030
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "security_ipfs_testing_node"
  }
}