package main

import (
	"io"
	"log"
	"net/http"

	"github.com/gomodule/redigo/redis"
)

func main() {
	http.HandleFunc("/pokemonwithoutredis", getPokemonWithOutRedis)
	http.HandleFunc("/pokemonwithredis", getPokemonWithRedis)

	log.Println("Server is running")
	log.Fatal(http.ListenAndServe(":4000", nil))
}

func getPokemonWithOutRedis(w http.ResponseWriter, r *http.Request) {
	pokemonName := r.URL.Query().Get("pokemon")

	client := http.DefaultClient

	req, err := http.NewRequest(http.MethodGet, "https://pokeapi.co/api/v2/pokemon/"+pokemonName, nil)
	if err != nil {
		log.Panic(err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Panic(err)
	}

	bytes, _ := io.ReadAll(res.Body)
	w.Header().Add("Content-Type", "application/json")
	w.Write(bytes)
}

func getPokemonWithRedis(w http.ResponseWriter, r *http.Request) {
	pokemonName := r.URL.Query().Get("pokemon")

	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Panic(err)
	}

	// Check Data From Redis (is't already caching?)
	reply, err := redis.Bytes(conn.Do("GET", pokemonName))
	if err == nil {
		// it's ready cache
		w.Header().Add("Content-Type", "application/json")
		w.Write(reply)
		return
	}

	// Create Request API pokeapi
	client := http.DefaultClient
	req, err := http.NewRequest(http.MethodGet, "https://pokeapi.co/api/v2/pokemon/"+pokemonName, nil)
	if err != nil {
		log.Panic(err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Panic(err)
	}

	bytes, _ := io.ReadAll(res.Body)

	// Save Response Data to Redis
	_, err = conn.Do("SET", pokemonName, string(bytes))
	if err != nil {
		log.Panic(err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(bytes)
}

// Testing Speed
// curl -o /dev/null -w "\n%{time_total} seconds\n" http://localhost:4000/pokemonwithoutredis?pokemon=pikachu
// curl -o /dev/null -w "\n%{time_total} seconds\n" http://localhost:4000/pokemonwithredis?pokemon=pikachu
