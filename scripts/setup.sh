#!/bin/bash

echo "🚀 Setting up development environment..."
echo "Checking prerequisites..."

# Check Docker
if ! command -v docker &> /dev/null; then
    echo "❌ Docker not found. Please install Rancher Desktop: https://rancherdesktop.io/"
    exit 1
else
    echo "✅ Docker is installed (version: $(docker --version))"
fi

# Check kubectl
if ! command -v kubectl &> /dev/null; then
    echo "❌ kubectl not found. Please install Rancher Desktop: https://rancherdesktop.io/"
    exit 1
else
    echo "✅ kubectl is installed (version: $(kubectl version --client --short 2>/dev/null || kubectl version --client))"
fi

# Check Go
if ! command -v go &> /dev/null; then
    echo "❌ Go not found. Please install Go from: https://golang.org/dl/"
    exit 1
else
    echo "✅ Go is installed (version: $(go version))"
fi

# Check and install Tilt if needed
if ! command -v tilt &> /dev/null; then
    echo "⚠️  Tilt not found. Installing..."
    curl -fsSL https://raw.githubusercontent.com/tilt-dev/tilt/master/scripts/install.sh | bash
    echo "✅ Tilt installed successfully"
else
    echo "✅ Tilt is already installed (version: $(tilt version))"
fi

echo "🎉 Development environment is ready!"