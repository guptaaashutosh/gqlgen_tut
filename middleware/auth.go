package middleware

import (
	"context"
	"fmt"
	"net/http"
)

func AuthenticateUser() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			fmt.Println("header : ",header)

			if header == "" {
				next.ServeHTTP(w, r)
				return
			}
	})
}
}

func UserMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authToken := r.Header.Get("Authorization")
			authRole := r.Header.Get("Role")

			var ctx = context.WithValue(r.Context(),"token",authToken)
			ctx = context.WithValue(ctx,"role",authRole)

			next.ServeHTTP(w, r.WithContext(ctx))
	})
}
}




