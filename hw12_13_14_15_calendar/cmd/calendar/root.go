package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/bubblesupreme/otus_go_homework/hw12_13_14_15_calendar/internal/server"
	internalgrpc "github.com/bubblesupreme/otus_go_homework/hw12_13_14_15_calendar/internal/server/grpc"

	"github.com/bubblesupreme/otus_go_homework/hw12_13_14_15_calendar/internal/app"
	internalhttp "github.com/bubblesupreme/otus_go_homework/hw12_13_14_15_calendar/internal/server/http"
	"github.com/bubblesupreme/otus_go_homework/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/bubblesupreme/otus_go_homework/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/bubblesupreme/otus_go_homework/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	log "github.com/bubblesupreme/otus_go_homework/hw12_13_14_15_calendar/internal/logger"
	homedir "github.com/mitchellh/go-homedir"
)

const (
	layoutTime = "01-02-2006-15-04-05"
)

var (
	cfgFile        string
	ErrStorageType = errors.New("storage type is not supported")
	ErrServerType  = errors.New("server type is not supported")
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "calendar",
	Short: "Calendar servis to manage events",
	Long: `The "Calendar" service is the most simplified service for storing calendar events and sending notifications.

	The service assumes the ability to:
		* add / update an event;
		* get a list of events for the day / week / month;
		* receive notification N days before the event.`,
	Run: run,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Calendar",
	Long:  `All software has versions. This is Calendar's`,
	Run: func(_ *cobra.Command, _ []string) {
		printVersion()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/calendar.json)")
	rootCmd.AddCommand(versionCmd)
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name "calendar.json".
		viper.AddConfigPath(home)
		viper.SetConfigName("calendar.json")
	}

	viper.SetConfigType("json")
	viper.AutomaticEnv()
	viper.BindEnv("dblogin", "POSTGRES_USER")
	viper.BindEnv("dbname", "POSTGRES_DB")
	viper.BindEnv("dbpassword", "POSTGRES_PASSWORD")
	viper.BindEnv("dbport", "POSTGRES_PORT")
	viper.BindEnv("dbhost", "POSTGRES_HOST")

	readConfig()

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config was changed: ", e.Name)
		readConfig()
	})
}

func readConfig() {
	if err := viper.ReadInConfig(); err == nil {
		log.Info("using config file: "+viper.ConfigFileUsed(), log.Fields{"settings": viper.AllSettings()})
	} else {
		log.Fatal("failed to read config file: "+err.Error(), nil)
	}
}

func run(_ *cobra.Command, _ []string) {
	initConfig()

	if err := func() error {
		config, err := NewConfig()
		if err != nil {
			log.Error("failed to read config", nil)
			return err
		}

		logF, err := configureLogger(config)
		if err != nil {
			log.Error("failed to configure logger", nil)
			return err
		}
		defer func() {
			if err := logF.Close(); err != nil {
				log.Fatal(fmt.Sprintf("failed to close log file: %s", err.Error()), nil)
			}
		}()

		storage, err := newStorage(config.Storage.Type, config.DataBase)
		if err != nil {
			log.Error("failed to create storage", log.Fields{"type": config.Storage.Type})
			return err
		}

		calendarApp := app.NewApp(storage)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		server, err := newServer(ctx, calendarApp, config.Server)
		if err != nil {
			log.Error("failed to create server", log.Fields{"type": config.Storage.Type})
			return err
		}

		go func() {
			signals := make(chan os.Signal, 1)
			signal.Notify(signals, syscall.SIGINT, syscall.SIGHUP)

			select {
			case <-ctx.Done():
				return
			case <-signals:
			}

			signal.Stop(signals)
			cancel()

			if err := server.Stop(); err != nil {
				log.Error("failed to stop http server: "+err.Error(), nil)
			}
		}()

		log.Info("calendar is running...", nil)
		if err := server.Start(); err != nil {
			log.Error("failed to start http server", nil)
			cancel()
			return err
		}

		return nil
	}(); err != nil {
		log.Fatal(err.Error(), nil)
	}
}

func newStorage(t string, c DBMSConf) (storage.Storage, error) {
	switch t {
	case "memory":
		return memorystorage.NewStorage(), nil
	case "sql":
		return sqlstorage.NewStorage(c.Login, c.Password, c.Port, c.Host, c.DBName)
	default:
		return nil, ErrStorageType
	}
}

func newServer(ctx context.Context, calendarApp *app.App, c ServerConf) (server.Server, error) {
	switch c.Type {
	case "grpc":
		return internalgrpc.NewServer(ctx, calendarApp, c.Port, c.Host)
	case "http":
		return internalhttp.NewServer(ctx, calendarApp, c.Port, c.Host)
	default:
		return nil, ErrServerType
	}
}

func configureLogger(c Config) (*os.File, error) {
	log.SetLevel(c.Logger.Level)

	fileName := fmt.Sprint("calendar", time.Now().Format(layoutTime), ".log")
	if err := os.MkdirAll(c.Logger.Path, os.ModePerm); err != nil {
		log.Error("failed to create log directory", log.Fields{
			"path": c.Logger.Path,
		})
		return nil, err
	}
	f, err := os.OpenFile(path.Join(c.Logger.Path, fileName), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o755)
	if err != nil {
		log.Error(fmt.Sprintf("failed to create log file: %s", err.Error()), log.Fields{
			"path":      c.Logger.Path,
			"file name": fileName,
		})
		return nil, err
	}
	// print to file and stdout
	log.SetOutput(io.MultiWriter(f, os.Stdout))

	return f, nil
}
