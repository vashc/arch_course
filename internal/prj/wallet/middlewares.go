package wallet

import (
	"arch_course/internal/prj"
	"log"

	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/golang-jwt/jwt"
)

// MiddlewareUserAuth is an authorization middleware
func MiddlewareUserAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := prj.ExtractToken(r)
		if err != nil {
			log.Printf("prj.ExtractToken error: %s\n", err.Error())
			code := http.StatusInternalServerError
			if errors.Is(err, prj.ErrUnathorizedUser) {
				code = http.StatusUnauthorized
			}

			http.Error(w, http.StatusText(code), code)
			return
		}

		//nolint:staticcheck // It's ok for now
		ctx := context.WithValue(r.Context(), prj.CtxAuthToken, token)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// MiddlewareUserCtx is a middleware for getting user context
func MiddlewareUserCtx(config *Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var err error
			tokenString := jwtauth.TokenFromHeader(r)

			var token *jwt.Token
			token, err = jwt.ParseWithClaims(
				tokenString,
				&JWTClaims{},
				func(token *jwt.Token) (interface{}, error) {
					// Algorithm type validation
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("%w: %v", prj.ErrWrongSignMethod, token.Header["alg"])
					}

					return []byte(config.JWTSecret), nil
				},
			)
			if err != nil {
				log.Printf("jwt.ParseWithClaims error: %s\n", err.Error())
				code := http.StatusInternalServerError
				http.Error(w, http.StatusText(code), code)
				return
			}

			if !token.Valid {
				log.Println("token is not valid")
				code := http.StatusInternalServerError
				http.Error(w, http.StatusText(code), code)
				return
			}

			claims, ok := token.Claims.(*JWTClaims)
			if !ok {
				log.Println("token.Claims.(*JWTClaims) is not ok")
				code := http.StatusInternalServerError
				http.Error(w, http.StatusText(code), code)
				return
			}

			//nolint:staticcheck // It's ok for now
			ctx := context.WithValue(r.Context(), prj.RequestParamUserID, claims.UserID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// MiddlewareUserPermission is a middleware to check if the logged user has permissions
// to read/write/delete requested user data
func MiddlewareUserPermission(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userID, ok := ctx.Value(prj.RequestParamUserID).(int64)
		if !ok {
			code := http.StatusUnauthorized
			http.Error(w, http.StatusText(code), code)
			return
		}

		requestedUserID, err := strconv.ParseInt(chi.URLParam(r, prj.RequestParamUserID), 10, 64)
		if err != nil {
			log.Printf("strconv.ParseInt error: %s\n", err.Error())
			code := http.StatusInternalServerError
			http.Error(w, http.StatusText(code), code)
			return
		}

		if userID != requestedUserID {
			code := http.StatusForbidden
			http.Error(w, http.StatusText(code), code)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
