#!/bin/bash

# Set Exit on fail param!!
set -e
# echo "Waiting for vault to start!!"

function check_vault_is_running() {
    echo "Waiting for vault to start!!"
    
    # Loop until a condition is met!!
    
    until response=$(curl -s http://localhost:8200/v1/sys/health?standbyok=true) && echo "$response" | grep "initialized"; do
        echo "Response: $response"
        sleep 5
    done
    
    
    echo "Vault is up and running at port 8200"
}

function check_cache_is_up() {
    echo "Checking if Redis is up..."
    
    until redis-cli -h localhost -p 6379 ping | grep -q "PONG"; do
        echo "Waiting for Redis to be available..."
        sleep 5
    done
    
    echo "Redis is up and running!"
}


check_cache_is_up
# check_vault_is_running