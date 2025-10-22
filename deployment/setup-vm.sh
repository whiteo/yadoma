#!/bin/bash

# Yadoma Agent - GCP VM Setup Script
# This script prepares a Google Cloud VM for running the Yadoma agent

set -e

echo "============================================"
echo "Yadoma Agent - GCP VM Setup Script"
echo "============================================"
echo ""

# Check if running as root
if [ "$EUID" -ne 0 ]; then
   echo "Please run as root (use sudo)"
   exit 1
fi

# Update system
echo "[1/6] Updating system packages..."
apt-get update
apt-get upgrade -y

# Install Docker if not already installed
if ! command -v docker &> /dev/null; then
    echo "[2/6] Installing Docker..."
    curl -fsSL https://get.docker.com -o get-docker.sh
    sh get-docker.sh
    rm get-docker.sh
    systemctl enable docker
    systemctl start docker
    echo "Docker installed successfully"
else
    echo "[2/6] Docker is already installed"
fi

# Verify Docker is running
if ! systemctl is-active --quiet docker; then
    echo "ERROR: Docker is not running"
    exit 1
fi

# Create yadoma deployment directory
echo "[3/6] Creating deployment directory..."
mkdir -p /opt/yadoma
chown -R root:root /opt/yadoma

# Create logs directory
mkdir -p /var/log/yadoma
chown -R root:root /var/log/yadoma

# Configure firewall for gRPC port 50051
echo "[4/6] Configuring firewall..."
if command -v ufw &> /dev/null; then
    ufw allow 50051/tcp comment 'Yadoma Agent gRPC'
    echo "Firewall configured (ufw)"
elif command -v firewall-cmd &> /dev/null; then
    firewall-cmd --permanent --add-port=50051/tcp
    firewall-cmd --reload
    echo "Firewall configured (firewalld)"
else
    echo "No firewall detected, skipping firewall configuration"
fi

# Set up log rotation
echo "[5/6] Setting up log rotation..."
cat > /etc/logrotate.d/yadoma-agent << 'EOF'
/var/log/yadoma/*.log {
    daily
    rotate 7
    compress
    delaycompress
    missingok
    notifempty
    create 0644 root root
    sharedscripts
    postrotate
        systemctl reload yadoma-agent > /dev/null 2>&1 || true
    endscript
}
EOF

# Display system information
echo "[6/6] System information:"
echo "  OS: $(lsb_release -d | cut -f2)"
echo "  Kernel: $(uname -r)"
echo "  Docker version: $(docker --version)"
echo "  Deployment path: /opt/yadoma"
echo "  Logs path: /var/log/yadoma"
echo ""

echo "============================================"
echo "Setup completed successfully!"
echo "============================================"
echo ""
echo "Next steps:"
echo "1. Set up GitHub secrets for CI/CD:"
echo "   - GCP_SSH_PRIVATE_KEY: Your private SSH key"
echo "   - GCP_VM_IP: VM external IP address"
echo "   - GCP_VM_USER: SSH username"
echo ""
echo "2. Push to master/main branch to trigger deployment"
echo ""
echo "3. Check agent status after deployment:"
echo "   sudo systemctl status yadoma-agent"
echo "   sudo journalctl -u yadoma-agent -f"
echo ""
