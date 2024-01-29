#!/bin/bash

while [[ $# -gt 0 ]]; do
    key="$1"
    case $key in
        -n|--name)
            name="$2"
            shift 2
            ;;
        *)
            echo "Unknown option or argument: $1"
            exit 1
            ;;
    esac
done

if [ -z "$name" ]; then
    read -p "Enter a cluster name: " name
fi

# Use the value of the --name parameter
echo "Value of the --name parameter: $name"