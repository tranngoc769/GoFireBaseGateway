package middleware

import (
	"context"
	"errors"
	"fmt"
	authUtil "go-firebase-gateway/common/auth"
	"go-firebase-gateway/common/log"
	"go-firebase-gateway/common/model"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/shaj13/go-guardian/v2/auth"
	"github.com/shaj13/go-guardian/v2/auth/strategies/basic"
	"github.com/shaj13/go-guardian/v2/auth/strategies/token"
	"github.com/shaj13/go-guardian/v2/auth/strategies/union"
	"github.com/shaj13/libcache"
	_ "github.com/shaj13/libcache/fifo"
)

var strategy union.Union
var tokenStrategy auth.Strategy
var cacheObj libcache.Cache

type Response struct {
	Token  string `json:"token"`
	UserId string `json:"user_id"`
}

func SetupGoGuardian() {
	cacheObj = libcache.FIFO.New(0)
	cacheObj.SetTTL(time.Minute * 10)
	cacheObj.RegisterOnExpired(func(key, _ interface{}) {
		cacheObj.Peek(key)
	})
	basicStrategy := basic.NewCached(validateUser, cacheObj)
	tokenStrategy = token.New(verifyJwe, cacheObj)
	strategy = union.New(tokenStrategy, basicStrategy)
}

func validateUser(ctx context.Context, r *http.Request, username, password string) (auth.Info, error) {
	log.Info("AuthMiddleware", "validateUser", "Executing Auth Middleware")

	if username == "goAPI" && password == "123456" {
		return auth.NewDefaultUser("medium", "1", nil, nil), nil
	}

	return nil, fmt.Errorf("Invalid credentials")
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Info("AuthMiddleware", "AuthMiddleware", "Executing Auth Middleware")
		_, user, err := strategy.AuthenticateRequest(c.Request)
		if err != nil {
			c.JSON(
				http.StatusUnauthorized,
				gin.H{
					"error": http.StatusText(http.StatusUnauthorized),
				},
			)
			c.Abort()
			return
		}
		c.Set("user", user)
		log.Info("AuthMiddleware", "Authenticated", user.GetUserName())
	}
}

func verifyJwe(ctx context.Context, r *http.Request, tokenString string) (auth.Info, time.Time, error) {
	accessTokenRes, err := authUtil.CheckAccessToken(tokenString)
	if err != nil {
		return nil, time.Time{}, err
	}
	accessToken, ok := accessTokenRes.(model.AccessToken)

	if !ok {
		return nil, time.Time{}, errors.New("get access token fail")
	}
	token, err := jwt.Parse(accessToken.JWT, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})

	if err != nil {
		return nil, time.Time{}, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		extension := make(map[string][]string)
		// extension["partner_code"] = []string{claims["partner_code"].(string)}
		// extension["level"] = []string{claims["level"].(string)}
		user := auth.NewDefaultUser(claims["sub"].(string), claims["id"].(string), nil, extension)

		return user, time.Time{}, nil
	}

	return nil, time.Time{}, fmt.Errorf("Invalid token")
}

func GetUserId(c *gin.Context) (interface{}, bool) {
	user, isExist := c.Get("user")
	return user.(auth.Info).GetID(), isExist
}

func GetUserLevel(c *gin.Context) (interface{}, bool) {
	user, isExist := c.Get("user")
	extension := user.(auth.Info).GetExtensions()
	return extension["level"][0], isExist
}

func GetPartnerCode(c *gin.Context) (interface{}, bool) {
	user, isExist := c.Get("user")
	extension := user.(auth.Info).GetExtensions()
	return extension["partner_code"][0], isExist
}

func isExpiredToken(f string) bool {
	i, err := strconv.ParseFloat(f, 32)
	if err != nil {
		return true
	}
	if float64(time.Now().Unix()) > i {
		return true
	} else {
		return false
	}
}
