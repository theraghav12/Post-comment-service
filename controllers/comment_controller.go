package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"post-comments-api/models"
	"post-comments-api/utils"
)

type CreateCommentRequest struct {
	PostID  uint    `json:"post_id"`
	Content string  `json:"content" binding:"required"`
	Author  *string `json:"author,omitempty"`
}

type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

func CreateComment(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}
	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var postID uint
	if idStr := c.Param("id"); idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
			return
		}
		postID = uint(id)
	} else if req.PostID > 0 {
		postID = req.PostID
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Post ID is required"})
		return
	}
	var post models.Post
	if err := utils.GetDB().First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	uid := userID.(uint)
	comment := models.Comment{
		PostID:  postID,
		UserID:  &uid,
		Author:  req.Author,
		Content: req.Content,
	}
	if err := utils.GetDB().Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}
	c.JSON(http.StatusCreated, comment)
}

func CreateCommentPublic(c *gin.Context) {
	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.PostID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Post ID is required"})
		return
	}
	var post models.Post
	if err := utils.GetDB().First(&post, req.PostID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	comment := models.Comment{
		PostID:  req.PostID,
		Author:  req.Author,
		Content: req.Content,
	}
	if err := utils.GetDB().Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}
	c.JSON(http.StatusCreated, comment)
}

func GetComments(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}
	// Pagination params
	page := 1
	pageSize := 10
	maxPageSize := 50
	if p := c.Query("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if ps := c.Query("page_size"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil && v > 0 {
			if v > maxPageSize {
				pageSize = maxPageSize
			} else {
				pageSize = v
			}
		}
	}
	offset := (page - 1) * pageSize
	var comments []models.Comment
	var total int64
	utils.GetDB().Model(&models.Comment{}).Where("post_id = ?", postID).Count(&total)
	if err := utils.GetDB().Where("post_id = ?", postID).Limit(pageSize).Offset(offset).Order("created_at ASC").Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}
	var resp []gin.H
	for _, comment := range comments {
		htmlContent, _ := utils.RenderMarkdown(comment.Content)
		resp = append(resp, gin.H{
			"id": comment.ID,
			"post_id": comment.PostID,
			"user_id": comment.UserID,
			"author": comment.Author,
			"content": comment.Content,
			"html_content": htmlContent,
			"created_at": comment.CreatedAt,
			"updated_at": comment.UpdatedAt,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"comments": resp,
		"pagination": gin.H{
			"page": page,
			"page_size": pageSize,
			"total": total,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

func UpdateComment(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}
	var comment models.Comment
	if err := utils.GetDB().First(&comment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}
	if comment.UserID == nil || *comment.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to update this comment"})
		return
	}
	var req UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	comment.Content = req.Content
	if err := utils.GetDB().Save(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment"})
		return
	}
	c.JSON(http.StatusOK, comment)
}

func DeleteComment(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}
	var comment models.Comment
	if err := utils.GetDB().First(&comment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}
	if comment.UserID == nil || *comment.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this comment"})
		return
	}
	if err := utils.GetDB().Delete(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted"})
}
