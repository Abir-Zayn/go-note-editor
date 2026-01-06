package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"go-note-editor/models"

	"github.com/gin-gonic/gin"
)

func Signup(c *gin.Context) {
	var req models.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_ANON_KEY")

	// Prepare signup payload for Supabase Auth
	payload := map[string]interface{}{
		"email":    req.Email,
		"password": req.Password,
		"data": map[string]string{
			"full_name": req.FullName,
		},
	}

	jsonPayload, _ := json.Marshal(payload)

	// Call Supabase Auth signup endpoint
	url := fmt.Sprintf("%s/auth/v1/signup", supabaseURL)
	httpReq, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("apikey", supabaseKey)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to auth service"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		var errResp map[string]interface{}
		json.Unmarshal(body, &errResp)
		c.JSON(resp.StatusCode, errResp)
		return
	}

	var authResp models.AuthResponse
	json.Unmarshal(body, &authResp)

	c.JSON(http.StatusOK, gin.H{
		"message": "Signup successful",
		"user": gin.H{
			"id":    authResp.User.ID,
			"email": authResp.User.Email,
		},
		"access_token":  authResp.AccessToken,
		"refresh_token": authResp.RefreshToken,
	})
}
