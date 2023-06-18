package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"youapp/entity"

	"github.com/dgrijalva/jwt-go"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(entity.NewErrorResponse(err.Error()))
			return
		}

		token := fields[1]
		claims := &entity.Claims{}

		tkn, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte("rahasia-perusahaan"), nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(entity.NewErrorResponse(err.Error()))
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(entity.NewErrorResponse(err.Error()))
			return
		}

		if !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(entity.NewErrorResponse(err.Error()))
			return
		}

		claims = tkn.Claims.(*entity.Claims)
		ctx := context.WithValue(r.Context(), "id", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
