package security

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"strings"

	"os"

	"time"

	// "repeatro/internal/models"
	// "repeatro/internal/repositories"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Security struct {
	PrivateKey      *ecdsa.PrivateKey
	PublicKey       *ecdsa.PublicKey
	ExpirationDelta time.Duration
}

type CustomClaims struct {
	UserID               uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims           // includes exp, nbf, iat, etc.
}

func ReadECDSAPrivateKey(path string) (*ecdsa.PrivateKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading private key file: %w", err)
	}

	block, _ := pem.Decode(data)
	if block == nil || block.Type != "EC PRIVATE KEY" {
		return nil, fmt.Errorf("invalid PEM block type for private key: %s", block.Type)
	}

	privKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parsing private key: %w", err)
	}
	return privKey, nil
}

// ReadECDSAPublicKey loads a public key from a PEM file
func ReadECDSAPublicKey(path string) (*ecdsa.PublicKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading public key file: %w", err)
	}

	block, _ := pem.Decode(data)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("invalid PEM block type for public key: %s", block.Type)
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parsing public key: %w", err)
	}

	pubKey, ok := pubInterface.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an ECDSA public key")
	}

	return pubKey, nil
}

func (s *Security) GetKyes() error {
	privateKey, err := ReadECDSAPrivateKey("../../../private.pem")
	if err != nil {
		panic(err)
		
	}

	publicKey, err := ReadECDSAPublicKey("../../../public.pem")
	if err != nil {
		panic(err)
	}

	s.PrivateKey = privateKey
	s.PublicKey = publicKey

	return nil
}

func (s *Security) EncodeString(input string, user_id uuid.UUID) (string, error) {
	claims := CustomClaims{
		UserID: user_id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.ExpirationDelta)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	token, err := unsignedToken.SignedString(s.PrivateKey)
	if err != nil {
		return "", nil
	}
	return token, nil
}

func (s *Security) DecodeToken(token string) (CustomClaims, error) {
	claimsToGet := &CustomClaims{}
	parser := jwt.NewParser(jwt.WithLeeway(0 * time.Second))
	tokenGot, err := parser.ParseWithClaims(token, claimsToGet, func(token *jwt.Token) (interface{}, error) {
		return &s.PublicKey, nil
	})
	if err != nil {
		return CustomClaims{}, err
	}

	if !tokenGot.Valid {
		return CustomClaims{}, fmt.Errorf("invalid token")
	}

	return *claimsToGet, nil
}

func (s *Security) validateToken(tokenString string, ctx *gin.Context) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return s.PublicKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok {
		return claims, nil
	}

	return nil, jwt.ErrTokenMalformed
}

func (s *Security) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}

		parts := strings.Split(authHeader, " ")
		fmt.Println("HERE1")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			return
		}

		tokenString := parts[1]
		claims, err := s.validateToken(tokenString, c)
		fmt.Println("HERE2")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.Set("userClaims", claims)
		c.Next()
	}
}
