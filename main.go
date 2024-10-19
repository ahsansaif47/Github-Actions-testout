package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	RedisC *redis.Client
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

	mux := setupRouter()

	println("Mux initialized sucessfully!!")
	_ = mux
	// log.Fatalln(http.ListenAndServe(":8080", mux))
}
