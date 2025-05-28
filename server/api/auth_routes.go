package api

import (
	"net/http"
	"net/netip"
	"strings"

	"github.com/cv711/odin-takehome/server/db"
	"github.com/cv711/odin-takehome/server/internal"
	"github.com/gin-gonic/gin"
)

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *API) auth(c *gin.Context) {
	var authRequest AuthRequest
	if err := c.BindJSON(&authRequest); err != nil {
		c.JSON(400, gin.H{
			"error": "invalid request",
		})
		return
	}
	if authRequest.Email == "" || authRequest.Password == "" {
		c.JSON(401, gin.H{
			"error": "unauthorized",
		})
		return
	}

	a.rateLimiter(c, authRequest)

	dbUser, err := a.db.LookupUser(c.Request.Context(), authRequest.Email)
	if err != nil {
		c.JSON(401, gin.H{
			"error": "unauthorized",
		})
		return
	}

	verified := internal.VerifyPassword(dbUser.PasswordHash, authRequest.Password)
	if !verified {
		c.JSON(401, gin.H{
			"error": "unauthorized",
		})
		return
	}

	c.SetSameSite(http.SameSiteStrictMode)

	dbUserID, err := dbUser.ID.Value()
	if err != nil {
		c.JSON(500, gin.H{
			"error": "something went wrong",
		})
		return
	}

	jwtToken, err := internal.GenerateJWTToken(dbUserID.(string))
	if err != nil {
		c.JSON(500, gin.H{
			"error": "something went wrong",
		})
		return
	}
	c.JSON(200, gin.H{
		"token": jwtToken,
	})
}

func (a *API) rateLimiter(c *gin.Context, authRequest AuthRequest) {
	remoteIP := strings.Split(c.ClientIP(), ":")[0]

	// Validate the IP address format
	parsedIP, err := netip.ParseAddr(remoteIP)
	if err != nil {
		a.log.Warn("Invalid IP address format", "ip", remoteIP, "error", err)
		c.JSON(400, gin.H{
			"error": "invalid IP address",
		})
		return
	}

	row, err := a.db.GetCounts(c.Request.Context(), db.GetCountsParams{
		RemoteIp: parsedIP.AsSlice(),
		Email:    authRequest.Email,
	})
	if err != nil {
		a.log.Error("Failed to get counts for rate limiting", "error", err)
		c.JSON(500, gin.H{
			"error": "internal server error",
		})
		return
	}

	if row.GlobalCount > 300 {
		c.JSON(429, gin.H{
			"error": "Too many attempts globally",
		})
		c.Abort()
		return
	}

	if row.IpCount > 25 {
		c.JSON(429, gin.H{
			"error": "Too many attempts from this IP",
		})
		c.Abort()
		return
	}

	if row.EmailCount > 10 {
		c.JSON(429, gin.H{
			"error": "Too many attempts for this email",
		})
		c.Abort()
		return
	}

	_, err = a.db.LogAttempt(c.Request.Context(), db.LogAttemptParams{
		Email:    authRequest.Email,
		RemoteIp: parsedIP.AsSlice(),
	})
	if err != nil {
		a.log.Error("Failed to log auth attempt", "error", err)
		c.JSON(500, gin.H{
			"error": "something went wrong",
		})
	}
}
