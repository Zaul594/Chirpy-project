package main

import (
	"log"
	"net/http"

	"github.com/bootdotdev/courses/learn-web-servers/create_read/internal/database"
	"github.com/go-chi/chi"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
}

func main() {
	const port = "8080"
	const filepathRoot = "."

	apiCfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
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
	ar.HandleFunc("/reset", apiCfg.resetHandler)
	ar.Post("/chirps", apiCfg.chirpCreateHandler)
	ar.Get("chirps", apiCfg.chirpCreateHandler)
	r.Mount("/api", apiCfg.chirpGetHandaler)

	corsMux := middlewareCors(r)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
