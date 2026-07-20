package main

import (
	"encoding/json"
	"log"
	"net/http"

	"osi-visualizer/osi"
)

func withCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		h(w, r)
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func handleLayers(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, osi.Layers)
}

func handleEncapsulate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "POST only"})
		return
	}
	var req osi.Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	if req.Message == "" {
		req.Message = "Hello"
	}
	writeJSON(w, http.StatusOK, map[string]any{"steps": osi.Encapsulate(req)})
}

func main() {
	http.HandleFunc("/api/layers", withCORS(handleLayers))
	http.HandleFunc("/api/encapsulate", withCORS(handleEncapsulate))

	addr := ":8080"
	log.Printf("OSI Visualizer backend listening on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
