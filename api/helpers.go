package api

import (
	"encoding/json"
	"net/http"

	"github.com/maxmindlin/swarm/data/mongodb"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, map[string]string{"error": msg})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// GenericCreateEndpoint is a generic shell for writing documents to MongoDB
func GenericCreateEndpoint(w http.ResponseWriter, doc interface{}, dbName string, collName string) {
	// Insert into DB
	s, err := mongodb.NewSession()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer s.Close()
	err = s.Insert(dbName, collName, doc)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, doc)
}
