package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
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

type DatabaseConfig struct {
	DatabaseHost string `json:"DATABASE_HOST"`
	DatabasePort string `json:"DATABASE_PORT"`
	DatabaseUser string `json:"DATABASE_USER"`
	DatabasePass string `json:"DATABASE_PASS"`
	DatabaseName string `json:"DATABASE_NAME"`
}

func (d *DatabaseConfig) GetConnectionString() string {
	// Get a secret. An empty string version gets the latest version of the secret.
	vaultURI := fmt.Sprintf("https://%s.vault.azure.net/", os.Getenv("KEY_VAULT_NAME"))

	// Create a credential using the NewDefaultAzureCredential type.
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}

	version := ""
	client, err := azsecrets.NewClient(vaultURI, cred, nil)
	if err != nil {
		log.Fatalf("failed to create the secret client: %v", err)
	}

	resp, err := client.GetSecret(context.TODO(), os.Getenv("SECRET_NAME"), version, nil)
	if err != nil {
		log.Fatalf("failed to get the secret: %v", err)
	}

	// Initialize a DatabaseConfig instance
	var config DatabaseConfig

	// Unmarshal the JSON string into the DatabaseConfig struct
	if err := json.Unmarshal([]byte(*resp.Value), &config); err != nil {
		log.Fatalf("failed to unmarshal secret value into DatabaseConfig: %v", err)
	}

	connectionString := fmt.Sprintf(
		"server=%s;user id=%s;password=%s;port=%s;database=%s;trustServerCertificate=true;encrypt=true",
		config.DatabaseHost,
		config.DatabaseUser,
		config.DatabasePass,
		config.DatabasePort,
		config.DatabaseName,
	)

	return connectionString
}
func main() {
	ctx := context.Background()
	var dbConfig DatabaseConfig
	connectionString := dbConfig.GetConnectionString()
	conn, err := sql.Open("mssql", connectionString)
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
