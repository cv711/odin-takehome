package api

import (
	"github.com/cv711/odin-takehome/server/db"
	"github.com/cv711/odin-takehome/server/internal"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type signupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *API) signup(c *gin.Context) {
	var req signupRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": "bad request",
		})
		return
	}

	if req.Email == "" || req.Password == "" {
		c.JSON(400, gin.H{
			"error": "bad request",
		})
		return
	}

	passworHash, err := internal.HashPassword(req.Password)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "internal server error",
		})
		return
	}

	// Check if the email already exists
	if user, err := a.db.LookupUser(c.Request.Context(), req.Email); err == nil && user.ID.Valid {
		c.JSON(400, gin.H{
			"error": "email already exists",
		})
		return
	}

	// Create the user
	if _, err := a.db.CreateUser(c.Request.Context(), db.CreateUserParams{
		Email:        req.Email,
		PasswordHash: passworHash,
	}); err != nil {
		a.log.Error("Failed to create user", "error", err)
		c.JSON(500, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "ok",
	})
}

func (a *API) getUser(c *gin.Context) {
	userID, err := a.getUserIDFromContext(c)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "bad request",
		})
		return
	}

	user, err := a.db.GetUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(401, gin.H{
			"error": "unauthorized",
		})
		return
	}

	id, err := user.ID.Value()
	if err != nil {
		c.JSON(500, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(200, gin.H{
		"id":         id,
		"email":      user.Email,
		"created_at": user.CreatedAt,
	})
}

func (a *API) getUserIDFromContext(c *gin.Context) (pgtype.UUID, error) {
	var userID pgtype.UUID
	if err := userID.Scan(c.GetString("current_user_id")); err != nil {
		return userID, err
	}
	return userID, nil
}
