package cmd

import (
	"context"
	"errors"
	"fmt"
	"github/malekradhouane/test-cdi/route"
	"github/malekradhouane/test-cdi/server"
	"github/malekradhouane/test-cdi/service"
	"github/malekradhouane/test-cdi/store"
	"github/malekradhouane/test-cdi/store/mysql"
	"net/http"
	"os"

	//import pg dialect
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var rootCmd = &cobra.Command{
	Use:  "ring-over",
	RunE: start,
}


func start(cmd *cobra.Command, args []string) error {
	viper.AutomaticEnv()
	logger := getLogger()
	logger.Debug("Running ring-over API...")

	ctx, cancel := context.WithCancel(context.Background())

	dbCfg := dbConfig()
	if dbCfg.DB == "" || dbCfg.User == "" || dbCfg.Password == "" {
		return onError(logger, cancel, "configuring the store", errors.New("DB name and credentials (user and password) are required"))
	}

	logger.Debug("Running ring-over API...")
	models := []interface{}{
		&store.Todo{},
	}
	db, err := mysql.NewClient(dbCfg, models)
	if err != nil {
		return onError(logger, cancel, "store client", err)
	}
	logger.Debug("connected to DB")



	logger.Debug("connecting to store",
		zap.String("host", dbCfg.Host),
		zap.Int("port", dbCfg.Port),
		zap.String("DB", dbCfg.DB),
		)


	logger.Debug("connected to DB")

	todoService := service.NewTodoService(db)
	todoActions := route.NewTodoActions(todoService)
	serverCfg := server.NewConfig()
	server := server.NewServer(serverCfg, todoActions)
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
		db.Shutdown()
		return onError(logger, cancel, "runtime", err)
	case <-ctx.Done():
		server.Shutdown(ctx)
		db.Shutdown()
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
	f.Int("db_port", 3306, "DB port")
	f.String("db_name", "", "DB name")
	f.String("db_user", "", "DB user")
	f.String("db_pwd", "", "DB password")
	viper.BindPFlags(f)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func dbConfig() *mysql.Config {
	viper.BindEnv("db_host", "db_host")
	viper.BindEnv("db_port", "db_port")
	viper.BindEnv("db_user", "db_user")
	viper.BindEnv("db_pwd", "db_pwd")
	viper.BindEnv("db_name", "db_name")

	return &mysql.Config{
		Host:     viper.GetString("db_host"),
		Port:     viper.GetInt("db_port"),
		DB:       viper.GetString("db_name"),
		User:     viper.GetString("db_user"),
		Password: viper.GetString("db_pwd"),
	}
}

func getLogger() *zap.Logger {
	var logger *zap.Logger
	logger, _ = zap.NewDevelopment()
	return logger
}
