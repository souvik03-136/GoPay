package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/souvik03-136/GoPay/internal/merrors"
)

// AuthService handles authentication logic.
type AuthService struct{}

// NewAuthService creates a new instance of AuthService.
func NewAuthService() *AuthService {
	return &AuthService{}
}

// GenerateJWT generates a new JWT token for a user.
func (s *AuthService) GenerateJWT(ctx *gin.Context, username string) {
	token, salt, err := GenerateToken(ctx, username)
	if err != nil {
		merrors.InternalServer(ctx, "Failed to generate JWT token")
		return
	}

	// Return the generated token and salt to the user
	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
		"salt":  salt,
	})
}

// ValidateJWT validates the provided JWT token from the request.
func (s *AuthService) ValidateJWT(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		merrors.Unauthorized(ctx, "Authorization header missing")
		return
	}

	// Extract the token from the "Bearer" prefix
	tokenString, err := extractToken(authHeader)
	if err != nil {
		merrors.Unauthorized(ctx, "Invalid authorization header format")
		return
	}

	// Extract salt from query parameters or headers (adjust based on your logic)
	salt := ctx.Query("salt")
	if salt == "" {
		merrors.BadRequest(ctx, "Salt is missing")
		return
	}

	// Validate the token
	claims, err := ValidateToken(ctx, tokenString, salt)
	if err != nil {
		merrors.Unauthorized(ctx, "Invalid JWT token")
		return
	}

	// Return valid token information
	ctx.JSON(http.StatusOK, gin.H{
		"message":  "Token is valid",
		"username": claims.Subject,
		"expires":  time.Unix(claims.ExpiresAt, 0),
	})
}

// RefreshJWT handles token refresh requests
func (s *AuthService) RefreshJWT(ctx *gin.Context) {
	// Extract the current token's username or claims
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		merrors.Unauthorized(ctx, "Authorization header missing")
		return
	}

	tokenString, err := extractToken(authHeader)
	if err != nil {
		merrors.Unauthorized(ctx, "Invalid authorization header format")
		return
	}

	// Extract salt from query parameters or headers
	salt := ctx.Query("salt")
	if salt == "" {
		merrors.BadRequest(ctx, "Salt is missing")
		return
	}

	// Validate the current token
	claims, err := ValidateToken(ctx, tokenString, salt)
	if err != nil {
		merrors.Unauthorized(ctx, "Invalid JWT token")
		return
	}

	// Generate a new token
	newToken, err := RefreshToken(ctx, claims.Subject)
	if err != nil {
		merrors.InternalServer(ctx, "Failed to refresh the JWT token")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"new_token": newToken,
	})
}

// extractToken extracts the token from the "Bearer <token>" format.
func extractToken(authHeader string) (string, error) {
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", errors.New("invalid authorization header format")
	}
	return parts[1], nil
}
