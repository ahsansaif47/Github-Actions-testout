#!/bin/bash

# Set Exit on fail param!!
set -x
echo "Preloading data into services!!"


function install_vault_cli() {
    curl -L https://releases.hashicorp.com/vault/1.15.0/vault_1.15.0_linux_amd64.zip -o vault.zip
    unzip vault.zip
    sudo mv vault /usr/local/bin/
    vault --version
}

# This function is for populating the vault store.
# Not the most ideal setup, but done is far better than perfect!!
function preload_vault_secrets() {
    
    # Enable KV secrets engine at the specified path
    vault secrets enable -path=super-secret kv || echo "KV already enabled, skipping..."
    
    # Store DB_URL
    if vault kv put -mount=super-secret foo super-secret=$1; then
        echo "Successfully stored super-secret"
    else
        echo "Failed to store super-secret"
    fi
    
    
    
    vault kv list super-secret
}

install_vault_cli
preload_vault_secrets