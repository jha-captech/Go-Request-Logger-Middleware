package main

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jha-captech/Go-Request-Logger-Middleware/middleware"
)

func main() {
	runStdLib()
	// runGin()
}

func runStdLib() {
	logger := slog.Default()

	mux := http.NewServeMux()

	stack := middleware.CreateMiddlewareStack(middleware.LoggerColorMiddleware(logger))

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]string{"message": "Hello, world!"}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	server := http.Server{
		Addr:    "localhost:8080",
		Handler: stack(mux),
	}

	println("listening...")
	log.Fatalln(server.ListenAndServe())
}

func runGin() {
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		println("in handler")
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	log.Fatalln(r.Run("localhost:8080")) // listen and serve on 0.0.0.0:8080
}
