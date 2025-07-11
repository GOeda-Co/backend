package security

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"net/http"
	"strings"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	grpcsso "github.com/tomatoCoderq/stats/internal/clients/sso/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type contextKey string
const userContextKey contextKey = "authUser"

func GetUserIdFromContext(ctx *gin.Context) (uuid.UUID, error) {
	userClaims, exists := ctx.Get("userClaims")
	if !exists {
		return uuid.UUID{}, fmt.Errorf("user claims do not exist")
	}

	claimsMap, ok := userClaims.(jwt.MapClaims)
	if !ok {
		return uuid.UUID{}, fmt.Errorf("cannot convert claims to map")
	}

	userIdString, ok := claimsMap["uid"].(string)
	if !ok {
		return uuid.UUID{}, fmt.Errorf("cannot get user_id from claims map")
	}

	userId, err := uuid.Parse(userIdString)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("cannot parse uuid")
	}

	return userId, nil
}

type AuthUser struct {
	ID uuid.UUID
	email string
	IsAdmin bool
}


type Security struct {
	PrivateKey      string
	PublicKey       *ecdsa.PublicKey
	ExpirationDelta time.Duration
}

type CustomClaims struct {
	UserID               uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims           // includes exp, nbf, iat, etc.
}

func (s *Security) validateToken(tokenString string) (jwt.MapClaims, error) {
	fmt.Println(s.PrivateKey)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.PrivateKey), nil
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
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			return
		}

		tokenString := parts[1]
		claims, err := s.validateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.Set("userClaims", claims)
		c.Next()
	}
}

func (s *Security) IsAdminMiddleware(client grpcsso.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := GetUserIdFromContext(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		isAdmin, err := client.IsAdmin(c, id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if !isAdmin {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "not admin"})
			return
		}

		c.Set("IsAdmin", isAdmin)
		c.Next()
	}
}

func (s *Security) AuthUnaryInterceptor(ssoClient *grpcsso.Client) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
		}
		authHeaders := md["authorization"]
		if len(authHeaders) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "authorization token is not supplied")
		}
		
		token := strings.TrimPrefix(authHeaders[0], "Bearer ")
		jwtClaimsMap, err := s.validateToken(token)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		}
		
		uid, ok := jwtClaimsMap["uid"].(string)
		if !ok || uid == "" {
			return nil, status.Errorf(codes.Unauthenticated, "invalid user ID in token")
		}
		
		uidUUID, err := uuid.Parse(uid)
		if err != nil {
			fmt.Println("TOK", token)
			return nil, status.Errorf(codes.Internal, "failed during parsing user ID: %v", err)
		}


		email, ok := jwtClaimsMap["email"].(string)
		if !ok || email == "" {
			return nil, status.Errorf(codes.Unauthenticated, "invalid user ID in token")
		}

		// isAdmin, ok := jwtClaimsMap["is_admin"].(bool)
		// if !ok {
		// 	return nil, status.Errorf(codes.Unauthenticated, "invalid user ID in token")
		// }

		authUser := AuthUser {
			ID: uidUUID,
			email: email,
		}	

		ctx = context.WithValue(ctx, userContextKey, authUser)

		return handler(ctx, req)
	}
}


// func New(
//     log *slog.Logger,
//     appSecret string,
//     permProvider PermissionProvider,
// ) func(next http.Handler) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		authHeader := c.GetHeader("Authorization")
// 		if authHeader == "" {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
// 			return
// 		}

// 		parts := strings.Split(authHeader, " ")
// 		fmt.Println("HERE1")
// 		if len(parts) != 2 || parts[0] != "Bearer" {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
// 			return
// 		}

// 		tokenString := parts[1]
// 		claims, err := s.validateToken(tokenString, c)
// 		fmt.Println("HERE2")
// 		if err != nil {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
// 			return
// 		}

// 		isAdmin, err := permProvider.IsAdmin(r.Context(), claims.UID)
//             if err != nil {
//                 log.Error("failed to check if user is admin", sl.Err(err))

//                 ctx := context.WithValue(r.Context(), errorKey, ErrFailedIsAdminCheck)
//                 next.ServeHTTP(w, r.WithContext(ctx))

//                 return
//             }

// 		c.Set("userClaims", claims)
// 		c.Next()
// 	}
// }
