package api

import (
	"encoding/json"
	"net/http"

	"github.com/maxmindlin/swarm/data/mongodb"

	"github.com/maxmindlin/swarm/model"
	"gopkg.in/mgo.v2/bson"
)

const (
	collection = "articles"
	db         = "swarm"
)

// CreateArticleEndpoint creates a DB entry for an article object
func CreateArticleEndpoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// Get data from json payload and decode to article model
	var article model.Article
	if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Insert into DB
	s, err := mongodb.NewSession()
	defer s.Close()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	article.ID = bson.NewObjectId()
	err = s.Insert(db, collection, article)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, article)
}
