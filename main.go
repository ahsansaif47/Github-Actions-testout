package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hashicorp/vault/api"
)

var (
	RedisC  *redis.Client
	VaultC  *api.Client
	VaultKV *api.KVv2
)

func Handler_foo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("foo"))
}

func Handler_bar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("bar"))
}

func connectToRedis() error {
	ctx := context.Background()
	// Set up Redis options, adjust according to your Redis configuration

	redis_addr := os.Getenv("REDIS_ADDR")

	fmt.Println("The redis address is: ", redis_addr)

	rdb := redis.NewClient(&redis.Options{
		Addr:     redis_addr, // Redis server address
		Password: "",         // Password, if any
		DB:       0,          // Default DB
	})

	// Ping the Redis server
	status := rdb.Ping(ctx)
	err := status.Err()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %v", err)
	}

	fmt.Println("Successfully connected to Redis")
	RedisC = rdb
	return nil
}

// Connect to Vault
func connectToVault() error {

	vault_addr := os.Getenv("VAULT_ADDR")
	vault_token := os.Getenv("VAULT_TOKEN")

	fmt.Printf("Vault Address: %s, Vault Token: %s", vault_addr, vault_token)

	config := api.DefaultConfig()
	config.Address = vault_addr

	client, err := api.NewClient(config)
	if err != nil {
		return fmt.Errorf("failed to create Vault client: %v", err)
	}

	client.SetToken(vault_token)

	VaultC = client

	fmt.Println("Successfully connected to Vault")
	return nil
}

// Add a key-value pair to Vault
func addSecretToVault(mount, key string, value map[string]interface{}) error {
	_, err := VaultC.KVv2(mount).Put(context.Background(), key, value)
	if err != nil {
		return fmt.Errorf("failed to add secret to Vault: %s", err.Error())
	}

	fmt.Printf("Successfully added secret '%s'\n", key)
	return nil
}

// Get a key-value pair from Vault
func getSecretFromVault(mount, key string) (map[string]interface{}, error) {
	secret, err := VaultC.KVv2(mount).Get(context.Background(), key)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve secret: %v", err)
	}

	fmt.Printf("Successfully retrieved secret '%v' for '%s'\n", secret.Data, key)
	return secret.Data, nil
}

// Add a key-value pair to Redis
func addKeyValue(key string, value string, expiration time.Duration) error {
	ctx := context.Background()

	err := RedisC.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set key-value pair: %v", err)
	}

	fmt.Printf("Successfully added key '%s' with value '%s'\n", key, value)
	return nil
}

// Delete a key-value pair from Redis
func deleteKeyValue(key string) error {
	ctx := context.Background()
	err := RedisC.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete key '%s': %v", key, err)
	}
	fmt.Printf("Successfully deleted key '%s'\n", key)
	return nil
}

// Delete a key-value pair from Redis
func getKeyValue(key string) error {
	ctx := context.Background()
	val, err := RedisC.Get(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to get key '%s'", key)
	}
	fmt.Printf("Successfully retrieved key value pair '%s' : '%s'\n", key, val)
	return nil
}

// type Server struct {
// 	Mux *http.ServeMux
// }

// var Srv Server = Server{}

func setupRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/h1", Handler_foo)
	mux.HandleFunc("/h2", Handler_bar)
	return mux
}

func main() {

	err := connectToRedis()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Redis connection established successfully.")
	}

	addKeyValue("foo", "bar", time.Duration(0))

	getKeyValue("foo")

	deleteKeyValue("foo")

	// Connect to Vault
	err = connectToVault()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if err := addSecretToVault("secret", "foo", map[string]interface{}{"super-secret": "bar"}); err != nil {
		fmt.Println("Error:", err)
		return
	}

	if _, err := getSecretFromVault("secret", "foo"); err != nil {
		fmt.Println("Error: ", err.Error())
	}

	// Getting one of the vaules populated earlier!!
	if _, err := getSecretFromVault("super-secret", "foo"); err != nil {
		fmt.Println("Error: ", err.Error())
	}

	// mux := setupRouter()

	// println("Mux initialized sucessfully!!")
	// _ = mux
	// log.Fatalln(http.ListenAndServe(":8080", mux))
}
