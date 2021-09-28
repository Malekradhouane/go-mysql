package cmd

import (
	"context"
	"fmt"
	"github/malekradhouane/test-cdi/auth"
	"github/malekradhouane/test-cdi/route"
	"github/malekradhouane/test-cdi/server"
	"github/malekradhouane/test-cdi/service"
	"github/malekradhouane/test-cdi/store"
	"github/malekradhouane/test-cdi/store/mongo"
	"net/http"
	"os"

	//import pg dialect
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var rootCmd = &cobra.Command{
	Use:  "data-impact",
	RunE: start,
}


func start(cmd *cobra.Command, args []string) error {
	viper.AutomaticEnv()
	logger := getLogger()
	logger.Debug("Running data-impact API...")
	models := []interface{}{
		&store.User{},
	}
	db := mongo.ConnectDB(models)
	dbCfg := dbConfig()

	ctx, cancel := context.WithCancel(context.Background())


	logger.Debug("connecting to store",
		zap.String("host", dbCfg.Host),
		zap.Int("port", dbCfg.Port),
		zap.String("DB", dbCfg.DB),
		)


	logger.Debug("connected to DB")

	userService := service.NewUserService(db)
	userActions := route.NewUserActions(userService)

	authCfg := auth.NewConfig()
	authMiddleware, _ := auth.NewAuthMiddleware(authCfg, db)
	serverCfg := server.NewConfig(authMiddleware)


	server := server.NewServer(serverCfg, userActions)
	errChan := make(chan error)
	go func() {
		logger.Debug("running the server")
		if err := server.Run(); err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		server.Shutdown(ctx)
		db.Client.Disconnect(ctx)
		return onError(logger, cancel, "runtime", err)
	case <-ctx.Done():
		server.Shutdown(ctx)
		db.Client.Disconnect(ctx)
		cancel()
		return nil
	}

}

func onError(log *zap.Logger, cancel context.CancelFunc, failingStep string, err error) error {
	log.Error(fmt.Sprintf("%s error", failingStep),
		zap.Error(err),
	)
	cancel()
	os.Exit(1)
	return nil
}

//Execute launches the app
func Execute() {
	viper.AutomaticEnv()
	f := rootCmd.Flags()
	f.String("db_host", "localhost", "DB host")
	f.Int("db_port", 27017, "DB port")
	f.String("db_name", "", "DB name")
	f.String("mongo_uri", "", "DB URI")
	f.String("db_user", "", "DB user")
	f.String("db_pwd", "", "DB password")
	viper.BindPFlags(f)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func dbConfig() *mongo.Config {
	viper.BindEnv("db_host", "db_host")
	viper.BindEnv("db_port", "db_port")
	viper.BindEnv("db_user", "db_user")
	viper.BindEnv("mongo_uri", "mongo_uri")
	viper.BindEnv("db_pwd", "db_pwd")
	viper.BindEnv("db_name", "db_name")
	viper.BindEnv("db_socket", "db_socket")
	viper.BindEnv("env", "env")

	return &mongo.Config{
		Host:     viper.GetString("db_host"),
		Port:     viper.GetInt("db_port"),
		URI:      viper.GetString("mongo_uri"),
		DB:       viper.GetString("db_name"),
		User:     viper.GetString("db_user"),
		Password: viper.GetString("db_pwd"),
		Socket:   viper.GetString("db_socket"),
		Env:      viper.GetString("env"),
	}
}

func getLogger() *zap.Logger {
	var logger *zap.Logger
	logger, _ = zap.NewDevelopment()
	return logger
}
