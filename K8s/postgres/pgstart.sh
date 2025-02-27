#!/bin/bash

# Start Minikube
minikube start

# Create the Persistent Volume (PV)
cat <<EOF | kubectl apply -f pg-volume.yaml
EOF

# Create the Persistent Volume Claim (PVC)
cat <<EOF | kubectl apply -f pg-pvclaim.yaml
EOF

# Create the Postgres Service
cat <<EOF | kubectl apply -f pg-service.yaml
EOF

# Create the Postgres Secret
cat <<EOF | kubectl apply -f pg-secrets.yaml
EOF

# Create the Postgres ConfigMap
cat <<EOF | kubectl apply -f pg-config.yaml
EOF

# Create the Postgres Deployment
cat <<EOF | kubectl apply -f pg-depl.yaml
EOF

echo "PostgreSQL deployment completed."