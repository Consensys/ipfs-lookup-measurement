output "monitor_ip" {
  description = "Public IP address of monitor"
  value       = aws_instance.ipfs-testing-monitor.public_ip
}

output "node_1_ip" {
  description = "Public IP address of node 1"
  value       = aws_instance.ipfs-testing-node-1.public_ip
}

output "node_2_ip" {
  description = "Public IP address of node 2"
  value       = aws_instance.ipfs-testing-node-2.public_ip
}

output "node_3_ip" {
  description = "Public IP address of node 3"
  value       = aws_instance.ipfs-testing-node-3.public_ip
}

output "node_4_ip" {
  description = "Public IP address of node 4"
  value       = aws_instance.ipfs-testing-node-4.public_ip
}

output "node_5_ip" {
  description = "Public IP address of node 5"
  value       = aws_instance.ipfs-testing-node-5.public_ip
}
