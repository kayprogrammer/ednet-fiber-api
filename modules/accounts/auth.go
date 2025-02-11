package accounts

import (
	"context"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/kayprogrammer/ednet-fiber-api/config"
	"github.com/kayprogrammer/ednet-fiber-api/ent"
	"github.com/kayprogrammer/ednet-fiber-api/ent/token"
	"github.com/kayprogrammer/ednet-fiber-api/ent/user"
)

type AccessTokenPayload struct {
	UserId   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	jwt.RegisteredClaims
}

type RefreshTokenPayload struct {
	Data string `json:"data"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userId uuid.UUID, userName string) string {
	cfg := config.GetConfig()
	expirationTime := time.Now().Add(time.Duration(cfg.AccessTokenExpireMinutes) * time.Minute)
	payload := AccessTokenPayload{
		UserId:   userId,
		Username: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	// Create the JWT string
	tokenString, err := token.SignedString(cfg.SecretKeyByte)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		log.Fatal("Error Generating Access token: ", err)
	}
	return tokenString
}

func GenerateRefreshToken() string {
	cfg := config.GetConfig()
	expirationTime := time.Now().Add(time.Duration(cfg.RefreshTokenExpireMinutes) * time.Minute)
	payload := RefreshTokenPayload{
		Data: config.GetRandomString(10),
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	// Create the JWT string
	tokenString, err := token.SignedString(cfg.SecretKeyByte)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		log.Fatal("Error Generating Refresh token: ", err)
	}
	return tokenString
}

func DecodeAccessToken(db *ent.Client, ctx context.Context, tokenStr string) (*ent.User, *string) {
	cfg := config.GetConfig()
	claims := &AccessTokenPayload{}

	tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return cfg.SecretKeyByte, nil
	})
	tokenErr := "Auth Token is Invalid or Expired!"
	if err != nil {
		return nil, &tokenErr
	}
	if !tkn.Valid {
		return nil, &tokenErr
	}

	// Fetch User model object
	userId := claims.UserId
	user, _ := db.User.Query().Where(user.ID(userId), user.HasTokensWith(token.Access(tokenStr))).
		Only(ctx)

	if user == nil {
		return nil, &tokenErr
	}
	return user, nil
}

func DecodeRefreshToken(db *ent.Client, ctx context.Context, tokenStr string) *ent.User {
	cfg := config.GetConfig()

	claims := &RefreshTokenPayload{}
	tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return cfg.SecretKeyByte, nil
	})
	if err != nil {
		return nil
	}
	if !tkn.Valid {
		log.Println("Invalid Refresh Token")
		return nil
	}
	user, _ := db.User.Query().Where(user.HasTokensWith(token.Refresh(tokenStr))).
		Only(ctx)
	if user == nil {
		return nil
	}
	return user
}