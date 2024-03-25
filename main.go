package main

import (
	"database/sql" // for connecting with the Postgre
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq" // import postgre drivers
)

// declare the global variable
var db *sql.DB

func main() {
	// Connect to PostgreSQL
	var err error
	db, err = sql.Open("postgres", "postgresql://username:password@localhost/dbname?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	// when main function exists db connection closed
	defer db.Close()
	fmt.Println("Connected to PostgreSQL")
	// Create a new router
	r := mux.NewRouter()
	// Define API endpoints
	r.HandleFunc("/deleteRecord", deleteRecordHandler).Methods("DELETE")

	fmt.Println("Server is listening on port 8080...")
	// Start the server
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
func deleteRecordHandler(w http.ResponseWriter, r *http.Request) {
	// Decode JSON payload
	var data map[string]string
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Extract ID from payload
	id, ok := data["id"]
	if !ok {
		http.Error(w, "ID not provided in payload", http.StatusBadRequest)
		return
	}
	// Delete record from PostgreSQL
	_, err = db.Exec("DELETE FROM your_table WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Send success response with HTTP status 200
	response := map[string]string{"message": "Deleted successfully"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
