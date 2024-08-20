package repository

import (
	"context"
	"database/sql"
	"log"
	"testing"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/mrofisr/azure-devops/internal/model"
	"github.com/testcontainers/testcontainers-go/modules/mssql"
)

func TestPostRepository_Create(t *testing.T) {
	ctx := context.Background()
	dbPassword := "yourStrongPassword1234$"

	container, err := mssql.Run(ctx,
		"mcr.microsoft.com/mssql/server:2022-latest",
		mssql.WithAcceptEULA(),
		mssql.WithPassword(dbPassword),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}
	dbURL, err := container.ConnectionString(ctx, "trustServerCertificate=true", "encrypt=true")
	if err != nil {
		log.Fatalf("failed to get connection string: %s", err)
	}
	conn, err := sql.Open("mssql", dbURL)
	if err != nil {
		log.Fatalf("failed to open connection: %s", err)
	}
	defer conn.Close()
	p := NewPostRepository(conn)
	postToCreate := model.Post{
		Title:   "Test Title",
		Content: "Test Content",
	}
	t.Run("create-table", func(t *testing.T) {
		query := `
			IF NOT EXISTS (SELECT * FROM sysobjects WHERE name='tblPosts' AND xtype='U')
			BEGIN
				CREATE TABLE dbo.tblPosts (
					id INT IDENTITY(1,1) PRIMARY KEY,
					title VARCHAR(255) NOT NULL,
					content TEXT NOT NULL,
					created_at DATETIME NOT NULL DEFAULT GETDATE(),
					updated_at DATETIME NOT NULL DEFAULT GETDATE()
				);
			END
		`
		_, err := conn.Exec(query)
		if err != nil {
			log.Fatal(err)
		}
	})
	t.Run("insert-post", func(t *testing.T) {
		err := p.Create(ctx, postToCreate)
		if err != nil {
			t.Errorf("failed to create post: %s", err)
		}
	})
	t.Run("find-all-post", func(t *testing.T) {
		posts, err := p.FindAll(ctx)
		if err != nil {
			t.Errorf("failed to get all posts: %s", err)
		}
		if len(posts) != 1 {
			t.Errorf("expected 1 post, got %d", len(posts))
		}
	})
	t.Run("find-post-by-id", func(t *testing.T) {
		post, err := p.FindByID(ctx, 1)
		if err != nil {
			t.Errorf("failed to get post by id: %s", err)
		}
		if post.ID != 1 {
			t.Errorf("expected post id 1, got %d", post.ID)
		}
	})
	t.Run("update-post", func(t *testing.T) {
		err := p.Update(ctx, model.Post{
			ID:      1,
			Title:   "Updated Title",
			Content: "Updated Content",
		})
		if err != nil {
			t.Errorf("failed to update post: %s", err)
		}
		post, err := p.FindByID(ctx, 1)
		if err != nil {
			t.Errorf("failed to get post by id: %s", err)
		}
		if post.Title != "Updated Title" {
			t.Errorf("expected post title Updated Title, got %s", post.Title)
		}
		if post.Content != "Updated Content" {
			t.Errorf("expected post content Updated Content, got %s", post.Content)
		}
	})
	t.Run("delete-post", func(t *testing.T) {
		err := p.Delete(ctx, 1)
		if err != nil {
			t.Errorf("failed to delete post: %s", err)
		}
		posts, err := p.FindAll(ctx)
		if err != nil {
			t.Errorf("failed to get all posts: %s", err)
		}
		if len(posts) != 0 {
			t.Errorf("expected 0 post, got %d", len(posts))
		}
	})
	t.Run("count-post", func(t *testing.T) {
		count, err := p.Count(ctx)
		if err != nil {
			t.Errorf("failed to count post: %s", err)
		}
		if count != 0 {
			t.Errorf("expected 0 post, got %d", count)
		}
	})
	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	})
}
