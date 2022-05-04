package server

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ory/viper"
	"github/malekradhouane/test-cdi/route"
)

//Server is the http layer of the app
type Server struct {
	Cfg         *Config
	TodoActions *route.TodoActions
}

//Config server config
type Config struct {
	Port           string
}

//NewConfig constructs a new config
func NewConfig() *Config {
	viper.BindEnv("PORT", "PORT")
	port := viper.GetString("PORT")
	if port == "" {
		port = "8080"
	}
	return &Config{
		Port:           port,
	}
}

func NewServer(cfg *Config,
	td *route.TodoActions,
) *Server {
	return &Server{Cfg: cfg,
		TodoActions: td,
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


	todolists := r.Group("/todolists")
	{
		todolists.GET("", s.TodoActions.Get)
		todolists.POST("",  s.TodoActions.Create)
		todolists.DELETE("/:id", s.TodoActions.Delete)

	}

	return r.Run(":" + s.Cfg.Port)
}

//Shutdown shutdowns the server
func (s *Server) Shutdown(ctx context.Context) {
	//Clean up here
}
