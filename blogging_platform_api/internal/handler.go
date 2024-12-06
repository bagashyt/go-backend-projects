package internal

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-backend-projects/blogging_platform_api/utils"
)

type Handler struct {
	store BlogStore
}

var (
	w http.ResponseWriter
	r *http.Request
)

func NewHandler(store BlogStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) BlogRoutes(router *gin.RouterGroup) {
	blogs := router.Group("/blogs")

	blogs.Use()
	{
		blogs.GET("/", h.handleGetBlogs)
		blogs.GET("/:id", h.handleGetBlogById)
	}

}

func (h *Handler) handleGetBlogs(c *gin.Context) {

	blogs, err := h.store.GetBlogs()
	if err != nil {

		utils.WriteError(c, http.StatusInternalServerError, err)
		// c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.WriteJSON(c, http.StatusOK, blogs)
	// c.JSON(http.StatusOK, blogs)
	return
}

func (h *Handler) handleGetBlogById(c *gin.Context) {
	vars := c.Param("id")

	BlogID, err := strconv.Atoi(vars)
	if err != nil {
		utils.WriteError(c, http.StatusBadRequest, fmt.Errorf("invalid blog ID"))
		return
	}

	blog, err := h.store.GetBlogById(BlogID)
	if err != nil {
		utils.WriteError(c, http.StatusInternalServerError, err)
		return
	}
	fmt.Printf("blog.Content: %v\n", blog.Content)
	utils.WriteJSON(c, http.StatusOK, blog)
}

// func (h *Handler) handleCreateBlog(w http.ResponseWriter, r *http.Request) {
// 	var blog BlogPost
// 	if err := utils.ParseJSON(r, &blog); err != nil {
// 		// utils.WriteError(w, http.StatusBadRequest, err)
// 		return
// 	}

// 	err := h.store.CreateBlog(blog)
// 	if err != nil {
// 		// utils.WriteError(w, http.StatusInternalServerError, err)
// 		return
// 	}

// 	//utils.WriteJSON(w, http.StatusCreated, blog)
// }
