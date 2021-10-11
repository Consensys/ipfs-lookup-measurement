#!/usr/bin/env sh

set -e
: "${HOST_NAME:=node-undefined}"

# Start loki
echo "      host: $HOST_NAME" >> promtail-local-config.yaml
./promtail-linux-amd64 -config.file=promtail-local-config.yaml &
# Start grafana
# Start agent
./agent &
./go-ipfs/cmd/ipfs/ipfs init
./go-ipfs/cmd/ipfs/ipfs config --json Discovery.MDNS '{"Enabled": false}'
./go-ipfs/cmd/ipfs/ipfs daemon > /app/all.log 2>&1
