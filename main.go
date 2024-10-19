package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
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
	var ctx = context.Background()
	// Set up Redis options, adjust according to your Redis configuration
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // Password, if any
		DB:       0,                // Default DB
	})

	// Ping the Redis server
	status := rdb.Ping(ctx)
	err := status.Err()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %v", err)
	}

	fmt.Println("Successfully connected to Redis")
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

	mux := setupRouter()

	log.Fatalln(http.ListenAndServe(":8080", mux))
}
