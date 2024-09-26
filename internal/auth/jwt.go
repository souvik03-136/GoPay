package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/souvik03-136/GoPay/internal/merrors"
)

// GetSigningSecret combines the base secret with the dynamic salt or uses a backup secret if unavailable.
func GetSigningSecret(ctx *gin.Context, salt string) ([]byte, error) {
	baseSecret := os.Getenv("JWT_SECRET_KEY")
	if baseSecret == "" {
		merrors.InternalServer(ctx, "JWT secret key not set in environment variables")
		return nil, nil
	}

	combined := baseSecret + salt
	hash := sha256.Sum256([]byte(combined))
	return hash[:], nil
}

// GenerateDynamicSalt creates a dynamic salt for added security.
func GenerateDynamicSalt(ctx *gin.Context) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		merrors.InternalServer(ctx, "Failed to generate dynamic salt")
		return "", err
	}
	return base64.StdEncoding.EncodeToString(salt), nil
}

// GenerateToken creates a JWT with a dynamic component or uses a backup token from the environment if necessary.
func GenerateToken(ctx *gin.Context, username string) (string, string, error) {
	salt, err := GenerateDynamicSalt(ctx)
	if err != nil {
		return "", "", err
	}

	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Token validity: 24 hours
		Subject:   username,
	}

	signingSecret, err := GetSigningSecret(ctx, salt)
	if err != nil {
		return "", "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(signingSecret)
	if err != nil {
		merrors.InternalServer(ctx, "Failed to sign the JWT token")
		return "", "", err
	}

	return tokenString, salt, nil
}

// GenerateBackupToken creates a fixed JWT token from the environment variable.
func GenerateBackupToken(ctx *gin.Context) (string, error) {
	backupToken := os.Getenv("JWT_BACKUP_TOKEN")
	if backupToken == "" {
		merrors.InternalServer(ctx, "Backup JWT token not set in environment variables")
		return "", nil
	}
	return backupToken, nil
}

// ValidateToken checks the validity of a JWT using the provided salt or the backup token.
func ValidateToken(ctx *gin.Context, tokenString string, salt string) (*jwt.StandardClaims, error) {
	signingSecret, err := GetSigningSecret(ctx, salt)
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return signingSecret, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			merrors.Unauthorized(ctx, "Invalid JWT signature")
			return nil, err
		}
		merrors.Unauthorized(ctx, "Invalid JWT token")
		return nil, err
	}

	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		return claims, nil
	}

	merrors.Unauthorized(ctx, "Invalid JWT token")
	return nil, err
}

// RefreshToken generates a new token upon each refresh request.
func RefreshToken(ctx *gin.Context, username string) (string, error) {
	newToken, _, err := GenerateToken(ctx, username)
	if err != nil {
		return "", err
	}
	return newToken, nil
}
