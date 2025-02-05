package internal

type BlogStore interface {
	CreateBlog(BlogPost) error
	GetBlogs() ([]*BlogPost, error)
	GetBlogById(id int) (*BlogPost, error)
}

type BlogPost struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Category  string `json:"category"`
	Tags      string `json:"tags"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
