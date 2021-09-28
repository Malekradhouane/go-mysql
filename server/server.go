package server

import (
	"context"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ory/viper"
	"github/malekradhouane/test-cdi/route"
)

//Server is the http layer of the app
type Server struct {
	Cfg         *Config
	UserActions *route.UserActions
}

//Config server config
type Config struct {
	Port           string
	AuthMiddleware *jwt.GinJWTMiddleware
}

//NewConfig constructs a new config
func NewConfig(authMiddleware *jwt.GinJWTMiddleware) *Config {
	viper.BindEnv("PORT", "PORT")
	port := viper.GetString("PORT")
	if port == "" {
		port = "8080"
	}
	return &Config{
		Port:           port,
		AuthMiddleware: authMiddleware,
	}
}

func NewServer(cfg *Config,
	uas *route.UserActions,
) *Server {
	return &Server{Cfg: cfg,
		UserActions: uas,
	}
}

//Run setup the app with all dependencies
func (s *Server) Run() error {
	r := gin.Default()
	corsCfg := cors.DefaultConfig()
	corsCfg.AddAllowHeaders("Authorization")
	corsCfg.AllowAllOrigins = true
	r.Use(cors.New(corsCfg))
	r.GET("/status", route.GetStatus)
	auth := r.Group("/api/auth")
	{
		auth.POST("", s.Cfg.AuthMiddleware.LoginHandler)

	}
	auth.Use(s.Cfg.AuthMiddleware.MiddlewareFunc())

	users := r.Group("/api/users")
	{
		users.POST("", s.UserActions.CreateUser)
		users.GET("/:id", s.Cfg.AuthMiddleware.MiddlewareFunc(), func(ctx *gin.Context) {
			if ctx.Param("id") == "list" {
				s.UserActions.ListUsers(ctx)
				return
			}
			s.UserActions.GetUser(ctx)
		})
		users.DELETE("/:id", s.Cfg.AuthMiddleware.MiddlewareFunc(), s.UserActions.DeleteUser)
		users.PATCH("/:id", s.Cfg.AuthMiddleware.MiddlewareFunc(), s.UserActions.UpdateUser)

	}

	return r.Run(":" + s.Cfg.Port)
}

//Shutdown shutdowns the server
func (s *Server) Shutdown(ctx context.Context) {
	//Clean up here
}
