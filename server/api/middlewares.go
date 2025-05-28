package api

import (
	"strings"

	"github.com/cv711/odin-takehome/server/internal"
	"github.com/gin-gonic/gin"
)

func (a *API) authRoute(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(401, gin.H{
			"error": "unauthorized",
		})
		c.Abort()
		return
	}

	jwtToken := strings.TrimPrefix(tokenString, "Bearer ")
	if jwtToken == "" {
		c.JSON(401, gin.H{
			"error": "unauthorized",
		})
		c.Abort()
		return
	}
	// Validate the JWT token
	claims, err := internal.ValidateJWTToken(jwtToken)
	if err != nil {
		a.log.Debug("Failed to get validate JWT: " + err.Error())
		c.JSON(401, gin.H{
			"error": "unauthorized",
		})
		c.Abort()
		return
	}

	// Check if the token issuer is valid
	if claims.Issuer != internal.TokenIssuer {
		c.JSON(401, gin.H{
			"error": "unauthorized",
		})
		c.Abort()
		return
	}

	// Store the user ID in the context
	c.Set("current_user_id", claims.Subject)
}
