#!/bin/bash

if ! command -v kubectl &> /dev/null; then
    echo "kubectl not found. Please install kubectl."
    exit 1
fi

if ! command -v kind &> /dev/null; then
    echo "KinD not found. Please install KinD."
    exit 1
fi

if ! command -v flux &> /dev/null; then
    echo "Flux not found. Please install Flux."
    exit 1
fi

if [[ ! $GITHUB_REPO || ! $GITHUB_USER || ! $GITHUB_REPO ]]; then
    echo "No GITHUB environment variables found"
    exit 1
fi

clusters_path="$(readlink -f "$0" | awk -F'scripts/start.sh' '{print $1}')"
clusters=$(ls "${clusters_path}clusters/" | xargs -n 1 basename | tr '\n' ' ')
for cluster in ${clusters}; do
    echo "-------------------------------------------"
    echo "Cluster: ${cluster}"
    context="kind-${cluster}"
    kind create cluster --name="$cluster"
    cluster_status=$(kubectl cluster-info 2>&1)
    if [[ "$cluster_status" == *"error"* || "$cluster_status" == *"Error"* ]]; then
        echo "An error occurred while creating the KinD cluster:"
        echo "$cluster_status"
        exit 1
    fi
    attempts=3
    elapsed_attempts=0
    is_ok=false
    while [ $elapsed_attempts -lt $attempts ]; do
        nodes_status=$(kubectl --context="$context" get nodes --output=jsonpath='{.items[*].status.conditions[?(@.type=="Ready")].status}')
        if [[ "$nodes_status" == *"True"* ]]; then
            is_ok=true
            break
        fi
        sleep 5
        ((elapsed_attempts++))
    done
    if [ $is_ok ]; then
        flux bootstrap github \
            --context="$context" \
            --owner="$GITHUB_USER" \
            --repository="$GITHUB_REPO" \
            --branch="main" \
            --path="infra/clusters/$cluster"
    fi
done