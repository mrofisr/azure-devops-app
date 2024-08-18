package repository

import (
	"context"
	"database/sql"
	"log"

	"github.com/mrofisr/azure-devops/internal/model"
)

type PostRepository interface {
	FindAll(ctx context.Context) ([]model.Post, error)
	FindByID(ctx context.Context, id int) (model.Post, error)
	Create(ctx context.Context, post model.Post) error
	Update(ctx context.Context, post model.Post) error
	Delete(ctx context.Context, id int) error
	Count(ctx context.Context) (int, error)
}

type postRepository struct {
	conn *sql.DB
}

func (p *postRepository) FindAll(ctx context.Context) ([]model.Post, error) {
	query := "SELECT * FROM dbo.tblPosts"
	rows, err := p.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []model.Post
	for rows.Next() {
		var post model.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}
func (p *postRepository) FindByID(ctx context.Context, id int) (model.Post, error) {
	query := "SELECT * FROM dbo.tblPosts WHERE id = ?"
	row := p.conn.QueryRowContext(ctx, query, id)
	var post model.Post
	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return model.Post{}, err
	}
	return post, nil
}
func (p *postRepository) Create(ctx context.Context, post model.Post) error {
	query := "INSERT INTO dbo.tblPosts (title, content, created_at, updated_at) VALUES (?, ?, GETDATE(), GETDATE())"
	_, err := p.conn.ExecContext(ctx, query, post.Title, post.Content)
	if err != nil {
		return err
	}
	return nil
}
func (p *postRepository) Update(ctx context.Context, post model.Post) error {
	query := "UPDATE dbo.tblPosts SET title = ?, content = ?, updated_at = GETDATE() WHERE id = ?"
	_, err := p.conn.ExecContext(ctx, query, post.Title, post.Content, post.ID)
	if err != nil {
		return err
	}
	return nil
}
func (p *postRepository) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM dbo.tblPosts WHERE id = ?"
	_, err := p.conn.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
func (p *postRepository) Count(ctx context.Context) (int, error) {
	query := "SELECT COUNT(*) FROM dbo.tblPosts"
	row := p.conn.QueryRowContext(ctx, query)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func NewPostRepository(conn *sql.DB) PostRepository {
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
	return &postRepository{conn}
}
