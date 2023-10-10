package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/Zaul594/Chirspy-project/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
	jwtSecret      string
	api_key        string
}

func main() {
	godotenv.Load(".env")

	const port = "8080"
	const filepathRoot = "."

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}
	api_key := os.Getenv("API_KEY")
	if api_key == "" {
		log.Fatal("API_KEY environment variable is not set")
	}

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parsed()
	if dbg != nil && *dbg {
		err := db.ResetDB()
		if err != nil {
			log.Fatal(err)
		}
	}

	apiCfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
		jwtSecret:      jwtSecret,
		api_key:        api_key,
	}

	r := chi.NewRouter()
	r.Handle("/app/assets/", apiCfg.metrIncMiddleWare(http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot)))))
	r.Handle("/app", apiCfg.metrIncMiddleWare(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))

	adm := chi.NewRouter()
	adm.Get("/metrics", apiCfg.admMetricHandler)
	r.Mount("/admin", adm)

	//api router
	ar := chi.NewRouter()
	ar.Get("/healthz", readyHanler)
	ar.Get("/metrics", apiCfg.metrHandler)
	ar.Get("/reset", apiCfg.resetHandler)

	ar.Post("/chirps", apiCfg.chirpCreateHandler)
	ar.Get("/chirps", apiCfg.chirpGetHandaler)
	ar.Get("/chirps/{chirpID}", apiCfg.GetOneChirpHandler)
	ar.Delete("/chirps/{chirpID}", apiCfg.deleteChirpyHandler)

	ar.Post("/login", apiCfg.handlerLogin)
	ar.Post("/users", apiCfg.createUserHandler)
	ar.Put("/users", apiCfg.userUpdateHandler)
	ar.Post("/polka/webhooks", apiCfg.upgradeHandler)

	ar.Post("/refresh", apiCfg.handlerRefresh)
	ar.Post("/revoke", apiCfg.handlerRevoke)
	r.Mount("/api", ar)

	corsMux := middlewareCors(r)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
