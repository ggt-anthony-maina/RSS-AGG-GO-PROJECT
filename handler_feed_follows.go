package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ggt-anthony-maina/rssagg/internal/database"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func(apiCfg *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User){
	 type parameters struct{
		FeedID uuid.UUID `json:"feed_id"`
}

	 decoder := json.NewDecoder(r.Body)

	 params := parameters{}
	 err := decoder.Decode(&params)
	 if err != nil{
		 respondWithError(w, 400, fmt.Sprintf("Error parsing JSON %v", err))
		 return
	 }

	 feedFellow, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID: user.ID,
		FeedID: params.FeedID,
 })

	 if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't Create feed follow: %v", err))
		return
	 }

	respondWithJSON(w, 200, databaseFeedFollowToFeedFollow(feedFellow))
}

//get all feed follows
func(apiCfg *apiConfig) handlerGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User){
	feedFollows, err := apiCfg.DB.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
	   respondWithError(w, 400, fmt.Sprintf("Couldn't Get feeds: %v", err))
	   return
	}

   respondWithJSON(w, 200, databaseFeedFollowTsoFeedFollows(feedFollows))
}


///function to delete feed_follow
func(apiCfg *apiConfig) handlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User){
	 feedFollowIDStr := chi.URLParam(r, "feedFollowID")
	 feedFollowID , err := uuid.Parse(feedFollowIDStr)
	 if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't parse feed follow id "))
		return
	 }
	 err = apiCfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID: feedFollowID,
		UserID: user.ID,
	 })

	 if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couln't delete feed follow: %v ", err))
		return
	 }
	 respondWithJSON(w, 200, struct{}{})
}