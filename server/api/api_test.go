package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"os"
	"strings"
	"testing"

	"github.com/cv711/odin-takehome/server/db"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)
	server := setupServer(t)
	router := gin.New()
	router = server.setupRoutes(router)
	cleanupDB(t, server.db)

	user1_email := "user1@example.com"
	user1_password := "password"

	t.Cleanup(func() { cleanupDB(t, server.db) })

	t.Run("health endpoint", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/health", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
	})

	t.Run("signup endpoint", func(t *testing.T) {
		w := httptest.NewRecorder()

		payload, _ := json.Marshal(signupRequest{
			Email:    user1_email,
			Password: user1_password,
		})
		req, _ := http.NewRequest("POST", "/api/signup", strings.NewReader(string(payload)))
		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
	})

	t.Run("signup endpoint should fail when password is missing", func(t *testing.T) {
		w := httptest.NewRecorder()

		payload, _ := json.Marshal(signupRequest{
			Email: user1_email,
		})
		req, _ := http.NewRequest("POST", "/api/signup", strings.NewReader(string(payload)))
		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
	})

	t.Run("signup endpoint should fail if user exists", func(t *testing.T) {
		w := httptest.NewRecorder()

		payload, _ := json.Marshal(signupRequest{
			Email:    user1_email,
			Password: user1_password,
		})
		req, _ := http.NewRequest("POST", "/api/signup", strings.NewReader(string(payload)))
		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
	})

	t.Run("auth endpoint", func(t *testing.T) {
		w := httptest.NewRecorder()
		payload, _ := json.Marshal(map[string]string{
			"email":    user1_email,
			"password": user1_password,
		})
		req, _ := http.NewRequest("POST", "/api/auth", strings.NewReader(string(payload)))
		req.RemoteAddr = "192.168.1.100:12345"
		router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		var response map[string]interface{}
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		assert.NotEmpty(t, response["token"], "Expected token in response")
	})

	t.Run("auth endpoint should fail with wrong password", func(t *testing.T) {
		w := httptest.NewRecorder()
		payload, _ := json.Marshal(map[string]string{
			"email":    user1_email,
			"password": "wrongpassword",
		})
		req, _ := http.NewRequest("POST", "/api/auth", strings.NewReader(string(payload)))
		req.RemoteAddr = "192.168.1.100:12345"
		router.ServeHTTP(w, req)
		assert.Equal(t, 401, w.Code)
		var response map[string]interface{}
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		assert.Equal(t, "unauthorized", response["error"], "Expected 'unauthorized' error message")
	})

	t.Run("get user endpoint", func(t *testing.T) {
		w := httptest.NewRecorder()
		// First, authenticate to get a token
		authPayload, _ := json.Marshal(map[string]string{
			"email":    user1_email,
			"password": user1_password,
		})
		authReq, _ := http.NewRequest("POST", "/api/auth", strings.NewReader(string(authPayload)))
		authReq.RemoteAddr = "192.168.1.100:12345"
		router.ServeHTTP(w, authReq)
		assert.Equal(t, 200, w.Code)
		var authResponse map[string]interface{}
		if err := json.NewDecoder(w.Body).Decode(&authResponse); err != nil {
			t.Fatalf("Failed to decode auth response: %v", err)
		}
		token, ok := authResponse["token"].(string)
		if !ok || token == "" {
			t.Fatal("Expected token in auth response")
		}
		// Now, use the token to get user info
		w = httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/user", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		var userResponse map[string]interface{}
		if err := json.NewDecoder(w.Body).Decode(&userResponse); err != nil {
			t.Fatalf("Failed to decode user response: %v", err)
		}
		assert.Equal(t, user1_email, userResponse["email"], "Expected email to match")
	})

	t.Run("get user endpoint should fail without token", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/user", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, 401, w.Code)
		var response map[string]interface{}
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		assert.Equal(t, "unauthorized", response["error"], "Expected 'unauthorized' error message")
	})

	t.Run("get user endpoint should fail with invalid token", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/user", nil)
		req.Header.Set("Authorization", "Bearer invalidtoken")
		router.ServeHTTP(w, req)
		assert.Equal(t, 401, w.Code)
		var response map[string]interface{}
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		assert.Equal(t, "unauthorized", response["error"], "Expected 'unauthorized' error message")
	})

	t.Run("auth endpoint should fail when limit is reached for the same email", func(t *testing.T) {
		payload, _ := json.Marshal(map[string]string{
			"email":    user1_email,
			"password": user1_password,
		})
		var w *httptest.ResponseRecorder
		for i := 0; i < 11; i++ {
			w = httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/auth", strings.NewReader(string(payload)))
			req.RemoteAddr = "192.168.1.100:12345"
			router.ServeHTTP(w, req)
		}
		assert.Equal(t, 429, w.Code)
		cleanupDB(t, server.db)
	})

	t.Run("auth endpoint should fail when limit is reached for the same IP", func(t *testing.T) {
		var w *httptest.ResponseRecorder
		for i := 0; i < 26; i++ {
			payload, _ := json.Marshal(map[string]string{
				"email":    fmt.Sprintf("user-%d@example.com", i),
				"password": user1_password,
			})

			w = httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/auth", strings.NewReader(string(payload)))
			req.RemoteAddr = "192.168.1.100:12345"
			router.ServeHTTP(w, req)
		}
		assert.Equal(t, 429, w.Code)
		cleanupDB(t, server.db)
	})

	t.Run("auth endpoint should fail when global limit is reached", func(t *testing.T) {
		var w *httptest.ResponseRecorder
		addr, _ := netip.ParseAddr("192.168.1.1")
		for i := 0; i < 302; i++ {
			payload, _ := json.Marshal(map[string]string{
				"email":    fmt.Sprintf("user-%d@example.com", i),
				"password": user1_password,
			})

			w = httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/auth", strings.NewReader(string(payload)))
			addr = addr.Next()
			req.RemoteAddr = fmt.Sprintf("%s:123", addr.String())
			router.ServeHTTP(w, req)
		}
		assert.Equal(t, 429, w.Code)
		cleanupDB(t, server.db)
	})
}

func setupServer(t *testing.T) *API {
	t.Helper()
	ctx := context.Background()

	log := slog.Default()

	dbPool := db.NewPool(ctx, log)
	if dbPool == nil {
		log.Error("Failed to create database pool")
		os.Exit(1)
	}

	internalDB := db.New(dbPool)

	return NewAPI(log, internalDB)
}

func cleanupDB(t *testing.T, db *db.Queries) {
	t.Helper()
	ctx := context.Background()

	// Clean up the database by deleting all users
	if err := db.DeleteAllUsers(ctx); err != nil {
		t.Fatalf("Failed to clean up database: %v", err)
	}
}
