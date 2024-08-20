package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/mrofisr/azure-devops/internal/handler"
	"github.com/mrofisr/azure-devops/internal/repository"
	"github.com/mrofisr/azure-devops/internal/router"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	ctx := context.Background()
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;trustServerCertificate=true;encrypt=true", os.Getenv("DATABASE_HOST"), os.Getenv("DATABASE_USER"), os.Getenv("DATABASE_PASS"), os.Getenv("DATABASE_PORT"), os.Getenv("DATABASE_NAME"))
	conn, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatal("Error creating connection pool: " + err.Error())
	}
	defer conn.Close()
	postRepo := repository.NewPostRepository(conn)
	postHandler := handler.NewPostHandler(postRepo)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "hello world"}`))
		w.WriteHeader(http.StatusOK)
	})
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		if err := conn.PingContext(ctx); err != nil {
			log.Fatal(err.Error())
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "pong"}`))
		w.WriteHeader(http.StatusOK)
	})
	r.Mount("/post", router.PostRouter(postHandler))
	fmt.Println("Server is running on port http://localhost:8080")
	fmt.Println("Press CTRL + C to exit")
	log.Fatal(http.ListenAndServe(":8080", r))
}
