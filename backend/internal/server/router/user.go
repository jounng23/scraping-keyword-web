package router

import (
	"errors"
	"net/http"
	"scraping-keyword-web/backend/pkg/config"
	"scraping-keyword-web/backend/pkg/db"
	strutil "scraping-keyword-web/backend/pkg/utils/strutils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

var JWT_SECRET_KEY = []byte(config.ServerConfig().JWTSecretKey)
var JWT_EXPIRATION = 24 * time.Hour

func (r *router) signin(c *gin.Context) {
	var userRequestBody SignInRequestBody
	_ = c.BindJSON(&userRequestBody)

	user, err := r.userRepo.GetUserByUsername(c.Request.Context(), userRequestBody.Username)
	if err != nil {
		log.Error().Msgf("failed to get user by authenticated info due to %v, username = %s", err.Error(), userRequestBody.Username)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if err = strutil.ComparePasswords(user.Password, userRequestBody.Password); err != nil {
		log.Error().Msg("failed to get user by authenticated info due to invalid password")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	token, err := generateToken(user.UserID, user.Username)
	if err != nil {
		log.Error().Msgf("failed to gen token due to %v, username = %s", err.Error(), userRequestBody.Username)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, User{Username: user.Username, Token: token})
}

func (r *router) signup(c *gin.Context) {
	var userRequestBody SignInRequestBody
	_ = c.BindJSON(&userRequestBody)
	newUser, err := r.userRepo.CreateUser(c.Request.Context(), db.User{
		Username: userRequestBody.Username,
		Password: userRequestBody.Password,
	})
	if err != nil {
		log.Error().Msgf("failed to create user due to %v, username = %s", err.Error(), userRequestBody.Username)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	token, err := generateToken(newUser.UserID, newUser.Username)
	if err != nil {
		log.Error().Msgf("failed to gen token due to %v, username = %s", err.Error(), userRequestBody.Username)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, User{Username: newUser.Username, Token: token})
}

func (r *router) verify(c *gin.Context) {
	claims, err := verifyTokenFromRequest(c.Request)
	if err != nil {
		log.Error().Msgf("failed to verify token due to %v", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, claims)
}

func generateToken(userID, username string) (string, error) {
	expirationTime := time.Now().Add(JWT_EXPIRATION)
	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWT_SECRET_KEY)
}

func verifyTokenFromRequest(req *http.Request) (claims *Claims, err error) {
	token, err := req.Cookie("token")
	if err != nil {
		log.Error().Msgf("failed to collect token due to %v", err.Error())
		return
	}
	return verifyAndParseJWTToken(token.Value)
}

func verifyAndParseJWTToken(tokenString string) (claims *Claims, err error) {
	claims = &Claims{}
	tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return JWT_SECRET_KEY, nil
	})
	if err != nil {
		return
	}
	if !tkn.Valid {
		return nil, errors.New("invalid token")
	}
	return
}
