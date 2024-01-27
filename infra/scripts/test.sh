#!/bin/bash

clusters_path="$(readlink -f "$0" | awk -F'scripts/test.sh' '{print $1}')"
clusters=$(ls "${clusters_path}clusters/" | xargs -n 1 basename | tr '\n' ' ')
for cluster in ${clusters}; do
    echo "-------------------------------------------"
    echo "Cluster: ${cluster}"
    context="kind-${cluster}"
    kubectl --context="$context" -n ingress-nginx port-forward svc/ingress-nginx-controller 8080:80 > /dev/null &
    pid=$!
    attempts=3
    elapsed_time=0
    is_ok=false
    while [ $elapsed_time -lt $attempts ]; do
        if nc -z localhost 8080 2>/dev/null; then
            is_ok=true
            break
        fi
        sleep 5
        ((elapsed_attempts++))
    done

    if [ $is_ok ]; then
        curl -H "Host: podinfo.${cluster}" http://localhost:8080
    else
        echo "Timeout error."
    fi
    if ps -p $pid > /dev/null; then
        kill $pid > /dev/null 2>&1
        sleep 1
    fi
done
