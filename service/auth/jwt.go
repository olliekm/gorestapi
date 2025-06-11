package auth

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/olliekm/gorestapi/config"
	"github.com/olliekm/gorestapi/types"
	"github.com/olliekm/gorestapi/utils"
	"golang.org/x/net/context"
)

type contextKey string

const UserKey contextKey = "userID"

func CreateJWT(secret []byte, userID int) (string, error) {
	expiration := time.Second * time.Duration(config.Envs.JWTExpirationInSeconds)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.Itoa(userID),
		"expiresAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func WithJWTAuth(handlerFunc http.HandlerFunc, store types.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the token
		tokenString := getTokenFromRequest(r)

		// Validate token
		token, err := validateToken(tokenString)
		if err != nil {
			permissionDenied(w)
			return
		}

		if !token.Valid {
			permissionDenied(w)
			return
		}

		// Fetch user ID from token
		claims := token.Claims.(jwt.MapClaims)
		str := claims["userID"].(string)
		userID, _ := strconv.Atoi(str)

		u, err := store.GetUserByID(userID)
		if err != nil {
			permissionDenied(w)
			return
		}

		// set context with user ID
		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, u.ID)
		r = r.WithContext(ctx)

		handlerFunc(w, r)
	}
}

func getTokenFromRequest(r *http.Request) string {
	token := r.Header.Get("Authorization")
	if token == "" {
		return ""
	}
	return token
}

func validateToken(t string) (*jwt.Token, error) {
	return jwt.Parse(t, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.NewValidationError("unexpected signing method", jwt.ValidationErrorSignatureInvalid)
		}

		return []byte(config.Envs.JWTSecret), nil
	})
}

func permissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
}

func GetUserIDFromContext(ctx context.Context) int {
	userID, ok := ctx.Value(UserKey).(int)
	if !ok {
		return -1
	}
	return userID
}
