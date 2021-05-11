package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/OttoWBitt/NEWT/common"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

		tokenString := req.Header.Get("Authorization")
		if len(tokenString) == 0 {
			erro := "Missing Authorization Header"
			common.RenderResponse(res, &erro, http.StatusUnauthorized)
			return
		}

		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

		claims, err := common.DecodeJwt(tokenString)
		if err != nil {
			erro := "Error verifying JWT token: " + err.Error()
			common.RenderResponse(res, &erro, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(req.Context(), "user", claims)
		next.ServeHTTP(res, req.WithContext(ctx))
	})
}
