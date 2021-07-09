package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/faceit/test/config"
	"github.com/faceit/test/logger"
	"github.com/faceit/test/notifier"
	"github.com/faceit/test/queue"
	"github.com/faceit/test/services/hasher"
	"github.com/faceit/test/services/password"
	"github.com/faceit/test/services/user"
	"github.com/faceit/test/store"
	countryhandler "github.com/faceit/test/web/country"
	"github.com/faceit/test/web/middleware"
	userhandler "github.com/faceit/test/web/user"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	glog "github.com/google/logger"

	_ "github.com/lib/pq" // postgres driver import
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	context, cancel := context.WithCancel(context.Background())

	go func() {
		oscall := <-c
		log.Printf("system call: %+v \n", oscall)
		cancel()
	}()

	err := loadService(context)
	if err != nil {
		log.Printf("APPLICATION ERROR: %s\n", err.Error())

		return
	}
}

func loadService(ctx context.Context) error {
	cfg, err := config.New()
	if err != nil {
		return err
	}

	log, err := initLog(cfg.Logger())
	if err != nil {
		return err
	}

	postgresClient, err := initDBClient(cfg.DB())
	if err != nil {
		return err
	}

	defer func() {
		err := postgresClient.Close()
		if err != nil {
			log.Errorf(ctx, "failed to close DB, error: %s", err)
		}
	}()

	userStore := store.NewUser(postgresClient, &sql.TxOptions{Isolation: sql.LevelDefault})
	passwordStore := store.NewPassword(postgresClient)
	countryStore := store.NewCountry(postgresClient)

	hasher := hasher.New()
	password := password.New(passwordStore, hasher)
	user := user.New(userStore, hasher)

	notifier := initNotifier(cfg.Notifier(), log)
	queue := queue.New(notifier)

	router := mux.NewRouter().StrictSlash(true)
	middleware := middleware.New(log)
	userhandler.NewHandler(router, log, middleware, user, countryStore, password, hasher, *queue)
	countryhandler.NewHandler(router, log, middleware, countryStore)

	server := &http.Server{
		Addr:    cfg.Service().Port,
		Handler: handlers.CORS()(router),
	}

	errCh := make(chan error)

	go startServer(ctx, log, server, errCh)

	select {
	case <-ctx.Done():
		log.Infof(context.Background(), "stopping service")

		err = server.Shutdown(ctx)
		if errors.Is(err, context.Canceled) {
			return nil
		}
		if err != nil {
			return err
		}

		queue.Closed()
	case err = <-errCh:
		log.Errorf(ctx, "service has been stopped with error: %s", err.Error())
		return err
	}

	queue.Closed()

	return nil

}

func initLog(cfg config.Logger) (logger.Logger, error) {
	f, err := os.Create(fmt.Sprintf("%d_%s", time.Now().UTC().Unix(), cfg.FilePath))
	if err != nil {
		return logger.Logger{}, fmt.Errorf("failed to create log file, error: %w", err)
	}

	l := glog.Init("test", cfg.Verbose, cfg.SystemLog, f)

	return logger.New(l), nil
}

func initDBClient(cfg config.DB) (*sql.DB, error) {
	url := fmt.Sprintf(cfg.PatternURL, cfg.Host, cfg.Port, cfg.UserName, cfg.Password, cfg.Name, cfg.SSLMode)

	db, err := sql.Open(cfg.UserName, url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database, %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("no connection to database, %w", err)
	}

	return db, nil
}

func initNotifier(cfg config.Notifier, l logger.Logger) *notifier.Notifier {
	var consumers []string

	consumers = append(consumers, cfg.OnCreate()...)
	consumers = append(consumers, cfg.OnUpdate()...)
	consumers = append(consumers, cfg.OnDelete()...)

	return notifier.New(nil, consumers, l)
}

func startServer(ctx context.Context, l logger.Logger, server *http.Server, errCh chan<- error) {
	l.Infof(ctx, "starting HTTP listener...")
	err := server.ListenAndServe()
	if err != nil {
		l.Errorf(ctx, "failed to listen and serve, error: %s", err)
		errCh <- err
	}

	l.Infof(ctx, "stopping HTTP listener...")
}
