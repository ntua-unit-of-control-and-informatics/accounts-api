package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	auth "euclia.xyz/accounts-api/authentication"
	models "euclia.xyz/accounts-api/models"
)

type IAuth interface {
	AuthMiddleware(h http.HandlerFunc) http.HandlerFunc
}

type AuthImplementation struct {
}

func (a *AuthImplementation) AuthMiddleware(h http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiAll := r.Header.Get("Authorization")
		if apiAll == "" {
			err := models.ErrorReport{
				Message: "Not authorized",
				Status:  http.StatusUnauthorized,
			}
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(err)
			return
		}
		apiKeyAr := strings.Split(apiAll, " ")
		authType := apiKeyAr[0]
		if authType == "Bearer" {
			apiKey := apiKeyAr[1]
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			_, err := auth.Verifier.Verify(ctx, apiKey)
			if err != nil {
				err := models.ErrorReport{
					Message: err.Error(),
					Status:  http.StatusUnauthorized,
				}
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(err)
				return
			}
		}

		h.ServeHTTP(w, r)
	})
}
