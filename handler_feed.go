package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ggt-anthony-maina/rssagg/internal/database"
	"github.com/google/uuid"
)

func(apiCfg *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User){
	 type parameters struct{
		Name string `json:"name"`
		URL string `json:"url`
}

	 decoder := json.NewDecoder(r.Body)

	 params := parameters{}
	 err := decoder.Decode(&params)
	 if err != nil{
		 respondWithError(w, 400, fmt.Sprintf("Error parsing JSON %v", err))
		 return
	 }

	 feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: params.Name,
		Url: params.URL,
		UserID: user.ID,
	 })

	 if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't Create feed: %v", err))
		return
	 }

	respondWithJSON(w, 200, databaseFeedToFeed(feed))
}

//get all feeds 
func(apiCfg *apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request){
	feeds, err := apiCfg.DB.GetFeeds(r.Context())
	if err != nil {
	   respondWithError(w, 400, fmt.Sprintf("Couldn't Get feeds: %v", err))
	   return
	}

   respondWithJSON(w, 200, databaseFeedsToFeeds(feeds))
}


///function to Get a user
func(apiCfg *apiConfig) handlerGetFeed(w http.ResponseWriter, r *http.Request, user database.User){
	
	respondWithJSON(w, 200, databaseUserToUser(user))
}