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
    vault secrets enable -version=2 -path=super-secret kv || echo "KV already enabled, skipping..."
    
    # # Store DB_URL
    if vault kv put -mount=super-secret foo super-secret=BAar; then
        echo "Successfully stored super-secret"
        vault kv get -mount=super-secret foo
    else
        echo "Failed to store super-secret" || exit
    fi
    
    # vault write super-secret/foo super-secret=bar
    vault kv get -mount=super-secret foo
    
    
    
    vault kv list super-secret
    
    echo "Listing secrets in detial :))"
    vault secrets list -detailed
    
}

install_vault_cli
preload_vault_secrets