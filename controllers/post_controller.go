package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"post-comments-api/models"
	"post-comments-api/utils"
)

type CreatePostRequest struct {
	Title   string  `json:"title" binding:"required"`
	Content string  `json:"content" binding:"required"`
	Author  *string `json:"author,omitempty"`
}

type UpdatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func CreatePost(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid := userID.(uint)
	post := models.Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  &uid,
		Author:  req.Author,
	}
	if err := utils.GetDB().Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}
	c.JSON(http.StatusCreated, post)
}

func CreatePostPublic(c *gin.Context) {
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	post := models.Post{
		Title:   req.Title,
		Content: req.Content,
		Author:  req.Author,
	}
	if err := utils.GetDB().Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}
	c.JSON(http.StatusCreated, post)
}

func GetPosts(c *gin.Context) {
	var posts []models.Post
	if err := utils.GetDB().Preload("Comments").Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}
	var resp []gin.H
	for _, post := range posts {
		htmlContent, _ := utils.RenderMarkdown(post.Content)
		resp = append(resp, gin.H{
			"id": post.ID,
			"user_id": post.UserID,
			"author": post.Author,
			"title": post.Title,
			"content": post.Content,
			"html_content": htmlContent,
			"created_at": post.CreatedAt,
			"updated_at": post.UpdatedAt,
			"comments": post.Comments,
		})
	}
	c.JSON(http.StatusOK, resp)
}

func GetPost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}
	var post models.Post
	if err := utils.GetDB().Preload("Comments").First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	htmlContent, _ := utils.RenderMarkdown(post.Content)
	c.JSON(http.StatusOK, gin.H{
		"id": post.ID,
		"user_id": post.UserID,
		"author": post.Author,
		"title": post.Title,
		"content": post.Content,
		"html_content": htmlContent,
		"created_at": post.CreatedAt,
		"updated_at": post.UpdatedAt,
		"comments": post.Comments,
	})
}

func UpdatePost(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}
	var post models.Post
	if err := utils.GetDB().First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	if post.UserID == nil || *post.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to update this post"})
		return
	}
	var req UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Content != "" {
		post.Content = req.Content
	}
	if err := utils.GetDB().Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}
	c.JSON(http.StatusOK, post)
}

func DeletePost(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}
	var post models.Post
	if err := utils.GetDB().First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	if post.UserID == nil || *post.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this post"})
		return
	}
	if err := utils.GetDB().Delete(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Post deleted"})
}
