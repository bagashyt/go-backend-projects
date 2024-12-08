package internal

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-backend-projects/blogging_platform_api/utils"
)

type Handler struct {
	store BlogStore
}

func NewHandler(store BlogStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) BlogRoutes(router *gin.RouterGroup) {
	blogs := router.Group("/blogs")

	blogs.Use()
	{
		blogs.GET("/", h.handleGetBlogs)
		blogs.GET("/:id", h.handleGetBlogById)
		blogs.POST("/", h.handleCreateBlog)
	}

}

func (h *Handler) handleGetBlogs(c *gin.Context) {

	blogs, err := h.store.GetBlogs()
	if err != nil {

		utils.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(c, http.StatusOK, blogs)
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

func (h *Handler) handleCreateBlog(c *gin.Context) {
	var blog BlogPost
	if err := utils.ParseJSON(c, &blog); err != nil {
		log.Fatalf(err.Error())
		utils.WriteError(c, http.StatusBadRequest, err)
		return
	}

	err := h.store.CreateBlog(blog)
	if err != nil {
		log.Fatalf(err.Error())
		utils.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(c, http.StatusCreated, blog)
}
