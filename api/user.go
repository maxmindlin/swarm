package api

import (
	"encoding/json"
	"net/http"

	"github.com/maxmindlin/swarm/model"
	"gopkg.in/mgo.v2/bson"
)

// CreateUserEndpoint inserts a new user profile into db
func CreateUserEndpoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// Get json payload from request and decode it to a user model
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Insert into DB
	user.ID = bson.NewObjectId()
	GenericCreateEndpoint(w, user, db, "users")
}
