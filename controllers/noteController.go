package controllers

import (
	"net/http"

	"go-note-editor/initializers"
	"go-note-editor/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateNote(c *gin.Context) {
	var req models.CreateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	note := models.Note{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID,
	}

	if err := initializers.DB.Create(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create note"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"note": note})
}

func GetAllNotes(c *gin.Context) {
	userIDStr, _ := c.Get("userID")
	userID, _ := uuid.Parse(userIDStr.(string))

	var notes []models.Note
	if err := initializers.DB.Where("user_id = ?", userID).Find(&notes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"notes": notes, "count": len(notes)})
}

func GetNoteByID(c *gin.Context) {
	id := c.Param("id")
	userIDStr, _ := c.Get("userID")
	userID, _ := uuid.Parse(userIDStr.(string))

	var note models.Note
	if err := initializers.DB.Where("id = ? AND user_id = ?", id, userID).First(&note).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"note": note})
}

func UpdateNote(c *gin.Context) {
	id := c.Param("id")
	userIDStr, _ := c.Get("userID")
	userID, _ := uuid.Parse(userIDStr.(string))

	var note models.Note
	if err := initializers.DB.Where("id = ? AND user_id = ?", id, userID).First(&note).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	var req models.CreateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	note.Title = req.Title
	note.Content = req.Content

	if err := initializers.DB.Save(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update note"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"note": note})
}

func DeleteNote(c *gin.Context) {
	id := c.Param("id")
	userIDStr, _ := c.Get("userID")
	userID, _ := uuid.Parse(userIDStr.(string))

	var note models.Note
	if err := initializers.DB.Where("id = ? AND user_id = ?", id, userID).First(&note).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	if err := initializers.DB.Delete(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete note"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note deleted successfully"})
}
