package main

import (
	"blackjack-trainer/deck"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type ServerState struct {
	mu   sync.Mutex
	Shoe *deck.Shoe
}

var state = ServerState{}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/start", startHandler)
	mux.HandleFunc("GET /api/deal", dealHandler)
	mux.HandleFunc("GET /api/count", countHandler)

	handlerWithCors := enableCORS(mux)
	port := ":8080"
	fmt.Println("Server is running on port", port)
	if err := http.ListenAndServe(port, handlerWithCors); err != nil {
		fmt.Println("Error starting server:", err)
	}

}

func startHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	state.mu.Lock()
	defer state.mu.Unlock()

	state.Shoe = deck.NewShoe(6) // Create a new shoe with 6 decks
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "New shoe created"}`))
}

func dealHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	state.mu.Lock()
	defer state.mu.Unlock()

	if state.Shoe == nil {
		http.Error(w, "Shoe not initialized. Please start a new game.", http.StatusBadRequest)
		return
	}

	card, ok := state.Shoe.Deal()
	if !ok {
		http.Error(w, "No more cards in the shoe", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(card)
}

func countHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	state.mu.Lock()
	defer state.mu.Unlock()

	if state.Shoe == nil {
		http.Error(w, "Shoe not initialized. Please start a new game.", http.StatusBadRequest)
		return
	}

	count := state.Shoe.RunningCount
	json.NewEncoder(w).Encode(map[string]int{"running_count": count})
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "application/json") // Forces JSON responses global

		// Short-circuit execution if it is just a pre-flight browser check
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
