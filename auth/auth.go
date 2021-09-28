package auth

import (
	"github/malekradhouane/test-cdi/api"
	"github/malekradhouane/test-cdi/errs"
	"github/malekradhouane/test-cdi/store"
	"time"
	"golang.org/x/crypto/bcrypt"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

//Config auth config
type Config struct {
	realm string
	key   string
}

//NewConfig constructs a new auth config
func NewConfig() *Config {
	return &Config{
		realm: "data-impact",
		key:   "secret", //To be changed
	}
}

//NewAuthMiddleware creates the middleware
func NewAuthMiddleware(cfg *Config, users store.UserStore) (*jwt.GinJWTMiddleware, error) {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:            cfg.realm,
		Key:              []byte(cfg.realm),
		Timeout:          time.Hour,
		MaxRefresh:       time.Hour * 24,
		SigningAlgorithm: "HS256",
		IdentityKey:      "identity",
		IdentityHandler:  idHandler,
		Authenticator:    authenticator(users),
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(api.AuthenticatedUser); ok {
				return jwt.MapClaims{
					"email":    v.Email,
					"id":       v.ID,
					"identity": v.ID,
				}
			}
			return jwt.MapClaims{}
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if _, ok := data.(*api.AuthenticatedUser); ok {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	if err != nil {
		return nil, err
	}

	return authMiddleware, nil
}

func idHandler(c *gin.Context) interface{} {
	return ExtractAuthenticated(c)
}

func authenticator(users store.UserStore) func(*gin.Context) (interface{}, error) {
	return func(c *gin.Context) (interface{}, error) {
		var login *api.Login
		if err := c.ShouldBind(&login); err != nil {
			return nil, jwt.ErrMissingLoginValues
		}
		user, err := users.Authenticate(c.Request.Context(), login)
		if err != nil {
			if err == errs.ErrNoSuchEntity || err == bcrypt.ErrMismatchedHashAndPassword {
				return nil, jwt.ErrFailedAuthentication
			}
			return nil, err
		}

		return api.AuthenticatedUser{ Email: user.Email}, nil
	}
}

//ExtractAuthenticated extracts authenticated user from the request
func ExtractAuthenticated(c *gin.Context) *api.AuthenticatedUser {
	claims := jwt.ExtractClaims(c)

	if claims != nil && claims["id"] != nil && claims["email"] != nil {
		return &api.AuthenticatedUser{
			ID:    claims["id"].(string),
			Email: claims["email"].(string),
		}
	}
	return nil
}
