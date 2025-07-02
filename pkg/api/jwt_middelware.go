package api

import (
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/valyala/fasthttp"
)

var JWTSecret string

func JwtMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		header := string(ctx.Request.Header.Peek("Authorization"))
		if !strings.HasPrefix(header, "Bearer ") {
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			ctx.SetBodyString(`{"error": "Unauthorized"}`)
			return
		}
		tokenSting := strings.TrimPrefix(header, "Bearer ")
		token, err := jwt.Parse(tokenSting, func(token *jwt.Token) (interface{}, error) {
			return []byte(JWTSecret), nil
		})
		if err != nil || !token.Valid {
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			ctx.SetBodyString(`{"error": "Unauthorized"}`)
			return
		}
		next(ctx)
	}
}
