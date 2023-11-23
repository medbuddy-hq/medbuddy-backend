package middleware

import (
	"github.com/gin-gonic/gin"
	"medbuddy-backend/internal/constant"
	"medbuddy-backend/internal/model"
	"medbuddy-backend/utility"
	"net/http"
	"strings"
)

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

// Practitioner is the middleware for admin-only endpoints
func Practitioner() gin.HandlerFunc {
	return func(c *gin.Context) {

		token := getToken(c.Request)
		if token == "" {
			rd := utility.BuildErrorResponse(http.StatusUnauthorized, constant.StatusFailed,
				constant.ErrUnauthorized, "no token specified", nil)
			c.JSON(http.StatusUnauthorized, rd)
			return
		}

		claims, err := VerifyToken(token)
		if err != nil {
			rd := utility.BuildErrorResponse(http.StatusUnauthorized, constant.StatusFailed,
				constant.ErrUnauthorized, err.Error(), nil)
			c.JSON(http.StatusUnauthorized, rd)
			return
		}

		// Check if request is from a practitioner
		if claims.Role != constant.Roles[constant.Practitioner] {
			rd := utility.BuildErrorResponse(http.StatusUnauthorized, constant.StatusFailed,
				constant.ErrUnauthorized, "cannot access this endpoint", nil)
			c.JSON(http.StatusUnauthorized, rd)
			return
		}

		// Set details from token into context and execute next handler
		c.Set("user info", &model.ContextInfo{
			ID:    claims.ID,
			Role:  claims.Role,
			Email: claims.Email,
		})
		c.Next()
	}
}

// Patient is the middleware for patient endpoints
func Patient() gin.HandlerFunc {
	return func(c *gin.Context) {

		token := getToken(c.Request)
		if token == "" {
			rd := utility.BuildErrorResponse(http.StatusUnauthorized, constant.StatusFailed,
				constant.ErrUnauthorized, "no token specified", nil)
			c.JSON(http.StatusUnauthorized, rd)
			return
		}

		claims, err := VerifyToken(token)
		if err != nil {
			rd := utility.BuildErrorResponse(http.StatusUnauthorized, constant.StatusFailed,
				constant.ErrUnauthorized, err.Error(), nil)
			c.JSON(http.StatusUnauthorized, rd)
			return
		}

		// Check if request is from a patient
		if claims.Role != constant.Roles[constant.Patient] {
			rd := utility.BuildErrorResponse(http.StatusUnauthorized, constant.StatusFailed,
				constant.ErrUnauthorized, "cannot access this endpoint", nil)
			c.JSON(http.StatusUnauthorized, rd)
			return
		}

		// Set details from token in context and execute next handler
		c.Set("user info", &model.ContextInfo{
			ID:    claims.ID,
			Role:  claims.Role,
			Email: claims.Email,
		})
		c.Next()
	}
}

// Generic is the middleware for patient endpoints
func Generic() gin.HandlerFunc {
	return func(c *gin.Context) {

		token := getToken(c.Request)
		if token == "" {
			rd := utility.BuildErrorResponse(http.StatusUnauthorized, constant.StatusFailed,
				constant.ErrUnauthorized, "no token specified", nil)
			c.JSON(http.StatusUnauthorized, rd)
			return
		}

		claims, err := VerifyToken(token)
		if err != nil {
			rd := utility.BuildErrorResponse(http.StatusUnauthorized, constant.StatusFailed,
				constant.ErrUnauthorized, err.Error(), nil)
			c.JSON(http.StatusUnauthorized, rd)
			return
		}

		// Set details from token in context and execute next handler
		c.Set("user info", &model.ContextInfo{
			ID:    claims.ID,
			Role:  claims.Role,
			Email: claims.Email,
		})
		c.Next()
	}
}

// getToken contains logic to fetch token from headers
func getToken(r *http.Request) (token string) {
	auth := r.Header.Get("Authorization")
	hToken := r.Header.Get("Token") // header token
	if auth == "" {
		if hToken == "" {
			return
		} else {
			token = hToken
		}
	} else {
		// Split Authorization to get bearer token
		strs := strings.Split(auth, " ")
		if len(strs) > 1 {
			token = strs[1]
		} else {
			token = auth
		}
	}

	return
}
