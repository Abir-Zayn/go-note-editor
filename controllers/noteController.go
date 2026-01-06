package controllers

import (
	"go-note-editor/initializers"
	"go-note-editor/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateNote creates a new note
func CreateNote(c *gin.Context) {
	var note models.Note

	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	note.UserID = userID.(uint)

	if err := initializers.DB.Create(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create note"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"note": note})
}

// GetAllNotes retrieves all notes for the authenticated user
func GetAllNotes(c *gin.Context) {
	var notes []models.Note

	userID, _ := c.Get("userID")

	if err := initializers.DB.Where("user_id = ?", userID).Find(&notes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"notes": notes, "count": len(notes)})
}

// GetNoteByID retrieves a single note by ID
func GetNoteByID(c *gin.Context) {
	id := c.Param("id")
	var note models.Note

	userID, _ := c.Get("userID")

	if err := initializers.DB.Where("id = ? AND user_id = ?", id, userID).First(&note).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"note": note})
}

// UpdateNote updates an existing note
func UpdateNote(c *gin.Context) {
	id := c.Param("id")
	var note models.Note

	userID, _ := c.Get("userID")

	// Check if note exists and belongs to user
	if err := initializers.DB.Where("id = ? AND user_id = ?", id, userID).First(&note).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	var updateData models.Note
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields
	note.Title = updateData.Title
	note.Content = updateData.Content

	if err := initializers.DB.Save(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update note"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"note": note})
}

// DeleteNote deletes a note (soft delete)
func DeleteNote(c *gin.Context) {
	id := c.Param("id")
	var note models.Note

	userID, _ := c.Get("userID")

	// Check if note exists and belongs to user
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

// SearchNotes searches notes by title or content
func SearchNotes(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	var notes []models.Note
	userID, _ := c.Get("userID")

	searchPattern := "%" + query + "%"
	if err := initializers.DB.Where("user_id = ? AND (title ILIKE ? OR content ILIKE ?)",
		userID, searchPattern, searchPattern).Find(&notes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search notes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"notes": notes, "count": len(notes)})
}
