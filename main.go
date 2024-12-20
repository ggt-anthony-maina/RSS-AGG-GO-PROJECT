package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ggt-anthony-maina/rssagg/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)
type apiConfig struct{
   DB *database.Queries
}
func main(){
   feed , err :=urlToFeed("https://wagslane.dev/index.xml")
   if err != nil{
      log.Fatal(err)
   }
   fmt.Println(feed)
   fmt.Println("Hello world")

   godotenv.Load(".env")

   portString := os.Getenv("PORT")
   if portString == ""{
	log.Fatal("PORT is not found in the environment")
   }

   dbUrl := os.Getenv("DB_URL")
   if dbUrl == ""{
	log.Fatal("DB_URL is not found in the environment")
   }

   conn, err := sql.Open("postgres", dbUrl)
   if err != nil{
      log.Fatal("Can't connect to the database:", err)
   }

   db := database.New(conn)
   apiCfg := apiConfig{
      DB: db,
   }

   go startScraping(db, 10, time.Minute)
   router := chi.NewRouter()

  router.Use(cors.Handler(cors.Options{
    AllowedOrigins:   []string{"https://*", "http://*"},
    AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
    ExposedHeaders:   []string{"Link"},
    AllowCredentials: false,
    MaxAge:           300, // Maximum value not ignored by any of major browsers
  }))

  v1Router := chi.NewRouter()

  v1Router.Get("/healthz", handlerReadiness)
  v1Router.Get("/err",handlerErr)
  v1Router.Post("/users",  apiCfg.handlerCreateUser)
  v1Router.Get("/users", apiCfg.middlewareAuth((apiCfg.handlerGetUser)))

  v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
  v1Router.Get("/feeds",apiCfg.middlewareAuth(apiCfg.handlerGetFeed) )
  v1Router.Get("/feeds/all", apiCfg.handlerGetFeeds)

  v1Router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))
  v1Router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollows))
  v1Router.Delete("/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollow))

  v1Router.Get("/posts", apiCfg.middlewareAuth(apiCfg.handlerGetPostsForUser))

  router.Mount("/v1", v1Router)
   srv := &http.Server{
	  Handler: router,
	  Addr: ":" + portString,
   }

   log.Printf("Server running on port %v", portString)
   err = srv.ListenAndServe()
   if err != nil{
	log.Fatal(err)
   }

   fmt.Println("PORT", portString)
}