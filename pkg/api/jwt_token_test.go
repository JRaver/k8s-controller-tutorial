package api

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
)

type TokenResponse struct {
	Token string `json:"token"`
}

func TestTokenHandler_IssuesValidJWT(t *testing.T) {
	JWTSecret = "test-secret"
	ctx := &fasthttp.RequestCtx{}
	TokenHandler(ctx)

	var tokenResp TokenResponse
	err := json.Unmarshal(ctx.Response.Body(), &tokenResp)
	require.NoError(t, err)
	require.NotEmpty(t, tokenResp.Token)

	// Parse token
	parsed, err := jwt.Parse(tokenResp.Token, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWTSecret), nil
	})
	require.NoError(t, err)
	require.True(t, parsed.Valid)
}

func TestJWTMiddleware_ValidAndInvalidToken(t *testing.T) {
	JWTSecret = "test-secret"
	claims := jwt.MapClaims{
		"sub": "testuser",
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(JWTSecret))
	require.NoError(t, err)

	// Valid token
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.Set("Authorization", "Bearer "+tokenStr)
	called := false
	mw := JwtMiddleware(func(ctx *fasthttp.RequestCtx) { called = true })
	mw(ctx)
	require.True(t, called, "middleware should call next handler for valid token")

	// Invalid token
	ctx2 := &fasthttp.RequestCtx{}
	ctx2.Request.Header.Set("Authorization", "Bearer invalidtoken")
	called = false
	mw(ctx2)
	require.False(t, called, "middleware should not call next handler for invalid token")
	require.Equal(t, fasthttp.StatusUnauthorized, ctx2.Response.StatusCode())

	// Expired token
	expiredClaims := jwt.MapClaims{
		"sub": "testuser",
		"exp": time.Now().Add(-time.Hour).Unix(),
	}
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	expiredStr, err := expiredToken.SignedString([]byte(JWTSecret))
	require.NoError(t, err)
	ctx3 := &fasthttp.RequestCtx{}
	ctx3.Request.Header.Set("Authorization", "Bearer "+expiredStr)
	called = false
	mw(ctx3)
	require.False(t, called, "middleware should not call next handler for expired token")
	require.Equal(t, fasthttp.StatusUnauthorized, ctx3.Response.StatusCode())
}
