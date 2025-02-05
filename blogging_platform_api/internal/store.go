package internal

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {

	return &Store{db: db}
}

func (s *Store) CreateBlog(blog BlogPost) error {
	_, err := s.db.Exec(context.TODO(), "INSERT INTO blogs (title, content, category, tags, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)", blog.Title, blog.Content, blog.Category, blog.Tags, blog.CreatedAt, blog.UpdatedAt)
	if err != nil {
		log.Fatalf("Unable to create Blog error: %s", err)
		return err
	}

	return nil
}

func (s *Store) GetBlogById(id int) (*BlogPost, error) {
	row, err := s.db.Query(context.TODO(), "SELECT * FROM blogs WHERE blog_id = $1", id)
	if err != nil {
		log.Fatalf("Unable to Get Blog by ID, error: %s", err)
		return nil, err
	}

	b := new(BlogPost)
	for row.Next() {
		b, err = scanRowIntoBlog(row)
		if err != nil {
			return nil, err
		}
	}

	return b, nil

}

func (s *Store) GetBlogs() ([]*BlogPost, error) {
	query := `SELECT * FROM blogs`
	rows, err := s.db.Query(context.TODO(), query)
	if err != nil {
		log.Fatalf("Unable to Get Blogs error: %s", err)
		return nil, err
	}

	blogs := make([]*BlogPost, 0)
	for rows.Next() {
		p, err := scanRowIntoBlog(rows)
		if err != nil {
			return nil, err
		}

		blogs = append(blogs, p)
	}

	return blogs, nil
}

func scanRowIntoBlog(rows pgx.Rows) (*BlogPost, error) {
	blog := new(BlogPost)

	err := rows.Scan(
		&blog.ID,
		&blog.Title,
		&blog.Content,
		&blog.Category,
		&blog.Tags,
		&blog.CreatedAt,
		&blog.UpdatedAt,
	)
	if err != nil {
		log.Fatalf(err.Error())
		return nil, err
	}

	return blog, nil
}
