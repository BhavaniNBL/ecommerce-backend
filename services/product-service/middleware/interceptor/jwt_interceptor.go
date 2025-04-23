package interceptor

import (
	"context"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Secret key used to sign JWT (should be loaded from config/env in production)
var jwtSecret = []byte("your_super_secret_key")

// Claims defines the structure for JWT claims
type Claims struct {
	UserID   string `json:"user_id"`
	UserType string `json:"user_type"`
	jwt.RegisteredClaims
}

// VerifyToken parses and validates the JWT token string
func VerifyToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, status.Error(codes.Unauthenticated, "Invalid token")
}

// AuthInterceptor is a gRPC interceptor for authenticating JWTs
func AuthInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "Missing metadata")
		}

		authHeader := md["authorization"]
		if len(authHeader) == 0 {
			return nil, status.Error(codes.Unauthenticated, "Authorization token not provided")
		}

		tokenString := strings.TrimPrefix(authHeader[0], "Bearer ")
		claims, err := VerifyToken(tokenString)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "Invalid token: %v", err)
		}

		// Add user info to context for downstream handlers
		ctx = context.WithValue(ctx, "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "user_type", claims.UserType)

		return handler(ctx, req)
	}
}
